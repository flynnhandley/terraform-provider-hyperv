provider "hyperv" {
  hypervisor = "HV01.contoso.local"
  username   = "hv_administrator"
  password   = "Password1234!"
}

resource "hyperv_virtual_machine" "db" {
  vm_name              = "db"
  cpu                  = 4
  ram                  = 4096
  switch               = "external_switch"
  disable_network_boot = true
  path                 = "C:\\ClusterStorage\\VMs"
  wait_for_ip          = true

  storage_disk {
    name             = "boot"
    diff_parent_path = "C:\\ClusterStorage\\VHDs\\server2012r2-0.1.0.vhdx"
  }

  provisioner "remote-exec" {
    inline = ["powershell.exe -executionpolicy unrestricted Write-Host Hello World!"]

    connection {
      type     = "winrm"
      user     = "Administrator"
      password = "Password1234!"
      timeout  = "5m"
    }
  }
}
