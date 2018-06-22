package hyperv

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

// VM ...
type VM struct {
	Name                string
	SwitchName          string
	BootVHDUri          string
	MAC                 string
	CPU                 int
	Generation          int
	RAMMB               int64
	Path                string
	Switch              string
	DisablePXE          string
	VlanID              int
	EnableSecureBoot    bool
	SecureBootTemplate  string
	DisableNetworkBoot  bool
	EnableDynamicMemory bool
	NetworkAdapters     []NetworkAdapter
	StorageDisks        []Disk
}

// NetworkAdapter ...
type NetworkAdapter struct {
	Name       string
	SwitchName string
	VlanID     string
	MAC        string
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
		Name:                d.Get("name").(string),
		RAMMB:               int64(d.Get("ram").(int)),
		Generation:          d.Get("generation").(int),
		CPU:                 d.Get("cpu").(int),
		Path:                d.Get("path").(string),
		NetworkAdapters:     GetNetworkAdapters(d),
		Switch:              d.Get("switch").(string),
		VlanID:              d.Get("vlan_id").(int),
		DisableNetworkBoot:  d.Get("disable_network_boot").(bool),
		EnableSecureBoot:    d.Get("enable_secure_boot").(bool),
		SecureBootTemplate:  d.Get("secure_boot_template").(string),
		EnableDynamicMemory: d.Get("enable_dynamic_memory").(bool),
	}

	if val, ok := d.GetOk("mac"); ok {
		vm.MAC = val.(string)
	}

	disks, err := GetDisks(d)
	log.Printf("NEW VM")
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
	log.Printf("GET DISKS")
	if vL, ok := d.GetOk("storage_disk"); ok {
		log.Printf("OK")
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

func WaitForIP(d *schema.ResourceData, hvDriver Driver, vm VM) error {

	if val, ok := d.GetOk("wait_for_ip"); ok {
		log.Printf("[DEBUG] Provision set, waiting for IP address")
		for _, v := range val.([]interface{}) {
			var ip string
			var index = 0
			var err error
			wfIP := v.(map[string]interface{})
			adapterName := wfIP["adapter_name"].(string)
			timeOut := wfIP["timeout"].(int)

			for ip, err = hvDriver.GetVirtualMachineNetworkAdapterAddress(vm.Name, adapterName); (ip == "") || strings.HasPrefix(ip, "169") || strings.HasPrefix(ip, "fe80"); {

				if err != nil {
					return err
				}

				log.Printf("[DEBUG] Waiting for IP" + adapterName)
				log.Printf("[DEBUG] IP is:" + ip)
				log.Printf("[DEBUG] adapter name is:" + adapterName)

				if index > (timeOut * 6) {
					return errors.New("could not find IP address before timeout expires, minutes: " + strconv.Itoa(timeOut))
				}

				time.Sleep(time.Second * 10)

				ip, _ = hvDriver.GetVirtualMachineNetworkAdapterAddress(vm.Name, adapterName)
				index++
			}
			d.SetConnInfo(map[string]string{"host": ip})
		}
	}
	return nil
}
