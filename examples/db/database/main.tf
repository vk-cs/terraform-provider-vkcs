resource "vkcs_db_database" "mysql-db" {
  name        = "testdb"
  dbms_id     = vkcs_db_instance.mysql.id
  charset     = "utf8"
  collate     = "utf8_general_ci"
}
