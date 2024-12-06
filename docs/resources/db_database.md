---
subcategory: "Databases"
layout: "vkcs"
page_title: "vkcs: vkcs_db_database"
description: |-
  Manages a db database.
---

# vkcs_db_database

Provides a db database resource. This can be used to create, modify and delete db databases.

## Example Usage

```terraform
resource "vkcs_db_database" "mysql_db" {
  name    = "testdb"
  dbms_id = vkcs_db_instance.mysql.id
  charset = "utf8"
  collate = "utf8_general_ci"
}
```
## Argument Reference
- `dbms_id` **required** *string* &rarr;  ID of the instance or cluster that database is created for.

- `name` **required** *string* &rarr;  The name of the database. Changing this creates a new database.

- `charset` optional *string* &rarr;  Type of charset used for the database. Changing this creates a new database.

- `collate` optional *string* &rarr;  Collate option of the database.  Changing this creates a new database.

- `vendor_options` optional &rarr;  <br>**New since v0.5.5**.
  - `force_deletion` optional *boolean* &rarr;  Whether to try to force delete the database. Some datastores restricts regular database deletion in some circumstances but provides force deletion for that cases.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `dbms_type` *string* &rarr;  Type of dbms for the database, can be "instance" or "cluster".

- `id` *string* &rarr;  ID of the resource.



## Import

Databases can be imported using the `dbms_id/name`

```shell
terraform import vkcs_db_database.mydb 67691f3e-a4d1-443e-b1e9-717f505cc458/mydbname
```

After the import you can use ```terraform show``` to view imported fields and write their values to your .tf file.

You should at least add following fields to your .tf file:

`name, dbms_id`
