resource "vkcs_compute_instance" "basic" {
  name              = "basic-tf-example"
  # AZ and flavor are mandatory
  availability_zone = "GZ1"
  flavor_name       = "Basic-1-2-20"
  # Use block_device to specify instance disk to get full control
  # of it in the future
  block_device {
    source_type           = "image"
    uuid                  = data.vkcs_images_image.debian.id
    destination_type      = "volume"
    volume_size           = 10
    delete_on_termination = true
  }
  # Specify at least one network to do not depend on project assets
  network {
    uuid = vkcs_networking_network.app.id
  }
  # Specify required security groups if yoo do not want `default` one
  security_groups = [
    vkcs_networking_secgroup.admin.id
  ]
  # If your configuration also defines a network for the instance,
  # ensure it is attachec to a router before creating of the instance
  depends_on = [
    vkcs_networking_router_interface.app
  ]
}
