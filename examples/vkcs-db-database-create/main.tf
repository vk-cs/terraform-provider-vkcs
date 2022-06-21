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
}

resource "vkcs_db_database" "db-database" {
  name        = "testdb"
  dbms_id     = vkcs_db_instance.db-instance.id
  charset     = "utf8"
  collate     = "utf8_general_ci"
}
