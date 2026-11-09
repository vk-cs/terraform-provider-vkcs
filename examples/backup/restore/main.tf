data "vkcs_backup_checkpoints" "source" {
  count   = var.enable_restore ? 1 : 0
  plan_id = var.backup_plan_id
}

locals {
  available_vm_backups = flatten([
    for resource in try(data.vkcs_backup_checkpoints.source[0].resources, []) : [
      for backup in resource.backups : {
        checkpoint_id      = backup.checkpoint_id
        source_resource_id = resource.id
      } if resource.resource_type == "OS::Nova::Server" && backup.status == "available"
    ]
  ])
}

resource "vkcs_backup_restore" "vm" {
  count              = var.enable_restore ? 1 : 0
  checkpoint_id      = try(local.available_vm_backups[0].checkpoint_id, "")
  source_resource_id = try(local.available_vm_backups[0].source_resource_id, "")
  restore_name       = "restored-vm-example"
  restore_net_id     = var.restore_net_id
  restore_keypair    = var.restore_keypair
}

data "vkcs_backup_restore" "vm" {
  count      = var.enable_restore ? 1 : 0
  restore_id = try(vkcs_backup_restore.vm[0].id, "")
}

resource "vkcs_compute_instance" "restored_vm" {
  count      = var.enable_restore ? 1 : 0
  name       = try(vkcs_backup_restore.vm[0].restore_name, "restored-vm-example")
  restore_id = try(vkcs_backup_restore.vm[0].id, "")
}

output "restore_resources_status" {
  value = try(data.vkcs_backup_restore.vm[0].resources_status, {})
}

output "restored_vm_id" {
  value = try(vkcs_compute_instance.restored_vm[0].id, "")
}
