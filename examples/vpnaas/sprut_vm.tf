# Create vm in sprut net with default, ssh and ping security groups.
# This is only necessary to verify that the vpn tunnel is configured correctly.
resource "vkcs_networking_secgroup" "ping_sprut" {
  name = "ping-sprut-tf-example"
  sdn  = local.sdn_sprut
}

resource "vkcs_networking_secgroup_rule" "ping_rule_sprut" {
  description       = "Ping rule in sprut"
  security_group_id = vkcs_networking_secgroup.ping_sprut.id
  direction         = "ingress"
  protocol          = "icmp"
  remote_ip_prefix  = "0.0.0.0/0"
  sdn               = local.sdn_sprut
}

resource "vkcs_networking_port" "sprut_vm_port" {
  name                         = "port_for_sprut_vm"
  network_id                   = vkcs_networking_network.sprut.id
  full_security_groups_control = true
  security_group_ids = [
    vkcs_networking_secgroup.ping_sprut.id
  ]

  depends_on = [
    vkcs_networking_subnet.sprut
  ]
}

# Create vm in sprut net
resource "vkcs_compute_instance" "vm_sprut" {
  name              = "sprut-vm-tf-example"
  availability_zone = "ME1"
  flavor_name       = "STD3-1-2"
  block_device {
    source_type           = "image"
    uuid                  = data.vkcs_images_image.image.id
    destination_type      = "volume"
    volume_size           = 11
    delete_on_termination = true
  }

  network {
    port = vkcs_networking_port.sprut_vm_port.id
  }

  # Wait for the route to be added into subnet
  depends_on = [
    vkcs_networking_subnet_route.static_sprut
  ]
}
