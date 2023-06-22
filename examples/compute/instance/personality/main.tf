resource "vkcs_compute_instance" "basic" {
  name              = "personality-tf-example"
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
  # config_drive must be enabled to use personality
  config_drive = true
  personality {
    file    = "/opt/app/config.json"
    content = jsonencode({ "foo" : "bar" })
  }
  depends_on = [
    vkcs_networking_router_interface.app
  ]
}
