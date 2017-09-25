Terraform Provider for Microsoft Hyper-V
==================


<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.8 (to build the provider plugin)
-	Powershell 5

Building The Provider (Windows)
---------------------

Clone repository to: `$GOPATH/src/github.com/flynnhandley/terraform-provider-hyperv`

```powershell
$ mkdir $GOPATH/src/github.com/terraform-providers -force; cd $GOPATH/src/github.com/terraform-providers
$ git clone git@github.com:terraform-providers/terraform-provider-google
```

Enter the provider directory and build the provider

```powershell
$ cd $GOPATH/src/github.com/github.com/flynnhandley/terraform-provider-hyperv
$ go build -o terraform-provider-hyperv.exe
```

Using the provider
----------------------

```
provider "hyperv" {
  hypervisor          = "hv01.dns.address"
  username            = "hyperv_user"
  password            = "${var.hv_password}"
}

resource "hyperv_virtual_switch" "dmz" {
  name  = "dmz"
}

resource "hyperv_virtual_switch" "internet" {
  name  = "internet"
}

resource "hyperv_virtual_machine" "web" {
  vm_name         = "web"
  boot_vhd_uri    = "http://address.to.repository/server2012r2-1.2.3.vhdx",
  switch_name     = "${hyperv_virtual_switch.internet.name}"
  cpu             = 2,
  ram_mb          = 1024,
  provision       = true,
  network_adapter {
    name = "internal",
    switch_name = "${hyperv_virtual_switch.dmz.name}",
  },

    provisioner "remote-exec" {
    inline = ["powershell.exe -executionpolicy unrestricted write-host Provisioning complete"],

    connection {
      type      = "winrm"
      user      = "vm_username"
      password  = "${var.vm_password}"
      timeout   = "5m"
    }
  }
}

resource "hyperv_virtual_machine" "db" {
  vm_name         = "db"
  boot_vhd_uri    = "http://address.to.repository/server2012r2-1.2.3.vhdx",
  switch_name     = "${hyperv_virtual_switch.db.name}"
  cpu             = 2,
  ram_mb          = 1024,
  provision       = true,
    provisioner "remote-exec" {
    inline = ["powershell.exe -executionpolicy unrestricted write-host Provisioning complete"],

    connection {
      type      = "winrm"
      user      = "vm_username"
      password  = "${var.vm_password}"
      timeout   = "5m"
    }
  }
}

// Set as environment variables using TF_VAR_name
variable "hv_password" {
  type    = "string"
}

variable "vm_password" {
  type    = "string"
}

```
