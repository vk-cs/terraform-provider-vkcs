resource "vkcs_compute_servergroup" "example" {
  name     = "servergroup-tf-example"
  policies = ["soft-anti-affinity"]
}

resource "vkcs_compute_instance" "servergroup" {
  name              = "servergroup-tf-example"
  availability_zone = "MS1"
  flavor_name       = "Basic-1-2-20"

  block_device {
    source_type           = "image"
    uuid                  = data.vkcs_images_image.debian.id
    destination_type      = "volume"
    volume_size           = 10
    volume_type           = "ceph-ssd"
    delete_on_termination = true
  }

  network {
    uuid = vkcs_networking_network.app.id
  }

  security_group_ids = [
    vkcs_networking_secgroup.admin.id
  ]

  # The server is placed into the group at creation time and can be moved
  # between groups later without recreation.
  server_group_id = vkcs_compute_servergroup.example.id

  depends_on = [
    vkcs_networking_router_interface.app
  ]
}
