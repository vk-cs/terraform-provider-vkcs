resource "random_password" "iceberg_owner" {
  length           = 20
  min_lower        = 1
  min_upper        = 1
  min_numeric      = 1
  min_special      = 1
  override_special = "!?%#/()-+*"
}

resource "random_password" "iceberg_trino" {
  length           = 20
  min_lower        = 1
  min_upper        = 1
  min_numeric      = 1
  min_special      = 1
  override_special = "!?%#/()-+*"
}

resource "random_password" "trino_example" {
  length           = 20
  min_lower        = 1
  min_upper        = 1
  min_numeric      = 1
  min_special      = 1
  override_special = "!?%#/()-+*"
}
