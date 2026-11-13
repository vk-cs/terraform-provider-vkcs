data "vkcs_backup_checkpoints" "plan" {
  count   = var.enable_checkpoints ? 1 : 0
  plan_id = var.backup_plan_id
}

output "legal_hold_checkpoint_ids" {
  value = flatten([
    for resource in try(data.vkcs_backup_checkpoints.plan[0].resources, []) : [
      for backup in resource.backups : backup.checkpoint_id
      if try(backup.object_lock.type, "") == "legalhold"
    ]
  ])
}
