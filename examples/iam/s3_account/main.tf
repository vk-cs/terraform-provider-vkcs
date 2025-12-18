resource "vkcs_iam_s3_account" "s3_account" {
  name        = "tf-example-s3-account"
  description = "S3 account created by Terraform example"
}

output "credentials" {
  value = {
    access_key = vkcs_iam_s3_account.s3_account.access_key
    secret_key = vkcs_iam_s3_account.s3_account.secret_key
  }
  sensitive = true
}
