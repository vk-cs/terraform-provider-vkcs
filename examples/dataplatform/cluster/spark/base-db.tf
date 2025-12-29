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

  size        = 100
  volume_type = "ceph-ssd"
  disk_autoexpand {
    autoexpand    = true
    max_disk_size = 1000
  }

  network {
    uuid = module.network.networks[0].id
  }

}

resource "vkcs_db_database" "postgres_db" {
  name    = "spark"
  dbms_id = vkcs_db_instance.db_instance.id
}

resource "vkcs_db_user" "postgres_user" {
  name     = "spark"
  password = random_password.spark.result

  dbms_id = vkcs_db_instance.db_instance.id

  vendor_options {
    skip_deletion = true
  }

  databases = [vkcs_db_database.postgres_db.name]
}

