package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/containers"
)

func TestAccKeyManagerContainer_basic(t *testing.T) {
	var container containers.Container
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerContainerBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerExists(
						"vkcs_keymanager_container.container_1", &container),
					resource.TestCheckResourceAttrPtr("vkcs_keymanager_container.container_1", "name", &container.Name),
					resource.TestCheckResourceAttrPtr("vkcs_keymanager_container.container_1", "type", &container.Type),
					resource.TestCheckResourceAttr("vkcs_keymanager_container.container_1", "secret_refs.#", "3"),
				),
			},
		},
	})
}

func TestAccKeyManagerContainer_acls(t *testing.T) {
	var container containers.Container
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerContainerAcls(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerExists(
						"vkcs_keymanager_container.container_1", &container),
					resource.TestCheckResourceAttrPtr("vkcs_keymanager_container.container_1", "name", &container.Name),
					resource.TestCheckResourceAttrPtr("vkcs_keymanager_container.container_1", "type", &container.Type),
					resource.TestCheckResourceAttr("vkcs_keymanager_container.container_1", "secret_refs.#", "3"),
					resource.TestCheckResourceAttr("vkcs_keymanager_container.container_1", "acl.0.read.0.project_access", "false"),
					resource.TestCheckResourceAttr("vkcs_keymanager_container.container_1", "acl.0.read.0.users.#", "2"),
				),
			},
		},
	})
}

func TestAccKeyManagerContainer_certificate_type(t *testing.T) {
	var container containers.Container
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerContainerCertificateType(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerExists(
						"vkcs_keymanager_container.container_1", &container),
					resource.TestCheckResourceAttrPtr("vkcs_keymanager_container.container_1", "name", &container.Name),
					resource.TestCheckResourceAttrPtr("vkcs_keymanager_container.container_1", "type", &container.Type),
					resource.TestCheckResourceAttr("vkcs_keymanager_container.container_1", "secret_refs.#", "3"),
				),
			},
		},
	})
}

func TestAccKeyManagerContainer_acls_update(t *testing.T) {
	var container containers.Container
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerContainerAcls(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerExists(
						"vkcs_keymanager_container.container_1", &container),
					resource.TestCheckResourceAttrPtr("vkcs_keymanager_container.container_1", "name", &container.Name),
					resource.TestCheckResourceAttrPtr("vkcs_keymanager_container.container_1", "type", &container.Type),
					resource.TestCheckResourceAttr("vkcs_keymanager_container.container_1", "secret_refs.#", "3"),
					resource.TestCheckResourceAttr("vkcs_keymanager_container.container_1", "acl.0.read.0.project_access", "false"),
					resource.TestCheckResourceAttr("vkcs_keymanager_container.container_1", "acl.0.read.0.users.#", "2"),
				),
			},
			{
				Config: testAccKeyManagerContainerAclsUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerExists(
						"vkcs_keymanager_container.container_1", &container),
					resource.TestCheckResourceAttrPtr("vkcs_keymanager_container.container_1", "name", &container.Name),
					resource.TestCheckResourceAttrPtr("vkcs_keymanager_container.container_1", "type", &container.Type),
					resource.TestCheckResourceAttr("vkcs_keymanager_container.container_1", "secret_refs.#", "3"),
					resource.TestCheckResourceAttr("vkcs_keymanager_container.container_1", "acl.0.read.0.project_access", "true"),
					resource.TestCheckResourceAttr("vkcs_keymanager_container.container_1", "acl.0.read.0.users.#", "0"),
				),
			},
		},
	})
}

func testAccCheckContainerDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(configer)
	kmClient, err := config.KeyManagerV1Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS KeyManager client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_keymanager_container" {
			continue
		}
		_, err = containers.Get(kmClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Container (%s) still exists", rs.Primary.ID)
		}
		if _, ok := err.(gophercloud.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckContainerExists(n string, container *containers.Container) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(configer)
		kmClient, err := config.KeyManagerV1Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating VKCS KeyManager client: %s", err)
		}

		var found *containers.Container

		found, err = containers.Get(kmClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*container = *found

		return nil
	}
}

const testAccKeyManagerContainer = `
resource "vkcs_keymanager_secret" "certificate_1" {
  name                 = "certificate"
  payload              = "certificate"
  secret_type          = "certificate"
  payload_content_type = "text/plain"
}

resource "vkcs_keymanager_secret" "private_key_1" {
  name                 = "private_key"
  payload              = "private_key"
  secret_type          = "private"
  payload_content_type = "text/plain"
}

resource "vkcs_keymanager_secret" "intermediate_1" {
  name                 = "intermediate"
  payload              = "intermediate"
  secret_type          = "certificate"
  payload_content_type = "text/plain"
}
`

func testAccKeyManagerContainerBasic() string {
	return fmt.Sprintf(`
%s

resource "vkcs_keymanager_container" "container_1" {
  name = "generic"
  type = "generic"

  secret_refs {
    name       = "certificate"
    secret_ref = vkcs_keymanager_secret.certificate_1.secret_ref
  }

  secret_refs {
    name       = "private_key"
    secret_ref = vkcs_keymanager_secret.private_key_1.secret_ref
  }

  secret_refs {
    name       = "intermediates"
    secret_ref = vkcs_keymanager_secret.intermediate_1.secret_ref
  }
}
`, testAccKeyManagerContainer)
}

func testAccKeyManagerContainerAcls() string {
	return fmt.Sprintf(`
%s

resource "vkcs_keymanager_container" "container_1" {
  name = "generic"
  type = "generic"

  secret_refs {
    name       = "certificate"
    secret_ref = vkcs_keymanager_secret.certificate_1.secret_ref
  }

  secret_refs {
    name       = "private_key"
    secret_ref = vkcs_keymanager_secret.private_key_1.secret_ref
  }

  secret_refs {
    name       = "intermediates"
    secret_ref = vkcs_keymanager_secret.intermediate_1.secret_ref
  }

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
`, testAccKeyManagerContainer)
}

func testAccKeyManagerContainerAclsUpdate() string {
	return fmt.Sprintf(`
%s

resource "vkcs_keymanager_container" "container_1" {
  name = "generic"
  type = "generic"

  secret_refs {
    name       = "certificate"
    secret_ref = vkcs_keymanager_secret.certificate_1.secret_ref
  }

  secret_refs {
    name       = "private_key"
    secret_ref = vkcs_keymanager_secret.private_key_1.secret_ref
  }

  secret_refs {
    name       = "intermediates"
    secret_ref = vkcs_keymanager_secret.intermediate_1.secret_ref
  }

  acl {
    read {}
  }
}
`, testAccKeyManagerContainer)
}

func testAccKeyManagerContainerCertificateType() string {
	return fmt.Sprintf(`
%s

resource "vkcs_keymanager_container" "container_1" {
  name = "generic"
  type = "certificate"

  secret_refs {
    name       = "certificate"
    secret_ref = vkcs_keymanager_secret.certificate_1.secret_ref
  }

  secret_refs {
    name       = "private_key"
    secret_ref = vkcs_keymanager_secret.private_key_1.secret_ref
  }

  secret_refs {
    name       = "intermediates"
    secret_ref = vkcs_keymanager_secret.intermediate_1.secret_ref
  }
}
`, testAccKeyManagerContainer)
}
