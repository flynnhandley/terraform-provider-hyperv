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
