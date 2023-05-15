---
subcategory: "Databases"
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
- `backup_id` **required** *string* &rarr;  The UUID of the backup.

- `description` optional *string* &rarr;  The description of the backup


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `created` *string* &rarr;  Backup creation timestamp

- `datastore` *object* &rarr;  Object that represents datastore of backup

- `dbms_id` *string* &rarr;  ID of the backed up instance or cluster

- `dbms_type` *string* &rarr;  Type of dbms of the backup, can be "instance" or "cluster".

- `id` *string* &rarr;  ID of the resource.

- `location_ref` *string* &rarr;  Location of backup data on backup storage

- `meta` *string* &rarr;  Metadata of the backup

- `name` *string* &rarr;  The name of the backup.

- `size` *number* &rarr;  Backup's volume size

- `updated` *string* &rarr;  Timestamp of backup's last update

- `wal_size` *number* &rarr;  Backup's WAL volume size


