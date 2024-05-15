resource "vkcs_compute_instance" "front_worker" {
  name      = "front-worker-${count.index}-tf-example"
  flavor_id = data.vkcs_compute_flavor.basic.id
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

data "vkcs_networking_secgroup" "default_secgroup" {
  name = "default"
}
