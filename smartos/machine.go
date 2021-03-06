package smartos

import (
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Machine struct {
	NodeName string
	ID       *uuid.UUID `json:"uuid,omitempty"`
	Alias    string     `json:"alias,omitempty"`
	Autoboot *bool      `json:"autoboot,omitempty"`
	Brand    string     `json:"brand,omitempty"`
	CPUCap   *uint32    `json:"cpu_cap,omitempty"`
	/*
		CPUShares                  uint32             `json:"cpu_shares,omitempty"`
	*/
	CustomerMetadata       map[string]string `json:"customer_metadata,omitempty"`
	SetCustomerMetadata    map[string]string `json:"set_customer_metadata,omitempty"`    // for updates
	RemoveCustomerMetadata []string          `json:"remove_customer_metadata,omitempty"` // for updates

	Disks []Disk `json:"disks,omitempty"`

	/*
		DelegateDataset            bool               `json:"delegate_dataset,omitempty"`
		DNSDomain                  string             `json:"dns_domain,omitempty"`
		FirewallEnabled            bool               `json:"firewall_enabled,omitempty"`
		Hostname                   string             `json:"hostname,omitempty"`
	*/
	ImageUUID *uuid.UUID `json:"image_uuid,omitempty"`
	/*
		InternalMetadata           map[string]string  `json:"internal_metadat,omitempty"`
		InternalMetadataNamespaces map[string]string  `json:"internal_metadata_namespaces,omitempty"`
		IndestructableDelegated    bool               `json:"indestructible_delegated,omitempty"`
		IndestructableZoneRoot     bool               `json:"indestructible_zoneroot,omitempty"`
	*/
	KernelVersion     string  `json:"kernel_version,omitempty"`
	MaintainResolvers *bool   `json:"maintain_resolvers,omitempty"`
	MaxPhysicalMemory *uint32 `json:"max_physical_memory,omitempty"`
	/*
		MaxSwap           uint32             `json:"max_swap,omitempty"`
	*/
	NetworkInterfaces []NetworkInterface `json:"nics,omitempty"`
	Quota             *uint32            `json:"quota,omitempty"`
	RAM               *uint32            `json:"ram,omitempty"`
	Resolvers         []string           `json:"resolvers,omitempty"`
	VirtualCPUCount   *uint32            `json:"vcpus,omitempty"`
	State             string             `json:"state,omitempty"`
	PrimaryIP         string             `json:"-"`

	Metadata map[string]string `json:"-"`
}

func (m *Machine) UpdatePrimaryIP() {
	m.PrimaryIP = ""
	for _, networkInterface := range m.NetworkInterfaces {
		if networkInterface.IsPrimary != nil {
			m.PrimaryIP = networkInterface.IPAddress
			break
		}
	}
}

func (m *Machine) UpdateMetadata() {
	metadata := map[string]string{}
	prefix := "terraform:"
	for k, v := range m.CustomerMetadata {
		if strings.HasPrefix(k, prefix) {
			metadata[strings.TrimPrefix(k, prefix)] = v
			delete(m.CustomerMetadata, k)
		}
	}
	m.Metadata = metadata
}

func newBool(value bool) *bool {
	n := value
	return &n
}

func newUint32(value uint32) *uint32 {
	n := value
	return &n
}

func (m *Machine) LoadFromSchema(d *schema.ResourceData) error {

	m.Alias = d.Get("alias").(string)
	m.Brand = d.Get("brand").(string)

	if NodeName, ok := d.GetOk("node_name"); ok {
		m.NodeName = NodeName.(string)
	}

	if iid, ok := d.GetOk("image_uuid"); ok {
		uuid, _ := uuid.Parse(iid.(string))
		m.ImageUUID = &uuid
	}

	if autoboot, ok := d.GetOk("autoboot"); ok {
		m.Autoboot = newBool(autoboot.(bool))
	}

	if cpuCap, ok := d.GetOk("cpu_cap"); ok {
		m.CPUCap = newUint32(uint32(cpuCap.(int)))
	}

	customerMetaData := map[string]string{}
	for k, v := range d.Get("customer_metadata").(map[string]interface{}) {
		customerMetaData[k] = v.(string)
	}
	m.CustomerMetadata = customerMetaData

	metadata := map[string]string{}
	for k, v := range d.Get("metadata").(map[string]interface{}) {
		metadata[k] = v.(string)
	}
	m.Metadata = metadata

	if disks, ok := d.GetOk("disks"); ok {
		m.Disks, _ = getDisks(disks)
	}

	if kernelVersion, ok := d.GetOk("kernel_version"); ok {
		m.KernelVersion = kernelVersion.(string)
	}

	if maxPhysicalMemory, ok := d.GetOk("max_physical_memory"); ok {
		m.MaxPhysicalMemory = newUint32(uint32(maxPhysicalMemory.(int)))
	}

	if maintainResolvers, ok := d.GetOk("maintain_resolvers"); ok {
		m.MaintainResolvers = newBool(maintainResolvers.(bool))
	}

	if nics, ok := d.GetOk("nics"); ok {
		m.NetworkInterfaces, _ = getNetworkInterfaces(nics)
	}

	if quota, ok := d.GetOk("quota"); ok {
		m.Quota = newUint32(uint32(quota.(int)))
	}

	if ram, ok := d.GetOk("ram"); ok {
		m.RAM = newUint32(uint32(ram.(int)))
	}

	if vcpus, ok := d.GetOk("vcpus"); ok {
		m.VirtualCPUCount = newUint32(uint32(vcpus.(int)))
	}

	var resolvers []string
	for _, resolver := range d.Get("resolvers").([]interface{}) {
		resolvers = append(resolvers, resolver.(string))
	}
	m.Resolvers = resolvers

	return nil
}

func (m *Machine) SaveToSchema(d *schema.ResourceData) error {
	d.Set("primary_ip", m.PrimaryIP)
	d.Set("id", m.ID.String())
	d.Set("node_name", m.NodeName)

	// We update the metadata in case machine provisioning pushed data there.
	d.Set("metadata", m.Metadata)

	if m.PrimaryIP != "" {
		log.Printf("Machine saved to schema with primary IP: '%s'", m.PrimaryIP)
		d.SetConnInfo(map[string]string{
			"type": "ssh",
			"host": m.PrimaryIP,
		})
	}

	return nil
}

func (m *Machine) setCustomerMetadata(key string, value interface{}) {
	if m.SetCustomerMetadata == nil {
		m.SetCustomerMetadata = make(map[string]string)
	}

	m.SetCustomerMetadata[key] = value.(string)
}

func (m *Machine) removeCustomerMetadata(key string) {
	m.RemoveCustomerMetadata = append(m.RemoveCustomerMetadata, key)
}

func stringsAreEqual(a interface{}, b interface{}) bool {
	return a.(string) == b.(string)
}

type NetworkInterface struct {
	/*
		AllowDHCPSpoofing     bool     `json:"allow_dhcp_spoofing,omitempty"`
	*/
	AllowIPSpoofing        bool `json:"allow_ip_spoofing"`
	AllowMACSpoofing       bool `json:"allow_mac_spoofing"`
	AllowRestrictedTraffic bool `json:"allow_restricted_traffic"`
	/*
		BlockedOutgoingPorts  []uint16 `json:"blocked_outgoing_ports,omitempty"`
	*/
	Gateways    []string `json:"gateways,omitempty"`
	Interface   string   `json:"interface,omitempty"`
	IPAddresses []string `json:"ips,omitempty"`
	IPAddress   string   `json:"ip,omitempty"`
	/*
		HardwareAddress       string   `json:"mac,omitempty"`
	*/
	Model        string `json:"model,omitempty"`
	Tag          string `json:"nic_tag,omitempty"`
	IsPrimary    *bool  `json:"primary,omitempty"`
	VirtualLANID uint16 `json:"vlan_id,omitempty"`

	/*
		VRRP Support
	*/
	VRRPVRID      *uint32 `json:"vrrp_vrid,omitempty"`
	VRRPPrimaryIP string  `json:"vrrp_primary_ip,omitempty"`
}

func getNetworkInterfaces(d interface{}) ([]NetworkInterface, error) {
	networkInterfaceDefinitions := d.([]interface{})

	var networkInterfaces []NetworkInterface

	for _, nid := range networkInterfaceDefinitions {
		networkInterfaceDefinition := nid.(map[string]interface{})

		allowRestrictedTraffic := false
		if value, ok := networkInterfaceDefinition["allow_restricted_traffic"].(bool); ok {
			allowRestrictedTraffic = value
		}

		allowIPSpoofing := false
		if value, ok := networkInterfaceDefinition["allow_ip_spoofing"].(bool); ok {
			allowIPSpoofing = value
		}

		allowMACSpoofing := false
		if value, ok := networkInterfaceDefinition["allow_mac_spoofing"].(bool); ok {
			allowMACSpoofing = value
		}

		var gateways []string
		for _, gateway := range networkInterfaceDefinition["gateways"].([]interface{}) {
			gateways = append(gateways, gateway.(string))
		}

		interfaceName := networkInterfaceDefinition["interface"].(string)

		var ips []string
		for _, ip := range networkInterfaceDefinition["ips"].([]interface{}) {
			ips = append(ips, ip.(string))
		}

		nicTag := networkInterfaceDefinition["nic_tag"].(string)

		var vlanID uint16
		if vlanIDCheck, ok := networkInterfaceDefinition["vlan_id"].(int); ok {
			vlanID = uint16(vlanIDCheck)
		}

		model := ""
		if m, ok := networkInterfaceDefinition["model"].(string); ok {
			model = m
		}

		var vrrpVRID *uint32
		if vrrpVRIDCheck, ok := networkInterfaceDefinition["vrrp_vrid"]; ok {
			// BUGBUG: for some reason vrrp_vrid is always found above but when
			// it's not specified, it comes back as zero.  This value is checked and
			// used as a "not present" value.
			vrid := vrrpVRIDCheck.(int)
			if vrid > 0 {
				vrrpVRID = newUint32(uint32(vrid))
			}
		}

		vrrpPrimaryIP := ""
		if m, ok := networkInterfaceDefinition["vrrp_primary_ip"].(string); ok {
			vrrpPrimaryIP = m
		}

		networkInterface := NetworkInterface{
			AllowRestrictedTraffic: allowRestrictedTraffic,
			AllowIPSpoofing:        allowIPSpoofing,
			AllowMACSpoofing:       allowMACSpoofing,
			Interface:              interfaceName,
			IPAddresses:            ips,
			Tag:                    nicTag,
			Gateways:               gateways,
			VirtualLANID:           vlanID,
			Model:                  model,
			VRRPVRID:               vrrpVRID,
			VRRPPrimaryIP:          vrrpPrimaryIP,
		}

		networkInterfaces = append(networkInterfaces, networkInterface)
	}

	return networkInterfaces, nil
}

type Disk struct {
	Boot        bool       `json:"boot,omitempty"`
	Compression string     `json:"compression,omitempty"`
	ImageUUID   *uuid.UUID `json:"image_uuid,omitempty"`
	ImageSize   uint32     `json:"image_size,omitempty"`
	Model       string     `json:"model,omitempty"`
	Size        *uint32    `json:"size,omitempty"`
}

func getDisks(d interface{}) ([]Disk, error) {
	diskDefinitions := d.([]interface{})

	var disks []Disk

	for _, dd := range diskDefinitions {
		diskDefinition := dd.(map[string]interface{})
		disk := Disk{}

		if b, ok := diskDefinition["boot"]; ok {
			disk.Boot = b.(bool)
		}

		if c, ok := diskDefinition["compression"]; ok {
			disk.Compression = c.(string)
		}

		if iid, ok := diskDefinition["image_uuid"]; ok {
			iid2, _ := uuid.Parse(iid.(string))
			if iid2 != uuid.Nil {
				disk.ImageUUID = &iid2
			} else {
				disk.ImageUUID = nil
			}
		}

		if is, ok := diskDefinition["image_size"]; ok {
			disk.ImageSize = uint32(is.(int))
		}

		if m, ok := diskDefinition["model"]; ok {
			disk.Model = m.(string)
		}

		if sz, ok := diskDefinition["size"]; ok {
			size := uint32(sz.(int))
			if size > 0 {
				disk.Size = &size
			}
		}

		disks = append(disks, disk)
	}

	return disks, nil
}
