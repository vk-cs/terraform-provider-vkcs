resource "vkcs_vpnaas_endpoint_group" "allowed_hosts" {
  name = "allowed-hosts-tf-example"
  type = "cidr"
  endpoints = [
    "10.2.0.0/24",
    "10.3.0.0/24",
  ]
}
