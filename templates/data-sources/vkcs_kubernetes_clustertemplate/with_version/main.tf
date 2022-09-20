data "vkcs_kubernetes_clustertemplate" "example_template_by_version" {
  version = "1.20.4"
}

output "example_template_id" {
  value = "${data.vkcs_kubernetes_clustertemplate.example_template_by_version.id}"
}
