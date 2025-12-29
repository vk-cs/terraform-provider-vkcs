module "network" {
  source = "https://github.com/vk-cs/terraform-vkcs-network/archive/refs/tags/v0.0.3.zip//terraform-vkcs-network-0.0.3"

  tags             = ["tf-example"]
  name             = "trino-tf-example"
  sdn              = "sprut"
  external_network = "internet"

  networks = [{
    subnets = [{
      cidr = "192.168.199.0/24"
    }]
  }]
}
