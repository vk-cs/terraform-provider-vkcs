---
subcategory: "Kubernetes"
layout: "vkcs"
page_title: "vkcs: vkcs_kubernetes_security_policy_templates_v2"
description: |-
  Use this data source to retrieve information about available VKCS Kubernetes cluster security policy templates.
---

# vkcs_kubernetes_security_policy_templates_v2

Use this data source to retrieve information about available VKCS Kubernetes cluster security policy templates.

## Example Usage
```terraform
data "vkcs_kubernetes_security_policy_templates_v2" "policies" {}
```

## Argument Reference
- `region` optional *string* &rarr;  The region for which to retrieve security policy templates. Defaults to provider's `region`.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  A synthetic identifier set to "policy_templates". This data source does not have a natural ID.

- `security_policies`  *set* &rarr;  A set of available security policy templates. Each template contains an `id`, `name`, `description`, `settings_description`, and `version`.
    - `description` *string* &rarr;  Brief description of the template's purpose.

    - `id` *string* &rarr;  Unique identifier of the security policy template.

    - `name` *string* &rarr;  Name of the security policy template (e.g., `k8sallowedrepos`).

    - `settings_description` *string* &rarr;  Base64-encoded JSON schema defining the configurable parameters of the policy.

    - `version` *string* &rarr;  Version of the security policy template.



