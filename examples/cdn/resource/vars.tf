resource "random_string" "random" {
  length = 5
  upper = false
  special = false
}

locals {
  cname = "tfexample-resource-${random_string.random.result}.vk.com"
}
