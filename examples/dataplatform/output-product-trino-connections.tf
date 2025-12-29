output "trino_connections" {
  value = {
    for c in data.vkcs_dataplatform_product.trino.configs.connections :
    c.plug => c.is_required
  }
}
