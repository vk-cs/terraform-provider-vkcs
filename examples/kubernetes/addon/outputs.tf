output "installed-addons" {
  value = [
    for a in data.vkcs_kubernetes_addons.cluster-addons.addons : a.name
    if a.installed
  ]
  description = "K8s cluster installed addons"

  depends_on = [
    vkcs_kubernetes_addon.kube-prometheus-stack
  ]
}
