package keymanager_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccKeyManagerSecret_importBasic(t *testing.T) {
	resourceName := "vkcs_keymanager_secret.secret_1"

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerSecretBasic,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccKeyManagerSecret_importACLs(t *testing.T) {
	resourceName := "vkcs_keymanager_secret.secret_1"

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerSecretAcls,
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"payload_content_encoding"},
			},
		},
	})
}
