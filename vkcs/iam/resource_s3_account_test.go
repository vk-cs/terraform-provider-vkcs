package iam_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/iam/s3accounts"
)

func TestAccIAMS3Account_basic(t *testing.T) {
	var s3Account s3accounts.S3Account

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccIAMS3AccountBasic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3AccountExists("vkcs_iam_s3_account.basic", &s3Account),
					resource.TestCheckResourceAttr("vkcs_iam_s3_account.basic", "name", "tfacc-s3-account-basic"),
					resource.TestCheckResourceAttrSet("vkcs_iam_s3_account.basic", "access_key"),
					resource.TestCheckResourceAttrSet("vkcs_iam_s3_account.basic", "account_id"),
					resource.TestCheckResourceAttrSet("vkcs_iam_s3_account.basic", "account_name"),
					resource.TestCheckResourceAttrSet("vkcs_iam_s3_account.basic", "created_at"),
					resource.TestCheckResourceAttrSet("vkcs_iam_s3_account.basic", "secret_key"),
				),
			},
			acctest.ImportStep("vkcs_iam_s3_account.basic", "secret_key"),
		},
	})
}

func TestAccIAMS3Account_full(t *testing.T) {
	var s3Account s3accounts.S3Account

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccIAMS3AccountFull),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3AccountExists("vkcs_iam_s3_account.full", &s3Account),
					resource.TestCheckResourceAttr("vkcs_iam_s3_account.full", "name", "tfacc-s3-account-full"),
					resource.TestCheckResourceAttr("vkcs_iam_s3_account.full", "description", "S3 Account created by acceptance test"),
					resource.TestCheckResourceAttrSet("vkcs_iam_s3_account.full", "access_key"),
					resource.TestCheckResourceAttrSet("vkcs_iam_s3_account.full", "account_id"),
					resource.TestCheckResourceAttrSet("vkcs_iam_s3_account.full", "account_name"),
					resource.TestCheckResourceAttrSet("vkcs_iam_s3_account.full", "created_at"),
					resource.TestCheckResourceAttrSet("vkcs_iam_s3_account.full", "secret_key"),
				),
			},
			acctest.ImportStep("vkcs_iam_s3_account.full", "secret_key"),
		},
	})
}

func testAccCheckS3AccountExists(n string, s3Account *s3accounts.S3Account) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("ID is not set")
		}

		opts := clients.ConfigOpts{}
		config, err := opts.LoadAndValidate()
		if err != nil {
			return fmt.Errorf("Error authenticating clients from environment: %s", err)
		}

		client, err := config.IAMServiceUsersV1Client(config.GetRegion())
		if err != nil {
			return fmt.Errorf("Error creating VKCS IAM Service Users client: %s", err)
		}

		found, err := s3accounts.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		*s3Account = *found

		return nil
	}
}

const testAccIAMS3AccountBasic = `
resource "vkcs_iam_s3_account" "basic" {
  name        = "tfacc-s3-account-basic"
}

`

const testAccIAMS3AccountFull = `
resource "vkcs_iam_s3_account" "full" {
  name        = "tfacc-s3-account-full"
  description = "S3 Account created by acceptance test"
}
`
