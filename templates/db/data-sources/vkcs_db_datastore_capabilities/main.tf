data "vkcs_db_datastore_capabilities" "postgres_caps" {
  datastore_name       = data.vkcs_db_datastore.postgres
  datastore_version_id = local.pg_v14_version_id
}

output "postgresql_capabilities" {
  value       = data.vkcs_db_datastore_capabilities.postgres_caps.capabilities
  description = "Available capabilities of the latest version of PostgreSQL datastore."
}
