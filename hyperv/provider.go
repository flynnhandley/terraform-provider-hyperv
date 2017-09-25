package hyperv

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {

	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"hypervisor": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HYPERV_HYPERVISOR", ""),
			},

			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HYPERV_USERNAME", ""),
			},

			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HYPERV_PASSWORD", ""),
			},

			"clustered": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"hyperv_virtual_machine": resourceHypervVM(),
			"hyperv_virtual_switch":  resourceHypervVirtualSwitch(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(data *schema.ResourceData) (interface{}, error) {

	config := Config{
		Username:   data.Get("username").(string),
		Password:   data.Get("password").(string),
		Hypervisor: data.Get("hypervisor").(string),
	}

	drv, err := config.GetDriver()

	return drv, err
}
