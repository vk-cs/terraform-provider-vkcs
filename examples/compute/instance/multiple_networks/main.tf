resource "vkcs_compute_instance" "multiple_networks" {
  name              = "multiple-networks-tf-example"
  availability_zone = "GZ1"
  flavor_name       = "Basic-1-2-20"
  block_device {
    source_type           = "image"
    uuid                  = data.vkcs_images_image.debian.id
    destination_type      = "volume"
    volume_size           = 10
    delete_on_termination = true
  }
  # Autocreate a new port in 'app' network
  network {
    uuid = vkcs_networking_network.app.id
  }
  # Use previously created port
  # This does not change security groups associated with the port
  # Also this changes DNS name of the port
  network {
    port = vkcs_networking_port.persistent_etcd.id
  }
  # Attach 'admin' security group to autocreated port
  # This does not associate the group to 'persistent' port
  security_group_ids = [
    vkcs_networking_secgroup.admin.id
  ]
  depends_on = [
    vkcs_networking_router_interface.app,
    vkcs_networking_router_interface.db
  ]
}
