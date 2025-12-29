resource "random_password" "iceberg_owner" {
  length           = 20
  min_lower        = 1
  min_upper        = 1
  min_numeric      = 1
  min_special      = 1
  override_special = "!?%#/()-+*"
}
