---
subcategory: "Databases"
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
- `id` **required** *string* &rarr;  The id of the database in form "dbms_id/db_name".

- `charset` optional *string* &rarr;  Type of charset used for the database.

- `collate` optional *string* &rarr;  Collate option of the database.

- `dbms_id` optional *string* &rarr;  ID of the instance or cluster that database is created for.

- `name` optional *string* &rarr;  The name of the database.


## Attributes Reference
No additional attributes are exported.

