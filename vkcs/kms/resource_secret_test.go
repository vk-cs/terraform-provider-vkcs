package kms_test

import (
	"fmt"
	"testing"

	acctest_helper "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccKMSSecretResource_basic(t *testing.T) {
	name := "tfacc-kms-secret-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		CheckDestroy:             testAccKMSCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccKMSSecretResourceBasic, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_kms_secret.secret", "id", name),
					resource.TestCheckResourceAttr("vkcs_kms_secret.secret", "data_json", `{"password":"initial"}`),
					resource.TestCheckResourceAttr("vkcs_kms_secret.secret", "delete_protection", "true"),
					resource.TestCheckResourceAttr("vkcs_kms_secret.secret", "version", "1"),
					resource.TestCheckResourceAttrSet("vkcs_kms_secret.secret", "created_time"),
				),
			},
			{
				Config: fmt.Sprintf(testAccKMSSecretResourceConfig, name, "updated", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_kms_secret.secret", "id", name),
					resource.TestCheckResourceAttr("vkcs_kms_secret.secret", "data_json", `{"password":"updated"}`),
					resource.TestCheckResourceAttr("vkcs_kms_secret.secret", "delete_protection", "false"),
					resource.TestCheckResourceAttr("vkcs_kms_secret.secret", "version", "2"),
				),
			},
			{
				Config: fmt.Sprintf(testAccKMSSecretResourceConfig, name+"-replacement", "replacement", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_kms_secret.secret", "id", name+"-replacement"),
					resource.TestCheckResourceAttr("vkcs_kms_secret.secret", "version", "1"),
				),
			},
		},
	})
}

const testAccKMSSecretResourceBasic = `
resource "vkcs_kms_secret" "secret" {
  path      = %q
  data_json = jsonencode({ password = "initial" })
}
`

const testAccKMSSecretResourceConfig = `
resource "vkcs_kms_secret" "secret" {
  path              = %q
  data_json         = jsonencode({ password = %q })
  delete_protection = %t
}
`
