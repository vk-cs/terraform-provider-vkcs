resource "vkcs_compute_instance" "basic" {
  name = "basic-tf-example"

  # AZ and flavor are mandatory
  availability_zone = "GZ1"

  flavor_name = "Basic-1-2-20"

  # Use block_device to specify instance disk to get full control
  # of it in the future
  block_device {
    # Set boot_index to mark root device if multiple
    # block devices are specified
    boot_index       = 0
    source_type      = "volume"
    uuid             = vkcs_blockstorage_volume.bootable.id
    destination_type = "volume"
    # Omitting delete_on_termination (or setting it to false)
    # allows you to manage previously created volume after instance deletion
    delete_on_termination = true
  }

  block_device {
    boot_index            = -1
    source_type           = "volume"
    uuid                  = vkcs_blockstorage_volume.data.id
    destination_type      = "volume"
    delete_on_termination = true
  }

  # Specify at least one network to not depend on project assets
  network {
    uuid = vkcs_networking_network.app.id
  }

  # Specify required security groups if you do not want `default` one
  security_group_ids = [
    vkcs_networking_secgroup.admin.id
  ]

  # If your configuration also defines a network for the instance,
  # ensure it is attached to a router before creating of the instance
  depends_on = [
    vkcs_networking_router_interface.app
  ]
}
