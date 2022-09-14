---
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
- `id` **String** (***Required***) The UUID of the instance

- `region` **String** (*Optional*) The region in which to obtain the Compute client. If omitted, the `region` argument of the provider is used.

- `user_data` **String** (*Optional*) The user data added when the server was created.


## Attributes Reference
- `id` **String** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `user_data` **String** See Argument Reference above.

- `access_ip_v4` **String** The first IPv4 address assigned to this server.

- `access_ip_v6` **String** The first IPv6 address assigned to this server.

- `availability_zone` **String** The availability zone of this server.

- `flavor_id` **String** The flavor ID used to create the server.

- `flavor_name` **String** The flavor name used to create the server.

- `image_id` **String** The image ID used to create the server.

- `image_name` **String** The image name used to create the server.

- `key_pair` **String** The name of the key pair assigned to this server.

- `metadata` <strong>Map of </strong>**String** A set of key/value pairs made available to the server.

- `name` **String** The name of the server.

- `network`  An array of maps, detailed below.
  - `fixed_ip_v4` **String** The IPv4 address assigned to this network port.

  - `fixed_ip_v6` **String** The IPv6 address assigned to this network port.

  - `mac` **String** The MAC address assigned to this network interface.

  - `name` **String** The name of the network

  - `port` **String** The port UUID for this network

  - `uuid` **String** The UUID of the network

- `power_state` **String** VM state

- `security_groups` <strong>Set of </strong>**String** An array of security group names associated with this server.

- `tags` <strong>Set of </strong>**String** A set of string tags for the instance.


