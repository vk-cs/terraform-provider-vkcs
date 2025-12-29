resource "vkcs_dataplatform_cluster" "example" {
  name        = "tf-basic-clickhouse"
  description = "tf-basic-clickhouse"
  network_id  = vkcs_networking_network.db.id
  subnet_id   = vkcs_networking_subnet.db.id

  product_name    = "clickhouse"
  product_version = "24.3.0"

  availability_zone = "GZ1"

  configs = {
    users = [
      {
        username = "vkdata"
        # Example only. Do not use in production.
        # Sensitive values must be provided securely and not stored in manifests.
        password = "Test_p#ssword-12-3"
        role     = "dbOwner"
      }
    ]

    maintenance = {
      start = "0 22 * * *"

      backup = {
        full = {
          keep_count = 10
          start      = "0 22 * * *"
        }

        incremental = {
          keep_count = 1
          start      = "0 22 * * *"
        }
      }
    }

    warehouses = [
      {
        name  = "db_customer"
        users = ["vkdata"]
      }
    ]
  }

  pod_groups = [
    {
      count = 1
      name  = "clickhouseKeeper"
      resource = {
        cpu_request = 0.5
        ram_request = 1
      }
      volumes = {
        data = {
          storage_class_name = "ceph-ssd"
          storage            = "5"
          count              = 1
        }
      }
    },
    {
      count = 1
      name  = "clickhouse"
      resource = {
        cpu_request = 2
        ram_request = 4
      }
      volumes = {
        data = {
          storage_class_name = "ceph-ssd"
          storage            = "5"
          count              = 1
        }
      }
    }
  ]

  depends_on = [vkcs_networking_router_interface.db]
}
