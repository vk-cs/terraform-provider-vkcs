resource "vkcs_dataplatform_cluster" "basic_trino" {
  name              = "tf-basic-trino"
  description       = "tf-basic-description"
  product_name      = "trino"
  product_version   = "0.468.1"
  availability_zone = "GZ1"

  network_id = vkcs_networking_network.app.id
  subnet_id  = vkcs_networking_subnet.app.id

  configs = {

    maintenance = {
      start = "0 22 * * *"
      crontabs = [
        {
          name  = "maintenance"
          start = "0 19 * * *"
          settings = [
            {
              alias = "duration"
              value = "3600"
            }
          ]
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
                value = ""
              },
              {
                alias = "hostname"
                value = ""
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
                value = ""
              },
              {
                alias = "s3_folder"
                value = ""
              },
              {
                alias = "catalog"
                value = ""
              }# ,
              # {
              #   alias = "parquet.use-bloom-filter"
              #   value = "true"
              # },
              # {
              #   alias = "parquet.ignore-statistics"
              #   value = "false"
              # },
              # {
              #   alias = "parquet.writer.validation-percentage"
              #   value = "5"
              # },
              # {
              #   alias = "parquet.writer.page-size"
              #   value = "1MB"
              # },
              # {
              #   alias = "parquet.writer.page-value-count"
              #   value = "80000"
              # },
              # {
              #   alias = "parquet.writer.block-size"
              #   value = "128MB"
              # },
              # {
              #   alias = "parquet.writer.batch-size"
              #   value = "10000"
              # },
              # {
              #   alias = "parquet.max-read-block-row-count"
              #   value = "8192"
              # },
              # {
              #   alias = "parquet.small-file-threshold"
              #   value = "10MB"
              # },
              # {
              #   alias = "parquet.experimental.vectorized-decoding.enabled"
              #   value = "true"
              # },
              # {
              #   alias = "iceberg.register-table-procedure.enabled"
              #   value = "false"
              # },
              # {
              #   alias = "iceberg.add-files-procedure.enabled"
              #   value = "false"
              # }
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

  depends_on = [vkcs_dataplatform_cluster.basic_iceberg]
}
