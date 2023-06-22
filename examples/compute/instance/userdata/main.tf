resource "vkcs_compute_instance" "basic" {
  name              = "basic-tf-example"
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
  user_data = <<EOF
    #cloud-config
    package_upgrade: true
    packages:
      - nginx
    runcmd:
      - systemctl start nginx
  EOF
  depends_on = [
    vkcs_networking_router_interface.app
  ]
}
