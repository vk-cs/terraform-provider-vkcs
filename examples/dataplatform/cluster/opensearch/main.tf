resource "vkcs_dataplatform_cluster" "opensearch" {
  name        = "opensearch-tf-example"
  description     = "OpenSearch example instance."
  product_name    = "opensearch"
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
