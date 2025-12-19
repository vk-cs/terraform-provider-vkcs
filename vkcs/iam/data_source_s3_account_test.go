package iam_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccIAMS3AccountDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccIAMS3AccountDataSourceBasic, map[string]string{"TestAccIAMS3AccountDataSourceBase": testAccIAMS3AccountDataSourceBase}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.vkcs_iam_s3_account.basic", "account_id", "vkcs_iam_s3_account.base", "account_id"),
					resource.TestCheckResourceAttrPair("data.vkcs_iam_s3_account.basic", "account_name", "vkcs_iam_s3_account.base", "account_name"),
					resource.TestCheckResourceAttrPair("data.vkcs_iam_s3_account.basic", "access_key", "vkcs_iam_s3_account.base", "access_key"),
					resource.TestCheckResourceAttrPair("data.vkcs_iam_s3_account.basic", "description", "vkcs_iam_s3_account.base", "description"),
					resource.TestCheckResourceAttrPair("data.vkcs_iam_s3_account.basic", "name", "vkcs_iam_s3_account.base", "name"),
				),
			},
		},
	})
}

func TestAccIAMS3AccountDataSource_queryByName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccIAMS3AccountDataSourceQueryByName, map[string]string{"TestAccIAMS3AccountDataSourceQueryByNameBase": testAccIAMS3AccountDataSourceQueryByNameBase}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.vkcs_iam_s3_account.basic", "id", "vkcs_iam_s3_account.s3_account_1", "id"),
					resource.TestCheckResourceAttrPair("data.vkcs_iam_s3_account.basic", "name", "vkcs_iam_s3_account.s3_account_1", "name"),
				),
			},
		},
	})
}

const testAccIAMS3AccountDataSourceBase = `
resource "vkcs_iam_s3_account" "base" {
  name = "tfacc-s3-account-base"
}
`

const testAccIAMS3AccountDataSourceBasic = `
{{ .TestAccIAMS3AccountDataSourceBase }}

data "vkcs_iam_s3_account" "basic" {
  id = vkcs_iam_s3_account.base.id
}
`

const testAccIAMS3AccountDataSourceQueryByNameBase = `
resource "vkcs_iam_s3_account" "s3_account_1" {
  name = "tfacc-s3-account-1"
}

resource "vkcs_iam_s3_account" "s3_account_2" {
  name = "tfacc-s3-account-2"
}
`

const testAccIAMS3AccountDataSourceQueryByName = `
{{ .TestAccIAMS3AccountDataSourceQueryByNameBase }}

data "vkcs_iam_s3_account" "basic" {
  name = vkcs_iam_s3_account.s3_account_1.name
}
`
