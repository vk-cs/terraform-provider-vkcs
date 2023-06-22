resource "vkcs_networking_port" "port_1" {
  name           = "port-tf-example"
  network_id     = vkcs_networking_network.app.id
  admin_state_up = "true"
}

resource "vkcs_networking_secgroup" "secgroup_1" {
  name        = "secgroup-tf-example"
  description = "Security group example"
}

resource "vkcs_networking_port_secgroup_associate" "port_1" {
  port_id = vkcs_networking_port.port_1.id
  security_group_ids = [
    vkcs_networking_secgroup.secgroup_1.id
  ]
}

resource "vkcs_compute_instance" "basic" {
  name = "basic-tf-example"

  availability_zone = "GZ1"
  flavor_name       = "Basic-1-2-20"

  block_device {
    source_type           = "image"
    uuid                  = data.vkcs_images_image.debian.id
    destination_type      = "volume"
    volume_size           = 10
    delete_on_termination = true
  }

  # Only port's security groups (i.e. secgroup-tf-example) will be applied
  network {
    port = vkcs_networking_port.port_1.id
  }
  # If port is used, this security group will not be applied
  security_groups = [
    vkcs_networking_secgroup.admin.name,
  ]

  depends_on = [
    vkcs_networking_router_interface.app
  ]
}

