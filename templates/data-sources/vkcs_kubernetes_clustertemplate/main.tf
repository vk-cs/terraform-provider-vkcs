data "vkcs_kubernetes_clustertemplate" "example_template" {
  name = "clustertemplate_1"
}

output "example_template_id" {
  value = "${data.vkcs_kubernetes_clustertemplate.example_template.id}"
}
