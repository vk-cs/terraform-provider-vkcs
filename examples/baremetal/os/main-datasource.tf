data "vkcs_baremetal_os" "ubuntu" {
  name      = "ubuntu"
  version   = "24.04"
  raid_type = "RAID1"
}

output "flavor_output" {
  value = {
    id        = data.vkcs_baremetal_os.ubuntu.id
    name      = data.vkcs_baremetal_os.ubuntu.name
    version   = data.vkcs_baremetal_os.ubuntu.version
    raid_type = data.vkcs_baremetal_os.ubuntu.raid_type
  }
}
