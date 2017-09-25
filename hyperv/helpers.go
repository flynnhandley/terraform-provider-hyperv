package hyperv

import (
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

// VM ...
type VM struct {
	Name              string
	SwitchName        string
	BootVHDUri        string
	CPU               int
	Generation        int
	RAMMB             int64
	DiskPath          string
	DisableSecureBoot string
	VLANID            string
	DisablePXE        string
	StorageDataDisks  []interface{}
	NetworkAdapters   []VMNetworkAdapter
}

// VMStorageDataDisk ...
type VMStorageDataDisk struct {
	Name         string
	Size         string
	format       string
	CreateOption string
}

// VMNetworkAdapter ...
type VMNetworkAdapter struct {
	Name       string
	SwitchName string
	VlanID     string
}

// Switch ..
type Switch struct {
	Hypervisor string
	Username   string
	Password   string
	Clustered  bool
	Driver     Driver
}

// GetVM ...
func GetVM(d *schema.ResourceData) (VM, error) {

	vm := VM{
		Name:       d.Get("vm_name").(string),
		SwitchName: d.Get("switch_name").(string),
		BootVHDUri: d.Get("boot_vhd_uri").(string),
		RAMMB:      int64(d.Get("ram_mb").(int)),
		Generation: d.Get("generation").(int),
		CPU:        d.Get("cpu").(int),
	}

	if vL, ok := d.GetOk("network_adapter"); ok {
		var adapters []VMNetworkAdapter
		for _, v := range vL.([]interface{}) {
			network := v.(map[string]interface{})
			adapterName, _ := network["name"].(string)
			adapterSwitchName, _ := network["switch_name"].(string)
			adapterVlanID, _ := network["vlan_id"].(int)

			adapter := VMNetworkAdapter{
				Name:       adapterName,
				SwitchName: adapterSwitchName,
				VlanID:     strconv.Itoa(adapterVlanID),
			}
			adapters = append(adapters, adapter)
		}
		vm.NetworkAdapters = adapters
	}

	return vm, nil
}

// GetSwitch ...
func GetSwitch() (Switch, error) {

	return Switch{}, nil
}
