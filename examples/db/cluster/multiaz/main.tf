resource "vkcs_db_cluster" "cluster" {
  name               = "multiaz-cluster-tf-example"
  availability_zones = ["GZ1", "MS1"]
  cluster_size       = 3
  flavor_id          = data.vkcs_compute_flavor.basic.id
  volume_size        = 10
  volume_type        = "ceph-ssd"
  datastore {
    version = "16"
    type    = "postgresql_multiaz"
  }
  network {
    uuid = vkcs_networking_network.db.id
  }

  depends_on = [
    vkcs_networking_router_interface.db
  ]
}

data "vkcs_networking_port" "vrrp_port" {
  id = vkcs_db_cluster.cluster.vrrp_port_id
}

output "cluster_ip" {
  value       = data.vkcs_networking_port.vrrp_port.all_fixed_ips[0]
  description = "IP address of the cluster."
}
