resource "vkcs_networking_network" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name       = "subnet_1"
  cidr       = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_sharedfilesystem_sharenetwork" "sharenetwork_1" {
  name              = "test_sharenetwork"
  description       = "test share network with security services"
  neutron_net_id    = "${vkcs_networking_network.network_1.id}"
  neutron_subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
}

resource "vkcs_sharedfilesystem_share" "share_1" {
  name             = "nfs_share"
  description      = "test share description"
  share_proto      = "NFS"
  size             = 1
  share_network_id = "${vkcs_sharedfilesystem_sharenetwork.sharenetwork_1.id}"
}

resource "vkcs_sharedfilesystem_share_access" "share_access_1" {
  share_id     = "${vkcs_sharedfilesystem_share.share_1.id}"
  access_type  = "ip"
  access_to    = "192.168.199.10"
  access_level = "rw"
}
