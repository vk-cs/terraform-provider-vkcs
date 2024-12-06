resource "vkcs_db_database" "mysql_db_1" {
  name    = "testdb_1"
  dbms_id = vkcs_db_instance.mysql.id
  charset = "utf8"
  collate = "utf8_general_ci"
}

resource "vkcs_db_database" "mysql_db_2" {
  name    = "testdb_2"
  dbms_id = vkcs_db_instance.mysql.id
  charset = "utf8"
  collate = "utf8_general_ci"
}
