resource "vkcs_sharedfilesystem_sharenetwork" "sharenetwork" {
  name                = "test_sharenetwork"
  neutron_net_id      = "${vkcs_networking_network.sfs.id}"
  neutron_subnet_id   = "${vkcs_networking_subnet.sfs.id}"
}

resource "vkcs_sharedfilesystem_share" "share" {
  name             = "nfs_share"
  description      = "test share description"
  share_proto      = "NFS"
  share_type       = "default_share_type"
  size             = 1
  share_network_id = "${vkcs_sharedfilesystem_sharenetwork.sharenetwork.id}"
}

resource "vkcs_sharedfilesystem_share_access" "share_access_1" {
  share_id     = "${vkcs_sharedfilesystem_share.share.id}"
  access_type  = "ip"
  access_to    = "192.168.199.10"
  access_level = "rw"
}

resource "vkcs_sharedfilesystem_share_access" "share_access_2" {
  share_id     = "${vkcs_sharedfilesystem_share.share.id}"
  access_type  = "ip"
  access_to    = "192.168.199.11"
  access_level = "rw"
}
