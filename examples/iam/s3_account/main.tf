resource "vkcs_iam_s3_account" "s3_account" {
  name        = "tf-example-s3-account"
  description = "S3 account created by Terraform example"
}

output "access_key" {
  value = vkcs_iam_s3_account.s3_account.access_key
}

output "secret_key" {
  value     = vkcs_iam_s3_account.s3_account.secret_key
  sensitive = true
}
