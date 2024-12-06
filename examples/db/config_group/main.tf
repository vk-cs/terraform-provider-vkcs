resource "vkcs_db_config_group" "mysql_80" {
  name = "db-config-group"
  datastore {
    type    = "mysql"
    version = "8.0"
  }
  values = {
    activate_all_roles_on_login : "true"
    autocommit : "1"
    block_encryption_mode : "test"
    innodb_segment_reserve_factor : "0.53"
  }
  description = "db-config-group-description"
}
