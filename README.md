

Terraform Provider for Microsoft Hyper-V
==================

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.8 (to build the provider plugin)
-	Powershell 5

# hyperv_virtual_machine

Create a virtual machine.

## Example Usage with Differencing Disks (Recommended)

```hcl
provider "hyperv" {
  hypervisor            = "HV01.contoso.local"
  username              = "hv_administrator"
  password              = "Password1234!"
}

resource "hyperv_virtual_switch" "application_switch" {
  name  = "application_switch"
}

resource "hyperv_virtual_machine" "web" {
  vm_name               = "web"
  cpu                   = 2
  ram_mb                = 1024
  switch                = "external_switch"
  disable_network_boot  = true
  path                  = "C:\\ClusterStorage\\VMs"

  network_adapter {
    name = "inside",
    switch_name = "inside",
  },

  storage_disk {
    name                  = "boot"
    diff_parent_path      = "C:\\ClusterStorage\\VHDs\\server2012r2-0.1.0.vhdx"
  }
}

resource "hyperv_virtual_machine" "db" {
  vm_name               = "db"
  cpu                   = 4
  ram_mb                = 4096
  switch                = "external_switch"
  disable_network_boot  = true
  path                  = "C:\\ClusterStorage\\VMs"

  network_adapter {
    name = "inside",
    switch_name = "inside",
  },

  storage_disk {
    name                  = "boot"
    diff_parent_path      = "C:\\ClusterStorage\\VHDs\\server2012r2-0.1.0.vhdx"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the virtual machine resource.
* `switch` - (Required) The name of the virtual switch that the virtual machine will connect to.
* `path` - (Optional) The directory to store the files for the new virtual machine, defaults to `C:\HyperV`.
* `processors` - (Optional) The number of Virtual Processors to assign to the virtual machine, defaults to `2`.
* `vlan_id` - (Optional) Specifies the virtual LAN identifier of the virtual machine network adapter.
* `mac` - (Optional) The mac address of the virtual machines network adapter.
* `disable_network_boot` - (Optional) Removes network devices from the BIOS boot sequence, this saves time when multiple adapters are used.
* `disable_secure_boot` - (Optional) Specifies whether to disable secure boot.
* `generation` - (Optional) Specifies the generation, as an integer, for the virtual machine. The values that are valid are 1 and 2, defaults to `2`
* `ram` - (Optional) An integer representing the amount of RAM allocated to the VM in MB, defaults to `2048`
* `wait_for_ip` - (Optional) Polls the HyperV host for the guest VM's IP address. If an IP address is detected, the host property of the provisioners connection is set. The wait_for_ip block is documented below.
* `storage_disk` - (Optional) The storage_disk block is documented below.
* `network_adapter` - (Optional) The network_adapter block is documented below.

`wait_for_ip` Supports the following:
* `timeout` - (Optional) Time in minutes to wait for an IP, an error will be generated if an IPv4 address is not detected after `timeout` minutes, defaults to `5`
* `adapter_name` - (Required) The name of the network adapter to query, defaults to `Network Adapter`

**The use of wait_for_ip is discouraged* it increases the chance of createing a *zombie resource**

`storage_disk` supports the following:

* `name` - (Required) Name of the VHDX, this must be unique.
* `image_path` - (Optional) A local directory on the HyperV host containing a VHDX image. For example *C:\\ClusterStorage\\VHDs\windows2012r2-0.0.1.vhdx*. `image_path` cannot be used with `image_url` or `diff_parent_path`
* `image_url` - (Optional) Specifies an image to download. For example *http://aritfactory.local/generic/windows2012r2-0.0.1.vhdx*. `image_url` cannot be used with `image_path` or `diff_parent_path`.
* `diff_parent_path` - (Optional) A local directory on the HyperV host containing a VHDX image, this image will be used as the parent and a new differencing disk will created. `diff_parent_path` cannot be used with `image_url` or `image_path`
* `size` - (Optional) The size in MB when creating a new vhdx. A new vhdx will be created if neither  `diff_parent_path`, `image_url` or `image_path` are specified, default is `50`

`network_adapter` supports the following:
* `name` - (Required) The name of the network adapter.
* `switch_name` - (Required) The virtual switch to connect the network adapter to,
* `vlan_id` - (Optional) Specifies the virtual LAN identifier of the virtual machine network adapter.
* `mac` - (Optional) he MAC address of this network adapter, must not contain any special characters.

## Attributes Reference

The following attributes are exported:

* `id` - The virtual machine ID.


## Example Using a Provisioner With  `wait_for_ip` Enabled (Not recommended)

```hcl
resource "hyperv_virtual_machine" "db" {
  vm_name               = "db"
  cpu                   = 4
  ram_mb                = 4096
  switch                = "external_switch"
  disable_network_boot  = true
  path                  = "C:\\ClusterStorage\\VMs"
  wait_for_ip           = true

  storage_disk {
    name                  = "boot"
    diff_parent_path      = "C:\\ClusterStorage\\VHDs\\server2012r2-0.1.0.vhdx"
  }

  provisioner "remote-exec" {
    inline = ["powershell.exe -executionpolicy unrestricted Write-Host Hello World!"],

    connection {
      type      = "winrm"
      user      = "Administrator"
      password  = "Password1234!"
      timeout   = "5m"
    }
  }
}
```

## Example Using a Provisioner with `wait_for_ip` Disabled  (recommended)
In order for this to work, you will need to provide the VM with either a static IP or DHCP binding and DNS entry.

```hcl
resource "hyperv_virtual_machine" "db" {
  vm_name               = "db"
  cpu                   = 4
  ram_mb                = 4096
  switch                = "external_switch"
  disable_network_boot  = true
  path                  = "C:\\ClusterStorage\\VMs"

  storage_disk {
    name                  = "boot"
    diff_parent_path      = "C:\\ClusterStorage\\VHDs\\server2012r2-0.1.0.vhdx"
  }

  provisioner "remote-exec" {
    inline = ["powershell.exe -executionpolicy unrestricted Write-Host Hello World!"],

    connection {
      type      = "winrm"
      host      = "db.contoso.local"
      user      = "Administrator"
      password  = "Password1234!"
      timeout   = "5m"
    }
  }
}
```