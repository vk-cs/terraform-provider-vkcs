package kms_test

import (
	"fmt"
	"testing"

	acctest_helper "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccKMSSecretsDataSource_basic(t *testing.T) {
	name := "tfacc-kms-secret-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		CheckDestroy:             testAccKMSCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccKMSSecretsDataSourceBasic, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vkcs_kms_secrets.secrets", "id"),
					resource.TestCheckResourceAttrSet("data.vkcs_kms_secrets.secrets", "secrets_count"),
					resource.TestCheckTypeSetElemAttr("data.vkcs_kms_secrets.secrets", "secrets.*", name),
				),
			},
		},
	})
}

const testAccKMSSecretsDataSourceBasic = `
resource "vkcs_kms_secret" "secret" {
  path              = %q
  delete_protection = false
  data_json         = jsonencode({ password = "test-password" })
}

data "vkcs_kms_secrets" "secrets" {
  depends_on = [vkcs_kms_secret.secret]
}
`
