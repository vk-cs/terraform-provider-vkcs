data "vkcs_dataplatform_templates" "cluster-templates" {
}

resource "vkcs_dataplatform_cluster" "basic-opensearch" {
  cluster_template_id = data.vkcs_dataplatform_templates.cluster-templates.templates[index(data.vkcs_dataplatform_templates.cluster-templates.*.product_name, "spark")].id
  name                = "tf-basic-spark"
  network_id          = "d2fad739-1b10-4dc8-9b2c-c246d7a7cc69"
  subnet_id           = "3a744943-fcc1-4a85-a96b-3dc4fff71885"
  product_name        = "spark"
  product_version     = "3.5.1"

  availability_zone = "UD2"
  configs = {
    settings = [
      {
        alias = "sparkproxy.spark_version"
        value = "spark-py-3.5.1:v3.5.1.2"
      }
    ]
    users = [
      {
        user     = "user"
        password = "somepa55word!"
      }
    ]
    maintenance = {
      start = "0 0 1 * *"
    }
    warehouses = [
      {
        name = "db_customer"
        users = [
          "user"
        ]
        connections = [
          {
            name = "s3_int"
            plug = "s3-int"
            settings = [
              {
                alias = "s3_bucket"
                value = "really-cool-bucket"
              },
              {
                alias = "s3_folder"
                value = "folder"
              }
            ]
          },
          {
            name = "postgres"
            plug = "postgresql"
            settings = [
              {
                alias = "db_name"
                value = "db"
              },
              {
                alias = "hostname"
                value = "database.com:5432"
              },
              {
                alias = "username"
                value = "db"
              },
              {
                alias = "password"
                value = "db"
              }
            ]
          }
        ]
      }
    ]
  }
  pod_groups = [
    {
      count = 1
      resource = {
        cpu_request = "8"
        ram_request = "8"
      }
      pod_group_template_id = "6a8e3515-d0d6-40f9-826e-e33dbe141485"
    },
    {
      count = 1
      resource = {
        cpu_request = "0.5"
        ram_request = "0.5"
      }
      volumes = [
        {
          type = "data"
          storage_class_name = "ceph"
          storage = "5"
          count = 1
        }
      ]
      pod_group_template_id = "498c4bf6-1e3e-4a06-b8b7-d60337e85dc1"
    }
  ]
}
