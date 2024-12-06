resource "vkcs_backup_plan" "backup_plan" {
  name               = "backup-plan-tf-example"
  provider_name      = "cloud_servers"
  incremental_backup = true
  # Create full backup every Monday at 04:00 MSK
  # Incremental backups are created each other day at the same time
  schedule = {
    date = ["Mo"]
    time = "04:00+03"
  }
  full_retention = {
    max_full_backup = 25
  }
  instance_ids = [vkcs_compute_instance.basic.id]
}
