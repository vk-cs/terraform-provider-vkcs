---
subcategory: "Virtual Machines"
layout: "vkcs"
page_title: "vkcs: vkcs_compute_instance"
description: |-
  Get information on an VKCS Instance
---

# vkcs_compute_instance

Use this data source to get the details of a running server

## Example Usage

```terraform
data "vkcs_compute_instance" "instance" {
  # Randomly generated UUID, for demonstration purposes
  id = "2ba26dc6-a12d-4889-8f25-794ea5bf4453"
}
```

## Argument Reference
- `id` **required** *string* &rarr;  The UUID of the instance

- `region` optional *string* &rarr;  The region in which to obtain the Compute client. If omitted, the `region` argument of the provider is used.

- `user_data` optional *string* &rarr;  The user data added when the server was created.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `access_ip_v4` *string* &rarr;  The first IPv4 address assigned to this server.

- `availability_zone` *string* &rarr;  The availability zone of this server.

- `flavor_id` *string* &rarr;  The flavor ID used to create the server.

- `flavor_name` *string* &rarr;  The flavor name used to create the server.

- `image_id` *string* &rarr;  The image ID used to create the server.

- `image_name` *string* &rarr;  The image name used to create the server.

- `key_pair` *string* &rarr;  The name of the key pair assigned to this server.

- `metadata` *map of* *string* &rarr;  A set of key/value pairs made available to the server.

- `name` *string* &rarr;  The name of the server.

- `network`  &rarr;  An array of maps, detailed below.
  - `fixed_ip_v4` *string* &rarr;  The IPv4 address assigned to this network port.

  - `mac` *string* &rarr;  The MAC address assigned to this network interface.

  - `name` *string* &rarr;  The name of the network

  - `port` *string* &rarr;  The port UUID for this network

  - `uuid` *string* &rarr;  The UUID of the network

- `power_state` *string* &rarr;  VM state

- `security_groups` *set of* *string* &rarr;  An array of security group names associated with this server.

- `tags` *set of* *string* &rarr;  A set of string tags for the instance.


