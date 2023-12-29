resource "vkcs_db_instance" "postgrespro_enterprise_1c" {
  name = "basic-tf-example"
  # AZ, flavor and datastore are mandatory
  availability_zone = "GZ1"
  flavor_id         = data.vkcs_compute_flavor.basic.id
  datastore {
    type    = "postgrespro_enterprise_1c"
    version = "12"
  }
  # Specify at least one network to not depend on project assets
  # Specify required security groups if you do not want `default` one
  network {
    uuid            = vkcs_networking_network.db.id
    security_groups = [vkcs_networking_secgroup.admin.id]
  }
  # Specify volume type, size and autoexpand options
  size        = 8
  volume_type = "ceph-ssd"
  disk_autoexpand {
    autoexpand    = true
    max_disk_size = 1000
  }
  # Specify required db capabilities
  capabilities {
    name = "node_exporter"
    settings = {
      "listen_port" : "9100"
    }
  }
  # Enable cloud monitoring
  cloud_monitoring_enabled = true
  # If your configuration also defines a network for the db instance,
  # ensure it is attached to a router before creating of the instance
  depends_on = [
    vkcs_networking_router_interface.db,
  ]
}

resource "vkcs_db_instance" "postgrespro_enterprise_1c_replica" {
  name = "basic-tf-example-replica"
  # AZ, flavor and datastore are mandatory
  availability_zone = "GZ1"
  flavor_id         = data.vkcs_compute_flavor.basic.id
  datastore {
    type    = "postgrespro_enterprise_1c"
    version = "12"
  }

  replica_of  = vkcs_db_instance.postgrespro_enterprise_1c.id

  network {
    uuid            = vkcs_networking_network.db.id
    security_groups = [vkcs_networking_secgroup.admin.id]
  }

  size        = 8
  volume_type = "ceph-ssd"
  

  capabilities {
    name = "node_exporter"
    settings = {
      "listen_port" : "9100"
    }
  }

  cloud_monitoring_enabled = true

  depends_on = [
    vkcs_networking_router_interface.db,
  ]
}
