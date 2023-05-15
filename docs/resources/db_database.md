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
resource "vkcs_db_instance" "db-instance" {
  name        = "db-instance"

  availability_zone = "GZ1"

  datastore {
    type    = "mysql"
    version = "5.7"
  }

  flavor_id   = data.vkcs_compute_flavor.db.id

  size        = 8
  volume_type = "ceph-ssd"

  network {
    uuid = vkcs_networking_network.db.id
  }

  depends_on = [
    vkcs_networking_router_interface.db
  ]
}

resource "vkcs_db_database" "db-database" {
  name        = "testdb"
  dbms_id     = vkcs_db_instance.db-instance.id
  charset     = "utf8"
  collate     = "utf8_general_ci"
}
```
## Argument Reference
- `dbms_id` **required** *string* &rarr;  ID of the instance or cluster that database is created for.

- `name` **required** *string* &rarr;  The name of the database. Changing this creates a new database.

- `charset` optional *string* &rarr;  Type of charset used for the database. Changing this creates a new database.

- `collate` optional *string* &rarr;  Collate option of the database.  Changing this creates a new database.


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
