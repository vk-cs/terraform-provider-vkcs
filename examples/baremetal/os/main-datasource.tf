data "vkcs_baremetal_os" "ubuntu" {
  name    = "ubuntu"
  version = "24.04"
}

output "flavor_output" {
  value = {
    id      = data.vkcs_baremetal_os.ubuntu.id
    name    = data.vkcs_baremetal_os.ubuntu.name
    version = data.vkcs_baremetal_os.ubuntu.version
  }
}
