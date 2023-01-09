resource "vkcs_db_config_group" "db-config-group" {
    name = "db-config-group"
    datastore {
        type = "mysql"
        version = "8.0"
    }
    values = {
        activate_all_roles_on_login : "true"
        autocommit : "1"
        block_encryption_mode : "test"
        innodb_segment_reserve_factor : "0.53"
    }
    description = "db-config-group-description"
}


resource "vkcs_db_instance" "db-instance" {
    name = "db-instance"

    availability_zone = "GZ1"
    
    datastore {
        type = "mysql"
        version = "8.0"
    }
    
    configuration_id = vkcs_db_config_group.db-config-group.id
    network {
      uuid = vkcs_networking_network.db.id
    }
    flavor_id = data.vkcs_compute_flavor.db.id
    volume_type = "ceph-ssd"
    size = 8

    depends_on = [
        vkcs_networking_router_interface.db
    ]
}
