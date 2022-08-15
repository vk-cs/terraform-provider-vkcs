---
layout: "vkcs"
page_title: "vkcs: vkcs_db_backup"
description: |-
  Get information on a db backup.
---

# vkcs_db_backup

Use this data source to get the information on a db backup resource.

## Example Usage

```terraform
data "vkcs_db_backup" "db-backup" {
  backup_id = "d27fbf1a-373a-479c-b951-31041756f289"
}
```

## Argument Reference
- `backup_id` **String** (***Required***) The UUID of the backup.

- `description` **String** (*Optional*) The description of the backup


## Attributes Reference
- `backup_id` **String** See Argument Reference above.

- `description` **String** See Argument Reference above.

- `created` **String** Backup creation timestamp

- `datastore` **Object** Object that represents datastore of backup

- `dbms_id` **String** ID of the backed up instance or cluster

- `dbms_type` **String** Type of dbms of the backup, can be "instance" or "cluster".

- `id` **String** ID of the resource.

- `location_ref` **String** Location of backup data on backup storage

- `meta` **String** Metadata of the backup

- `name` **String** The name of the backup.

- `size` **Number** Backup's volume size

- `updated` **String** Timestamp of backup's last update

- `wal_size` **Number** Backup's WAL volume size


