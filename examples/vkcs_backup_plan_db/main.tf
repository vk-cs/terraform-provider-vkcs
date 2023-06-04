resource "vkcs_backup_plan" "backup_plan" {
  name          = "backup-plan-tf-example"
  provider_name = "dbaas"
  schedule = {
    every_hours = 12
  }
  full_retention = {
    max_full_backup = 25
  }
  incremental_backup = false
  instance_ids       = [vkcs_db_instance.basic.id]
}
