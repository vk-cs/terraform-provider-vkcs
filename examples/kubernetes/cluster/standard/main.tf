resource "vkcs_kubernetes_cluster" "k8s_cluster" {
  name                = "k8s-standard-cluster"
  cluster_type        = "standard"
  cluster_template_id = data.vkcs_kubernetes_clustertemplate.k8s_31.id
  master_flavor       = data.vkcs_compute_flavor.basic.id
  master_count        = 1

  labels = {
    cloud_monitoring         = "true"
    kube_log_level           = "2"
    clean_volumes            = "true"
    master_volume_size       = "100"
    cluster_node_volume_type = "ceph-ssd"
  }

  availability_zone   = "MS1"
  network_id          = vkcs_networking_network.app.id
  subnet_id           = vkcs_networking_subnet.app.id
  floating_ip_enabled = true

  sync_security_policy = true
  # If your configuration also defines a network for the instance,
  # ensure it is attached to a router before creating of the instance
  depends_on = [
    vkcs_networking_router_interface.app,
  ]
}
