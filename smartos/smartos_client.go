package smartos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"regexp"
	"io/ioutil"

	"github.com/google/uuid"
	"golang.org/x/crypto/ssh"
)

type SmartOSClient struct {
	hosts           map[string]interface{}
	user            string
	clients         map[string]*ssh.Client
	agentConnection net.Conn
	authMethods     []ssh.AuthMethod
	sshkey		string
}

func (c *SmartOSClient) Connect(nodeName string) error {
	var err error = nil

	if c.clients[nodeName] != nil {
		return nil
	}

        keyfile, err := ioutil.ReadFile(c.sshkey)
        if err != nil {
                log.Println("SSH: Can't read key: ", err.Error())
        }

        keyparser, err := ssh.ParsePrivateKey(keyfile)
        if err != nil {
                log.Println("SSH: Can't parse key: ", err.Error())
        }

	
	config := &ssh.ClientConfig{
		User:            c.user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(keyparser),},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	log.Printf("SSH: Connecting to host: %s at %s", nodeName, c.hosts[nodeName].(string))
	connectString := fmt.Sprintf("%s:22", c.hosts[nodeName].(string))
	client, err := ssh.Dial("tcp", connectString, config)
	if err != nil {
		log.Println("SSH: Connection failed: ", err.Error())
		return err
	}

	c.clients[nodeName] = client

	log.Println("SSH: Connected successfully")
	return nil
}

func (c *SmartOSClient) Close(nodeName string) {
	if c.clients[nodeName] != nil {
		c.clients[nodeName].Close()
		c.clients[nodeName] = nil
	}
}

func (c *SmartOSClient) CreateMachine(nodeName string, machine *Machine) (*uuid.UUID, error) {
	log.Printf("Creating machine on node: %s", nodeName)

	err := c.Connect(nodeName)
	if err != nil {
		return nil, err
	}

	session, err := c.clients[nodeName].NewSession()
	if err != nil {
		return nil, err
	}

	defer session.Close()

	// Ensure the image has been imported
	if machine.ImageUUID != nil && *machine.ImageUUID != uuid.Nil {
		log.Printf("Ensuring image with UUID %s has been imported", machine.ImageUUID.String())
		err = c.ImportRemoteImage(nodeName, *machine.ImageUUID)
		if err != nil {
			log.Fatalln("Failed to import image for machine.  Error: ", err.Error())
			return nil, err
		}
	} else if machine.Brand == "joyent" || machine.Brand == "lx" {
		log.Fatalln("No image specifiec for OS VM.")
		return nil, fmt.Errorf("No image specifiec for OS VM.")
	}

	// Ensure any disk images are imported
	for _, disk := range machine.Disks {
		if disk.ImageUUID != nil && *disk.ImageUUID != uuid.Nil {
			err = c.ImportRemoteImage(nodeName, *disk.ImageUUID)
			if err != nil {
				log.Fatalf("Failed to import disk image: %s (Error: %s)", disk.ImageUUID.String(), err.Error())
				return nil, err
			}
		}
	}

	json, err := json.Marshal(machine)
	if err != nil {
		log.Fatalln("Failed to create JSON for machine.  Error: ", err.Error())
		return nil, err
	}

	log.Println("JSON: ", string(json))

	session.Stdin = bytes.NewReader(json)

	var b bytes.Buffer
	session.Stderr = &b

	log.Println("SSH execute: vmadm create")
	err = session.Run("vmadm create")
	if err != nil {
		return nil, fmt.Errorf("remote command vmadm failed.  Error: %s (%s)", err, b.String())
	}

	output := b.String()
	log.Printf("Returned data: %s", output)

	re := regexp.MustCompile("Successfully created VM ([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})")
	matches := re.FindStringSubmatch(output)

	if len(matches) != 2 {
		return nil, fmt.Errorf("unrecognized response from vmadm: %s", output)
	}

	log.Println("Matched regex: ", matches[1])
	uuid, err := uuid.Parse(matches[1])
	if err != nil {
		return nil, err
	}

	return &uuid, nil
}

func (c *SmartOSClient) GetMachine(nodeName string, id uuid.UUID) (*Machine, error) {
	err := c.Connect(nodeName)
	if err != nil {
		return nil, err
	}

	session, err := c.clients[nodeName].NewSession()
	if err != nil {
		return nil, err
	}

	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b

	var stderr bytes.Buffer
	session.Stderr = &stderr

	log.Println("SSH execute: vmadm get", id.String())
	err = session.Run("vmadm get " + id.String())
	if err != nil {
		return nil, fmt.Errorf("remote command vmadm failed.  Error: %s (%s)", err, stderr.String())
	}

	outputBytes := b.Bytes()

	output := string(outputBytes)
	log.Printf("Returned data: %s", output)

	var machine Machine
	err = json.Unmarshal(outputBytes, &machine)
	if err != nil {
		log.Printf("Failed to parse returned JSON: %s", err)
		return nil, err
	}

	machine.NodeName = nodeName
	machine.UpdatePrimaryIP()
	machine.UpdateMetadata()

	return &machine, nil
}

func (c *SmartOSClient) UpdateMachine(nodeName string, machine *Machine) error {
	err := c.Connect(nodeName)
	if err != nil {
		return err
	}

	session, err := c.clients[nodeName].NewSession()
	if err != nil {
		return err
	}

	defer session.Close()

	json, err := json.Marshal(machine)
	if err != nil {
		log.Fatalln("Failed to create JSON for machine.  Error: ", err.Error())
	}

	log.Println("JSON: ", string(json))

	session.Stdin = bytes.NewReader(json)

	var b bytes.Buffer
	session.Stderr = &b

	log.Println("SSH execute: vmadm update" + machine.ID.String())
	err = session.Run("vmadm update " + machine.ID.String())
	if err != nil {
		return fmt.Errorf("remote command vmadm failed.  Error: %s (%s)", err, b.String())
	}

	output := b.String()
	log.Printf("Returned data: %s", output)

	return nil
}

func (c *SmartOSClient) DeleteMachine(nodeName string, id uuid.UUID) error {
	err := c.Connect(nodeName)
	if err != nil {
		return err
	}

	session, err := c.clients[nodeName].NewSession()
	if err != nil {
		return err
	}

	defer session.Close()

	var b bytes.Buffer
	session.Stderr = &b

	log.Println("SSH execute: vmadm delete ", id.String())
	err = session.Run("vmadm delete " + id.String())
	if err != nil {
		return fmt.Errorf("remote command vmadm failed.  Error: %s", err)
	}

	output := b.String()
	log.Printf("Returned data: %s", output)

	re := regexp.MustCompile("Successfully deleted VM ([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})")
	matches := re.FindStringSubmatch(output)

	if len(matches) != 2 {
		return fmt.Errorf("unrecognized response from vmadm: %s", output)
	}

	return nil
}

func (c *SmartOSClient) GetLocalImage(nodeName string, name string, version string) (*Image, error) {
	err := c.Connect(nodeName)
	if err != nil {
		return nil, err
	}

	session, err := c.clients[nodeName].NewSession()
	if err != nil {
		return nil, err
	}

	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b

	var stderr bytes.Buffer
	session.Stderr = &stderr

	command := fmt.Sprintf("imgadm list -j name=%s version=%s", name, version)
	err = session.Run(command)
	if err != nil {
		return nil, fmt.Errorf("remote command vmadm failed.  Error: %s (%s)", err, stderr.String())
	}

	outputBytes := b.Bytes()

	output := string(outputBytes)
	log.Printf("Returned data: %s", output)

	var images []map[string]interface{}
	err = json.Unmarshal(outputBytes, &images)
	if err != nil {
		log.Printf("Failed to parse returned JSON: %s", err)
		return nil, err
	}

	if len(images) == 0 {
		return nil, nil
	}

	imageInfo := images[0]
	manifest := imageInfo["manifest"].(map[string]interface{})

	image := Image{}
	image.Name = manifest["name"].(string)
	image.Version = manifest["version"].(string)

	imageID, err := uuid.Parse(manifest["uuid"].(string))
	if err != nil {
		log.Printf("Failed to parse uuid: %s", err)
		return nil, err
	}
	image.ID = &imageID

	return &image, nil
}

func (c *SmartOSClient) FindRemoteImage(nodeName string, name string, version string) (*Image, error) {
	err := c.Connect(nodeName)
	if err != nil {
		return nil, err
	}

	session, err := c.clients[nodeName].NewSession()
	if err != nil {
		return nil, err
	}

	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b

	var stderr bytes.Buffer
	session.Stderr = &stderr

	command := fmt.Sprintf("imgadm avail -j name=%s version=%s", name, version)
	err = session.Run(command)
	if err != nil {
		return nil, fmt.Errorf("remote command vmadm failed.  Error: %s (%s)", err, stderr.String())
	}

	outputBytes := b.Bytes()

	output := string(outputBytes)
	log.Printf("Returned data: %s", output)

	var images []map[string]interface{}
	err = json.Unmarshal(outputBytes, &images)
	if err != nil {
		log.Printf("Failed to parse returned JSON: %s", err)
		return nil, err
	}

	if len(images) == 0 {
		return nil, nil
	}

	imageInfo := images[0]
	manifest := imageInfo["manifest"].(map[string]interface{})

	image := Image{}
	image.Name = manifest["name"].(string)
	image.Version = manifest["version"].(string)

	imageID, err := uuid.Parse(manifest["uuid"].(string))
	if err != nil {
		log.Printf("Failed to parse uuid: %s", err)
		return nil, err
	}
	image.ID = &imageID

	return &image, nil
}

func (c *SmartOSClient) ImportRemoteImage(nodeName string, uuid uuid.UUID) error {
	err := c.Connect(nodeName)
	if err != nil {
		return err
	}

	session, err := c.clients[nodeName].NewSession()
	if err != nil {
		return err
	}

	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b

	var stderr bytes.Buffer
	session.Stderr = &stderr

	log.Printf("Importing image with UUID: %s\n", uuid.String())

	command := fmt.Sprintf("imgadm import %s", uuid.String())
	err = session.Run(command)
	if err != nil {
		return fmt.Errorf("remote command vmadm failed.  Error: %s (%s)", err, stderr.String())
	}

	outputBytes := b.Bytes()

	output := string(outputBytes)
	log.Printf("Returned data: %s", output)

	return nil
}

func (c *SmartOSClient) GetImage(nodeName string, name string, version string) (*Image, error) {
	image, err := c.GetLocalImage(nodeName, name, version)
	if err != nil {
		return nil, err
	}

	if image == nil {
		image, err = c.FindRemoteImage(nodeName, name, version)
		if err != nil {
			return nil, err
		}

		err = c.ImportRemoteImage(nodeName, *image.ID)
		if err != nil {
			return nil, err
		}
	}

	return image, nil
}
