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
  disk_autoexpand {
    autoexpand    = true
    max_disk_size = 1000
  }

  network {
    uuid = vkcs_networking_network.db.id
  }

  capabilities {
    name = "node_exporter"
    settings = {
      "listen_port" : "9100"
    }
  }

  depends_on = [
    vkcs_networking_router_interface.db
  ]
}

resource "vkcs_db_instance" "db-replica" {
  name        = "db-instance-replica"
  datastore {
    type    = "mysql"
    version = "5.7"
  }
  replica_of  = vkcs_db_instance.db-instance.id

  flavor_id   = data.vkcs_compute_flavor.db.id

  size        = 8
  volume_type = "ceph-ssd"

  network {
    uuid = vkcs_networking_network.db.id
  }
}
