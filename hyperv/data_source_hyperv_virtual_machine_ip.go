package hyperv

import (
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceHypervVirtualMachineIP() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceHypervVirtualMachineIPRead,

		Schema: map[string]*schema.Schema{
			"virtual_machine_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceHypervVirtualMachineIPRead(d *schema.ResourceData, meta interface{}) error {
	hvDriver := meta.(Driver)
	// vmID := d.Get("virtual_machine_id").(string)
	vmID := "Asdf"
	ip := ""

	log.Printf("[DEBUG] Provision set, waiting for IP address")
	for ip, _ := hvDriver.GetVirtualMachineNetworkAdapterAddress(vmID); ip == ""; {
		log.Printf("[DEBUG] Waiting for IP")
		time.Sleep(time.Second * 10)
		ip, _ = hvDriver.GetVirtualMachineNetworkAdapterAddress(vmID)
	}

	ip, _ = hvDriver.GetVirtualMachineNetworkAdapterAddress(vmID)
	d.SetConnInfo(map[string]string{"host": ip})

	return nil
}
