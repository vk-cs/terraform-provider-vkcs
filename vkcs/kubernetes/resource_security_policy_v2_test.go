package kubernetes_test

import (
	"errors"
	"fmt"
	"testing"

	acctest_helper "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/secpolicies"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func TestAccKubernetesSecurityPolicyV2_basic(t *testing.T) {
	t.Parallel()

	clusterName := "tfacc-secpol-basic-v2-" + acctest_helper.RandStringFromCharSet(3, acctest_helper.CharSetAlphaNum)
	clusterConfig := acctest.AccTestRenderConfig(testAccKubernetesClusterV2Base, map[string]string{
		"TestAccKubernetesNetworkingBase": testAccKubernetesNetworkingBase,
		"ClusterName":                     clusterName,
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesSecurityPolicyV2Basic, map[string]string{
					"TestAccKubernetesClusterV2Base": clusterConfig,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesSecurityPolicyV2Exists("vkcs_kubernetes_cluster_v2.base"),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_cluster_v2.base", "cluster_id", "vkcs_kubernetes_cluster_v2.base", "id"),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_cluster_v2.base", "security_policy_template_id", "data.vkcs_kubernetes_security_policy_template_v2.policy_template", "id"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.base", "namespace", "*"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.base", "enabled", "true"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.base", "policy_settings", `{"cpu":"1000m","excludedNamespaces":["prometheus"],"memory":"2Gi"}`),

					resource.TestCheckResourceAttrSet("vkcs_kubernetes_cluster_v2.base", "id"),
					resource.TestCheckResourceAttrSet("vkcs_kubernetes_cluster_v2.base", "region"),
				),
			},
			acctest.ImportStep("vkcs_kubernetes_cluster_v2.base"),
		},
	})
}

func TestAccKubernetesSecurityPolicyV2_full(t *testing.T) {
	t.Parallel()

	clusterName := "tfacc-secpol-full-v2-" + acctest_helper.RandStringFromCharSet(3, acctest_helper.CharSetAlphaNum)
	clusterConfig := acctest.AccTestRenderConfig(testAccKubernetesClusterV2Base, map[string]string{
		"TestAccKubernetesNetworkingBase": testAccKubernetesNetworkingBase,
		"ClusterName":                     clusterName,
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesSecurityPolicyV2Full, map[string]string{
					"TestAccKubernetesClusterV2Base": clusterConfig,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesSecurityPolicyV2Exists("vkcs_kubernetes_security_policy_v2.full"),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_security_policy_v2.full", "cluster_id", "vkcs_kubernetes_cluster_v2.base", "id"),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_security_policy_v2.full", "security_policy_template_id", "data.vkcs_kubernetes_security_policy_template_v2.policy_template", "id"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_security_policy_v2.full", "namespace", "test-namespace"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_security_policy_v2.full", "enabled", "false"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_security_policy_v2.full", "policy_settings", `{"cpu":"2000m","excludedNamespaces":["monitoring","logging"],"memory":"4Gi"}`),

					resource.TestCheckResourceAttrSet("vkcs_kubernetes_security_policy_v2.full", "id"),
					resource.TestCheckResourceAttrSet("vkcs_kubernetes_security_policy_v2.full", "region"),
				),
			},
			acctest.ImportStep("vkcs_kubernetes_security_policy_v2.full"),
		},
	})
}

func TestAccKubernetesSecurityPolicyV2_update(t *testing.T) {
	t.Parallel()

	clusterName := "tfacc-secpol-upd-v2-" + acctest_helper.RandStringFromCharSet(3, acctest_helper.CharSetAlphaNum)
	clusterConfig := acctest.AccTestRenderConfig(testAccKubernetesClusterV2Base, map[string]string{
		"TestAccKubernetesNetworkingBase": testAccKubernetesNetworkingBase,
		"ClusterName":                     clusterName,
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesSecurityPolicyV2UpdateOld, map[string]string{
					"TestAccKubernetesClusterV2Base": clusterConfig,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesSecurityPolicyV2Exists("vkcs_kubernetes_security_policy_v2.update"),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_security_policy_v2.update", "cluster_id", "vkcs_kubernetes_cluster_v2.base", "id"),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_security_policy_v2.update", "security_policy_template_id", "data.vkcs_kubernetes_security_policy_template_v2.policy_template", "id"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_security_policy_v2.update", "namespace", "*"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_security_policy_v2.update", "enabled", "true"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_security_policy_v2.update", "policy_settings", `{"cpu":"1000m","excludedNamespaces":["prometheus"],"memory":"2Gi"}`),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesSecurityPolicyV2UpdateNew, map[string]string{
					"TestAccKubernetesClusterV2Base": clusterConfig,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesSecurityPolicyV2Exists("vkcs_kubernetes_security_policy_v2.update"),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_security_policy_v2.update", "cluster_id", "vkcs_kubernetes_cluster_v2.base", "id"),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_security_policy_v2.update", "security_policy_template_id", "data.vkcs_kubernetes_security_policy_template_v2.policy_template", "id"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_security_policy_v2.update", "namespace", "updated-namespace"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_security_policy_v2.update", "enabled", "false"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_security_policy_v2.update", "policy_settings", `{"cpu":"4000m","excludedNamespaces":["monitoring"],"memory":"8Gi"}`),
				),
			},
			acctest.ImportStep("vkcs_kubernetes_security_policy_v2.update"),
		},
	})
}

func testAccCheckKubernetesSecurityPolicyV2Exists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Kubernetes security policy not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("id is not set")
		}

		opts := clients.ConfigOpts{}
		config, err := opts.LoadAndValidate()
		if err != nil {
			return fmt.Errorf("Error authenticating clients from environment: %s", err)
		}

		client, err := config.ManagedK8SClient(acctest.OsRegionName)
		if err != nil {
			return fmt.Errorf("error creating Kubernetes API client: %s", err)
		}

		_, err = secpolicies.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			if errutil.IsNotFound(err) {
				return errors.New("Kubernetes security policy not found")
			}
			return err
		}

		return nil
	}
}

const testAccKubernetesSecurityPolicyV2Basic = `
{{ .TestAccKubernetesClusterV2Base }}

data "vkcs_kubernetes_security_policy_template_v2" "policy_template" {
  name    = "k8scontainerrequests"
  version = "1.0.0"
}

resource "vkcs_kubernetes_security_policy_v2" "base" {
  cluster_id                  = vkcs_kubernetes_cluster_v2.base.id
  security_policy_template_id = data.vkcs_kubernetes_security_policy_template_v2.policy_template.id
  namespace                   = "*"
  policy_settings             = jsonencode({
    cpu                = "1000m"
    memory             = "2Gi"
    excludedNamespaces = ["prometheus"]
  })
}
`

const testAccKubernetesSecurityPolicyV2Full = `
{{ .TestAccKubernetesClusterV2Base }}

data "vkcs_kubernetes_security_policy_template_v2" "policy_template" {
  name    = "k8scontainerrequests"
  version = "1.0.0"
}

resource "vkcs_kubernetes_security_policy_v2" "full" {
  cluster_id                  = vkcs_kubernetes_cluster_v2.base.id
  security_policy_template_id = data.vkcs_kubernetes_security_policy_template_v2.policy_template.id
  namespace                   = "test-namespace"
  enabled                     = false
  policy_settings             = jsonencode({
    cpu                = "2000m"
    memory             = "4Gi"
    excludedNamespaces = ["monitoring", "logging"]
  })
}
`

const testAccKubernetesSecurityPolicyV2UpdateOld = `
{{ .TestAccKubernetesClusterV2Base }}

data "vkcs_kubernetes_security_policy_template_v2" "policy_template" {
  name    = "k8scontainerrequests"
  version = "1.0.0"
}

resource "vkcs_kubernetes_security_policy_v2" "update" {
  cluster_id                  = vkcs_kubernetes_cluster_v2.base.id
  security_policy_template_id = data.vkcs_kubernetes_security_policy_template_v2.policy_template.id
  namespace                   = "*"
  enabled                     = true
  policy_settings             = jsonencode({
    cpu                = "1000m"
    memory             = "2Gi"
    excludedNamespaces = ["prometheus"]
  })
}
`

const testAccKubernetesSecurityPolicyV2UpdateNew = `
{{ .TestAccKubernetesClusterV2Base }}

data "vkcs_kubernetes_security_policy_template_v2" "policy_template" {
  name    = "k8scontainerrequests"
  version = "1.0.0"
}

resource "vkcs_kubernetes_security_policy_v2" "update" {
  cluster_id                  = vkcs_kubernetes_cluster_v2.base.id
  security_policy_template_id = data.vkcs_kubernetes_security_policy_template_v2.policy_template.id
  namespace                   = "updated-namespace"
  enabled                     = false
  policy_settings             = jsonencode({
    cpu                = "4000m"
    memory             = "8Gi"
    excludedNamespaces = ["monitoring"]
  })
}
`
