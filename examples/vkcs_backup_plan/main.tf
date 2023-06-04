resource "vkcs_backup_plan" "backup_plan" {
  name          = "backup-plan-tf-example"
  provider_name = "cloud_servers"
  schedule = {
    date = ["Tu"]
    time = "11:12+03"
  }
  full_retention = {
    max_full_backup = 25
  }
  incremental_backup = true
  instance_ids       = [vkcs_compute_instance.basic.id]
}
