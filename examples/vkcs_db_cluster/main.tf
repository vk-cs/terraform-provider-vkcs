resource "vkcs_db_cluster" "db-cluster" {
  name        = "db-cluster"

  availability_zone = "GZ1"
  datastore {
    type    = "postgresql"
    version = "12"
  }

  cluster_size = 3

  flavor_id   = data.vkcs_compute_flavor.db.id

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
  id = "${vkcs_db_cluster.db-cluster.loadbalancer_id}"
}

data "vkcs_networking_port" "loadbalancer-port" {
  port_id = "${data.vkcs_lb_loadbalancer.loadbalancer.vip_port_id}"
}

# Use this to connect to the cluster
output "cluster_ips" {
  value = "${data.vkcs_networking_port.loadbalancer-port.all_fixed_ips}"
  description = "IP addresses of the cluster."
}
