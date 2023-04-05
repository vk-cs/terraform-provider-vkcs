// Workaround for Terraform version < 1.4 because there is no
// terraform_data resource
resource "vkcs_networking_secgroup" "base" {
  name = "test"
  description = "Dummy security group for saving public DNS zone name"
}

locals {
  zone_name = "example-${vkcs_networking_secgroup.base.id}.com"
}
