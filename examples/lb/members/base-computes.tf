resource "vkcs_compute_instance" "front_worker" {
  count           = 2
  name            = "front-worker-${count.index}-tf-example"
  flavor_id       = data.vkcs_compute_flavor.basic.id
  security_groups = ["default", "admin-tf-example", "http-tf-example"]
  image_id        = data.vkcs_images_image.debian.id

  network {
    uuid        = vkcs_networking_network.app.id
  }
  depends_on = [
    vkcs_networking_secgroup.admin,
    vkcs_networking_secgroup.http,
    vkcs_networking_router_interface.app
   ]
}
