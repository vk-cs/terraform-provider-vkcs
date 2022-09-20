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

resource "vkcs_sharedfilesystem_securityservice" "securityservice_1" {
  name        = "security"
  description = "created by terraform"
  type        = "active_directory"
  server      = "192.168.199.10"
  dns_ip      = "192.168.199.10"
  domain      = "example.com"
  user        = "joinDomainUser"
  password    = "s8cret"
}

resource "vkcs_sharedfilesystem_sharenetwork" "sharenetwork_1" {
  name              = "test_sharenetwork_secure"
  description       = "share the secure love"
  neutron_net_id    = "${vkcs_networking_network.network_1.id}"
  neutron_subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
  security_service_ids = [
	"${vkcs_sharedfilesystem_securityservice.securityservice_1.id}",
  ]
}

resource "vkcs_sharedfilesystem_share" "share_1" {
  name             = "cifs_share"
  share_proto      = "CIFS"
  size             = 1
  share_network_id = "${vkcs_sharedfilesystem_sharenetwork.sharenetwork_1.id}"
}

resource "vkcs_sharedfilesystem_share_access" "share_access_1" {
  share_id     = "${vkcs_sharedfilesystem_share.share_1.id}"
  access_type  = "user"
  access_to    = "windows"
  access_level = "ro"
}

resource "vkcs_sharedfilesystem_share_access" "share_access_2" {
  share_id     = "${vkcs_sharedfilesystem_share.share_1.id}"
  access_type  = "user"
  access_to    = "linux"
  access_level = "rw"
}
