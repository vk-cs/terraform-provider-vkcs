resource "vkcs_kubernetes_cluster_v2" "cluster" {
  name        = "base-infra"
  description = "An example of a Kubernetes cluster v2 created via Terraform"
  version     = "v1.31.4"
  
  cluster_type = "standard"
  master_count = 1
  
  master_flavor = data.vkcs_compute_flavor.master.id
  
  network_id  = vkcs_networking_network.k8s-network.id
  subnet_id   = vkcs_networking_subnet.k8s-subnet.id
  
  availability_zones = ["GZ1"]
  
  external_network_id = data.vkcs_networking_network.extnet.id
  
  network_plugin = "calico"
  pods_ipv4_cidr = "10.100.0.0/16"
  
  enable_public_ip = true
  
  labels = {
    managed-by  = "terraform"
  }
  
  insecure_registries = ["registry.example.com:5000"]
  
  loadbalancer_subnet_id = vkcs_networking_subnet.k8s-subnet.id
  loadbalancer_allowed_cidrs = ["10.0.0.0/8", "192.168.0.0/16"]
}
