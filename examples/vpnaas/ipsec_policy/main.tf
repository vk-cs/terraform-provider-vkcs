resource "vkcs_vpnaas_ipsec_policy" "data_center" {
  name        = "database-key-policy-tf-example"
  description = "Policy that restricts remote working users to connect to our data ceneter over VPN"
  lifetime {
    units = "seconds"
    value = 3600
  }
}
