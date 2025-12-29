output "clickhouse_user_roles" {
  value = data.vkcs_dataplatform_product.clickhouse.configs.user_roles[*].name
}
