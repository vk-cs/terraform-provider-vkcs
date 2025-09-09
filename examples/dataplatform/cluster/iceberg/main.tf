resource "time_sleep" "wait_20_seconds" {
  depends_on = [vkcs_networking_router_interface.db]

  create_duration = "20s"
}

resource "vkcs_dataplatform_cluster" "basic_iceberg" {
  name        = "tf-basic-iceberg"
  description = "tf-basic-iceberg-description"
  network_id  = vkcs_networking_network.db.id
  subnet_id   = vkcs_networking_subnet.db.id

  product_name    = "iceberg-metastore"
  product_version = "17.2.0"

  availability_zone = "GZ1"
  configs = {
    maintenance = {
      start = "0 0 1 * *"
      backup = {
        full = {
          keep_time = 10
          start     = "0 0 1 * *"
        }
      }
    }
    warehouses = [
      {
        name = "metastore",
      }
    ]
    users = [
      {
        username = "vkdata"
        password = "Test_p#ssword-12-3"
        role     = "dbOwner"
      }
    ]
  }
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
      name  = "bouncer"
      count = 0
      # even though bouncer is disabled, we need to specify its resource request
      resource = {
        cpu_request = "0.2"
        ram_request = "1"
      }
    }
  ]

  # need to wait for network access to appear after creation of router_interface
  depends_on = [time_sleep.wait_20_seconds]
}
