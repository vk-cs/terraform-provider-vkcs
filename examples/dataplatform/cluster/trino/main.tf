locals {
  # In this example Iceberg has two entrypoints: one is direct to postgres, another one is to postgres via pgbouncer
  # Here we do not bother what connection we get for Trino
  # In real life select appropriate entrypoint by service attributes
  # See Data Platform guide for an example
  iceberg_host      = regex(".*@([^:/]+):([0-9]+).*", vkcs_dataplatform_cluster.iceberg.info.services[0].connection_string)[0]
  iceberg_port      = regex(".*@([^:/]+):([0-9]+).*", vkcs_dataplatform_cluster.iceberg.info.services[0].connection_string)[1]
  iceberg_host_port = "${local.iceberg_host}:${local.iceberg_port}"
}

resource "vkcs_dataplatform_cluster" "trino" {
  name            = "trino-tf-example"
  description     = "Trino example instance."
  product_name    = "trino"
  product_version = "0.468.1"

  network_id        = module.network.networks[0].id
  subnet_id         = module.network.networks[0].subnets[0].id
  availability_zone = "GZ1"

  # in order to create a trino in the same cluster as the iceberg
  stack_id = vkcs_dataplatform_cluster.iceberg.stack_id

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
    },
  ]
  configs = {
    users = [
      {
        username = "example"
        password = random_password.trino_example.result
      }
    ]
    warehouses = [
      {
        name = "example"
        connections = [
          {
            name = "iceberg"
            plug = "iceberg-metastore-int"
            settings = [
              {
                alias = "hostname"
                value = local.iceberg_host_port
              },
              {
                alias = "username"
                value = "trino"
              },
              {
                alias = "password"
                value = random_password.iceberg_trino.result
              },
              {
                alias = "db_name"
                value = "example"
              },
              {
                alias = "s3_bucket"
                value = "dataplatform-tf-example"
              },
              {
                alias = "s3_folder"
                # Just a unique folder to not mess up other examples
                value = module.network.router_id
              },
              {
                alias = "catalog"
                value = "iceberg"
              }
            ]
          }
        ]
      }
    ]
    maintenance = {
      start = "0 22 * * *"
      crontabs = [
        {
          name  = "maintenance"
          start = "0 19 * * *"
          settings = [
            {
              alias = "duration"
              # Overwrite default value
              value = "600"
            },
          ]
        }
      ]
    }
  }

  # If you create networking in the same bundle of resources with Data Platform resource
  # add dependency on corresponding vkcs_networking_router_interface resource.
  # However this is not required if you set up networking witth terraform-vkcs-network module.
  # depends_on = [vkcs_networking_router_interface.db]
}
