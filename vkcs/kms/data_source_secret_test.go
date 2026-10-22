package kms_test

import (
	"fmt"
	"testing"

	acctest_helper "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccKMSSecretDataSource_basic(t *testing.T) {
	name := "tfacc-kms-secret-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		CheckDestroy:             testAccKMSCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccKMSSecretDataSourceBasic, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vkcs_kms_secret.secret", "id", name),
					resource.TestCheckResourceAttr("data.vkcs_kms_secret.secret", "path", name),
					resource.TestCheckResourceAttrPair("data.vkcs_kms_secret.secret", "data_json", "vkcs_kms_secret.secret", "data_json"),
					resource.TestCheckResourceAttrSet("data.vkcs_kms_secret.secret", "created_time"),
					resource.TestCheckResourceAttrSet("data.vkcs_kms_secret.secret", "version"),
				),
			},
		},
	})
}

const testAccKMSSecretDataSourceBasic = `
resource "vkcs_kms_secret" "secret" {
  path              = %q
  delete_protection = false
  data_json         = jsonencode({ password = "test-password" })
}

data "vkcs_kms_secret" "secret" {
  path = vkcs_kms_secret.secret.path
}
`
