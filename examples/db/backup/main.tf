resource "vkcs_db_backup" "mysql-backup" {
  name = "mssql-backup"
  dbms_id = vkcs_db_instance.mysql.id
}
