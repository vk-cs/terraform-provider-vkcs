resource "vkcs_dataplatform_cluster" "example" {
  name            = "clickhouse-tf-example"
  description     = "ClickHouse example instance."
  product_name    = "clickhouse"
  product_version = "24.3.0"

  network_id        = vkcs_networking_network.db.id
  subnet_id         = vkcs_networking_subnet.db.id
  availability_zone = "GZ1"

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
    },
  ]
  configs = {
    users = [
      {
        username = "owner"
        password = random_password.clickhouse_owner.result
        role     = "dbOwner"
      },
    ]
    warehouses = [{
      name = "example"
    }]
    maintenance = {
      # Set start om maintenance the same as start of full backup.
      # Otherwise you get unpredictable behavior of interaction between
      # Terraform, VKCS Terraform provider and Data Platform API.
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
  }

  # If you create networking in the same bundle of resources with Data Platform resource
  # add dependency on corresponding vkcs_networking_router_interface resource.
  # However this is not required if you set up networking witth terraform-vkcs-network module.
  depends_on = [vkcs_networking_router_interface.db]
}
