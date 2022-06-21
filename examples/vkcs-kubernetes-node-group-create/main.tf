resource "vkcs_networking_network" "k8s" {
    name           = "k8s-net"
    admin_state_up = true
}

resource "vkcs_networking_subnet" "k8s-subnetwork" {
  name            = "k8s-subnet"
  network_id      = vkcs_networking_network.k8s.id
  cidr            = "192.168.199.0/24"
  ip_version      = 4
  dns_nameservers = ["8.8.8.8", "8.8.4.4"]
}

data "vkcs_networking_network" "extnet" {
  name = "ext-net"
}

resource "vkcs_networking_router" "k8s" {
  name                = "k8s-router"
  admin_state_up      = true
  external_network_id = data.vkcs_networking_network.extnet.id
}

resource "vkcs_networking_router_interface" "k8s" {
  router_id = vkcs_networking_router.k8s.id
  subnet_id = vkcs_networking_subnet.k8s-subnetwork.id
}

data "vkcs_compute_flavor" "k8s" {
  name = "Standard-2-4-50"
}

data "vkcs_kubernetes_clustertemplate" "ct" {
  version = "1.21.4"
}

resource "vkcs_kubernetes_cluster" "k8s-cluster" {
  depends_on = [
    vkcs_networking_router_interface.k8s,
  ]

  name                = "k8s-cluster"
  cluster_template_id = data.vkcs_kubernetes_clustertemplate.ct.id
  master_flavor       = data.vkcs_compute_flavor.k8s.id
  master_count        = 1

  network_id          = vkcs_networking_network.k8s.id
  subnet_id           = vkcs_networking_subnet.k8s-subnetwork.id
  floating_ip_enabled = true
  availability_zone   = "MS1"
  insecure_registries = ["1.2.3.4"]
}

resource "vkcs_kubernetes_node_group" "default_ng" {
    cluster_id = data.vkcs_kubernetes_cluster.k8s-cluster

    node_count = 1
    name = "default"
    max_nodes = 5
    min_nodes = 1

    labels {
        key = "env"
        value = "test"
    }

    labels {
        key = "disktype"
        value = "ssd"
    }
    
    taints {
        key = "taintkey1"
        value = "taintvalue1"
        effect = "PreferNoSchedule"
    }

    taints {
        key = "taintkey2"
        value = "taintvalue2"
        effect = "PreferNoSchedule"
    }
}
