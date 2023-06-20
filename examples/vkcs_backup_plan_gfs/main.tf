resource "vkcs_backup_plan" "backup_plan" {
  name          = "backup-plan-tf-example"
  provider_name = "cloud_servers"
  incremental_backup = false
  # Backup the instance three times in week at 23:00 (02:00 MSK next day)
  schedule = {
    date = ["Mo", "We", "Fr"]
    time = "23:00"
  }
  # Keep backups: one for every last four weeks, one for every month of the last year, one for last two years
  gfs_retention = {
    gfs_weekly  = 4
    gfs_monthly = 11
    gfs_yearly  = 2
  }
  instance_ids       = [vkcs_compute_instance.basic.id]
}
