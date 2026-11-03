data "vkcs_backup_checkpoints" "source" {
  plan_id = var.backup_plan_id
}

locals {
  available_vm_backups = flatten([
    for resource in data.vkcs_backup_checkpoints.source.resources : [
      for backup in resource.backups : {
        checkpoint_id      = backup.checkpoint_id
        source_resource_id = resource.id
      } if resource.resource_type == "OS::Nova::Server" && backup.status == "available"
    ]
  ])
}

resource "vkcs_backup_restore" "vm" {
  checkpoint_id      = local.available_vm_backups[0].checkpoint_id
  source_resource_id = local.available_vm_backups[0].source_resource_id
  restore_name       = "restored-vm-example"
  restore_net_id     = var.restore_net_id
  restore_keypair    = var.restore_keypair
}

data "vkcs_backup_restore" "vm" {
  restore_id = vkcs_backup_restore.vm.id
}

output "restore_resources_status" {
  value = data.vkcs_backup_restore.vm.resources_status
}
