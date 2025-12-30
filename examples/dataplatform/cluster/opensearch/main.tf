resource "vkcs_dataplatform_cluster" "example" {
  name        = "tf-basic-opensearch"
  description = "tf-basic-opensearch"
  network_id  = vkcs_networking_network.db.id
  subnet_id   = vkcs_networking_subnet.db.id

  product_name    = "opensearch"
  product_version = "2.15.0"

  availability_zone = "GZ1"

  configs = {
    users = [
      {
        username = "vkdata"
        # Example only. Do not use in production.
        # Sensitive values must be provided securely and not stored in manifests.
        password = "Test_p#ssword-12-3"
        role     = "administrator"
      }
    ]

    maintenance = {
      start = "0 22 * * *"

      backup = {
        full = {
          keepCount = 10
          start     = "0 22 * * *"
        }

        incremental = {
          keepCount = 1
          start     = "0 22 * * *"
        }
      }
    }

    warehouses = [
      {
        name  = "opensearch"
        users = ["vkdata"]
      }
    ]
  }

  pod_groups = [
    {
      count = 0
      name  = "bootstrap"
      resource = {
        cpu_request = 0.5
        ram_request = 1
        cpu_limit = 0.5
        ram_limit = 1
      }
    },
    {
      count = 1
      name  = "dashboards"
      resource = {
        cpu_request = 0.5
        ram_request = 1
        cpu_limit = 0.5
        ram_limit = 1
      }
    },
    {
      count = 3
      name  = "masters"
      resource = {
        cpu_request = 0.5
        ram_request = 2
        cpu_limit = 0.5
        ram_limit = 2
      }
      volumes = {
        data = {
          storage_class_name = "ceph-ssd"
          storage            = "30"
          count              = 3
        }
      }
    }
  ]

  depends_on = [vkcs_networking_router_interface.db]
}
