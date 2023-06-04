resource "vkcs_backup_plan" "backup_plan" {
  name          = "backup-plan-tf-example"
  provider_name = "cloud_servers"
  schedule = {
    date = ["Tu"]
    time = "08:12"
  }
  gfs_retention = {
    gfs_weekly  = 10
    gfs_monthly = 2
    gfs_yearly  = 1
  }
  incremental_backup = false
  instance_ids       = [vkcs_compute_instance.basic.id]
}
