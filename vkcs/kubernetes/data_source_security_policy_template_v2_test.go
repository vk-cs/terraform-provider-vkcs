package kubernetes_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccKubernetesSecPolicyTemplateV2Datasource_byID(t *testing.T) {
	policyName := "k8scontainerlimits"
	policyVersion := "1.0.0"
	testConfig := acctest.AccTestRenderConfig(testAccKubernetesSecPolicyTemplateV2ConfigByID, map[string]string{
		"DataSourceKubernetesSecPolicyTemplatesV2Config": testAccKubernetesSecPolicyTemplatesV2Config,
		"TargetPolicyName":    policyName,
		"TargetPolicyVersion": policyVersion,
	})
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vkcs_kubernetes_security_policy_template_v2.by_id", "name", policyName),
					resource.TestCheckResourceAttr("data.vkcs_kubernetes_security_policy_template_v2.by_id", "version", policyVersion),
					resource.TestCheckResourceAttrSet("data.vkcs_kubernetes_security_policy_template_v2.by_id", "id"),
					resource.TestCheckResourceAttrSet("data.vkcs_kubernetes_security_policy_template_v2.by_id", "description"),
					resource.TestCheckResourceAttrSet("data.vkcs_kubernetes_security_policy_template_v2.by_id", "settings_description"),
				),
			},
		},
	})
}

func TestAccKubernetesSecPolicyTemplateV2Datasource_byNameAndVersion(t *testing.T) {
	policyName := "k8scontainerlimits"
	policyVersion := "1.0.0"
	testConfig := acctest.AccTestRenderConfig(testAccKubernetesSecPolicyTemplateV2ConfigByNameAndVersion, map[string]string{
		"TargetPolicyName":    policyName,
		"TargetPolicyVersion": policyVersion,
	})
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vkcs_kubernetes_security_policy_template_v2.by_name_and_version", "name", policyName),
					resource.TestCheckResourceAttr("data.vkcs_kubernetes_security_policy_template_v2.by_name_and_version", "version", policyVersion),
					resource.TestCheckResourceAttrSet("data.vkcs_kubernetes_security_policy_template_v2.by_name_and_version", "id"),
					resource.TestCheckResourceAttrSet("data.vkcs_kubernetes_security_policy_template_v2.by_name_and_version", "description"),
					resource.TestCheckResourceAttrSet("data.vkcs_kubernetes_security_policy_template_v2.by_name_and_version", "settings_description"),
				),
			},
		},
	})
}

const testAccKubernetesSecPolicyTemplateV2ConfigByID = `
{{ .DataSourceKubernetesSecPolicyTemplatesV2Config }}

locals {
  target_policy_name    = "{{ .TargetPolicyName }}"
  target_policy_version = "{{ .TargetPolicyVersion }}"
  
  target_policy_id = [
    for template in data.vkcs_kubernetes_security_policy_templates_v2.basic.security_policies : 
    template.id 
    if template.name == local.target_policy_name && template.version == local.target_policy_version
  ][0]
}

data "vkcs_kubernetes_security_policy_template_v2" "by_id" {
  id = local.target_policy_id
}
`

const testAccKubernetesSecPolicyTemplateV2ConfigByNameAndVersion = `
data "vkcs_kubernetes_security_policy_template_v2" "by_name_and_version" {
  name    = "{{ .TargetPolicyName }}"
  version = "{{ .TargetPolicyVersion }}"
}
`
