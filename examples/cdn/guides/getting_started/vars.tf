resource "random_string" "random" {
  length  = 5
  upper   = false
  special = false
}

locals {
  cname = "tfguide-resource-${random_string.random.result}.vk.com"
}
