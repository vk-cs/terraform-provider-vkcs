resource "vkcs_compute_instance" "volumes_tags_externalip" {
  name              = "volumes-tags-externalip-tf-example"
  availability_zone = "GZ1"
  flavor_name       = "Basic-1-2-20"
  # Use previously created volume as root device
  block_device {
    # Set boot_index to mark root device if multiple
    # block devices are specified
    boot_index       = 0
    source_type      = "volume"
    uuid             = vkcs_blockstorage_volume.bootable.id
    destination_type = "volume"
    # Omitting delete_on_termination (or setting it to false)
    # allows you to manage previously created volume after instance deletion
  }
  # Gracefully shutdown instance before deleting to keep
  # data consistent on persistent volume
  stop_before_destroy = true
  # Add empty disk to use ii during the instance lifecycle
  block_device {
    source_type           = "blank"
    destination_type      = "volume"
    volume_size           = 20
    delete_on_termination = true
  }
  tags = ["tf-example"]
  # Create the instance in external network to get external IP automatically
  # instead of using FIP on private network
  network {
    name = "ext-net"
  }
}
