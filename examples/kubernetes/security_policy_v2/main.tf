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
