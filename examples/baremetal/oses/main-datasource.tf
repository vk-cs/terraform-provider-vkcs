data "vkcs_baremetal_oses" "main" {}

output "oses_output" {
  value = {
    oses = data.vkcs_baremetal_oses.main
  }
}
