resource "vkcs_dataplatform_cluster" "trino" {
  name            = "trino-tf-guide"
  description     = "Trino example instance from Data Platform guide."
  product_name    = "trino"
  product_version = "0.468.1"

  network_id        = vkcs_networking_network.db.id
  subnet_id         = vkcs_networking_subnet.db.id
  availability_zone = "GZ1"

  # In order to create Trino in the same cluster stack as the ClickHouse.
  stack_id = vkcs_dataplatform_cluster.clickhouse.stack_id
  # This argument must be the same for all products in the same cluster stack.
  floating_ip_pool = "auto"

  pod_groups = [
    {
      name  = "coordinator"
      count = 1
      resource = {
        cpu_request = "2.0"
        ram_request = "4.0"
      }
    },
    {
      name  = "worker"
      count = 1
      resource = {
        cpu_request = "2.0"
        ram_request = "4.0"
      }
    }
  ]
  configs = {
    users = [
      {
        username = "example"
        password = random_password.trino_example.result
      }
    ]
    warehouses = [{
      # For some Data Platform product value of the `name` argument has no sense
      # but is fixed. So you must set exactly this value.
      # Otherwise you get unpredictable behavior of interaction between
      # Terraform, VKCS Terraform provider and Data Platform API.
      name = "trino"
      connections = [
        {
          name = "clickhouse"
          plug = "clickhouse"
          settings = [
            {
              alias = "hostname"
              value = "${local.clickhouse_tcp.host}:${local.clickhouse_tcp.port}"
            },
            {
              alias = "username"
              value = "trino"
            },
            {
              alias = "password"
              value = random_password.clickhouse_trino.result
            },
            {
              alias = "ssl"
              value = "false"
            },
            {
              alias = "db_name"
              value = "clickhouse"
            },
            {
              alias = "catalog"
              value = "clickhouse"
            },
          ]
        },
      ]
    }]
    maintenance = {
      start = "0 22 * * *"
      crontabs = [
        {
          name  = "maintenance"
          start = "0 19 * * *"
        }
      ]
    }
  }
}
