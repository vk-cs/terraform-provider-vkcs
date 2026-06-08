package kubernetes_test

import (
	"fmt"
	"testing"

	acctest_helper "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/addons"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

// the _basic test matches the _full test
func TestAccKubernetesClusterAddonV2_basic(t *testing.T) {
	t.Parallel()

	clusterName := "tfacc-clad-basic-v2-" + acctest_helper.RandStringFromCharSet(2, acctest_helper.CharSetAlphaNum)
	clusterConfig := acctest.AccTestRenderConfig(testAccKubernetesClusterV2Base, map[string]string{
		"TestAccKubernetesNetworkingBase": testAccKubernetesNetworkingBase,
		"ClusterName":                     clusterName,
	})
	clusterWithNodeGroupConfig := acctest.AccTestRenderConfig(testAccKubernetesNodeGroupV2Base, map[string]string{
		"TestAccKubernetesNodeGroupV2ClusterBase": clusterConfig,
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesClusterAddonV2Base, map[string]string{
					"TestAccKubernetesClusterAddonV2ClusterBase": clusterWithNodeGroupConfig,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesClusterAddonV2Exists("vkcs_kubernetes_cluster_addon_v2.base"),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_cluster_addon_v2.base", "cluster_id", "vkcs_kubernetes_cluster_v2.base", "id"),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_cluster_addon_v2.base", "addon_id", "data.vkcs_kubernetes_addon_v2.ingress-nginx", "addon_id"),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_cluster_addon_v2.base", "addon_version_id", "data.vkcs_kubernetes_addon_v2.ingress-nginx", "id"),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_cluster_addon_v2.base", "values", "data.vkcs_kubernetes_addon_v2.ingress-nginx", "values_template"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_addon_v2.base", "addon_name", "ingress-nginx"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_addon_v2.base", "namespace", "ingress-nginx"),

					resource.TestCheckResourceAttrSet("vkcs_kubernetes_cluster_addon_v2.base", "id"),
					resource.TestCheckResourceAttrSet("vkcs_kubernetes_cluster_addon_v2.base", "region"),
					resource.TestCheckResourceAttrSet("vkcs_kubernetes_cluster_addon_v2.base", "created_at"),
					resource.TestCheckResourceAttrSet("vkcs_kubernetes_cluster_addon_v2.base", "updated_at"),
					resource.TestCheckResourceAttrSet("vkcs_kubernetes_cluster_addon_v2.base", "status"),
				),
			},
			acctest.ImportStep("vkcs_kubernetes_cluster_addon_v2.base"),
		},
	})
}

func testAccCheckKubernetesClusterAddonV2Exists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Kubernetes cluster addon not found: %s", n)
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
			return fmt.Errorf("Error creating Kubernetes API client: %s", err)
		}

		_, err = addons.GetClusterAddon(client, rs.Primary.ID).Extract()
		if err != nil {
			if errutil.IsNotFound(err) {
				return fmt.Errorf("Kubernetes cluster addon not found")
			}
			return err
		}

		return nil
	}
}

const testAccKubernetesClusterAddonV2Base = `
{{ .TestAccKubernetesClusterAddonV2ClusterBase }}

data "vkcs_kubernetes_addon_v2" "ingress-nginx" {
  name    = "ingress-nginx"
  version = "4.12.1"
}

resource "vkcs_kubernetes_cluster_addon_v2" "base" {
  cluster_id       = vkcs_kubernetes_cluster_v2.base.id
  addon_id         = data.vkcs_kubernetes_addon_v2.ingress-nginx.addon_id
  addon_name       = "ingress-nginx"
  addon_version_id = data.vkcs_kubernetes_addon_v2.ingress-nginx.id
  namespace        = "ingress-nginx"
  values           = data.vkcs_kubernetes_addon_v2.ingress-nginx.values_template
}
`
