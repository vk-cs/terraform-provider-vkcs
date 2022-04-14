---
layout: "vkcs"
page_title: "vkcs: db_user"
subcategory: ""
description: |-
  Manages a db user.
---

# vkcs\_db\_user

Provides a db user resource. This can be used to create, modify and delete db user.

## Example Usage

```terraform

resource "vkcs_db_user" "myuser" {
  name        = "myuser"
  password    = "password"
  
  instance_id = example_db_instance_id
  
  databases   = [example_db_database_name, example_db_other_database_name]
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the user. Changing this creates a new user.

* `password` - (Required) The password of the user.

* `host` - IP address of the host that user will be accessible from.

* `dbms_id` - (Optional) ID of the instance or cluster that user is created for.

* `databases` - List of names of the databases, that user is created for.

Either `instance_id` or `dbms_id` must be configured.

## Import

Users can be imported using the `dbms_id/name`

```
$ terraform import vkcs_db_user.myuser my_dbms_id/myusername
```

After the import you can use ```terraform show``` to view imported fields and write their values to your .tf file.

You should at least add following fields to your .tf file:

`name, dbms_id, password`
