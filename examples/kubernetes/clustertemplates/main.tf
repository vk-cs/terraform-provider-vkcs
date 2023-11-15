data "vkcs_kubernetes_clustertemplates" "templates" {}

output "available_templates_by_name" {
  value = [
    for template in data.vkcs_kubernetes_clustertemplates.templates.cluster_templates :
    { name = template.name, version = template.version }
  ]
}
