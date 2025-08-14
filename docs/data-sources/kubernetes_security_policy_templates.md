---
subcategory: "Kubernetes"
layout: "vkcs"
page_title: "vkcs: vkcs_kubernetes_security_policy_templates"
description: |-
  Get information on a kubernetes security policy templates.
---

# vkcs_kubernetes_security_policy_templates

Provides a kubernetes security policy templates datasource. This can be used to get information about all available VKCS kubernetes security policy templates.

**New since v0.7.0**.

## Example Usage

```terraform
data "vkcs_kubernetes_security_policy_templates" "templates" {}
```

## Argument Reference
- `region` optional *string* &rarr;  The region in which to obtain the service client. If omitted, the `region` argument of the provider is used.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  Random identifier of the data source.

- `security_policy_templates`  *list* &rarr;  Available kubernetes security policy templates.
    - `created_at` *string* &rarr;  Template creation timestamp

    - `description` *string* &rarr;  Description of the security policy template.

    - `id` *string* &rarr;  ID of the template.

    - `name` *string* &rarr;  Name of the security policy template.

    - `settings_description` *string* &rarr;  Security policy settings description.

    - `updated_at` *string* &rarr;  Template update timestamp.

    - `version` *string* &rarr;  Version of the security policy template.



