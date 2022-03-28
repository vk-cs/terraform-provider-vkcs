terraform {
  required_providers {
    mcs = {
      source = "vk-cs/vkcs"
      version = "~> 0.1.0"
    }
  }
}


resource "vkcs_networking_network" "k8s" {
  name           = "k8s-net"
  admin_state_up = true
}


resource "vkcs_networking_subnet" "k8s-subnetwork" {
  name            = "k8s-subnet"
  network_id      = var.k8s-network-id
  cidr            = "10.100.0.0/16"
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
  router_id = var.k8s-router-id
  subnet_id = var.k8s-subnet-id
}


resource "vkcs_compute_keypair" "keypair" {
  name       = "default"
  public_key = file(var.public-key-file)
}


data "vkcs_compute_flavor" "k8s" {
  name = var.k8s-flavor
}

data "vkcs_kubernetes_clustertemplate" "ct" {
  version = "1.20.4"
}

resource "vkcs_kubernetes_cluster" "k8s-cluster" {

  master_flavor = var.new-master-flavor

  name = "k8s-cluster"
  cluster_template_id = data.vkcs_kubernetes_clustertemplate.ct.id
  master_count        = 1
  keypair = vkcs_compute_keypair.keypair.id
  network_id = vkcs_networking_network.k8s.id
  subnet_id = vkcs_networking_subnet.k8s-subnetwork.id
  floating_ip_enabled = true
  availability_zone = "MS1"
}