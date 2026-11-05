resource "vkcs_backup_plan" "checkpoint_source" {
  name               = "backup-plan-checkpoints-example"
  provider_name      = "cloud_servers"
  incremental_backup = false
  schedule = {
    every_hours = 24
  }
  full_retention = {
    max_full_backup = 10
  }
  instance_ids = [vkcs_compute_instance.basic.id]
}

data "vkcs_backup_checkpoints" "plan" {
  plan_id = vkcs_backup_plan.checkpoint_source.id
}

output "legal_hold_checkpoint_ids" {
  value = flatten([
    for resource in data.vkcs_backup_checkpoints.plan.resources : [
      for backup in resource.backups : backup.checkpoint_id
      if try(backup.object_lock.type, "") == "legalhold"
    ]
  ])
}
