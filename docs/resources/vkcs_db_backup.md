---
layout: "vkcs"
page_title: "vkcs: db_backup"
subcategory: ""
description: |-
  Manages a db backup.
---

# vkcs\_db\_backup

Provides a db backup resource. This can be used to create and delete db backup.
**New since v.0.1.4**.

## Example Usage

```terraform

resource "vkcs_db_instance" "db-instance" {
  name = "db-instance"

  datastore {
    type    = "postgresql"
    version = "13"
  }

  floating_ip_enabled = true

  flavor_id         = "c8c42890-1ae9-411f-8cce-42e2d7c9b7d0"
  availability_zone = "MS1"

  size        = 8
  volume_type = "ceph-ssd"
  disk_autoexpand {
    autoexpand    = true
    max_disk_size = 1000
  }

  network {
    uuid = "4a6883c9-dd84-488d-a48f-afa5e37d3e2f"
  }

  capabilities {
    name = "node_exporter"
  }

  capabilities {
    name = "postgres_extensions"
  }
}

resource "vkcs_db_backup" "db-backup" {
    name = "db-backup"
    dbms_id = vkcs_db_instance.db-instance.id
    description = "tf-backup"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the backup. Changing this creates a new backup

* `dbms_id` - (Required) ID of the instance or cluster, to create backup of.

* `description` - The description of the backup

* `container_prefix` - Prefix of S3 bucket (<prefix>-<project_id>) to store backup data. Default: databasebackups

## Attributes reference

The following attributes are exported:

* `location_ref` - Location of backup data on backup storage

* `created` - Backup creation timestamp

* `updated` - Timestamp of backup's last update

* `size` - Backup's volume size

* `wal_size` - Backup's WAL volume size

* `datastore` - Object that represents datastore of backup
    * `type` - (Required) Type of the datastore. Changing this creates a new instance.
    * `version` - (Required) Version of the datastore. Changing this creates a new instance.

* `meta` - Metadata of the backup

## Import

Backups can be imported using the `id`, e.g.

```
$ terraform import vkcs_db_backup.mybackup 67b86ce7-0924-48a6-8a18-683cfc6b4183
```

After the import you can use ```terraform show``` to view imported fields and write their values to your .tf file.
