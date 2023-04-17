---
layout: "vkcs"
page_title: "vkcs: vkcs_db_user"
description: |-
  Get information on a db user.
---

# vkcs_db_user

Use this data source to get the information on a db user resource.

## Example Usage

```terraform
data "vkcs_db_user" "db-user" {
  id = "d0438c95-1854-4ced-806c-c3b4125b586d/myuser"
}
```

## Argument Reference
- `id` **String** (***Required***) The id of the user in form "dbms_id/user_name".

- `databases` **String** (*Optional*) List of names of the databases, that user is created for.

- `dbms_id` **String** (*Optional*) ID of the instance or cluster that user is created for.

- `host` **String** (*Optional*) IP address of the host that user will be accessible from.

- `name` **String** (*Optional*) The name of the user. Changing this creates a new user.

- `password` **String** (*Optional*) The password of the user.


## Attributes Reference
No additional attributes are exported.

