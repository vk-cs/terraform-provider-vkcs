package kms_test

import (
	"fmt"
	"os"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kms/v1/keys"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kms/v1/secrets"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/kms"
)

func testAccKMSPreCheck(t *testing.T) {
	t.Helper()
	if os.Getenv("OS_PROJECT_ID") == "" {
		t.Fatal("OS_PROJECT_ID must be set for KMS acceptance tests")
	}
	opts := clients.ConfigOpts{}
	if _, err := opts.LoadAndValidate(); err != nil {
		t.Fatalf("Invalid KMS acceptance credentials: %v", err)
	}
}

// Verify deletion through the API, including resources created by data source tests.
func testAccKMSCheckDestroy(state *terraform.State) error {
	opts := clients.ConfigOpts{}
	config, err := opts.LoadAndValidate()
	if err != nil {
		return err
	}
	client, err := config.KMSV1Client(config.GetRegion())
	if err != nil {
		return err
	}
	for _, rs := range state.RootModule().Resources {
		switch rs.Type {
		case "vkcs_kms_key":
			_, err = keys.Get(client, rs.Primary.ID).Extract()
		case "vkcs_kms_secret":
			_, err = secrets.Get(client, rs.Primary.ID).Extract()
		default:
			continue
		}
		if err == nil {
			return fmt.Errorf("%s %q still exists", rs.Type, rs.Primary.ID)
		}
		if !errutil.IsNotFound(err) {
			return err
		}
	}
	return nil
}

func testAccKMSKeyConfig(name, keyType string, deletionAllowed bool) string {
	return fmt.Sprintf(`
resource "vkcs_kms_key" "key" {
  name = %q
  type = %q
  deletion_allowed = %t
}
`, name, keyType, deletionAllowed)
}

func TestAccKMSKeyResource_basic(t *testing.T) {
	name := sdkacctest.RandomWithPrefix("tf-acc-kms-key")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccKMSPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		CheckDestroy:             testAccKMSCheckDestroy,
		Steps: []resource.TestStep{
			{Config: fmt.Sprintf(`resource "vkcs_kms_key" "key" { name = %q }`, name), Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("vkcs_kms_key.key", "id", name),
				resource.TestCheckResourceAttr("vkcs_kms_key.key", "type", "aes256-gcm96"),
				resource.TestCheckResourceAttr("vkcs_kms_key.key", "deletion_allowed", "false"),
			)},
			{Config: testAccKMSKeyConfig(name, "aes256-gcm96", true), Check: resource.TestCheckResourceAttr("vkcs_kms_key.key", "deletion_allowed", "true")},
			{Config: testAccKMSKeyConfig(name+"-replacement", "aes128-gcm96", true), Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("vkcs_kms_key.key", "id", name+"-replacement"),
				resource.TestCheckResourceAttr("vkcs_kms_key.key", "type", "aes128-gcm96"),
			)},
		},
	})
}

func TestAccKMSKeyResource_keyTypes(t *testing.T) {
	for _, keyType := range []string{"aes128-gcm96", "aes256-gcm96", "chacha20-poly1305", "xchacha20-poly1305"} {
		t.Run(keyType, func(t *testing.T) {
			name := sdkacctest.RandomWithPrefix("tf-acc-kms-key")
			resource.Test(t, resource.TestCase{
				PreCheck:                 func() { testAccKMSPreCheck(t) },
				ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
				CheckDestroy:             testAccKMSCheckDestroy,
				Steps:                    []resource.TestStep{{Config: testAccKMSKeyConfig(name, keyType, true), Check: resource.TestCheckResourceAttr("vkcs_kms_key.key", "type", keyType)}},
			})
		})
	}
}

func TestAccKMSKeyDataSource_basic(t *testing.T) {
	name := sdkacctest.RandomWithPrefix("tf-acc-kms-key")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccKMSPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		CheckDestroy:             testAccKMSCheckDestroy,
		Steps: []resource.TestStep{{Config: testAccKMSKeyConfig(name, "aes256-gcm96", true) + `
data "vkcs_kms_key" "key" { name = vkcs_kms_key.key.name }
`, Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("data.vkcs_kms_key.key", "id", name),
			resource.TestCheckResourceAttrPair("data.vkcs_kms_key.key", "type", "vkcs_kms_key.key", "type"),
			resource.TestCheckResourceAttrPair("data.vkcs_kms_key.key", "deletion_allowed", "vkcs_kms_key.key", "deletion_allowed"),
		)}},
	})
}

func TestAccKMSKeysListDataSource_basic(t *testing.T) {
	name := sdkacctest.RandomWithPrefix("tf-acc-kms-key")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccKMSPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		CheckDestroy:             testAccKMSCheckDestroy,
		Steps: []resource.TestStep{{Config: testAccKMSKeyConfig(name, "aes256-gcm96", true) + `
data "vkcs_kms_keys_list" "keys" { depends_on = [vkcs_kms_key.key] }
`, Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet("data.vkcs_kms_keys_list.keys", "id"),
			resource.TestCheckResourceAttrSet("data.vkcs_kms_keys_list.keys", "keys_count"),
			resource.TestCheckTypeSetElemAttr("data.vkcs_kms_keys_list.keys", "keys.*", name),
		)}},
	})
}

// Both cryptographic data sources are exercised together to verify a real round trip.
func TestAccKMSKeyEncryptDecryptDataSources_basic(t *testing.T) {
	name := sdkacctest.RandomWithPrefix("tf-acc-kms-key")
	config := testAccKMSKeyConfig(name, "aes256-gcm96", true) + `
data "vkcs_kms_key_encrypt" "encrypted" {
  key = vkcs_kms_key.key.name
  plaintext = base64encode("acceptance-test plaintext")
}
data "vkcs_kms_key_decrypt" "decrypted" {
  key = vkcs_kms_key.key.name
  ciphertext = data.vkcs_kms_key_encrypt.encrypted.ciphertext
}
`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccKMSPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		CheckDestroy:             testAccKMSCheckDestroy,
		Steps: []resource.TestStep{{Config: config, Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet("data.vkcs_kms_key_encrypt.encrypted", "id"),
			resource.TestCheckResourceAttrSet("data.vkcs_kms_key_encrypt.encrypted", "ciphertext"),
			resource.TestCheckResourceAttrSet("data.vkcs_kms_key_decrypt.decrypted", "id"),
			resource.TestCheckResourceAttrPair("data.vkcs_kms_key_decrypt.decrypted", "plaintext", "data.vkcs_kms_key_encrypt.encrypted", "plaintext"),
		)}},
	})
}

func testAccKMSSecretConfig(path, password string, protected bool) string {
	return fmt.Sprintf(`
resource "vkcs_kms_secret" "secret" {
  path = %q
  data_json = jsonencode({ password = %q })
  delete_protection = %t
}
`, path, password, protected)
}

func TestAccKMSSecretResource_basic(t *testing.T) {
	name := sdkacctest.RandomWithPrefix("tf-acc-kms-secret")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccKMSPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		CheckDestroy:             testAccKMSCheckDestroy,
		Steps: []resource.TestStep{
			{Config: fmt.Sprintf(`
resource "vkcs_kms_secret" "secret" {
  path = %q
  data_json = jsonencode({ password = "initial" })
}
`, name), Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("vkcs_kms_secret.secret", "id", name),
				resource.TestCheckResourceAttr("vkcs_kms_secret.secret", "data_json", `{"password":"initial"}`),
				resource.TestCheckResourceAttr("vkcs_kms_secret.secret", "delete_protection", "true"),
				resource.TestCheckResourceAttr("vkcs_kms_secret.secret", "version", "1"),
				resource.TestCheckResourceAttrSet("vkcs_kms_secret.secret", "created_time"),
			)},
			{Config: testAccKMSSecretConfig(name, "updated", false), Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("vkcs_kms_secret.secret", "id", name),
				resource.TestCheckResourceAttr("vkcs_kms_secret.secret", "data_json", `{"password":"updated"}`),
				resource.TestCheckResourceAttr("vkcs_kms_secret.secret", "delete_protection", "false"),
				resource.TestCheckResourceAttr("vkcs_kms_secret.secret", "version", "2"),
			)},
			{Config: testAccKMSSecretConfig(name+"-replacement", "replacement", false), Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("vkcs_kms_secret.secret", "id", name+"-replacement"),
				resource.TestCheckResourceAttr("vkcs_kms_secret.secret", "version", "1"),
			)},
		},
	})
}

func TestAccKMSSecretDataSource_basic(t *testing.T) {
	name := sdkacctest.RandomWithPrefix("tf-acc-kms-secret")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccKMSPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		CheckDestroy:             testAccKMSCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccKMSSecretDataSourceBasic, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vkcs_kms_secret.secret", kms.EntityFieldID, name),
					resource.TestCheckResourceAttr("data.vkcs_kms_secret.secret", kms.SecretFieldPath, name),
					resource.TestCheckResourceAttrPair("data.vkcs_kms_secret.secret", kms.SecretFieldDataJSON, "vkcs_kms_secret.secret", kms.SecretFieldDataJSON),
					resource.TestCheckResourceAttrSet("data.vkcs_kms_secret.secret", kms.SecretFieldCreatedTime),
					resource.TestCheckResourceAttrSet("data.vkcs_kms_secret.secret", kms.SecretFieldVersion),
				),
			},
		},
	})
}

const testAccKMSSecretDataSourceBasic = `
resource "vkcs_kms_secret" "secret" {
  path = %q
  delete_protection = false
  data_json = jsonencode({ password = "test-password" })
}

data "vkcs_kms_secret" "secret" {
  path = vkcs_kms_secret.secret.path
}
`

func TestAccKMSSecretsListDataSource_basic(t *testing.T) {
	name := sdkacctest.RandomWithPrefix("tf-acc-kms-secret")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccKMSPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		CheckDestroy:             testAccKMSCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccKMSSecretsListDataSourceBasic, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vkcs_kms_secrets_list.secrets", kms.EntityFieldID),
					resource.TestCheckResourceAttrSet("data.vkcs_kms_secrets_list.secrets", kms.SecretsListFieldSecretsCount),
					resource.TestCheckTypeSetElemAttr("data.vkcs_kms_secrets_list.secrets", kms.SecretsListFieldSecrets+".*", name),
				),
			},
		},
	})
}

const testAccKMSSecretsListDataSourceBasic = `
resource "vkcs_kms_secret" "secret" {
  path = %q
  delete_protection = false
  data_json = jsonencode({ password = "test-password" })
}

data "vkcs_kms_secrets_list" "secrets" {
  depends_on = [vkcs_kms_secret.secret]
}
`
