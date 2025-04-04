resource "vkcs_compute_keypair" "generated_key" {
  name = "generated-key-tf-example"
}

resource "vkcs_compute_instance" "basic" {
  name              = "windows-password-tf-example"
  availability_zone = "ME1"
  flavor_name       = "STD2-2-8"
  key_pair          = vkcs_compute_keypair.generated_key.name

  block_device {
    source_type           = "image"
    uuid                  = data.vkcs_images_image.windows.id
    destination_type      = "volume"
    volume_size           = 50
    volume_type           = "high-iops-ha"
    delete_on_termination = true
  }

  network {
    uuid = vkcs_networking_network.app.id
  }

  security_group_ids = [
    vkcs_networking_secgroup.rdp.id,
  ]

  vendor_options {
    get_password_data = true
  }

  depends_on = [
    vkcs_networking_router_interface.app
  ]
}

output "windows_password" {
  value     = rsadecrypt(vkcs_compute_instance.basic.password_data, vkcs_compute_keypair.generated_key.private_key)
  sensitive = true
}
