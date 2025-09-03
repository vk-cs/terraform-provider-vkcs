resource "vkcs_iam_service_user" "service_user" {
  name        = "tf-example-service-user"
  description = "Service user created by Terraform example"
  role_names = [
    "mcs_admin_vm",
    "mcs_admin_network"
  ]
}

output "credentials" {
  value = {
    login    = vkcs_iam_service_user.service_user.login
    password = vkcs_iam_service_user.service_user.password
  }
  sensitive = true
}
