resource "vkcs_dataplatform_cluster" "spark" {
  name            = "spark-tf-example"
  description     = "Spark example instance."
  product_name    = "spark"
  product_version = "3.5.1"

  network_id        = vkcs_networking_network.db.id
  subnet_id         = vkcs_networking_subnet.db.id
  availability_zone = "GZ1"

  pod_groups = [
    {
      name  = "sparkconnect"
      count = 1
      # In fact this group requires more resorces than spark template defines
      # This are the minimum resource requirements to run Spark
      resource = {
        cpu_request = "10"
        ram_request = "10"
      }
    },
    {
      name  = "sparkhistory"
      count = 1
      resource = {
        cpu_request = "0.5"
        ram_request = "1"
      }
      volumes = {
        "data" = {
          storage_class_name = "ceph-ssd"
          storage            = "5"
          count              = 1
        }
      }
    },
    {
      name  = "sparkproxy"
      count = 1
      resource = {
        "cpu_request" : "0.5",
        "ram_request" : "0.5",
      },
    },
    {
      name  = "sparkintlogs"
      count = 1
      resource = {
        "cpu_request" : "0.5",
        "ram_request" : "0.5",
      },
    },
    {
      name  = "authservice"
      count = 1
      resource = {
        "cpu_request" : "0.5",
        "ram_request" : "0.5",
      },
    },
    {
      name  = "sparkjobs"
      count = 1
      resource = {
        "cpu_request" : "4.0",
        "ram_request" : "8.0",
      },
    },
  ]
  configs = {
    settings = [
      {
        alias = "sparkproxy.spark_version"
        value = "spark-py-3.5.1:v3.5.1.2"
      }
    ]
    warehouses = [{
      name = "spark"
      connections = [
        {
          name = "s3_int"
          plug = "s3-int"
          settings = [
            {
              alias = "s3_bucket"
              # Data Platform Spark requires unique bucket for this type of connection
              value = format("spark-tf-example-%s", vkcs_networking_router.router.id)
            },
            {
              alias = "s3_folder"
              value = "spark"
            }
          ]
        },
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
    }]
    maintenance = {
      start = "0 0 1 * *"
    }
  }

  # If you create networking in the same bundle of resources with Data Platform resource
  # add dependency on corresponding vkcs_networking_router_interface resource.
  # However this is not required if you set up networking witth terraform-vkcs-network module.
  depends_on = [vkcs_networking_router_interface.db]
}
