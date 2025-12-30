resource "vkcs_dataplatform_cluster" "kafka" {
  name        = "kafka-tf-example"
  description     = "Kafka example instance."
  product_name    = "trino"
  product_version = "0.468.1"

  network_id  = vkcs_networking_network.db.id
  subnet_id   = vkcs_networking_subnet.db.id
  availability_zone = "GZ1"

  configs = {
    maintenance = {
    }
    warehouses = [
    ]
    users = [
    ]
  }
  pod_groups = [
  ]

  depends_on = [vkcs_networking_router_interface.db]
}
