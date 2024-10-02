---
subcategory: "Databases"
layout: "vkcs"
page_title: "vkcs: vkcs_db_user"
description: |-
  Manages a db user.
---

# vkcs_db_user

Provides a db user resource. This can be used to create, modify and delete db user.

## Example Usage

```terraform
resource "vkcs_db_user" "mysql-user" {
  name        = "testuser"
  password    = "Test_p@ssword-12-3"

  dbms_id     = vkcs_db_instance.mysql.id

  databases   = [vkcs_db_database.mysql-db-1.name, vkcs_db_database.mysql-db-2.name]
}
```
## Argument Reference
- `dbms_id` **required** *string* &rarr;  ID of the instance or cluster that user is created for.

- `name` **required** *string* &rarr;  The name of the user. Changing this creates a new user.

- `password` **required** sensitive *string* &rarr;  The password of the user.

- `databases` optional *set of* *string* &rarr;  List of names of the databases, that user is created for.

- `host` optional *string* &rarr;  IP address of the host that user will be accessible from.

- `vendor_options` optional &rarr;  <br>**New since v0.5.3**.
  - `skip_deletion` optional *boolean* &rarr;  Boolean to control whether to user deletion should be skipped. If set to true, the resource will be removed from the state, but the remote object will not be deleted. This is useful for PostgreSQL, where users cannot be deleted from the API if they own database objects.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `dbms_type` *string* &rarr;  Type of dbms for the user, can be "instance" or "cluster".

- `id` *string* &rarr;  ID of the resource.



## Import

Users can be imported using the `dbms_id/name`

```shell
terraform import vkcs_db_user.myuser b29f9249-b0e0-43f2-9278-34ed8284a4dc/myusername
```

After the import you can use ```terraform show``` to view imported fields and write their values to your .tf file.

You should at least add following fields to your .tf file:

`name, dbms_id, password`
