output "trino_crontabs" {
  value = data.vkcs_dataplatform_product.trino.configs.crontabs[*]
}
