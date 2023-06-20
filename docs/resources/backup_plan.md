---
subcategory: "Backup"
layout: "vkcs"
page_title: "vkcs: vkcs_backup_plan"
description: |-
  Manages a backup plan resource within VKCS.
---

# vkcs_backup_plan

Manages a backup plan resource.

## Example Usage
### Create plan for compute instance with full_retention policy and with incremental backups
```terraform
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
```

### Create plan for compute instance with gfs_retention policy, without incremental backups, using UTC timezone
```terraform
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
```

### Create plan for db instance with full_retention policy, making backup every 12 hours
```terraform
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
```

## Argument Reference
- `incremental_backup` **required** *boolean* &rarr;  Whether incremental backups strategy should be used. If enabled, the schedule.date field must specify one day, on which full backup will be created. On other days, incremental backups will be created. <br>**Note:** This option may be enabled for only for 'cloud_servers' provider.

- `instance_ids` **required** *string* &rarr;  List of ids of instances to make backup for

- `name` **required** *string* &rarr;  Name of the backup plan

- `schedule` , ***required***
  - `date` optional *string* &rarr;  List of days when to perform backups. If incremental_backups is enabled, only one day should be specified

  - `every_hours` optional *number* &rarr;  Hour interval of backups, must be one of: 3, 12, 24. This field is incompatible with date/time fields

  - `time` optional *string* &rarr;  Time of backup in format hh:mm (for UTC timezone) or hh:mm+tz (for other timezones)


- `full_retention` , optional &rarr;  Parameters for full retention policy. If incremental_backup is enabled, policy specifies number of full backups stored. Incompatible with gfs_retention
  - `max_full_backup` **required** *number* &rarr;  Maximum number of backups


- `gfs_retention` , optional &rarr;  Parameters for gfs retention policy. If incremental_backup is enabled, policy specifies number of full backups stored. Incompatible with full_retention
  - `gfs_weekly` **required** *number* &rarr;  Number of weeks to store backups

  - `gfs_monthly` optional *number* &rarr;  Number of months to store backups

  - `gfs_yearly` optional *number* &rarr;  Number of years to store backups


- `provider_id` optional *string* &rarr;  ID of backup provider

- `provider_name` optional *string* &rarr;  Name of backup provider, must be one of: cloud_servers, dbaas

- `region` optional *string* &rarr;  The `region` to fetch availability zones from, defaults to the provider's `region`.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource



## Import

Keypairs can be imported using the `name`, e.g.
```shell
terraform import vkcs_backup_plan.mybackupplan 5dfe75cb-a00f-4bc8-9551-bd38f64747e7
```
