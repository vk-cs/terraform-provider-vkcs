resource "vkcs_vpnaas_endpoint_group" "subnet_hosts" {
  type      = "subnet"
  endpoints = [vkcs_networking_subnet.app.id]
  sdn       = "neutron"
}
