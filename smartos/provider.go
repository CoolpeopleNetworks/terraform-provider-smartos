package smartos

import (
	"net"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema:         providerSchema(),
		ResourcesMap:   providerResources(),
		DataSourcesMap: providerDataSources(),
		ConfigureFunc:  providerConfigure,
	}
}

// List of supported configuration fields for the provider.
// More info in https://github.com/hashicorp/terraform/blob/v0.6.6/helper/schema/schema.go#L29-L142
func providerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"hosts": {
			Type:        schema.TypeMap,
			Required:    true,
			Description: "Host addresses of the SmartOS global zone.",
		},
		"user": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "User to authenticate with.",
		},
		"sshkey": {
			Type:        schema.TypeString,
                        Required:    true,
                        Description: "User's private SSH key.",
                },
	}
}

func providerResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"smartos_machine": resourceMachine(),
	}
}

func providerDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"smartos_image": datasourceImage(),
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	sshSocket := os.Getenv("SSH_AUTH_SOCK")
	agentConnection, err := net.Dial("unix", sshSocket)
	if err != nil {
		return nil, err
	}

	authMethods := []ssh.AuthMethod{}
	authMethods = append(authMethods, ssh.PublicKeysCallback(agent.NewClient(agentConnection).Signers))

	client := SmartOSClient{
		hosts:           d.Get("hosts").(map[string]interface{}),
		user:            d.Get("user").(string),
		agentConnection: agentConnection,
		authMethods:     authMethods,
		clients:         make(map[string]*ssh.Client),
		sshkey:          d.Get("sshkey").(string),
	}

	return &client, nil
}
