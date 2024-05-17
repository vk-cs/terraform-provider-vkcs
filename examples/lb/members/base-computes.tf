resource "vkcs_compute_instance" "front_worker" {
  count     = 2
  name      = "front-worker-${count.index}-tf-example"
  flavor_id = data.vkcs_compute_flavor.basic.id

  block_device {
    source_type      = "image"
    uuid             = data.vkcs_images_image.debian.id
    destination_type = "volume"
    volume_size      = 10
    # Must be set to delete volume after instance deletion
    # Otherwise you get "orphaned" volume with terraform
    delete_on_termination = true
  }

  security_group_ids = [
    vkcs_networking_secgroup.admin.id,
    vkcs_networking_secgroup.http.id
  ]

  network {
    uuid = vkcs_networking_network.app.id
  }

  depends_on = [
    vkcs_networking_router_interface.app
  ]
}
