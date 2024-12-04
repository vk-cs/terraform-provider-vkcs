---
subcategory: "Cloud Monitoring"
layout: "vkcs"
page_title: "vkcs: vkcs_cloud_monitoring"
description: |-
  Manages a cloud monitoring for the given image within VKCS.
---

# vkcs_cloud_monitoring

Receives settings for cloud monitoring for the `vkcs_compute_instance'.

~> **Attention:**
If you create a virtual machine with cloud monitoring enabled, then take a disk snapshot and create a new instance from
it,
monitoring will also be enabled on the new one. If you then delete the `vkcs_cloud_monitoring` resource,
the monitoring service user will be deleted along with it, causing cloud monitoring to stop working.

## Example Usage

```terraform
resource "vkcs_cloud_monitoring" "basic" {
  image_id = data.vkcs_images_image.debian.id
}
```

## Argument Reference
- `image_id` **required** *string* &rarr;  ID of the image to create cloud monitoring for.

- `region` optional *string* &rarr;  The region in which to obtain the service client. If omitted, the `region` argument of the provider is used.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource.

- `script` *string* &rarr;  Shell script of the cloud monitoring.

- `service_user_id` *string* &rarr;  ID of the service monitoring user.



~> **Note:**
You can use this resource for multiple compute instances in the same project with the same image.

~> **Note:**
Monitoring script may be in bash or powershell format, depending on the OS.
