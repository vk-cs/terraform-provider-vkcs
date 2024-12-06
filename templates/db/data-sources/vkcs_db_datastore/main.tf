data "vkcs_db_datastore" "datastore" {
  name = "mysql"
}

output "mysql_versions" {
  value       = data.vkcs_db_datastore.datastore.versions
  description = "List of versions of MySQL that are available within VKCS."
}
