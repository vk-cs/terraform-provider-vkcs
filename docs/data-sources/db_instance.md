---
subcategory: "Databases"
layout: "vkcs"
page_title: "vkcs: vkcs_db_instance"
description: |-
  Get information on a db instance.
---

# vkcs_db_instance

Use this data source to get the information on a db instance resource.

## Example Usage

```terraform
data "vkcs_db_instance" "db-instance" {
  id = "e7da2869-2ae2-4900-99e3-a44fec2b11ac"
}
```

## Argument Reference
- `id` **required** *string* &rarr;  The id of the instance.

- `backup_schedule` optional &rarr;  Object that represents configuration of PITR backup. This functionality is available only for postgres datastore. **New since v0.1.4**.
  - `interval_hours` **required** *number* &rarr;  Time interval between backups, specified in hours. Available values: 3, 6, 8, 12, 24.

  - `keep_count` **required** *number* &rarr;  Number of backups to be stored.

  - `name` **required** *string* &rarr;  Name of the schedule.

  - `start_hours` **required** *number* &rarr;  Hours part of timestamp of initial backup.

  - `start_minutes` **required** *number* &rarr;  Minutes part of timestamp of initial backup.

- `datastore` optional &rarr;  Object that represents datastore of the instance.
  - `type` **required** *string* &rarr;  Type of the datastore.

  - `version` **required** *string* &rarr;  Version of the datastore.

- `flavor_id` optional *string* &rarr;  The ID of flavor for the instance.

- `hostname` optional *string* &rarr;  The hostname of the instance.

- `ip` optional *string* &rarr;  IP address of the instance.

- `name` optional *string* &rarr;  The name of the instance.

- `region` optional *string* &rarr;  Region of the resource.

- `status` optional *string* &rarr;  Instance status.

- `volume` optional &rarr;  Object that describes volume of the instance.
  - `size` **required** *number* &rarr;  Size of the instance volume.

  - `used` **required** *number* &rarr;  Size of the used volume space.

  - `volume_id` **required** *string* &rarr;  ID of the instance volume.

  - `volume_type` **required** *string* &rarr;  Type of the instance volume.


## Attributes Reference
No additional attributes are exported.

