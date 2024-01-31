package kubernetes_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccKubernetesSecurityPolicyTemplateDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesSecurityPolicyTemplateDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vkcs_kubernetes_security_policy_template.template", "name", "k8sallowedrepos"),
					resource.TestCheckResourceAttr("data.vkcs_kubernetes_security_policy_template.template", "version", "1.0.0"),
				),
			},
		},
	})
}

const testAccKubernetesSecurityPolicyTemplateDataSourceBasic = `
data "vkcs_kubernetes_security_policy_template" "template" {
	name = "k8sallowedrepos"
}
`
