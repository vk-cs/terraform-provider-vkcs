resource "vkcs_db_instance" "db_instance" {
  name = "db-instance"

  availability_zone = "GZ1"

  datastore {
    type    = "postgresql"
    version = "11"
  }

  flavor_id = data.vkcs_compute_flavor.db.id

  size        = 10
  volume_type = "ceph-ssd"
  disk_autoexpand {
    autoexpand    = true
    max_disk_size = 1000
  }
  wal_volume {
    size        = 10
    volume_type = "ceph-ssd"
  }

  wal_disk_autoexpand {
    autoexpand    = true
    max_disk_size = 20
  }

  network {
    uuid = vkcs_networking_network.db.id
  }

  depends_on = [
    vkcs_networking_router_interface.db
  ]
}
