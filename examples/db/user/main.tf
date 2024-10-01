resource "vkcs_db_user" "mysql-user" {
  name        = "testuser"
  password    = "Test_p@ssword-12-3"

  dbms_id     = vkcs_db_instance.mysql.id

  databases   = [vkcs_db_database.mysql-db-1.name, vkcs_db_database.mysql-db-2.name]
}
