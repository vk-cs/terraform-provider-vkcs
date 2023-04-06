data "vkcs_db_datastore_parameters" "mysql_params" {
	datastore_name = data.vkcs_db_datastore.mysql
	datastore_version_id = data.vkcs_db_datastore.mysql.versions[0].id
}

output "mysql_parameters" {
	value = data.vkcs_db_datastore_parameters.mysql_params.parameters
	description = "Available configuration parameters of the latest version of MySQL datastore."
}
