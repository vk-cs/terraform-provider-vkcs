resource "vkcs_compute_instance" "front_worker" {
  name      = "front-worker-tf-example"
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
    data.vkcs_networking_secgroup.default_secgroup.id,
    vkcs_networking_secgroup.admin.id,
    vkcs_networking_secgroup.http.id
  ]
  image_id = data.vkcs_images_image.debian.id

  network {
    uuid = vkcs_networking_network.app.id
  }
}
