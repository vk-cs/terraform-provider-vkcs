resource "vkcs_db_instance" "pg_with_backup" {
  name              = "pg-with-backup-tf-example"
  availability_zone = "GZ1"
  flavor_id         = data.vkcs_compute_flavor.basic.id

  datastore {
    type    = "postgresql"
    version = "16"
  }

  network {
    uuid = vkcs_networking_network.db.id
  }

  size        = 8
  volume_type = "ceph-ssd"
  disk_autoexpand {
    autoexpand    = true
    max_disk_size = 1000
  }

  backup_schedule {
    name           = "three_hours_backup_tf_example"
    start_hours    = 16
    start_minutes  = 20
    interval_hours = 3
    keep_count     = 3
  }

  depends_on = [
    vkcs_networking_router_interface.db,
  ]
}
