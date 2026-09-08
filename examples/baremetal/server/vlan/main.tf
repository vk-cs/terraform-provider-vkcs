resource "vkcs_baremetal_server" "server_vlan" {
  name              = "server-vlan"
  availability_zone = "GZ1"
  flavor_id         = data.vkcs_baremetal_flavor.minimal.id
  os_id             = data.vkcs_baremetal_os.ubuntu.id
  key_pair          = vkcs_compute_keypair.generated_key.name
  raid_type         = "RAID1"

  nic {
    name = "nic0"
    vlan {
      native     = true
      network_id = vkcs_networking_network.app.id
      subnet_id  = vkcs_networking_subnet.app.id
    }

    vlan {
      id         = 100
      network_id = vkcs_networking_network.db.id
      subnet_id  = vkcs_networking_subnet.db.id
    }
  }
}


