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
  name           = "k8s-net"
  admin_state_up = true
}

resource "vkcs_db_instance" "db-instance" {
  name        = "db-instance"

  datastore {
    type    = "postgresql"
    version = "10"
  }

  flavor_id   = data.vkcs_compute_flavor.db.id
  
  size        = 8
  volume_type = "ceph-ssd"
  network {
    uuid = vkcs_networking_network.db.id
  }

  root_enabled  = true
  root_password = var.db-root-user-pwd
}

output "root_user_password" {
  value = vkcs_db_instance.db-instance.root_password
}

