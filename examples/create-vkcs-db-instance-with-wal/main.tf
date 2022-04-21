terraform {
  required_providers {
    vkcs = {
      source  = "vk-cs/vkcs"
      version = "~> 0.1.0"
    }
  }
}

data "vkcs_compute_flavor" "db" {
  name = var.db-instance-flavor
}

resource "vkcs_networking_network" "db" {
  name           = "db-net"
  admin_state_up = true
}

resource "vkcs_db_instance" "db-instance" {
  name        = "db-instance"

  datastore {
    type    = "postgresql"
    version = "11"
  }

  public_access     = true

  flavor_id   = data.vkcs_compute_flavor.db.id

  size        = 10
  volume_type = "ceph-ssd"
  disk_autoexpand {
    autoexpand    = true
    max_disk_size = 1000
  }
  wal_volume {
    size          = 10
    volume_type   = "ceph-ssd"
    autoexpand    = true
    max_disk_size = 20
  }

  network {
    uuid = vkcs_networking_network.db.id
  }
}
