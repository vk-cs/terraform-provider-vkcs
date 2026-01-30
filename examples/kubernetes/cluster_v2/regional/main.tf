resource "vkcs_kubernetes_cluster_v2" "k8s_cluster" {
  name        = "k8s-regional-cluster"
  description = "An example of a regional Kubernetes cluster v2 created via Terraform"
  version     = "v1.34.2"

  cluster_type       = "regional"
  master_count       = 3
  availability_zones = ["GZ1", "MS1", "ME1"]
  master_flavor      = data.vkcs_compute_flavor.master.id

  network_id             = vkcs_networking_network.app.id
  subnet_id              = vkcs_networking_subnet.app.id
  loadbalancer_subnet_id = vkcs_networking_subnet.app.id
  network_plugin         = "calico"
  pods_ipv4_cidr         = "10.100.0.0/16"

  # If your configuration also defines a network for the instance,
  # ensure it is attached to a router before creating of the instance
  depends_on = [
    vkcs_networking_router_interface.app,
  ]
}
