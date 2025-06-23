resource "random_string" "random" {
  length  = 5
  upper   = false
  special = false
}

locals {
  s3_bucket = "tfexample-resource-${random_string.random.result}"
}
