package kms_test

import (
	"fmt"
	"testing"

	acctest_helper "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccKMSKeysDataSource_basic(t *testing.T) {
	name := "tfacc-kms-key-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		CheckDestroy:             testAccKMSCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccKMSKeysDataSourceBasic, name, "aes256-gcm96", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vkcs_kms_keys.keys", "id"),
					resource.TestCheckResourceAttrSet("data.vkcs_kms_keys.keys", "keys_count"),
					resource.TestCheckTypeSetElemAttr("data.vkcs_kms_keys.keys", "keys.*", name),
				),
			},
		},
	})
}

const testAccKMSKeysDataSourceBasic = testAccKMSKeyResourceConfig + `
data "vkcs_kms_keys" "keys" {
  depends_on = [vkcs_kms_key.key]
}
`
