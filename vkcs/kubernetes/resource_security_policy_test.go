package kubernetes_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	acctest_helper "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/containerinfra/v1/securitypolicies"
)

func TestAccKubernetesSecurityPolicy_basic_big(t *testing.T) {
	var policy securitypolicies.SecurityPolicy
	clusterName := "tfacc-basic-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesSecurityPolicyBasic, map[string]string{"TestAccKubernetesClusterBasic": acctest.AccTestRenderConfig(testAccKubernetesClusterBasic, map[string]string{"TestAccKubernetesNetworkingBase": testAccKubernetesNetworkingBase, "TestAccKubernetesClusterBase": testAccKubernetesClusterBase, "ClusterName": clusterName})}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesSecurityPolicyExists("vkcs_kubernetes_security_policy.basic", &policy),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_security_policy.basic", "cluster_id", "vkcs_kubernetes_cluster.basic", "id"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_security_policy.basic", "enabled", "true"),
				),
			},
			acctest.ImportStep("vkcs_kubernetes_security_policy.basic"),
		},
	})
}

func testAccCheckKubernetesSecurityPolicyExists(n string, policy *securitypolicies.SecurityPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("kubernetes security policy not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("id is not set")
		}

		config, err := clients.ConfigureFromEnv(context.Background())
		if err != nil {
			return fmt.Errorf("Error authenticating clients from environment: %s", err)
		}

		client, err := config.ContainerInfraV1Client(acctest.OsRegionName)
		if err != nil {
			return fmt.Errorf("error creating Kubernetes API client: %s", err)
		}

		found, err := securitypolicies.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found == nil {
			return errors.New("kubernetes security policy not found")
		}

		*policy = *found
		return nil
	}
}

const testAccKubernetesSecurityPolicyBasic = `
{{ .TestAccKubernetesClusterBasic }}
locals {
	policy_settings = {
	  "ranges" = [
		{
		  "min_replicas" = 1
		  "max_replicas" = 2
		}
	  ]
	}
  }

resource "vkcs_kubernetes_security_policy" "basic" {
    cluster_id = vkcs_kubernetes_cluster.basic.id
    enabled = true
    namespace = "*"
    policy_settings = jsonencode(local.policy_settings)
    security_policy_template_id = "bd90b09a-c814-4d3a-bbe7-b323e4252fa7"
}
`
