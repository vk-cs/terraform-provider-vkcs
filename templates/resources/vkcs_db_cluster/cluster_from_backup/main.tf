resource "vkcs_db_cluster" "mydb-cluster" {
  name        = "mydb-cluster"

  datastore {
    type    = "postgresql"
    version = "12"
  }

  cluster_size = 3

  flavor_id   = "9e931469-1490-489e-88af-29a289681c53"

  volume_size = 10
  volume_type = "MS1"

  network {
    uuid = "3ee9b184-3311-4d85-840b-7a9c48e7beac"
  }

  restore_point {
    backup_id = "backup_id"
  }
}
