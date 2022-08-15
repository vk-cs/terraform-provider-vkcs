---
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
- `id` **String** (***Required***) The id of the instance.

- `backup_schedule` (*Optional*) Object that represents configuration of PITR backup. This functionality is available only for postgres datastore. **New since v.0.1.4**.
  - `interval_hours` **Number** (***Required***) Time interval between backups, specified in hours. Available values: 3, 6, 8, 12, 24.

  - `keep_count` **Number** (***Required***) Number of backups to be stored.

  - `name` **String** (***Required***) Name of the schedule.

  - `start_hours` **Number** (***Required***) Hours part of timestamp of initial backup.

  - `start_minutes` **Number** (***Required***) Minutes part of timestamp of initial backup.

- `datastore` (*Optional*) Object that represents datastore of the instance.
  - `type` **String** (***Required***) Type of the datastore.

  - `version` **String** (***Required***) Version of the datastore.

- `flavor_id` **String** (*Optional*) The ID of flavor for the instance.

- `hostname` **String** (*Optional*) The hostname of the instance.

- `ip` **String** (*Optional*) IP address of the instance.

- `name` **String** (*Optional*) The name of the instance.

- `region` **String** (*Optional*) Region of the resource.

- `status` **String** (*Optional*) Instance status.

- `volume` (*Optional*) Object that describes volume of the instance.
  - `size` **Number** (***Required***) Size of the instance volume.

  - `used` **Number** (***Required***) Size of the used volume space.

  - `volume_id` **String** (***Required***) ID of the instance volume.

  - `volume_type` **String** (***Required***) Type of the instance volume.


## Attributes Reference
- `id` **String** See Argument Reference above.

- `backup_schedule`  See Argument Reference above.
  - `interval_hours` **Number** See Argument Reference above.

  - `keep_count` **Number** See Argument Reference above.

  - `name` **String** See Argument Reference above.

  - `start_hours` **Number** See Argument Reference above.

  - `start_minutes` **Number** See Argument Reference above.

- `datastore`  See Argument Reference above.
  - `type` **String** See Argument Reference above.

  - `version` **String** See Argument Reference above.

- `flavor_id` **String** See Argument Reference above.

- `hostname` **String** See Argument Reference above.

- `ip` **String** See Argument Reference above.

- `name` **String** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `status` **String** See Argument Reference above.

- `volume`  See Argument Reference above.
  - `size` **Number** See Argument Reference above.

  - `used` **Number** See Argument Reference above.

  - `volume_id` **String** See Argument Reference above.

  - `volume_type` **String** See Argument Reference above.


