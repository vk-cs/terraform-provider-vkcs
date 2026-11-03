data "vkcs_backup_checkpoints" "plan" {
  plan_id = var.backup_plan_id
}

output "legal_hold_checkpoint_ids" {
  value = flatten([
    for resource in data.vkcs_backup_checkpoints.plan.resources : [
      for backup in resource.backups : backup.checkpoint_id
      if try(backup.object_lock.type, "") == "legalhold"
    ]
  ])
}
