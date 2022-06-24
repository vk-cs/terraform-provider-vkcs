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

resource "vkcs_compute_keypair" "keypair" {
  name       = "default"
  public_key = file(var.public-key-file)
}

resource "vkcs_networking_network" "db" {
  name           = "db-net"
}

resource "vkcs_db_instance" "db-instance" {
  name        = "db-instance"

  datastore {
    type    = "mysql"
    version = "5.7"
  }

  keypair           = vkcs_compute_keypair.keypair.id
  public_access     = true

  flavor_id   = data.vkcs_compute_flavor.db.id

  size        = 10
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
}