data "vkcs_baremetal_flavors" "main" {}

output "flavors_output" {
  value = {
    flavors = data.vkcs_baremetal_flavors.main
  }
}
