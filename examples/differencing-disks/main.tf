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