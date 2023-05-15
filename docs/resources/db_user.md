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

resource "vkcs_db_user" "db-user" {
  name        = "testuser"
  password    = "SomePass1_"

  dbms_id     = vkcs_db_instance.db-instance.id

  databases   = [vkcs_db_database.db-database.name]
}
```
## Argument Reference
- `dbms_id` **required** *string* &rarr;  ID of the instance or cluster that user is created for.

- `name` **required** *string* &rarr;  The name of the user. Changing this creates a new user.

- `password` **required** sensitive *string* &rarr;  The password of the user.

- `databases` optional *string* &rarr;  List of names of the databases, that user is created for.

- `host` optional *string* &rarr;  IP address of the host that user will be accessible from.


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
