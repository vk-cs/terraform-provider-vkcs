resource "vkcs_dataplatform_cluster" "clickhouse" {
  name            = "clickhouse-tf-guide"
  description     = "ClickHouse example instance from Data Platform guide."
  product_name    = "clickhouse"
  product_version = "25.3.0"

  network_id        = module.network.networks[0].id
  subnet_id         = module.network.networks[0].subnets[0].id
  availability_zone = "GZ1"
  # Enable public access to simplify testing of the product.
  floating_ip_pool = "auto"

  pod_groups = [
    # Omit settings for clickhouseKeeper pod group to illustrate
    # how Data Platform handle this with default settings.
    # NOTE: If you omit settings for a pod group you cannot scale
    # the pod group later.
    # Increase ram_request and storage values for clickhouse pod group
    # against default settings.
    {
      count = 3
      name  = "clickhouse"
      resource = {
        cpu_request = "2.0"
        ram_request = "8.0"
      }
      volumes = {
        data = {
          storage_class_name = "ceph-ssd"
          storage            = "150"
          count              = 1
        }
      }
    },
  ]
  configs = {
    settings = [
      # Increase value of the setting against default one.
      {
        alias = "clickhouse.background_common_pool_size"
        value = 10
      },
    ]
    users = [
      {
        username = "owner"
        password = random_password.clickhouse_owner.result
        role     = "dbOwner"
      },
      {
        username = "trino"
        password = random_password.clickhouse_trino.result
        role     = "readOnly"
      },
    ]
    warehouses = [
      # Define database name.
      {
        name = "clickhouse"
      },
    ]
    maintenance = {
      # Set start om maintenance the same as start of full backup.
      # Otherwise you get unpredictable behavior of interaction between
      # Terraform, VKCS Terraform provider and Data Platform API.
      start = "0 1 * * 0"
      backup = {
        full = {
          keep_count = 5
          start      = "0 1 * * 0"
        }
        incremental = {
          keep_count = 7
          start      = "0 1 * * 1-6"
        }
      }
    }
  }

  # If you create networking in the same bundle of resources with Data Platform resource
  # add dependency on corresponding vkcs_networking_router_interface resource.
  # However this is not required if you set up networking witth terraform-vkcs-network module.
  # depends_on = [vkcs_networking_router_interface.db]
}
