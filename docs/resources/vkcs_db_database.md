---
layout: "vkcs"
page_title: "vkcs: db_database"
subcategory: ""
description: |-
  Manages a db database.
---

# vkcs\_db\_database (Resource)

Provides a db database resource. This can be used to create, modify and delete db databases.

## Example Usage

```terraform

resource "vkcs_db_database" "mydb" {
  name        = "mydb"
  instance_id = example_db_instance_id
  charset     = "utf8"
  collate     = "utf8_general_ci"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the database. Changing this creates a new database.

* `dbms_id` - (Optional) ID of the instance or cluster that database is created for.

* `charset` - Type of charset used for the database. Changing this creates a new database.

* `collate` - Collate option of the database.  Changing this creates a new database.

Either `instance_id` or `dbms_id` must be configured.

## Import

Databases can be imported using the `dbms_id/name`

```
$ terraform import vkcs_db_database.mydb my_dbms_id/mydbname
```

After the import you can use ```terraform show``` to view imported fields and write their values to your .tf file.

You should at least add following fields to your .tf file:

`name, dbms_id`