package hyperv

import (
	"errors"

	"github.com/hashicorp/terraform/helper/schema"
)

const ()

func resourceHypervVirtualSwitch() *schema.Resource {
	return &schema.Resource{

		Create: resourceHypervVirtualSwitchCreate,
		Read:   resourceHypervVirtualSwitchRead,
		Update: resourceHypervVirtualSwitchUpdate,
		Delete: resourceHypervVirtualSwitchDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "internal",
			},
		},
	}
}

func resourceHypervVirtualSwitchCreate(d *schema.ResourceData, meta interface{}) error {

	hvDriver := meta.(Driver)
	name := d.Get("name").(string)

	if exists, _ := hvDriver.GetVirtualSwitchID(map[string]string{"name": name}); exists != "" {
		return errors.New("Cannot create Switch: " + name + ". Already exists!")
	}

	switchType := d.Get("type").(string)

	id, err := hvDriver.CreateVirtualSwitch(name, switchType)

	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceHypervVirtualSwitchRead(d, meta)
}

func resourceHypervVirtualSwitchRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceHypervVirtualSwitchUpdate(d *schema.ResourceData, meta interface{}) error {

	resourceHypervVirtualSwitchDelete(d, meta)

	resourceHypervVirtualSwitchCreate(d, meta)

	return resourceHypervVirtualSwitchRead(d, meta)
}

func resourceHypervVirtualSwitchDelete(d *schema.ResourceData, meta interface{}) error {
	hvDriver := meta.(Driver)
	id := d.Id()

	return hvDriver.DeleteVirtualSwitch(id)
}
