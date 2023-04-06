data "vkcs_db_datastores" "datastores" {}

output "available_datastores" {
    value = data.vkcs_db_datastores.datastores.datastores
    description = "List of datastores that are available within VKCS."
}
