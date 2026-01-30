---
subcategory: "Kubernetes"
layout: "vkcs"
page_title: "vkcs: vkcs_kubernetes_security_policy_v2"
description: |-
  Provides a Kubernetes cluster security policy resource. This can be used to create, modify, and delete Kubernetes cluster security policy.
---

# vkcs_kubernetes_security_policy_v2

Provides a Kubernetes cluster security policy resource. This can be used to create, modify, and delete Kubernetes cluster security policy.

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

resource "vkcs_kubernetes_security_policy_v2" "replicalimits" {
  cluster_id                  = vkcs_kubernetes_cluster_v2.k8s_cluster.id
  enabled                     = true
  namespace                   = "*"
  policy_settings             = jsonencode(local.policy_settings)
  security_policy_template_id = data.vkcs_kubernetes_security_policy_template_v2.replicalimits.id
}
```
## Argument Reference
- `cluster_id` **required** *string* &rarr;  The ID of the Kubernetes cluster. **Forces replacement** on change.

- `namespace` **required** *string* &rarr;  Namespace to apply security policy to. Changing this updates the security policy.

- `policy_settings` **required** *string* &rarr;  Policy settings. Changing this updates the security policy.

- `security_policy_template_id` **required** *string* &rarr;  The ID of the security policy template. **Forces replacement** on change.

- `enabled` optional *boolean* &rarr;  Controls whether the security policy is enabled. Default is `true`. Changing this updates the security policy.

- `region` optional *string* &rarr;  The region in which to create the Kubernetes security policy. If omitted, the provider's `region` is used. **Forces replacement** on change.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  The ID of the cluster security policy.



## Import

Security policies can be imported using the `id`, e.g.

```shell
terraform import vkcs_kubernetes_security_policy_v2.sp 3DGdkgFCVhIRjYUgHCKQfU4bO6J
```
