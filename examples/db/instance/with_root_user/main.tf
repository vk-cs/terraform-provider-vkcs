resource "vkcs_db_instance" "db-instance" {
  name        = "db-instance"

  availability_zone = "GZ1"

  datastore {
    type    = "postgresql"
    version = "13"
  }

  flavor_id   = data.vkcs_compute_flavor.db.id
  
  size        = 8
  volume_type = "ceph-ssd"
  network {
    uuid = vkcs_networking_network.db.id
  }

  root_enabled  = true

  depends_on = [
    vkcs_networking_router_interface.db
  ]
}

output "root_user_password" {
  value     = vkcs_db_instance.db-instance.root_password
  sensitive = true
}
