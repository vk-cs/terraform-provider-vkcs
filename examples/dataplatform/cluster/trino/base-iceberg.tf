resource "vkcs_dataplatform_cluster" "iceberg" {
  name            = "iceberg-tf-example"
  description     = "Iceberg example instance for Trino example."
  product_name    = "iceberg-metastore"
  product_version = "17.2.0"

  network_id        = vkcs_networking_network.db.id
  subnet_id         = vkcs_networking_subnet.db.id
  availability_zone = "GZ1"

  pod_groups = []
  configs = {
    users = [
      {
        username = "owner"
        password = random_password.iceberg_owner.result
        role     = "dbOwner"
      },
      {
        username = "trino"
        password = random_password.iceberg_trino.result
        role     = "common"
      },
    ]
    warehouses = [
      {
        name = "example"
      }
    ]
    maintenance = {
      start = "0 0 1 * *"
      backup = {
        full = {
          keep_time = 10
          start     = "0 0 1 * *"
        }
      }
    }
  }
  depends_on = [vkcs_networking_router_interface.db]
}
