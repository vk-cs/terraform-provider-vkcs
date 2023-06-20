resource "vkcs_backup_plan" "backup_plan" {
  name          = "backup-plan-tf-example"
  provider_name = "dbaas"
  # Must be false since DBaaS does not support incremental backups
  incremental_backup = false
  # Backup database data every 12 hours since the next hour after the plan creation
  schedule = {
    every_hours = 12
  }
  full_retention = {
    max_full_backup = 25
  }
  instance_ids       = [vkcs_db_instance.basic.id]
}
