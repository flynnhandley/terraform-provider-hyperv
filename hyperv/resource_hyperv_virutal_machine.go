package hyperv

import (
	"errors"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceHypervVM() *schema.Resource {
	return &schema.Resource{

		Create: resourceHypervVMCreate,
		Read:   resourceHypervVMRead,
		Update: resourceHypervVMUpdate,
		Delete: resourceHypervVMDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			// Mandatory parameters
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"switch": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Optional parameters
			"path": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "C:\\HyperV",
			},

			"processors": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2,
			},

			"vlan_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"mac": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateMacAddress,
			},

			"disable_network_boot": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"disable_secure_boot": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"generation": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2,
			},

			"ram": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2048,
			},

			"wait_for_ip": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timeout": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  5,
						},
						"adapter_name": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "Network Adapter",
						},
					},
				},
			},

			"storage_disk": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"image_path": {
							Type:     schema.TypeString,
							Optional: true,
							ConflictsWith: []string{"storage_disk.image_url",
								"storage_disk.diff_parent_path",
								"storage_disk.size"},
						},
						"image_url": {
							Type:     schema.TypeString,
							Optional: true,
							ConflictsWith: []string{"storage_disk.image_path",
								"storage_disk.diff_parent_path",
								"storage_disk.size"},
						},

						"diff_parent_path": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"storage_disk.image_url", "storage_disk.image_path"},
						},

						"size": {
							Type:          schema.TypeInt,
							Optional:      true,
							Default:       50,
							ConflictsWith: []string{"storage_disk.image_url", "storage_disk.image_path"},
						},
					},
				},
			},
			"network_adapter": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"switch_name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"vlan_id": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						"mac": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceHypervVMCreate(d *schema.ResourceData, meta interface{}) error {
	var hvDriver = meta.(Driver)
	var id string
	var err error
	var vm VM

	// Make ResourceData esy to work with / do some validation
	if vm, err = NewVM(d); err != nil {
		return err
	}

	// Check if VM already exists (Outside of terraform)
	if existingID, err := hvDriver.GetVirtualMachineId(map[string]string{"vmName": vm.Name}); (existingID != "") || (err != nil) {
		if err != nil {
			return err
		}
		return errors.New("Cannot create VM: " + vm.Name + ". Already exists!")
	}

	// Create VM on HV
	log.Printf("[DEBUG] Creating VM: " + vm.Name)
	if id, err = hvDriver.CreateVirtualMachine(vm.Name, vm.Path, vm.RAMMB, vm.Switch, vm.Generation); err != nil {
		return err
	}

	// VM Created
	d.SetId(id)

	if vm.MAC != "" {
		if err = hvDriver.SetNetworkAdapterStaticMacAddress(vm.Name, "Network Adapter", vm.MAC); err != nil {
			return err
		}
	}

	// Attach network adapters
	if vm.NetworkAdapters != nil {
		for _, a := range vm.NetworkAdapters {
			log.Printf("[DEBUG] Creating network adapter with vmID: " + id + " ,name: " + a.Name + ",switchName: " + a.SwitchName + ",vlanID: " + a.VlanID)
			if err = hvDriver.AddVMNetworkAdapter(id, a.Name, a.SwitchName, a.VlanID); err != nil {
				return err
			}
		}
	}

	// Attach storage disks
	if vm.StorageDisks != nil {
		for _, d := range vm.StorageDisks {
			if d.DiffParentPath != "" {
				if _, err = hvDriver.NewDifferencingDisk(id, d.Name, d.DiffParentPath); err != nil {
					return err
				}
			} else if d.ImagePath != "" {
				if _, err = hvDriver.NewDiskFromImagePath(id, d.Name, d.ImagePath); err != nil {
					return err
				}
			} else if d.ImageURL != "" {
				if _, err = hvDriver.NewDiskFromImageURL(id, d.Name, d.ImageURL); err != nil {
					return err
				}
			} else {
				if _, err = hvDriver.NewVhd(id, d.Name, d.Size); err != nil {
					return err
				}
			}
		}
	}

	// Disable network boot (Removes network from boot order)
	if vm.DisableNetworkBoot {
		if err = hvDriver.DisableNetworkBoot(id); err != nil {
			return err
		}
	}

	// Start VM
	log.Printf("[DEBUG] Starting VM")
	if err = hvDriver.Start(vm.Name); err != nil {
		return err
	}

	// This approach is not recommended, instead, set the MAC address and use DHCP / DNS to communicate.
	if err = WaitForIp(d, hvDriver, vm); err != nil {
		return err
	}

	return resourceHypervVMRead(d, meta)
}

func resourceHypervVMRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceHypervVMUpdate(d *schema.ResourceData, meta interface{}) error {

	resourceHypervVMDelete(d, meta)

	resourceHypervVMCreate(d, meta)

	return resourceHypervVMRead(d, meta)
}

func resourceHypervVMDelete(d *schema.ResourceData, meta interface{}) error {

	hvDriver := meta.(Driver)
	var vmID = d.Id()

	if err := hvDriver.DeleteVirtualMachine(vmID); err != nil {
		return err
	}

	return nil
}
