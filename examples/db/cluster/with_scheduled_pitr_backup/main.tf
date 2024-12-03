resource "vkcs_db_cluster" "pg_with_backup" {
  name = "pg-with-backup-tf-example"

  availability_zone = "GZ1"
  datastore {
    type    = "postgresql"
    version = "16"
  }

  cluster_size = 3

  flavor_id = data.vkcs_compute_flavor.basic.id

  volume_size = 10
  volume_type = "ceph-ssd"

  network {
    uuid = vkcs_networking_network.db.id
  }

  backup_schedule {
    name           = "three_hours_backup_tf_example"
    start_hours    = 16
    start_minutes  = 20
    interval_hours = 3
    keep_count     = 3
  }

  depends_on = [
    vkcs_networking_router_interface.db
  ]
}
