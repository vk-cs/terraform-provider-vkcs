---
layout: "vkcs"
page_title: "Save Terraform state in VKCS Cloud Storage"
description: |-
  Save Terraform state in VKCS Cloud Storage
---

## Configure a bucket
The first thing you need is to [create an account and credentials in cloud storage](https://cloud.vk.com/docs/en/storage/s3/service-management/account-management).

You will also need a bucket to save the terraform state, you can do this by following one of these documentations: 
1) [Create bucket using terraform](https://cloud.vk.com/docs/en/tools-for-using-services/terraform/how-to-guides/aws)
2) [Create bucket using website or aws cli](https://cloud.vk.com/docs/en/storage/s3/service-management/buckets/create-bucket)

~> **Note:** If you decide to create a bucket using terraform, it is important to place terraform manifests in another folder.

## Configure terraform

To set up the configuration, you will need the region and domain of the VKCS Cloud Storage.
You can find them in the preparation steps of this [documentation](https://cloud.vk.com/docs/en/tools-for-using-services/terraform/how-to-guides/aws).

After that, in the main project, where we want to save terraform state in bucket, we need to add the following setting:
```terraform
terraform {
  backend "s3" {
    bucket     = "<bucket-name>"
    key        = "<path/to/tfstate>"
    region     = "<region>"
    access_key = "<public access key>"
    secret_key = "<secret key>"

    skip_credentials_validation = true
    skip_metadata_api_check     = true
    skip_requesting_account_id  = true
    skip_region_validation      = true
    skip_s3_checksum            = true
    endpoints = {
      s3 = "<domain>"
    }
  }
}
```
~> **Attention:** All arguments of terraform resources, including passwords, will be stored in the raw state as plain-text. 
[Read more about sensitive data in state](https://www.terraform.io/docs/language/state/sensitive-data.html)

Now you can use terraform as usual, and terraform state will be automatically saved to the bucket using the specified key.
