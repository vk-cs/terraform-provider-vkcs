data "vkcs_compute_flavor" "db" {
  name = "Standard-2-8-50"
}

resource "vkcs_db_instance" "db_instance" {
  name = "db-instance-tf-example"

  availability_zone = "GZ1"

  datastore {
    type    = "postgresql"
    version = "16"
  }

  flavor_id           = data.vkcs_compute_flavor.db.id
  floating_ip_enabled = true

  size        = 10
  volume_type = "ceph-ssd"
  disk_autoexpand {
    autoexpand    = true
    max_disk_size = 1000
  }

  network {
    uuid = vkcs_networking_network.db.id
  }

  depends_on = [
    vkcs_networking_router_interface.db
  ]
}

resource "vkcs_db_database" "postgres_db" {
  name    = "testdb_1"
  dbms_id = vkcs_db_instance.db_instance.id
}

resource "vkcs_db_user" "postgres_user" {
  name = "testuser"
  # Example only. Do not use in production.
  # Sensitive values must be provided securely and not stored in manifests.
  password = "Test_p#ssword-12-3"

  dbms_id = vkcs_db_instance.db_instance.id

  vendor_options {
    skip_deletion = true
  }

  databases = [vkcs_db_database.postgres_db.name]
}

