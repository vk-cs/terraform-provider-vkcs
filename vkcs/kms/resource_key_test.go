package kms_test

import (
	"fmt"
	"testing"

	acctest_helper "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccKMSKeyResource_basic(t *testing.T) {
	name := "tfacc-kms-key-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		CheckDestroy:             testAccKMSCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccKMSKeyResourceBasic, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_kms_key.key", "id", name),
					resource.TestCheckResourceAttr("vkcs_kms_key.key", "type", "aes256-gcm96"),
					resource.TestCheckResourceAttr("vkcs_kms_key.key", "deletion_allowed", "false"),
				),
			},
			{
				Config: fmt.Sprintf(testAccKMSKeyResourceConfig, name, "aes256-gcm96", true),
				Check:  resource.TestCheckResourceAttr("vkcs_kms_key.key", "deletion_allowed", "true"),
			},
			{
				Config: fmt.Sprintf(testAccKMSKeyResourceConfig, name+"-replacement", "aes128-gcm96", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_kms_key.key", "id", name+"-replacement"),
					resource.TestCheckResourceAttr("vkcs_kms_key.key", "type", "aes128-gcm96"),
				),
			},
		},
	})
}

func TestAccKMSKeyResource_keyTypes(t *testing.T) {
	for _, keyType := range []string{"aes128-gcm96", "aes256-gcm96", "chacha20-poly1305", "xchacha20-poly1305"} {
		t.Run(keyType, func(t *testing.T) {
			name := "tfacc-kms-key-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)
			resource.Test(t, resource.TestCase{
				PreCheck:                 func() { acctest.AccTestPreCheck(t) },
				ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
				CheckDestroy:             testAccKMSCheckDestroy,
				Steps: []resource.TestStep{
					{
						Config: fmt.Sprintf(testAccKMSKeyResourceConfig, name, keyType, true),
						Check:  resource.TestCheckResourceAttr("vkcs_kms_key.key", "type", keyType),
					},
				},
			})
		})
	}
}

const testAccKMSKeyResourceBasic = `
resource "vkcs_kms_key" "key" {
  name = %q
}
`

const testAccKMSKeyResourceConfig = `
resource "vkcs_kms_key" "key" {
  name             = %q
  type             = %q
  deletion_allowed = %t
}
`
