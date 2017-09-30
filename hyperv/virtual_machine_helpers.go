package hyperv

import (
	"errors"
	"log"
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
	Switch            string
	DisablePXE        string
	NetworkAdapters   []NetworkAdapter
	StorageDisks      []Disk
}

// NetworkAdapter ...
type NetworkAdapter struct {
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

// Disk ..
type Disk struct {
	ID             string
	Type           string
	Name           string
	Size           int64
	Path           string
	ImageURL       string
	ImagePath      string
	DiffParentPath string
}

// NewVM ...
func NewVM(d *schema.ResourceData) (VM, error) {

	vm := VM{
		Name:            d.Get("vm_name").(string),
		RAMMB:           int64(d.Get("ram_mb").(int)),
		Generation:      d.Get("generation").(int),
		CPU:             d.Get("cpu").(int),
		NetworkAdapters: GetNetworkAdapters(d),
		Switch:          d.Get("switch").(string),
	}

	disks, err := GetDisks(d)
	log.Printf("NEW VM NEW VM NEW VM NEW VM NEW VM NEW VM NEW VM NEW VM ")
	vm.StorageDisks = disks

	return vm, err
}

// GetSwitch ...
func GetSwitch() (Switch, error) {

	return Switch{}, nil
}

// GetNetworkAdapters ..
func GetNetworkAdapters(d *schema.ResourceData) []NetworkAdapter {

	if vL, ok := d.GetOk("network_adapter"); ok {
		var adapters []NetworkAdapter
		for _, v := range vL.([]interface{}) {
			network := v.(map[string]interface{})
			newAdapter := NetworkAdapter{}

			if val, ok := network["name"]; ok {
				newAdapter.Name = val.(string)
			}

			if val, ok := network["switch_name"]; ok {
				newAdapter.SwitchName = val.(string)
			}

			if _, ok := network["vlan_id"]; ok {
				newAdapter.VlanID = strconv.Itoa(network["vlan_id"].(int))
			}

			adapters = append(adapters, newAdapter)
		}
		return adapters
	}
	return nil
}

// GetDisks ..
func GetDisks(d *schema.ResourceData) ([]Disk, error) {
	log.Printf("GER DISKS GER DISKS GER DISKS GER DISKS GER DISKS GER DISKS GER DISKS GER DISKS GER DISKS ")
	if vL, ok := d.GetOk("storage_disk"); ok {
		log.Printf("ok ok ok ok OK OK OK OK OK")
		var disks []Disk
		for _, v := range vL.([]interface{}) {
			disk := v.(map[string]interface{})

			newDisk := Disk{}

			if val, ok := disk["name"]; ok {
				newDisk.Name = val.(string)
			}

			if val, ok := disk["type"]; ok {
				newDisk.Type = val.(string)
			}

			if val, ok := disk["image_path"]; ok {
				newDisk.ImagePath = val.(string)
			}

			if val, ok := disk["image_url"]; ok {
				newDisk.ImageURL = val.(string)
			}

			if val, ok := disk["diff_parent_path"]; ok {
				newDisk.DiffParentPath = val.(string)
			}

			if val, ok := disk["size"]; ok {
				newDisk.Size = int64(val.(int)) * 1024 * 1024 * 1024 // (GB to b)
			}

			switch newDisk.Type {
			case "image":
				if (newDisk.ImagePath == "") && (newDisk.ImageURL == "") && (newDisk.DiffParentPath == "") {
					return nil, errors.New("Either 'image_path' or 'image_url' or 'diff_parent_path' must be set when the 'type' parameter value is 'image'")
				}
			}
			log.Printf("APPENDING DISKS")
			disks = append(disks, newDisk)
		}
		log.Printf("RETURNING DISKS")
		return disks, nil
	}
	return nil, nil
}
