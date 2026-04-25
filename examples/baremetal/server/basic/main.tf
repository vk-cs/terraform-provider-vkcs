resource "vkcs_baremetal_server" "server" {
  name              = "tf-server-57b4d7f1"
  availability_zone = "GZ1"
  flavor_id         = data.vkcs_baremetal_flavor.minimal.id
  os_id             = data.vkcs_baremetal_os.ubuntu.id
  key_pair          = vkcs_compute_keypair.generated_key.name
  raid_type         = "RAID1"

  user_data = <<EOF
    #cloud-config
    package_upgrade: true
    packages:
      - nginx
    runcmd:
      - systemctl start nginx
  EOF

  nic {
    name = "nic0"
    vlan {
      native     = true
      network_id = vkcs_networking_network.app.id
      subnet_id  = vkcs_networking_subnet.app.id
    }
  }

  nic {
    name = "nic1"
    vlan {
      native     = true
      network_id = vkcs_networking_network.db.id
      subnet_id  = vkcs_networking_subnet.db.id
    }
  }
}
