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
		"hosts": &schema.Schema{
			Type:        schema.TypeMap,
			Required:    true,
			Description: "Host addresses of the SmartOS global zone.",
		},
		"user": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "User to authenticate with.",
		},
	}
}

func providerResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"smartos_machine": resourceMachine(),
		"smartcluster_machine": resourceMachine(),
	}
}

func providerDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"smartos_image": datasourceImage(),
		"smartcluster_image": datasourceImage(),
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	agentConnection, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		return nil, err
	}

	signers := agent.NewClient(agentConnection).Signers
	authMethods := []ssh.AuthMethod{
		ssh.PublicKeysCallback(signers),
	}

	client := SmartOSClient{
		hosts:           d.Get("hosts").(map[string]interface{}),
		user:            d.Get("user").(string),
		agentConnection: agentConnection,
		authMethods:     authMethods,
		clients:         make(map[string]*ssh.Client),
	}

	return &client, nil
}
