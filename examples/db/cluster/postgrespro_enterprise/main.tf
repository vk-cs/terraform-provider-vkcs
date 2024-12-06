resource "vkcs_db_cluster" "pg_cluster" {
  name = "pg-cluster"

  availability_zone = "GZ1"
  datastore {
    type    = "postgrespro_enterprise"
    version = "12"
  }

  cluster_size = 3

  flavor_id                = data.vkcs_compute_flavor.basic.id
  cloud_monitoring_enabled = true

  volume_size = 10
  volume_type = "ceph-ssd"

  network {
    uuid = vkcs_networking_network.db.id
  }

  depends_on = [
    vkcs_networking_router_interface.db
  ]
}

data "vkcs_lb_loadbalancer" "loadbalancer" {
  id = vkcs_db_cluster.pg_cluster.loadbalancer_id
}

data "vkcs_networking_port" "loadbalancer_port" {
  id = data.vkcs_lb_loadbalancer.loadbalancer.vip_port_id
}

output "cluster_ips" {
  value       = data.vkcs_networking_port.loadbalancer_port.all_fixed_ips
  description = "IP addresses of the cluster."
}
