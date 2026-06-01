---
subcategory: "Kubernetes"
layout: "vkcs"
page_title: "vkcs: vkcs_kubernetes_security_policy_template_v2"
description: |-
  Use this data source to retrieve information about specific VKCS Kubernetes cluster security policy template.
---

# vkcs_kubernetes_security_policy_template_v2

Use this data source to retrieve information about specific VKCS Kubernetes cluster security policy template.

The cluster security policy template can be identified by either:

* `id`
* the combination of `name` and `version`

These options are mutually exclusive.

## Example Usage
```terraform
data "vkcs_kubernetes_security_policy_template_v2" "replicalimits" {
  name    = "k8sreplicalimits"
  version = "1.0.0"
}
```

## Argument Reference
- `id` optional *string* &rarr;  The security policy template's ID. Used to fetch a specific security policy template.

- `name` optional *string* &rarr;  Name of the security policy template (e.g., `k8sallowedrepos`).

- `region` optional *string* &rarr;  The region for which to retrieve security policy templates. Defaults to provider's `region`.

- `version` optional *string* &rarr;  Version of the security policy template.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `description` *string* &rarr;  Brief description of the template's purpose.

- `settings_description` *string* &rarr;  Base64-encoded JSON schema defining the configurable parameters of the policy.


