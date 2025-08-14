---
subcategory: "Databases"
layout: "vkcs"
page_title: "vkcs: vkcs_db_backup"
description: |-
  Manages a db backup.
---

# vkcs_db_backup

Provides a db backup resource. This can be used to create and delete db backup.

**New since v0.1.4**.

## Example Usage

```terraform
resource "vkcs_db_backup" "mysql_backup" {
  name    = "mssql-backup"
  dbms_id = vkcs_db_instance.mysql.id
}
```
## Argument Reference
- `dbms_id` **required** *string* &rarr;  ID of the instance or cluster, to create backup of.

- `name` **required** *string* &rarr;  The name of the backup. Changing this creates a new backup

- `container_prefix` optional *string* &rarr;  Prefix of S3 bucket ([prefix] - [project_id]) to store backup data. Default: databasebackups

- `description` optional *string* &rarr;  The description of the backup

- `region` optional *string* &rarr;  The region in which to obtain the service client. If omitted, the `region` argument of the provider is used.<br>**New since v0.4.0**.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `created` *string* &rarr;  Backup creation timestamp

- `datastore`  *list* &rarr;  Object that represents datastore of backup
    - `type` *string* &rarr;  Version of the datastore. Changing this creates a new instance.

    - `version` *string* &rarr;  Type of the datastore. Changing this creates a new instance.


- `dbms_type` *string* &rarr;  Type of dbms for the backup, can be "instance" or "cluster".

- `id` *string* &rarr;  ID of the resource.

- `location_ref` *string* &rarr;  Location of backup data on backup storage

- `meta` *string* &rarr;  Metadata of the backup

- `size` *number* &rarr;  Backup's volume size

- `updated` *string* &rarr;  Timestamp of backup's last update

- `wal_size` *number* &rarr;  Backup's WAL volume size



## Import

Backups can be imported using the `id`, e.g.

```shell
terraform import vkcs_db_backup.mybackup 67b86ce7-0924-48a6-8a18-683cfc6b4183
```

After the import you can use ```terraform show``` to view imported fields and write their values to your .tf file.
