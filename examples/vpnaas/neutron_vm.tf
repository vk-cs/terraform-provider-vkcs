# Create vm in neutron net with default, ssh and ping security groups.
# This is only necessary to verify that the vpn tunnel is configured correctly.
resource "vkcs_networking_secgroup" "ping_neutron" {
  name = "ping-neutron-tf-example"
  sdn  = local.sdn_neutron
}

resource "vkcs_networking_secgroup_rule" "ping_rule_neutron" {
  description       = "Ping rule in neutron"
  security_group_id = vkcs_networking_secgroup.ping_neutron.id
  direction         = "ingress"
  protocol          = "icmp"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "vkcs_networking_secgroup" "ssh_neutron" {
  name = "ssh-neutron-tf-example"
  sdn  = local.sdn_neutron
}

resource "vkcs_networking_secgroup_rule" "ssh_rule_neutron" {
  description       = "SSH rule in neutron"
  security_group_id = vkcs_networking_secgroup.ssh_neutron.id
  direction         = "ingress"
  protocol          = "tcp"
  # Specify SSH port
  port_range_max = 22
  port_range_min = 22
  # Allow access from any sources
  remote_ip_prefix = "0.0.0.0/0"
}

resource "vkcs_networking_port" "neutron_vm_port" {
  name       = "port_for_neutron_vm"
  network_id = vkcs_networking_network.neutron.id
  full_security_groups_control = true
  security_group_ids = [
    vkcs_networking_secgroup.ssh_neutron.id,
    vkcs_networking_secgroup.ping_neutron.id
  ]

  depends_on = [
    vkcs_networking_subnet.neutron
  ]
}

resource "vkcs_compute_keypair" "key_pair" {
  name = "vpn_key_pair_tf_example"
}

# Create vm in neutron net
resource "vkcs_compute_instance" "vm_neutron" {
  name              = "neutron-vm-tf-example"
  availability_zone = "ME1"
  flavor_name       = "STD3-1-2"
  key_pair          = vkcs_compute_keypair.key_pair.name # need for connect to vm by ssh
  block_device {
    source_type           = "image"
    uuid                  = data.vkcs_images_image.image.id
    destination_type      = "volume"
    volume_size           = 11
    delete_on_termination = true
  }

  network {
    port = vkcs_networking_port.neutron_vm_port.id
  }
  network {
    uuid = data.vkcs_networking_network.extnet_neutron.id
  }
  security_group_ids = [
    vkcs_networking_secgroup.ssh_neutron.id,
    vkcs_networking_secgroup.ping_neutron.id
  ]

  # Wait for the route to be added into subnet
  depends_on = [
    vkcs_networking_subnet_route.static_neutron
  ]
}

resource "terraform_data" "ssh_with_ping" {
  triggers_replace = [
    vkcs_compute_instance.vm_neutron.id,
    vkcs_compute_instance.vm_sprut.id,
    vkcs_compute_instance.vm_neutron.network[1].fixed_ip_v4,
    vkcs_networking_port.sprut_vm_port.all_fixed_ips[0],
    vkcs_vpnaas_service.neutron.id,
    vkcs_vpnaas_service.sprut.id,
  ]

  provisioner "remote-exec" {
    inline = [
      "ping -c 1 -w 900 ${vkcs_networking_port.sprut_vm_port.all_fixed_ips[0]}"
    ]

    connection {
      type        = "ssh"
      user        = "centos"
      private_key = vkcs_compute_keypair.key_pair.private_key
      host        = vkcs_compute_instance.vm_neutron.network[1].fixed_ip_v4
      timeout     = "7m"
    }
  }

  depends_on = [
    vkcs_compute_instance.vm_sprut,
    vkcs_compute_instance.vm_neutron
  ]
}
