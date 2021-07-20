package smartos

import (
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMachine() *schema.Resource {
	return &schema.Resource{
		Create: resourceMachineCreate,
		Read:   resourceMachineRead,
		Update: resourceMachineUpdate,
		Delete: resourceMachineDelete,

		Schema: map[string]*schema.Schema{
			"serial_code": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"node_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"alias": {
				Type:     schema.TypeString,
				Required: true,
			},
			/*
				"archive_on_delete": {
					Type:     schema.TypeBool,
					Optional: true,
				},
			*/
			"autoboot": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			/*
				"billing_id": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"bhyve_extra_opts": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},
				"boot": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"bootrom": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},
			*/
			"brand": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cpu_cap": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			/*
				"cpu_shares": {
					Type:     schema.TypeInt,
					Optional: true,
				},
					"cpu_type": {
						Type:     schema.TypeString,
						Optional: true,
						ForceNew: true,
					},
			*/
			"customer_metadata": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			/*
				"delegate_dataset": {
					Type:     schema.TypeBool,
					Optional: true,
				},
			*/
			"disks": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"boot": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"compression": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"image_uuid": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"image_size": { // in MiB
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"model": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"size": { // in MiB
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			/*
				"disk_driver": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},
				"do_not_inventory": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"dns_domain": {
					Type:     schema.TypeString,
					Optional: true,
				},
				// "filesystems.*"
				"firewall_enabled": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"flexible_disk_size": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"fs_allowed": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"hostname": {
					Type:     schema.TypeString,
					Optional: true,
				},
			*/
			"image_uuid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			/*
				"internal_metadata": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: {
						Type: schema.TypeString,
					},
				},
				"internal_metadata_namespaces": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: {
						Type: schema.TypeString,
					},
				},
				"indestructible_delegated": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"indestructible_zoneroot": {
					Type:     schema.TypeBool,
					Optional: true,
				},
			*/
			"kernel_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			/*
				"limit_priv": {
					Type:     schema.TypeString,
					Optional: true,
				},
			*/
			"maintain_resolvers": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			/*
				"max_locked_memory": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"max_lwps": {
					Type:     schema.TypeInt,
					Optional: true,
				},
			*/
			"max_physical_memory": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			/*
				"max_swap": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"mdata_exec_timeout": {
					Type:     schema.TypeInt,
					Optional: true,
				},
			*/
			"nics": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_restricted_traffic": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"allow_ip_spoofing": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"allow_mac_spoofing": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"gateways": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"interface": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"ips": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"nic_tag": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"model": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"vlan_id": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"vrrp_vrid": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"vrrp_primary_ip": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			/*
				"nic_driver": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"nowait": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"owner_uuid": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"qemu_opts": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"qemu_extra_opts": {
					Type:     schema.TypeString,
					Optional: true,
				},
			*/
			"primary_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"quota": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ram": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"resolvers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			// "routes.*" - object
			/*
				"spice_opts": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"spice_password": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"spice_port": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"tmpfs": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"uuid": {
					Type:     schema.TypeString,
					Optional: true,
				},
			*/
			"vcpus": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			/*
				"vga": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"virtio_txburst": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"virtio_txtimer": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"vnc_password": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"vnc_port": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"zfs_data_compression": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"zfs_data_recsize": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"zfs_filesystem_limit": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"zfs_io_priority": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"zfs_root_compression": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"zfs_root_recsize": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"zfs_snapshot_limit": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"zlog_max_size": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"zpool": {
					Type:     schema.TypeString,
					Optional: true,
				},
			*/
		},
	}
}

func parseId(id string) (string, uuid.UUID, error) {
	var nodeName string = ""
	var uuidString string = ""

	parts := strings.Split(id, "/")
	if len(parts) != 2 {
		return "", uuid.Nil, fmt.Errorf("ID returned does not contain a /")
	}

	nodeName = parts[0]
	uuidString = parts[1]

	uuid_value, err := uuid.Parse(uuidString)
	if err != nil {
		log.Printf("uuid.Parse failed to parse %s", uuidString)
		return "", uuid.Nil, err
	}

	log.Printf("ID Parse - NodeName: %s, UUID: %s", nodeName, uuidString)
	return nodeName, uuid_value, nil
}

func createId(nodeName string, uuid uuid.UUID) string {
	return fmt.Sprintf("%s/%s", nodeName, uuid.String())
}

func resourceMachineCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("---------------- MachineCreate")
	d.SetId("")

	client := m.(*SmartOSClient)
	machine := Machine{}
	err := machine.LoadFromSchema(d)
	if err != nil {
		return err
	}

	uuid, err := client.CreateMachine(machine.NodeName, &machine)
	if err != nil {
		return err
	}

	d.SetId(createId(machine.NodeName, *uuid))

	err = resourceMachineRead(d, m)
	log.Printf("---------------- MachineCreate (COMPLETE)")
	return err
}

func resourceMachineRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("---------------- MachineRead")
	client := m.(*SmartOSClient)
	nodeName, uuid, err := parseId(d.Id())
	if err != nil {
		log.Printf("Failed to parse incoming ID [%s] - %s", d.Id(), err)
		return err
	}

	machine, err := client.GetMachine(nodeName, uuid)
	if err != nil {
		log.Printf("Failed to retrieve machine with ID %s.  Error: %s", d.Id(), err)
		return err
	}

	err = machine.SaveToSchema(d)
	log.Printf("---------------- MachineRead (COMPLETE)")
	return err
}

func resourceMachineUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("---------------- MachineUpdate")
	nodeName, machineId, err := parseId(d.Id())
	if err != nil {
		return err
	}

	d.Partial(true)

	machineUpdate := Machine{
		ID:       &machineId,
		NodeName: nodeName,
	}

	updatesRequired := false

	if d.HasChange("alias") && !d.IsNewResource() {
		_, newValue := d.GetChange("alias")

		machineUpdate.Alias = newValue.(string)
		updatesRequired = true
	}

	if d.HasChange("autoboot") && !d.IsNewResource() {
		_, newValue := d.GetChange("autoboot")

		machineUpdate.Autoboot = newBool(newValue.(bool))
		updatesRequired = true
	}

	if d.HasChange("cpu_cap") && !d.IsNewResource() {
		_, newValue := d.GetChange("cpu_cap")

		machineUpdate.CPUCap = newUint32(uint32(newValue.(int)))
		updatesRequired = true
	}

	if d.HasChange("customer_metadata") && !d.IsNewResource() {
		oldSchemaValue, newSchemaValue := d.GetChange("customer_metadata")
		oldMap := oldSchemaValue.(map[string]interface{})
		newMap := newSchemaValue.(map[string]interface{})

		var addItem func(key string, value interface{}) = machineUpdate.setCustomerMetadata
		var removeItem func(key string) = machineUpdate.removeCustomerMetadata

		if ReconcileMaps(oldMap, newMap, addItem, addItem, removeItem, stringsAreEqual) {
			updatesRequired = true
		}
	}

	if d.HasChange("maintain_resolvers") && !d.IsNewResource() {
		_, newValue := d.GetChange("maintain_resolvers")

		machineUpdate.MaintainResolvers = newBool(newValue.(bool))
		updatesRequired = true
	}

	if d.HasChange("max_physical_memory") && !d.IsNewResource() {
		_, newValue := d.GetChange("max_physical_memory")

		machineUpdate.MaxPhysicalMemory = newUint32(uint32(newValue.(int)))
		updatesRequired = true
	}

	if d.HasChange("quota") && !d.IsNewResource() {
		_, newValue := d.GetChange("quota")

		machineUpdate.Quota = newUint32(uint32(newValue.(int)))
		updatesRequired = true
	}

	if d.HasChange("resolvers") && !d.IsNewResource() {
		_, newSchemaValue := d.GetChange("resolvers")

		var resolvers []string
		for _, resolver := range newSchemaValue.([]interface{}) {
			resolvers = append(resolvers, resolver.(string))
		}
		machineUpdate.Resolvers = resolvers
		updatesRequired = true
	}

	if d.HasChange("nics") && !d.IsNewResource() {
		_, newSchemaValue := d.GetChange("nics")

		var nics []NetworkInterface
		for _, nic := range newSchemaValue.([]interface{}) {
			nics = append(nics, nic.(NetworkInterface))
		}
		machineUpdate.NetworkInterfaces = nics
		updatesRequired = true
	}

	if updatesRequired {
		client := m.(*SmartOSClient)

		err = client.UpdateMachine(nodeName, &machineUpdate)
		if err != nil {
			return err
		}
	}

	d.Partial(false)
	err = resourceMachineRead(d, m)
	log.Printf("---------------- MachineUpdate (COMPLETE)")
	return err
}

func resourceMachineDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("Request to delete machine with ID: %s\n", d.Id())

	client := m.(*SmartOSClient)
	nodeName, machineId, err := parseId(d.Id())
	if err != nil {
		return err
	}

	return client.DeleteMachine(nodeName, machineId)
}
