resource "vkcs_dataplatform_cluster" "basic-trino" {
  cluster_template_id = "17230566-bfaa-492c-b22c-33b500a17155"
  name                = "tf-basic-trino"
  network_id          = "d2fad739-1b10-4dc8-9b2c-c246d7a7cc69"
  subnet_id           = "3a744943-fcc1-4a85-a96b-3dc4fff71885"
  product_name        = "trino"
  product_version     = "0.468.0"

  availability_zone = "UD2"
  configs = {
    users = [
      {
        user     = "user"
        password = "somepa55word!"
      }
    ]
    maintenance = {
      start = "0 22 * * *"
    }
    warehouses = [
      {
        name = "db_customer"
        users = [
          "user"
        ]
      }
    ]
  }
  pod_groups = [
    {
      count = 1
      resource = {
        cpu_request = "4"
        ram_request = "6"
      }
      pod_group_template_id = "75cd99d3-7089-4372-8353-3ce675f55284"
    },
    {
      count = 0
      resource = {
        cpu_request = "4"
        ram_request = "6"
      }
      pod_group_template_id = "1da36325-1bc8-43eb-ab52-113332001cff"
    }
  ]
}