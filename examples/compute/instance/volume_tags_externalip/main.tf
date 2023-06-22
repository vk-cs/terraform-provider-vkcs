resource "vkcs_blockstorage_volume" "volume" {
  name              = "volume-tf-example"
  size              = 5
  image_id          = data.vkcs_images_image.debian.id
  availability_zone = "GZ1"
  volume_type       = "ceph-ssd"
}

resource "vkcs_compute_instance" "boot-from-volume" {
  name = "bootfromvolume"

  availability_zone = "GZ1"
  flavor_name       = "Basic-1-2-20"

  block_device {
    source_type           = "volume"
    uuid                  = vkcs_blockstorage_volume.volume.id
    destination_type      = "volume"
    boot_index            = 0
    delete_on_termination = true
  }

  network {
    uuid = vkcs_networking_network.app.id
  }

  security_groups = [
    vkcs_networking_secgroup.admin.name
  ]

  tags = ["tf-example"]

  depends_on = [
    vkcs_networking_router_interface.app
  ]
}

resource "vkcs_networking_floatingip" "fip" {
  pool = data.vkcs_networking_network.extnet.name
}

resource "vkcs_compute_floatingip_associate" "fip" {
  floating_ip = vkcs_networking_floatingip.fip.address
  instance_id = vkcs_compute_instance.boot-from-volume.id
}
