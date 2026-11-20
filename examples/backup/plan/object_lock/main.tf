resource "vkcs_backup_plan" "protected" {
  count              = var.enable_object_lock ? 1 : 0
  name               = "backup-plan-object-lock-example"
  provider_name      = "cloud_servers"
  incremental_backup = false
  schedule = {
    every_hours = 24
  }
  full_retention = {
    max_full_backup = 30
  }
  instance_ids = [var.instance_id]
  object_lock = {
    type = "compliance"
    days = 30
  }
}
