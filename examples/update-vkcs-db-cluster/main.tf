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
}

resource "vkcs_db_cluster" "db-cluster" {
  name        = "db-cluster"

  datastore {
    type    = "postgresql"
    version = "12"
  }

  cluster_size = 5

  flavor_id   = data.vkcs_compute_flavor.db.id

  volume_size = 12
  volume_type = "ceph-ssd"

  network {
    uuid = vkcs_networking_network.db.id
  }
}