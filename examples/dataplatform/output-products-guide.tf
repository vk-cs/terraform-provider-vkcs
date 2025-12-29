output "guide_products" {
  value = [
    for p in data.vkcs_dataplatform_products.products.products :
    p if strcontains(p.product_name, "click") || strcontains(p.product_name, "trino")
  ]
}
