resource "vkcs_compute_instance" "gateway" {
  name              = "gateway-tf-example"
  availability_zone = "GZ1"
  flavor_name       = "Basic-1-2-20"
  block_device {
    source_type           = "image"
    uuid                  = data.vkcs_images_image.debian.id
    destination_type      = "volume"
    volume_size           = 10
    delete_on_termination = true
  }
  network {
    uuid = vkcs_networking_network.app.id
  }
  depends_on = [
    vkcs_networking_router_interface.app
  ]
}

