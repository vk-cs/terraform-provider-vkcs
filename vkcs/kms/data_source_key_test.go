package kms_test

import (
	"fmt"
	"testing"

	acctest_helper "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccKMSKeyDataSource_basic(t *testing.T) {
	name := "tfacc-kms-key-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		CheckDestroy:             testAccKMSCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccKMSKeyDataSourceBasic, name, "aes256-gcm96", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vkcs_kms_key.key", "id", name),
					resource.TestCheckResourceAttrPair("data.vkcs_kms_key.key", "type", "vkcs_kms_key.key", "type"),
					resource.TestCheckResourceAttrPair("data.vkcs_kms_key.key", "deletion_allowed", "vkcs_kms_key.key", "deletion_allowed"),
				),
			},
		},
	})
}

const testAccKMSKeyDataSourceBasic = testAccKMSKeyResourceConfig + `
data "vkcs_kms_key" "key" {
  name = vkcs_kms_key.key.name
}
`
