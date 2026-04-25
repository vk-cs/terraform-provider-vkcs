data "vkcs_baremetal_server" "server" {
  id = vkcs_baremetal_server.server.id
}

output "server_output" {
  value = {
    id                = data.vkcs_baremetal_server.server.id
    name              = data.vkcs_baremetal_server.server.name
    region            = data.vkcs_baremetal_server.server.region
    availability_zone = data.vkcs_baremetal_server.server.availability_zone

    cpu_cores = data.vkcs_baremetal_server.server.cpu_cores
    cpu_types = data.vkcs_baremetal_server.server.cpu_types

    is_locked = data.vkcs_baremetal_server.server.is_locked

    local_disk_sizes = data.vkcs_baremetal_server.server.local_disk_sizes
    ram_megabytes    = data.vkcs_baremetal_server.server.ram_megabytes

    power_state = data.vkcs_baremetal_server.server.power_state
    status      = data.vkcs_baremetal_server.server.status

    tags = data.vkcs_baremetal_server.server.tags

    image_id   = data.vkcs_baremetal_server.server.image_id
    image_name = data.vkcs_baremetal_server.server.image_name
    os_type    = data.vkcs_baremetal_server.server.os_type

    raid_type = data.vkcs_baremetal_server.server.raid_type
    flavor_id = data.vkcs_baremetal_server.server.flavor_id

    target_boot_order = data.vkcs_baremetal_server.server.target_boot_order

    local_disks_info = data.vkcs_baremetal_server.server.local_disks_info
  }
}

output "boot_devices" {
  value = [
    for b in data.vkcs_baremetal_server.server.target_boot_order :
    b.boot_device_type
  ]
}

output "disks" {
  value = [
    for d in data.vkcs_baremetal_server.server.local_disks_info : {
      path  = d.path
      size  = d.size
      type  = d.type
      model = d.model
    }
  ]
}

