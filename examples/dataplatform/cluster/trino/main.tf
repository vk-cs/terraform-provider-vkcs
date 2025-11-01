resource "vkcs_dataplatform_cluster" "basic_trino" {
  name              = "tf-basic-trino"
  description       = "tf-basic-description"
  product_name      = "trino"
  product_version   = "0.468.1"
  availability_zone = "GZ1"

  network_id = vkcs_networking_network.app.id
  subnet_id  = vkcs_networking_subnet.app.id

  # in order to create a trino in the same cluster as the iceberg
  stack_id = vkcs_dataplatform_cluster.basic_iceberg.stack_id

  configs = {
    users = [
      {
        username = "vkdata"
        password = "Test_p#ssword-12-3"
      }
    ]

    maintenance = {
      start = "0 22 * * *"
      crontabs = [
        {
          name  = "maintenance"
          start = "0 19 * * *"
        }
      ]
    }

    warehouses = [
      {
        name = "trino"
        connections = [
          {
            name = "iceberg"
            plug = "iceberg-metastore-int"
            settings = [
              {
                alias = "db_name"
                value = "metastore"
              },
              {
                alias = "hostname"
                value = regex(".*@([^/]+).*", vkcs_dataplatform_cluster.basic_iceberg.info.services[0].connection_string)[0]
              },
              {
                alias = "username"
                value = vkcs_dataplatform_cluster.basic_iceberg.configs.users[0].username
              },
              {
                alias = "password"
                value = vkcs_dataplatform_cluster.basic_iceberg.configs.users[0].password
              },
              {
                alias = "s3_bucket"
                value = local.s3_bucket
              },
              {
                alias = "s3_folder"
                value = "s3_folder"
              },
              {
                alias = "catalog"
                value = "catalog"
              }
            ]
          }
        ]
      }
    ]
  }

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
}
