resource "vkcs_vpnaas_ike_policy" "data_center" {
  name           = "key-policy-tf-example"
  description    = "Policy that restricts remote working users to connect to our data ceneter over VPN"
  auth_algorithm = "sha256"
}
