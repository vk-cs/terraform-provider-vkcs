resource "vkcs_compute_instance" "restored_vm" {
  count              = var.enable_restore ? 1 : 0
  name               = var.instance_name
  checkpoint_id      = var.checkpoint_id
  source_resource_id = var.source_resource_id != "" ? var.source_resource_id : null
  flavor_id          = var.flavor_id
  key_pair           = var.key_pair

  network {
    uuid = var.network_id
  }
}

data "vkcs_backup_restore" "status" {
  count      = var.enable_restore ? 1 : 0
  restore_id = try(vkcs_compute_instance.restored_vm[0].restore_id, "")
}

output "restored_vm_id" {
  value = try(vkcs_compute_instance.restored_vm[0].id, "")
}

output "restore_resources_status" {
  value = try(data.vkcs_backup_restore.status[0].resources_status, {})
}
