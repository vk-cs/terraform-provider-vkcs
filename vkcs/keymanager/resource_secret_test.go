package keymanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/keymanager"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/secrets"
)

func TestAccKeyManagerSecret_basic(t *testing.T) {
	var secret secrets.Secret
	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerSecretBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretExists(
						"vkcs_keymanager_secret.secret_1", &secret),
					resource.TestCheckResourceAttrPtr("vkcs_keymanager_secret.secret_1", "name", &secret.Name),
					resource.TestCheckResourceAttrPtr("vkcs_keymanager_secret.secret_1", "secret_type", &secret.SecretType),
					resource.TestCheckResourceAttr("vkcs_keymanager_secret.secret_1", "payload", "foobar"),
				),
			},
		},
	})
}

func TestAccKeyManagerSecret_basicWithMetadata(t *testing.T) {
	var secret secrets.Secret
	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerSecretBasicWithMetadata,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretExists(
						"vkcs_keymanager_secret.secret_1", &secret),
					resource.TestCheckResourceAttrPtr("vkcs_keymanager_secret.secret_1", "name", &secret.Name),
					resource.TestCheckResourceAttrPtr("vkcs_keymanager_secret.secret_1", "secret_type", &secret.SecretType),
				),
			},
		},
	})
}

func TestAccKeyManagerSecret_updateMetadata(t *testing.T) {
	var secret secrets.Secret
	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerSecretBasicWithMetadata,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretExists(
						"vkcs_keymanager_secret.secret_1", &secret),
					resource.TestCheckResourceAttrPtr("vkcs_keymanager_secret.secret_1", "name", &secret.Name),
					resource.TestCheckResourceAttrPtr("vkcs_keymanager_secret.secret_1", "secret_type", &secret.SecretType),
					testAccCheckMetadataEquals("foo", "bar", &secret),
				),
			},
			{
				Config: testAccKeyManagerSecretUpdateMetadata,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretExists(
						"vkcs_keymanager_secret.secret_1", &secret),
					resource.TestCheckResourceAttrPtr("vkcs_keymanager_secret.secret_1", "name", &secret.Name),
					resource.TestCheckResourceAttrPtr("vkcs_keymanager_secret.secret_1", "secret_type", &secret.SecretType),
					testAccCheckMetadataEquals("foo", "update", &secret),
				),
			},
		},
	})
}

func TestAccKeyManagerUpdateSecret_payload(t *testing.T) {
	var secret secrets.Secret
	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerSecretNoPayload,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretExists(
						"vkcs_keymanager_secret.secret_1", &secret),
					resource.TestCheckResourceAttrPtr("vkcs_keymanager_secret.secret_1", "name", &secret.Name),
					resource.TestCheckResourceAttrPtr("vkcs_keymanager_secret.secret_1", "secret_type", &secret.SecretType),
					testAccCheckPayloadEquals("", &secret),
				),
			},
			{
				Config: testAccKeyManagerSecretUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretExists(
						"vkcs_keymanager_secret.secret_1", &secret),
					resource.TestCheckResourceAttrPtr("vkcs_keymanager_secret.secret_1", "name", &secret.Name),
					resource.TestCheckResourceAttrPtr("vkcs_keymanager_secret.secret_1", "secret_type", &secret.SecretType),
					testAccCheckPayloadEquals("updatedfoobar", &secret),
				),
			},
			{
				Config: testAccKeyManagerSecretUpdateWhitespace,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretExists(
						"vkcs_keymanager_secret.secret_1", &secret),
					resource.TestCheckResourceAttrPtr("vkcs_keymanager_secret.secret_1", "name", &secret.Name),
					resource.TestCheckResourceAttrPtr("vkcs_keymanager_secret.secret_1", "secret_type", &secret.SecretType),
					testAccCheckPayloadEquals("updatedfoobar", &secret),
				),
			},
			{
				Config: testAccKeyManagerSecretUpdateBase64,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretExists(
						"vkcs_keymanager_secret.secret_1", &secret),
					resource.TestCheckResourceAttrPtr("vkcs_keymanager_secret.secret_1", "name", &secret.Name),
					resource.TestCheckResourceAttrPtr("vkcs_keymanager_secret.secret_1", "secret_type", &secret.SecretType),
					testAccCheckPayloadEquals("base64foobar ", &secret),
				),
			},
		},
	})
}

func TestAccKeyManagerSecret_acls(t *testing.T) {
	var secret secrets.Secret
	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerSecretAcls,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretExists(
						"vkcs_keymanager_secret.secret_1", &secret),
					resource.TestCheckResourceAttrPtr("vkcs_keymanager_secret.secret_1", "name", &secret.Name),
					resource.TestCheckResourceAttrPtr("vkcs_keymanager_secret.secret_1", "secret_type", &secret.SecretType),
					resource.TestCheckResourceAttr("vkcs_keymanager_secret.secret_1", "acl.0.read.0.project_access", "false"),
					resource.TestCheckResourceAttr("vkcs_keymanager_secret.secret_1", "acl.0.read.0.users.#", "2"),
					resource.TestCheckResourceAttr("vkcs_keymanager_secret.secret_2", "acl.0.read.0.project_access", "true"),
					resource.TestCheckResourceAttr("vkcs_keymanager_secret.secret_2", "acl.0.read.0.users.#", "0"),
				),
			},
		},
	})
}

func TestAccKeyManagerSecret_acls_update(t *testing.T) {
	var secret secrets.Secret
	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerSecretAcls,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretExists(
						"vkcs_keymanager_secret.secret_1", &secret),
					resource.TestCheckResourceAttrPtr("vkcs_keymanager_secret.secret_1", "name", &secret.Name),
					resource.TestCheckResourceAttrPtr("vkcs_keymanager_secret.secret_1", "secret_type", &secret.SecretType),
					resource.TestCheckResourceAttr("vkcs_keymanager_secret.secret_1", "acl.0.read.0.project_access", "false"),
					resource.TestCheckResourceAttr("vkcs_keymanager_secret.secret_1", "acl.0.read.0.users.#", "2"),
					resource.TestCheckResourceAttr("vkcs_keymanager_secret.secret_2", "acl.0.read.0.project_access", "true"),
					resource.TestCheckResourceAttr("vkcs_keymanager_secret.secret_2", "acl.0.read.0.users.#", "0"),
				),
			},
			{
				Config: testAccKeyManagerSecretAclsUpdate1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretExists(
						"vkcs_keymanager_secret.secret_1", &secret),
					resource.TestCheckResourceAttrPtr("vkcs_keymanager_secret.secret_1", "name", &secret.Name),
					resource.TestCheckResourceAttrPtr("vkcs_keymanager_secret.secret_1", "secret_type", &secret.SecretType),
					resource.TestCheckResourceAttr("vkcs_keymanager_secret.secret_1", "acl.0.read.0.project_access", "false"),
					resource.TestCheckResourceAttr("vkcs_keymanager_secret.secret_1", "acl.0.read.0.users.#", "2"),
					resource.TestCheckResourceAttr("vkcs_keymanager_secret.secret_2", "acl.0.read.0.project_access", "false"),
					resource.TestCheckResourceAttr("vkcs_keymanager_secret.secret_2", "acl.0.read.0.users.#", "1"),
				),
			},
			{
				Config: testAccKeyManagerSecretAclsUpdate2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretExists(
						"vkcs_keymanager_secret.secret_1", &secret),
					resource.TestCheckResourceAttrPtr("vkcs_keymanager_secret.secret_1", "name", &secret.Name),
					resource.TestCheckResourceAttrPtr("vkcs_keymanager_secret.secret_1", "secret_type", &secret.SecretType),
					resource.TestCheckResourceAttr("vkcs_keymanager_secret.secret_1", "acl.0.read.0.project_access", "true"),
					resource.TestCheckResourceAttr("vkcs_keymanager_secret.secret_1", "acl.0.read.0.users.#", "0"),
					resource.TestCheckResourceAttr("vkcs_keymanager_secret.secret_2", "acl.0.read.0.project_access", "true"),
					resource.TestCheckResourceAttr("vkcs_keymanager_secret.secret_2", "acl.0.read.0.users.#", "0"),
				),
			},
		},
	})
}

func testAccCheckSecretDestroy(s *terraform.State) error {
	config := acctest.AccTestProvider.Meta().(clients.Config)
	kmClient, err := config.KeyManagerV1Client(acctest.OsRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS KeyManager client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_keymanager_secret" {
			continue
		}
		_, err = secrets.Get(kmClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Secret (%s) still exists", rs.Primary.ID)
		}
		if _, ok := err.(gophercloud.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckSecretExists(n string, secret *secrets.Secret) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acctest.AccTestProvider.Meta().(clients.Config)
		kmClient, err := config.KeyManagerV1Client(acctest.OsRegionName)
		if err != nil {
			return fmt.Errorf("Error creating VKCS KeyManager client: %s", err)
		}

		var found *secrets.Secret

		found, err = secrets.Get(kmClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*secret = *found

		return nil
	}
}

func testAccCheckPayloadEquals(payload string, secret *secrets.Secret) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.AccTestProvider.Meta().(clients.Config)
		kmClient, err := config.KeyManagerV1Client(acctest.OsRegionName)
		if err != nil {
			return fmt.Errorf("Error creating VKCS KeyManager client: %s", err)
		}

		opts := secrets.GetPayloadOpts{
			PayloadContentType: "text/plain",
		}

		uuid := keymanager.GetUUIDFromSecretRef(secret.SecretRef)
		secretPayload, _ := secrets.GetPayload(kmClient, uuid, opts).Extract()
		if string(secretPayload) != payload {
			return fmt.Errorf("Payloads do not match. Expected %v but got %v", payload, secretPayload)
		}
		return nil
	}
}

func testAccCheckMetadataEquals(key string, value string, secret *secrets.Secret) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.AccTestProvider.Meta().(clients.Config)
		kmClient, err := config.KeyManagerV1Client(acctest.OsRegionName)
		if err != nil {
			return fmt.Errorf("Error creating VKCS networking client: %s", err)
		}

		uuid := keymanager.GetUUIDFromSecretRef(secret.SecretRef)
		metadatum, err := secrets.GetMetadatum(kmClient, uuid, key).Extract()
		if err != nil {
			return err
		}
		if metadatum.Value != value {
			return fmt.Errorf("Metadata does not match. Expected %v but got %v", metadatum, value)
		}

		return nil
	}
}

const testAccKeyManagerSecretBasic = `
resource "vkcs_keymanager_secret" "secret_1" {
  algorithm = "aes"
  bit_length = 256
  mode = "cbc"
  name = "mysecret"
  payload = "foobar"
  payload_content_type = "text/plain"
  secret_type = "passphrase"
}`

const testAccKeyManagerSecretBasicWithMetadata = `
resource "vkcs_keymanager_secret" "secret_1" {
  algorithm = "aes"
  bit_length = 256
  mode = "cbc"
  name = "mysecret"
  payload = "foobar"
  payload_content_type = "text/plain"
  secret_type = "passphrase"
  metadata = {
    foo = "bar"
  }
}`

const testAccKeyManagerSecretUpdateMetadata = `
resource "vkcs_keymanager_secret" "secret_1" {
  algorithm = "aes"
  bit_length = 256
  mode = "cbc"
  name = "mysecret"
  payload = "foobar"
  payload_content_type = "text/plain"
  secret_type = "passphrase"
  metadata = {
    foo = "update"
  }
}`

const testAccKeyManagerSecretNoPayload = `
resource "vkcs_keymanager_secret" "secret_1" {
  algorithm = "aes"
  bit_length = 256
  mode = "cbc"
  name = "mysecret"
  secret_type = "passphrase"
}`

const testAccKeyManagerSecretUpdate = `
resource "vkcs_keymanager_secret" "secret_1" {
  algorithm = "aes"
  bit_length = 256
  mode = "cbc"
  name = "mysecret"
  payload = "updatedfoobar"
  payload_content_type = "text/plain"
  secret_type = "passphrase"
}`

const testAccKeyManagerSecretUpdateWhitespace = `
resource "vkcs_keymanager_secret" "secret_1" {
  algorithm = "aes"
  bit_length = 256
  mode = "cbc"
  name = "mysecret"
  payload = <<EOF
updatedfoobar
EOF
  payload_content_type = "text/plain"
  secret_type = "passphrase"
}`

const testAccKeyManagerSecretUpdateBase64 = `
resource "vkcs_keymanager_secret" "secret_1" {
  algorithm = "aes"
  bit_length = 256
  mode = "cbc"
  name = "mysecret"
  payload = base64encode("base64foobar ")
  payload_content_type = "application/octet-stream"
  payload_content_encoding = "base64"
  secret_type = "passphrase"
}`

const testAccKeyManagerSecretAcls = `
resource "vkcs_keymanager_secret" "secret_1" {
  algorithm = "aes"
  bit_length = 256
  mode = "cbc"
  name = "mysecret"
  payload = base64encode("base64foobar ")
  payload_content_type = "application/octet-stream"
  payload_content_encoding = "base64"
  secret_type = "passphrase"

  acl {
    read {
      project_access = false
      users = [
        "619e2ad074321cf246b03a89e95afee95fb26bb0b2d1fc7ba3bd30fcca25588a",
        "96b3ebddf275996285eae440e71227ba47c651be18391b0f2ebf1032ebae5dca",
      ]
    }
  }
}

resource "vkcs_keymanager_secret" "secret_2" {
  algorithm = "aes"
  bit_length = 256
  mode = "cbc"
  name = "mysecret"
  payload = "foobar"
  payload_content_type = "text/plain"
  secret_type = "passphrase"
}
`

const testAccKeyManagerSecretAclsUpdate1 = `
resource "vkcs_keymanager_secret" "secret_1" {
  algorithm = "aes"
  bit_length = 256
  mode = "cbc"
  name = "mysecret"
  payload = base64encode("base64foobar ")
  payload_content_type = "application/octet-stream"
  payload_content_encoding = "base64"
  secret_type = "passphrase"

  acl {
    read {
      project_access = false
      users = [
        "96b3ebddf275996285eae440e71227ba47c651be18391b0f2ebf1032ebae5dca",
        "619e2ad074321cf246b03a89e95afee95fb26bb0b2d1fc7ba3bd30fcca25588a",
      ]
    }
  }
}

resource "vkcs_keymanager_secret" "secret_2" {
  algorithm = "aes"
  bit_length = 256
  mode = "cbc"
  name = "mysecret"
  payload = "foobar"
  payload_content_type = "text/plain"
  secret_type = "passphrase"

  acl {
    read {
      project_access = false
      users = [
        "96b3ebddf275996285eae440e71227ba47c651be18391b0f2ebf1032ebae5dca",
      ]
    }
  }
}
`

const testAccKeyManagerSecretAclsUpdate2 = `
resource "vkcs_keymanager_secret" "secret_1" {
  algorithm = "aes"
  bit_length = 256
  mode = "cbc"
  name = "mysecret"
  payload = base64encode("base64foobar ")
  payload_content_type = "application/octet-stream"
  payload_content_encoding = "base64"
  secret_type = "passphrase"

  acl {
    read {
      project_access = true
    }
  }
}

resource "vkcs_keymanager_secret" "secret_2" {
  algorithm = "aes"
  bit_length = 256
  mode = "cbc"
  name = "mysecret"
  payload = "foobar"
  payload_content_type = "text/plain"
  secret_type = "passphrase"

  acl {
    read {
      project_access = true
    }
  }
}
`
