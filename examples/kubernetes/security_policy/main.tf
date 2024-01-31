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
