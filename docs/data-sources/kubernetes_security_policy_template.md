---
subcategory: "Kubernetes"
layout: "vkcs"
page_title: "vkcs: vkcs_kubernetes_security_policy_template"
description: |-
  Get information on a kubernetes security policy template.
---

# vkcs_kubernetes_security_policy_template

Provides a kubernetes security policy template datasource. This can be used to get information about an VKCS kubernetes security policy template.

**New since v0.7.0**.

## Example Usage

```terraform
data "vkcs_kubernetes_security_policy_template" "replicalimits" {
  name = "k8sreplicalimits"
}
```

## Argument Reference
- `name` **required** *string* &rarr;  Name of the security policy template.

- `region` optional *string* &rarr;  The region in which to obtain the service client. If omitted, the `region` argument of the provider is used.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `created_at` *string* &rarr;  Template creation timestamp

- `description` *string* &rarr;  Description of the security policy template.

- `id` *string* &rarr;  ID of the resource.

- `settings_description` *string* &rarr;  Security policy settings description.

- `updated_at` *string* &rarr;  Template update timestamp.

- `version` *string* &rarr;  Version of the security policy template.


