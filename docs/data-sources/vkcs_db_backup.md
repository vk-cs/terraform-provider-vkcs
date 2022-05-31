---
layout: "vkcs"
page_title: "vkcs: db_backup"
subcategory: ""
description: |-
  Get information on a db backup.
---

# vkcs\_db\_backup

Use this data source to get the information on a db backup resource.
**New since v.0.1.4**.

## Example Usage

```terraform

data "vkcs_db_backup" "db-backup" {
  id = "d27fbf1a-373a-479c-b951-31041756f289"
}
```
## Argument Reference

The following arguments are supported:

* `id` - (Required) The UUID of the backup.

## Attributes reference

The following attributes are exported:

* `name` - The name of the backup. Changing this creates a new backup

* `dbms_id` - ID of the instance or cluster, to create backup of

* `description` - The description of the backup

* `container_prefix` - Prefix of S3 bucket (<prefix>-<project_id>) to store backup data. Default: databasebackups

* `location_ref` - Location of backup data on backup storage

* `created` - Backup creation timestamp

* `updated` - Timestamp of backup's last update

* `size` - Backup's volume size

* `wal_size` - Backup's WAL volume size

* `datastore` - Object that represents datastore of backup
    * `type` - (Required) Type of the datastore. Changing this creates a new instance.
    * `version` - (Required) Version of the datastore. Changing this creates a new instance.

* `meta` - Metadata of the backup
