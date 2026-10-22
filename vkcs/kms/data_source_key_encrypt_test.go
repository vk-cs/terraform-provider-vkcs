package kms_test

import (
	"fmt"
	"testing"

	acctest_helper "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccKMSKeyEncryptDecryptDataSources_basic(t *testing.T) {
	name := "tfacc-kms-key-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)
	config := fmt.Sprintf(testAccKMSKeyEncryptDecryptDataSourcesBasic, name, "aes256-gcm96", true)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		CheckDestroy:             testAccKMSCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vkcs_kms_key_encrypt.encrypted", "id"),
					resource.TestCheckResourceAttrSet("data.vkcs_kms_key_encrypt.encrypted", "ciphertext"),
					resource.TestCheckResourceAttrSet("data.vkcs_kms_key_decrypt.decrypted", "id"),
					resource.TestCheckResourceAttrPair("data.vkcs_kms_key_decrypt.decrypted", "plaintext", "data.vkcs_kms_key_encrypt.encrypted", "plaintext"),
				),
			},
		},
	})
}

const testAccKMSKeyEncryptDecryptDataSourcesBasic = testAccKMSKeyResourceConfig + `
data "vkcs_kms_key_encrypt" "encrypted" {
  key       = vkcs_kms_key.key.name
  plaintext = base64encode("acceptance-test plaintext")
}

data "vkcs_kms_key_decrypt" "decrypted" {
  key        = vkcs_kms_key.key.name
  ciphertext = data.vkcs_kms_key_encrypt.encrypted.ciphertext
}
`
