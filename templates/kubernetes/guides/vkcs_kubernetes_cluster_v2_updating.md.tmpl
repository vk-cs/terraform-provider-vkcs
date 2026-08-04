---
layout: "vkcs"
page_title: "Upgrading a Kubernetes cluster version"
description: |-
  Upgrade the Kubernetes version of a Managed Kubernetes cluster with the VKCS provider using Terraform.
---

# Upgrading a Kubernetes cluster version

This guide explains how to upgrade the Kubernetes version of a Managed Kubernetes (mk8s)
`v2` cluster with Terraform. It assumes you already have a cluster managed by the
[`vkcs_kubernetes_cluster_v2`](../resources/kubernetes_cluster_v2) resource — see
[Getting Started with Kubernetes](./kubernetes_cluster_v2_getting_started) if you don't.

## How the upgrade works

An upgrade is triggered by changing the `version` argument of the cluster resource to a
**higher** supported version. When you apply the change, VK Cloud performs a rolling update:

1. **Master nodes are updated first**, one at a time (remove from cluster → update → verify →
   return to cluster → next node). Clusters with 3-5 master nodes stay available for the API
   during this phase; a single-master cluster has a control-plane outage while its master is
   updated.
2. **Worker node groups are updated next**, also as a rolling update. The number of nodes taken
   out of service simultaneously in each group is controlled by that group's
   `parallel_upgrade_chunk` (the percentage of nodes that may be unavailable during the
   upgrade). Workers upgrade to the cluster's new version automatically — the
   [`vkcs_kubernetes_node_group_v2`](../resources/kubernetes_node_group_v2) resource has no
   separate `version` argument.

Alongside Kubernetes itself, some cluster components are updated (for example CoreDNS,
Gatekeeper, Shell Operator).

## Before you upgrade

* **Only upgrades are allowed.** You cannot downgrade — the provider rejects a `version` lower
  than the current one.
* **Upgrade one minor version at a time.** Kubernetes does not support skipping minor versions
  (e.g. go `1.32.x` → `1.33.x` → `1.34.x`, not `1.32.x` → `1.34.x`). Patch upgrades within a
  minor version are fine.
* **Check the supported target versions** with the
  [`vkcs_kubernetes_versions_v2`](../data-sources/kubernetes_versions_v2) data source instead of
  guessing a version string.
* **Ensure spare capacity.** Autoscaling does **not** run during an upgrade, so the cluster must
  have enough spare master and worker nodes to absorb the workload from nodes that are
  temporarily unavailable. VK Cloud recommends 3-5 master nodes and keeping ~1% extra worker
  capacity on top of the nodes being drained.
* **Back up a customized CoreDNS Corefile.** During the upgrade CoreDNS is reset to the default
  Corefile. If you modified it, back it up first and reapply your changes afterward.
* **Add-ons are not upgraded automatically.** Installed add-ons keep their versions across a
  cluster upgrade. Upgrade them separately afterward (see
  [`vkcs_kubernetes_cluster_addon_v2`](../resources/kubernetes_cluster_addon_v2)).
* **Match your `kubectl`.** Keep the `kubectl` minor version within one of the new cluster
  version.

## Step 1. Find the target version

```terraform
data "vkcs_kubernetes_versions_v2" "available" {}

output "available_k8s_versions" {
  value = data.vkcs_kubernetes_versions_v2.available.k8s_versions
}
```

```shell
terraform apply
# available_k8s_versions = ["v1.31.4", "v1.32.1", "v1.33.3", "v1.34.2"]
```

## Step 2. Tune the worker rollout speed (optional)

`parallel_upgrade_chunk` on each node group defines the maximum percentage of that group's nodes
that may be upgraded at once. Lower it for a more cautious, slower rollout; raise it to finish
faster at the cost of more simultaneously-unavailable nodes. Make sure the remaining nodes can
carry the workload of those being drained.

```terraform
resource "vkcs_kubernetes_node_group_v2" "workers" {
  # ... other arguments unchanged ...
  parallel_upgrade_chunk = 20 # upgrade at most 20% of this group's nodes at a time
}
```

For example, in a 9-node group, `parallel_upgrade_chunk = 33` upgrades three nodes at a time;
`20` upgrades two; `10` upgrades one.

## Step 3. Bump the cluster version

Change `version` on the cluster resource to the target version:

```terraform
resource "vkcs_kubernetes_cluster_v2" "k8s" {
  name    = "k8s-getting-started"
  version = "v1.34.2" # was "v1.33.3"

  # ... all other arguments unchanged ...
}
```

## Step 4. Review the plan

```shell
terraform plan
```

Confirm that the plan shows an **in-place update** of `version` on
`vkcs_kubernetes_cluster_v2` (`~ version = "v1.33.3" -> "v1.34.2"`) and **not** a
destroy/recreate (`-/+`). Changing `version` alone updates the cluster in place; if you see a
replacement, you also changed a `Forces replacement` argument by mistake — revert that change.

```console
  ~ resource "vkcs_kubernetes_cluster_v2" "k8s" {
        id      = "3GXEzAgwMCg3r8nRS05mf9ndTq8"
        name    = "k8s-getting-started"
      ~ version = "v1.33.3" -> "v1.34.2"
        # (other attributes unchanged)
    }

Plan: 0 to add, 1 to change, 0 to destroy.
```

## Step 5. Apply

```shell
terraform apply
```

The upgrade runs master-first, then the worker node groups. Expect it to take a while — the
control plane typically takes several minutes, plus rolling time per node group that depends on
node count and `parallel_upgrade_chunk`. Terraform blocks until the rollout finishes.

## Step 6. Verify

```shell
export KUBECONFIG=$(pwd)/kubeconfig.yaml
kubectl get nodes
```

All nodes should report the new version in the `VERSION` column and be `Ready`:

```console
NAME                                    STATUS   ROLES    AGE   VERSION
k8s-getting-started-default-workers-0   Ready    <none>   1h    v1.34.2
k8s-getting-started-default-workers-1   Ready    <none>   1h    v1.34.2
```

You can also check the cluster resource:

```shell
terraform state show vkcs_kubernetes_cluster_v2.k8s | grep -E "version|status"
```

## Step 7. Post-upgrade tasks

* Reapply any customizations to the CoreDNS Corefile if you use a modified one.
* Upgrade add-ons that need to move to a version compatible with the new Kubernetes release.
* Upgrade your local `kubectl` if its minor version now differs from the cluster by more than one.

## Recommendations

* Test the upgrade on a staging cluster before production.
* Upgrade sequentially, one minor version at a time, verifying the cluster is healthy between
  steps.
* Do not let the cluster fall behind supported versions — unsupported versions are not covered
  by technical support. Check end-of-support dates in the
  [version support policy](https://cloud.vk.com/docs/en/kubernetes/k8s/concepts/versions/version-support).

## See also

* [Getting Started with Kubernetes](./kubernetes_cluster_v2_getting_started)
* [`vkcs_kubernetes_cluster_v2`](../resources/kubernetes_cluster_v2)
* [`vkcs_kubernetes_node_group_v2`](../resources/kubernetes_node_group_v2)
* [`vkcs_kubernetes_versions_v2`](../data-sources/kubernetes_versions_v2)
