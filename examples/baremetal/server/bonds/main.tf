resource "vkcs_baremetal_server" "server_bond" {
  name              = "server-bond"
  availability_zone = "GZ1"
  flavor_id         = data.vkcs_baremetal_flavor.minimal.id
  os_id             = data.vkcs_baremetal_os.ubuntu.id
  key_pair          = vkcs_compute_keypair.generated_key.name
  raid_type         = "RAID1"

  bond {
    name            = "bond0"
    interface_names = ["nic0"]
    vlan {
      native     = true
      network_id = vkcs_networking_network.app.id
      subnet_id  = vkcs_networking_subnet.app.id
    }
  }

  bond {
    name            = "bond1"
    interface_names = ["nic1"]

    vlan {
      native     = true
      network_id = vkcs_networking_network.db.id
      subnet_id  = vkcs_networking_subnet.db.id
    }
  }
}


