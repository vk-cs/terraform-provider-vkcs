resource "vkcs_db_instance" "db-instance" {
  name        = "db-instance"

  availability_zone = "GZ1"

  datastore {
    type    = "mysql"
    version = "5.7"
  }

  flavor_id   = data.vkcs_compute_flavor.db.id

  size        = 8
  volume_type = "ceph-ssd"  

  network {
    uuid = vkcs_networking_network.db.id
  }

  depends_on = [
    vkcs_networking_router_interface.db
  ]
}

resource "vkcs_db_database" "db-database" {
  name        = "testdb"
  dbms_id     = vkcs_db_instance.db-instance.id
  charset     = "utf8"
  collate     = "utf8_general_ci"
}

resource "vkcs_db_user" "db-user" {
  name        = "testuser"
  password    = "SomePass1_"

  dbms_id     = vkcs_db_instance.db-instance.id

  databases   = [vkcs_db_database.db-database.name]
}
