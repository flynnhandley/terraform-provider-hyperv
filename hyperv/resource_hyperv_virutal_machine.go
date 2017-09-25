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

			"switch_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"boot_vhd_uri": {
				Type:     schema.TypeString,
				Required: true,
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

			"disk_path": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "C:\\hyperv",
			},

			"disable_secure_boot": {
				Type:     schema.TypeBool,
				Optional: true,
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

			"storage_data_disk": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"size": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"format": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"create_option": {
							Type:     schema.TypeString,
							Optional: true,
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

	if vm, err = GetVM(d); err != nil {
		return err
	}

	if existingID, _ := hvDriver.GetVirtualMachineId(map[string]string{"vmName": vm.Name}); existingID != "" {
		return errors.New("Cannot create VM: " + vm.Name + ". Already exists!")
	}

	// Create VM on HV
	log.Printf("[DEBUG] Creating VM")
	if id, err = hvDriver.CreateVirtualMachine(vm.Name, "C:\\hyperv", vm.RAMMB, vm.SwitchName, vm.Generation); err != nil {
		return err
	}

	// Set CPU Count
	if err = hvDriver.SetVirtualMachineCpuCount(id, vm.CPU); err != nil {
		return err
	}

	// VM Created
	d.SetId(id)

	// Set VLAN (Applies to all adapters currently connected to the VM)
	if vm.VLANID != "" {
		log.Printf("[DEBUG] SETTING VLAN")

		if err = hvDriver.SetVirtualMachineVlanId(id, vm.VLANID); err != nil {

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

	// Download / Attach VHD
	log.Printf("[DEBUG] Downloading and attaching boot VHD")
	if _, err = hvDriver.AttachBootVHD(id, vm.BootVHDUri); err != nil {
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
