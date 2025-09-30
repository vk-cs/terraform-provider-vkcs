resource "vkcs_dataplatform_cluster" "basic_trino" {
  name              = "tf-basic-trino"
  description       = "tf-basic-description"
  network_id        = vkcs_networking_network.db.id
  subnet_id         = vkcs_networking_subnet.db.id
  product_name      = "trino"
  product_version   = "0.449.0"
  availability_zone = "GZ1"

  configs = {
    maintenance = {
      start = "0 0 1 * *"
    }

    warehouses = [
      {
        name = "trino"
        connections = [
          {
            name = "postgres"
            plug = "postgresql"
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
      name  = "coordinator"
      count = 1
      resource = {
        cpu_request = "4.0"
        ram_request = "16.0"
      }
    },
    {
      name  = "worker"
      count = 3
      resource = {
        cpu_request = "4.0"
        ram_request = "16.0"
      }
    }
  ]

}
