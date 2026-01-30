package kubernetes_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccKubernetesSecPolicyTemplatesV2Datasource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesSecPolicyTemplatesV2Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.vkcs_kubernetes_security_policy_templates_v2.base",
						"security_policies.#",
						regexp.MustCompile(`[1-9]\d*`),
					),
					resource.TestCheckResourceAttrSet(
						"data.vkcs_kubernetes_security_policy_templates_v2.base",
						"region",
					),
					resource.TestCheckResourceAttrSet(
						"data.vkcs_kubernetes_security_policy_templates_v2.base",
						"id",
					),
					resource.TestCheckResourceAttr(
						"data.vkcs_kubernetes_security_policy_templates_v2.base",
						"id",
						"policy_templates",
					),
				),
			},
		},
	})
}

const testAccKubernetesSecPolicyTemplatesV2Config = `
data "vkcs_kubernetes_security_policy_templates_v2" "base" {}
`
