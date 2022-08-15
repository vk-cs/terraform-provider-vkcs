---
layout: "vkcs"
page_title: "vkcs: vkcs_db_database"
description: |-
  Get information on a db database.
---

# vkcs_db_database

Use this data source to get the information on a db database resource.

## Example Usage

```terraform
data "vkcs_db_database" "db-database" {
  id = "325a2871-f311-45ac-ae91-bfef20fc768e/mydatabase"
}
```

## Argument Reference
- `id` **String** (***Required***) The id of the database in form "dbms_id/db_name".

- `charset` **String** (*Optional*) Type of charset used for the database.

- `collate` **String** (*Optional*) Collate option of the database.

- `dbms_id` **String** (*Optional*) ID of the instance or cluster that database is created for.

- `name` **String** (*Optional*) The name of the database.


## Attributes Reference
- `id` **String** See Argument Reference above.

- `charset` **String** See Argument Reference above.

- `collate` **String** See Argument Reference above.

- `dbms_id` **String** See Argument Reference above.

- `name` **String** See Argument Reference above.


