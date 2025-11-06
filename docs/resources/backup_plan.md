---
subcategory: "Backup"
layout: "vkcs"
page_title: "vkcs: vkcs_backup_plan"
description: |-
  Manages a backup plan resource within VKCS.
---

# vkcs_backup_plan

Manages a backup plan resource.

**New since v0.4.0**.

## Example Usage
### Incremental backup for compute instance
```terraform
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
```

### Full backup with GFS retention policy for compute instance
```terraform
resource "vkcs_backup_plan" "backup_plan" {
  name               = "backup-plan-tf-example"
  provider_name      = "cloud_servers"
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
  instance_ids = [vkcs_compute_instance.basic.id]
}
```

### Backup for db instance
```terraform
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
  instance_ids = [vkcs_db_instance.mysql.id]
}
```

## Argument Reference
- `incremental_backup` **required** *boolean* &rarr;  Whether incremental backups strategy should be used. If enabled, the schedule.date field must specify one day, on which full backup will be created. On other days, incremental backups will be created. <br>**Note:** This option may be enabled for only for 'cloud_servers' provider.

- `name` **required** *string* &rarr;  Name of the backup plan

- `schedule` ***required***
    - `date` optional *string* &rarr;  List of days when to perform backups. If incremental_backups is enabled, only one day should be specified

    - `every_hours` optional *number* &rarr;  Hour interval of backups, must be one of: 3, 12, 24. This field is incompatible with date/time fields

    - `time` optional *string* &rarr;  Time of backup in format hh:mm (for UTC timezone) or hh:mm+tz (for other timezones, e.g. 10:00+03 for MSK, 10:00-04 for ET)


- `backup_targets`  *list* &rarr;  List of backup targets specifying instance_id and volume_ids for each instance. Either backup_targets or instance_ids must be specified, but not both.<br>**New since v0.13.1**.
    - `instance_id` **required** *string* &rarr;  ID of the instance for which specific volumes are backed up.

    - `volume_ids` optional *set of* *string* &rarr;  List of volume IDs to back up for the instance. If no list is specified, backups will be created for all disks.


- `full_retention` optional &rarr;  Parameters for full retention policy. Specifies number of full backups stored. Incremental backups (if enabled) are not counted as full. Incompatible with gfs_retention
    - `max_full_backup` **required** *number* &rarr;  Maximum number of backups


- `gfs_retention` optional &rarr;  Parameters for gfs retention policy. Specifies number of full backups stored. Incremental backups (if enabled) are not counted as full. Incompatible with full_retention
    - `gfs_weekly` **required** *number* &rarr;  Number of weeks to store backups

    - `gfs_monthly` optional *number* &rarr;  Number of months to store backups

    - `gfs_yearly` optional *number* &rarr;  Number of years to store backups


- `instance_ids` optional *string* &rarr;  List of ids of instances to make backup for. Either backup_targets or instance_ids must be specified, but not both.

- `provider_id` optional *string* &rarr;  ID of backup provider

- `provider_name` optional *string* &rarr;  Name of backup provider, must be one of: cloud_servers, dbaas

- `region` optional *string* &rarr;  The `region` to fetch availability zones from, defaults to the provider's `region`. Changing this creates a new plan.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource



## Import

Backup plan can be imported using the `name`, e.g.
```shell
terraform import vkcs_backup_plan.mybackupplan 5dfe75cb-a00f-4bc8-9551-bd38f64747e7
```
