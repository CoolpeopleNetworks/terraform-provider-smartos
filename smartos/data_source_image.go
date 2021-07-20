package smartos

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceImage() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		Read:          datasourceImageReadRunc,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func datasourceImageReadRunc(d *schema.ResourceData, m interface{}) error {
	client := m.(*SmartOSClient)

	name := d.Get("name").(string)
	version := d.Get("version").(string)
	nodeName := d.Get("node_name").(string)

	var image *Image
	var err error

	image, err = client.GetLocalImage(nodeName, name, version)
	if err != nil {
		return err
	}

	if image == nil {
		image, err = client.FindRemoteImage(nodeName, name, version)
		if err == nil && image == nil {
			return fmt.Errorf("Image not found")
		}
	}

	if err != nil {
		log.Printf("Failed to retrieve image with name: %s, version: %s.  Error: %s", name, version, err)
		return err
	}

	d.SetId(image.ID.String())

	return nil
}
