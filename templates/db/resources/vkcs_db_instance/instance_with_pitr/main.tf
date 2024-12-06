resource "vkcs_db_instance" "db_instance" {
  name = "db-instance"

  datastore {
    type    = "postgresql"
    version = "13"
  }

  floating_ip_enabled = true

  flavor_id         = "9e931469-1490-489e-88af-29a289681c53"
  availability_zone = "MS1"

  size        = 8
  volume_type = "MS1"
  disk_autoexpand {
    autoexpand    = true
    max_disk_size = 1000
  }

  network {
    uuid = "3ee9b184-3311-4d85-840b-7a9c48e7beac"
  }

  capabilities {
    name = "node_exporter"
  }

  capabilities {
    name = "postgres_extensions"
  }

  backup_schedule {
    name           = three_hours_backup
    start_hours    = 16
    start_minutes  = 20
    interval_hours = 3
    keep_count     = 3
  }
}
