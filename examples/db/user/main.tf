resource "vkcs_db_user" "mysql_user" {
  name     = "testuser"
  password = "Test_p@ssword-12-3"

  dbms_id = vkcs_db_instance.mysql.id

  databases = [vkcs_db_database.mysql_db_1.name, vkcs_db_database.mysql_db_2.name]
}
