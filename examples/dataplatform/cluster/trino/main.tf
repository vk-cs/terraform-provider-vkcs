resource "vkcs_dataplatform_cluster" "basic_trino" {
  name              = "tf-basic-trino"
  description       = "tf-basic-description"
  network_id        = vkcs_networking_network.db.id
  subnet_id         = vkcs_networking_subnet.db.id
  product_name      = "trino"
  product_version   = "0.468.0"
  availability_zone = "GZ1"

  configs = {
    settings = [
      {
      alias = "coordinator.config.memory.heapHeadroomPerNode"
      value = "6GB"
      },
      {
      alias = "coordinator.config.query.maxMemoryPerNode"
      value = "6GB"
      },
      {
      alias = "worker.config.query.maxMemoryPerNode"
      value = "6GB"
      },
      {
      alias = "worker.config.memory.heapHeadroomPerNode"
      value = "6GB"
      },
    ]
    maintenance = {
      start    = "0 22 * * *"
      crontabs = [
        {
          name     = "maintenance"
          start    = "0 19 * * *"
          settings = [
            {
              alias = "duration"
              value = "3600"
            },
            {
              alias = "iceberg.expire-snapshots.min-retention"
              value = "7d"
            },
            {
              alias = "iceberg.remove-orphan-files.min-retention"
              value = "7d"
            },
            {
              alias = "iceberg.optimize.file-size-threshold"
              value = "128MB"
            },
            {
              alias = "iceberg.write.metadata.delete-after-commit.enabled"
              value = "true"
            },
            {
              alias = "iceberg.write.metadata.previous-versions-max"
              value = "10"
            }
          ]
        }
      ]
    }

    warehouses = [
      {
        name        = "trino"
        connections = [
          {
            name     = "postgres"
            plug     = "postgresql"
            settings = [
              {
                alias = "db_name"
                value = vkcs_db_database.postgres_db.name
              },
              {
                alias = "hostname"
                value = "${vkcs_db_instance.db_instance.ip[0]}:5432"
              },
              {
                alias = "username"
                value = vkcs_db_user.postgres_user.name
              },
              {
                alias = "password"
                value = vkcs_db_user.postgres_user.password
              }
            ]
          }
        ]
      }
    ]
  }

  pod_groups = [
    {
      name     = "coordinator"
      count    = 1
      resource = {
        cpu_request = "4.0"
        ram_request = "16.0"
      }
    },
    {
      name     = "worker"
      count    = 1
      resource = {
        cpu_request = "4.0"
        ram_request = "16.0"
      }
    }
  ]
}
