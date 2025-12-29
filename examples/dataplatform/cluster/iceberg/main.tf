resource "vkcs_dataplatform_cluster" "iceberg" {
  name            = "iceberg-tf-example"
  description     = "Iceberg example instance."
  product_name    = "iceberg-metastore"
  product_version = "17.2.0"

  network_id        = module.network.networks[0].id
  subnet_id         = module.network.networks[0].subnets[0].id
  availability_zone = "GZ1"

  pod_groups = [
    {
      name  = "postgres"
      count = 1
      resource = {
        cpu_request = "0.5"
        ram_request = "1"
      }
      volumes = {
        "data" = {
          storage_class_name = "ceph-ssd"
          storage            = "10"
          count              = 1
        }
        "wal" = {
          storage_class_name = "ceph-ssd"
          storage            = "10"
          count              = 1
        }
      }
    },
    {
      name = "bouncer"
      # bouncer could be enabled later
      count = 0
      # even though bouncer is disabled, we need to specify its resource request
      resource = {
        cpu_request = "0.2"
        ram_request = "1"
      }
    },
  ]
  configs = {
    users = [
      {
        username = "owner"
        password = random_password.iceberg_owner.result
        role     = "dbOwner"
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

  # If you create networking in the same bundle of resources with Data Platform resource
  # add dependency on corresponding vkcs_networking_router_interface resource.
  # However this is not required if you set up networking witth terraform-vkcs-network module.
  # depends_on = [vkcs_networking_router_interface.db]
}
