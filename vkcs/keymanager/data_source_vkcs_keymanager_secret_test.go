package keymanager_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccKeyManagerSecretDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerSecretV1DataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.vkcs_keymanager_secret.secret_1", "id",
						"vkcs_keymanager_secret.secret_1", "id"),
					resource.TestCheckResourceAttrPair(
						"data.vkcs_keymanager_secret.secret_2", "id",
						"vkcs_keymanager_secret.secret_2", "id"),
					resource.TestCheckResourceAttrPair(
						"data.vkcs_keymanager_secret.secret_1", "payload",
						"vkcs_keymanager_secret.secret_1", "payload"),
					resource.TestCheckResourceAttrPair(
						"data.vkcs_keymanager_secret.secret_2", "payload",
						"vkcs_keymanager_secret.secret_2", "payload"),
					resource.TestCheckResourceAttr(
						"data.vkcs_keymanager_secret.secret_1", "metadata.foo", "update"),
				),
			},
		},
	})
}

func TestAccKeyManagerSecretDataSource_acls(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerSecretV1DataSourceAcls,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.vkcs_keymanager_secret.secret_1", "id",
						"vkcs_keymanager_secret.secret_1", "id"),
					resource.TestCheckResourceAttrPair(
						"data.vkcs_keymanager_secret.secret_2", "id",
						"vkcs_keymanager_secret.secret_2", "id"),
					resource.TestCheckResourceAttrPair(
						"data.vkcs_keymanager_secret.secret_1", "payload",
						"vkcs_keymanager_secret.secret_1", "payload"),
					resource.TestCheckResourceAttrPair(
						"data.vkcs_keymanager_secret.secret_2", "payload",
						"vkcs_keymanager_secret.secret_2", "payload"),
					resource.TestCheckResourceAttr(
						"data.vkcs_keymanager_secret.secret_1", "metadata.foo", "update"),
					resource.TestCheckResourceAttr("data.vkcs_keymanager_secret.secret_1", "acl.0.read.0.project_access", "false"),
					resource.TestCheckResourceAttr("data.vkcs_keymanager_secret.secret_1", "acl.0.read.0.users.#", "1"),
					resource.TestCheckResourceAttr("data.vkcs_keymanager_secret.secret_2", "acl.0.read.0.project_access", "true"),
					resource.TestCheckResourceAttr("data.vkcs_keymanager_secret.secret_2", "acl.0.read.0.users.#", "0"),
				),
			},
		},
	})
}

func TestAccKeyManagerSecretDataSource_migrateToFramework(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.AccTestPreCheck(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"vkcs": {
						VersionConstraint: "0.3.0",
						Source:            "vk-cs/vkcs",
					},
				},
				Config: testAccKeyManagerSecretV1DataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.vkcs_keymanager_secret.secret_1", "id",
						"vkcs_keymanager_secret.secret_1", "id"),
					resource.TestCheckResourceAttrPair(
						"data.vkcs_keymanager_secret.secret_2", "id",
						"vkcs_keymanager_secret.secret_2", "id"),
					resource.TestCheckResourceAttrPair(
						"data.vkcs_keymanager_secret.secret_1", "payload",
						"vkcs_keymanager_secret.secret_1", "payload"),
					resource.TestCheckResourceAttrPair(
						"data.vkcs_keymanager_secret.secret_2", "payload",
						"vkcs_keymanager_secret.secret_2", "payload"),
					resource.TestCheckResourceAttr(
						"data.vkcs_keymanager_secret.secret_1", "metadata.foo", "update"),
				),
			},
			{
				ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
				Config:                   testAccKeyManagerSecretV1DataSourceBasic,
				PlanOnly:                 true,
			},
		},
	})
}

const testAccKeyManagerSecretV1DataSourceBasic = `
resource "vkcs_keymanager_secret" "secret_1" {
  algorithm   = "aes"
  bit_length  = 192
  mode        = "cbc"
  name        = "mysecret"
  payload     = "foobar"
  secret_type = "passphrase"
  payload_content_type = "text/plain"
  metadata = {
    foo = "update"
  }
}

resource "vkcs_keymanager_secret" "secret_2" {
  algorithm   = "aes"
  bit_length  = 256
  mode        = "cbc"
  name        = "mysecret"
  secret_type = "passphrase"
  payload     = "foo"
  expiration  = "3000-07-31T12:02:46Z"
  payload_content_type = "text/plain"
}

data "vkcs_keymanager_secret" "secret_1" {
  bit_length  = vkcs_keymanager_secret.secret_1.bit_length
  secret_type = "passphrase"
}

data "vkcs_keymanager_secret" "secret_2" {
  mode              = "cbc"
  secret_type       = "passphrase"
  expiration_filter = vkcs_keymanager_secret.secret_2.expiration
}
`

const testAccKeyManagerSecretV1DataSourceAcls = `
resource "vkcs_keymanager_secret" "secret_1" {
  algorithm   = "aes"
  bit_length  = 192
  mode        = "cbc"
  name        = "mysecret"
  payload     = "foobar"
  secret_type = "passphrase"
  payload_content_type = "text/plain"
  metadata = {
    foo = "update"
  }
  acl {
    read {
      project_access = false
      users = [
        "96b3ebddf275996285eae440e71227ba47c651be18391b0f2ebf1032ebae5dca",
      ]
    }
  }
}

resource "vkcs_keymanager_secret" "secret_2" {
  algorithm   = "aes"
  bit_length  = 256
  mode        = "cbc"
  name        = "mysecret"
  secret_type = "passphrase"
  payload     = "foo"
  expiration  = "3000-07-31T12:02:46Z"
  payload_content_type = "text/plain"
  acl {
    read {
      project_access = true
    }
  }
}

data "vkcs_keymanager_secret" "secret_1" {
  bit_length  = vkcs_keymanager_secret.secret_1.bit_length
  secret_type = "passphrase"
}

data "vkcs_keymanager_secret" "secret_2" {
  mode              = "cbc"
  secret_type       = "passphrase"
  expiration_filter = vkcs_keymanager_secret.secret_2.expiration
}
`
