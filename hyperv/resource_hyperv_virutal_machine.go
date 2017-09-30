package hyperv

import (
	"errors"
	"log"
	"time"

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
			"vm_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"switch": {
				Type:     schema.TypeString,
				Required: true,
			},

			"path": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "C:\\HyperV",
			},

			"cpu": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2,
			},

			"generation": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2,
			},

			"ram_mb": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2048,
			},

			"vlan_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"disable_pxe": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"provision": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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

	if vm, err = NewVM(d); err != nil {
		return err
	}

	if existingID, err := hvDriver.GetVirtualMachineId(map[string]string{"vmName": vm.Name}); (existingID != "") || (err != nil) {
		if err != nil {
			return err
		}
		return errors.New("Cannot create VM: " + vm.Name + ". Already exists!")
	}

	// Create VM on HV
	log.Printf("[DEBUG] Creating VM: " + vm.Name)
	if id, err = hvDriver.CreateVirtualMachine(vm.Name, "C:\\hyperv", vm.RAMMB, vm.Switch, vm.Generation); err != nil {
		return err
	}

	// VM Created
	d.SetId(id)

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
	log.Printf("[DEBUG] LETS COUNT THE DISKS LETS COUNT THE DISKS LETS COUNT THE DISKS LETS COUNT THE DISKS LETS COUNT THE DISKS LETS COUNT THE DISKS LETS COUNT THE DISKS ")
	if vm.StorageDisks != nil {
		log.Printf("[DEBUG] STORAGE DISKS NOT NULL [DEBUG] STORAGE DISKS NOT NULL [DEBUG] STORAGE DISKS NOT NULL [DEBUG] STORAGE DISKS NOT NULL [DEBUG] STORAGE DISKS NOT NULL [DEBUG] STORAGE DISKS NOT NULL ")
		for _, d := range vm.StorageDisks {

			log.Printf("[DEBUG] DDISKS ARE REAL")

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
				log.Printf("[DEBUG] Creating VHD")
				if _, err = hvDriver.NewVhd(id, d.Name, d.Size); err != nil {
					return err
				}
			}

		}
	}

	if err = hvDriver.SetVirtualMachineRemoveNetworkBoot(id); err != nil {
		return err
	}

	// Start VM
	log.Printf("[DEBUG] Starting VM")
	if err = hvDriver.Start(vm.Name); err != nil {
		return err
	}

	// If provision is set, wait for IP and then add to connection info
	if d.Get("provision").(bool) {
		log.Printf("[DEBUG] Provision set, waiting for IP address")
		for ip, _ := hvDriver.GetVirtualMachineNetworkAdapterAddress(vm.Name); ip == ""; {
			log.Printf("[DEBUG] Waiting for IP")
			time.Sleep(time.Second * 10)
			ip, _ = hvDriver.GetVirtualMachineNetworkAdapterAddress(vm.Name)
		}

		ip, _ := hvDriver.GetVirtualMachineNetworkAdapterAddress(vm.Name)
		d.SetConnInfo(map[string]string{"host": ip})
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
	var vmName string
	var err error

	hvDriver := meta.(Driver)
	var vmID = d.Id()

	if vmName, err = hvDriver.InvokeCommand("(Get-VM -ID "+vmID+").Name", nil); vmName == "" {
		return err
	}

	if err := hvDriver.DeleteVirtualMachine(vmID); err != nil {
		return err
	}

	// Wait for HV to delete VM and then clean up files
	cmd := "While(Get-VM -ID " + vmID + " -ErrorAction SilentlyContinue){Start-Sleep 1};if(Test-path 'C:\\HyperV\\" + vmName + ".vhdx'){Remove-Item 'C:\\HyperV\\" + vmName + ".vhdx' -force};if(Test-path 'C:\\HyperV\\" + vmName + "'){Remove-Item 'C:\\HyperV\\" + vmName + "' -recurse -force}"
	if _, err := hvDriver.InvokeCommand(cmd, nil); err != nil {
		return err
	}
	return nil
}
