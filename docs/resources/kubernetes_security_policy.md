---
subcategory: "Kubernetes"
layout: "vkcs"
page_title: "vkcs: vkcs_kubernetes_security_policy"
description: |-
  Manages a kubernetes cluster security policy.
---

# vkcs_kubernetes_security_policy

Provides a kubernetes cluster security policy resource. This can be used to create, modify and delete kubernetes security policies.

**New since v0.7.0**.

## Example Usage
```terraform
locals {
  policy_settings = {
    "ranges" = [
      {
        "min_replicas" = 1
        "max_replicas" = 2
      }
    ]
  }
}

resource "vkcs_kubernetes_security_policy" "replicalimits" {
  cluster_id                  = vkcs_kubernetes_cluster.k8s-cluster.id
  enabled                     = true
  namespace                   = "*"
  policy_settings             = jsonencode(local.policy_settings)
  security_policy_template_id = data.vkcs_kubernetes_security_policy_template.replicalimits.id
}
```
## Argument Reference
- `cluster_id` **required** *string* &rarr;  The ID of the kubernetes cluster. Changing this creates a new security policy.

- `namespace` **required** *string* &rarr;  Namespace to apply security policy to.

- `policy_settings` **required** *string* &rarr;  Policy settings.

- `security_policy_template_id` **required** *string* &rarr;  The ID of the security policy template. Changing this creates a new security policy.

- `enabled` optional *boolean* &rarr;  Controls whether the security policy is enabled. Default is true.

- `region` optional *string* &rarr;  The region in which to obtain the Container Infra client. If omitted, the `region` argument of the provider is used. Changing this creates a new security policy.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `created_at` *string* &rarr;  Creation timestamp

- `id` *string* &rarr;  ID of the resource

- `updated_at` *string* &rarr;  Update timestamp.



## Import

Security policies can be imported using the `id`, e.g.

```shell
terraform import vkcs_kubernetes_security_policy.sp 723bfe25-5b2b-4410-aba0-1c0ef6d1c8b0
```
