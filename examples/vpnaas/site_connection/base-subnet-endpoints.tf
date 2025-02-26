resource "vkcs_vpnaas_endpoint_group" "subnet_hosts" {
  type = "cidr"
  endpoints = [
    vkcs_networking_subnet.app.cidr
  ]
}
