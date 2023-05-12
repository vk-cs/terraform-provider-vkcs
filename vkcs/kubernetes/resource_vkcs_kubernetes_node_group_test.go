package kubernetes_test

import (
	"fmt"
	"strconv"
	"testing"

	sdk_acctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/containerinfra/v1/clusters"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/containerinfra/v1/nodegroups"
)

func nodeGroupFixture(name, flavorID string, count, max, min int, autoscaling bool) *nodegroups.CreateOpts {
	return &nodegroups.CreateOpts{
		Name:        name,
		FlavorID:    flavorID,
		NodeCount:   count,
		MaxNodes:    max,
		MinNodes:    min,
		Autoscaling: autoscaling,
	}
}

const nodeGroupResourceFixture = `
		%s

		resource "vkcs_kubernetes_node_group" "%[2]s" {
          cluster_id          = vkcs_kubernetes_cluster.%s.id
		  name                = "%[2]s"
		  flavor_id           = "%[4]s"
		  node_count          =  "%d"
		  max_nodes           =  "%d"
		  min_nodes           =  "%d"
		  autoscaling_enabled =  "%t"
		}`

func TestAccKubernetesNodeGroup_basic(t *testing.T) {
	var cluster clusters.Cluster
	var nodeGroup nodegroups.NodeGroup

	clusterName := "testcluster" + sdk_acctest.RandStringFromCharSet(8, sdk_acctest.CharSetAlphaNum)
	createClusterFixture := clusterFixture(clusterName, acctest.ClusterTemplateID, acctest.OsFlavorID,
		acctest.OsKeypairName, acctest.OsNetworkID, acctest.OsSubnetworkID, "MS1", 1)
	clusterResourceName := "vkcs_kubernetes_cluster." + clusterName

	nodeGroupName := "testng" + sdk_acctest.RandStringFromCharSet(8, sdk_acctest.CharSetAlphaNum)
	ngFixture := nodeGroupFixture(nodeGroupName, acctest.OsFlavorID, 1, 5, 1, false)
	nodeGroupResourceName := "vkcs_kubernetes_node_group." + nodeGroupName

	ngNodeCountScaleFixture := nodeGroupFixture(nodeGroupName, acctest.OsFlavorID, 2, 5, 1, false)
	ngPatchOptsFixture := nodeGroupFixture(nodeGroupName, acctest.OsFlavorID, 2, 4, 2, true)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheckKubernetes(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckKubernetesClusterDestroy, Steps: []resource.TestStep{
			{
				Config: testAccKubernetesNodeGroupBasic(clusterName, testAccKubernetesClusterBasic(createClusterFixture), ngFixture),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesClusterExists(clusterResourceName, &cluster),
					testAccCheckKubernetesNodeGroupExists(nodeGroupResourceName, clusterResourceName, &nodeGroup),
					checkNodeGroupAttrs(nodeGroupResourceName, ngFixture),
				),
			},
			{
				Config: testAccKubernetesNodeGroupBasic(clusterName, testAccKubernetesClusterBasic(createClusterFixture), ngNodeCountScaleFixture),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(nodeGroupResourceName, "node_count", strconv.Itoa(ngNodeCountScaleFixture.NodeCount)),
					testAccCheckKubernetesNodeGroupScaled(nodeGroupResourceName),
				),
			},
			{
				Config: testAccKubernetesNodeGroupBasic(clusterName, testAccKubernetesClusterBasic(createClusterFixture), ngPatchOptsFixture),
				Check: resource.ComposeTestCheckFunc(
					checkNodeGroupPatchAttrs(nodeGroupResourceName, ngPatchOptsFixture),
					testAccCheckKubernetesNodeGroupPatched(nodeGroupResourceName),
				),
			},
		},
	})
}

func checkNodeGroupAttrs(resourceName string, nodeGroup *nodegroups.CreateOpts) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if s.Empty() == true {
			return fmt.Errorf("state not updated")
		}

		checksStore := []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(resourceName, "name", nodeGroup.Name),
			resource.TestCheckResourceAttr(resourceName, "node_count", strconv.Itoa(nodeGroup.NodeCount)),
			resource.TestCheckResourceAttr(resourceName, "flavor_id", nodeGroup.FlavorID),
			resource.TestCheckResourceAttr(resourceName, "max_nodes", strconv.Itoa(nodeGroup.MaxNodes)),
			resource.TestCheckResourceAttr(resourceName, "min_nodes", strconv.Itoa(nodeGroup.MinNodes)),
			resource.TestCheckResourceAttr(resourceName, "autoscaling_enabled", strconv.FormatBool(nodeGroup.Autoscaling)),
		}

		return resource.ComposeTestCheckFunc(checksStore...)(s)
	}
}

func checkNodeGroupPatchAttrs(resourceName string, nodeGroup *nodegroups.CreateOpts) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if s.Empty() == true {
			return fmt.Errorf("state not updated")
		}

		checksStore := []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(resourceName, "max_nodes", strconv.Itoa(nodeGroup.MaxNodes)),
			resource.TestCheckResourceAttr(resourceName, "min_nodes", strconv.Itoa(nodeGroup.MinNodes)),
			resource.TestCheckResourceAttr(resourceName, "autoscaling_enabled", strconv.FormatBool(nodeGroup.Autoscaling)),
		}

		return resource.ComposeTestCheckFunc(checksStore...)(s)
	}
}

func testAccCheckKubernetesNodeGroupExists(n, clusterResourceName string, nodeGroup *nodegroups.NodeGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, found, err := getNgAndResource(n, s)
		if err != nil {
			return err
		}
		cluster, _, err := getClusterAndResource(clusterResourceName, s)
		if err != nil {
			return err
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		if found.UUID != rs.Primary.ID {
			return fmt.Errorf("node group not found")
		}

		if cluster.Primary.ID != rs.Primary.Attributes["cluster_id"] {
			return fmt.Errorf(
				"mismatched cluster id in node_group; expected %s, but got %s",
				cluster.Primary.ID, rs.Primary.Attributes["cluster_id"])
		}

		*nodeGroup = *found

		return nil
	}
}

func testAccCheckKubernetesNodeGroupScaled(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, found, err := getNgAndResource(n, s)
		if err != nil {
			return err
		}

		if strconv.Itoa(found.NodeCount) != rs.Primary.Attributes["node_count"] {
			return fmt.Errorf("mismatched node_count")
		}
		return nil
	}
}

func testAccCheckKubernetesNodeGroupPatched(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, found, err := getNgAndResource(n, s)
		if err != nil {
			return err
		}

		if strconv.Itoa(found.MaxNodes) != rs.Primary.Attributes["max_nodes"] {
			return fmt.Errorf("mismatched max_nodes")
		}
		if strconv.Itoa(found.MinNodes) != rs.Primary.Attributes["min_nodes"] {
			return fmt.Errorf("mismatched min_nodes")
		}
		if strconv.FormatBool(found.Autoscaling) != rs.Primary.Attributes["autoscaling_enabled"] {
			return fmt.Errorf("mismatched autoscaling")
		}
		return nil
	}
}

func getNgAndResource(n string, s *terraform.State) (*terraform.ResourceState, *nodegroups.NodeGroup, error) {
	rs, ok := s.RootModule().Resources[n]
	if !ok {
		return nil, nil, fmt.Errorf("node group not found: %s", n)
	}

	config := acctest.AccTestProvider.Meta().(clients.Config)
	containerInfraClient, err := config.ContainerInfraV1Client(acctest.OsRegionName)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating container infra client: %s", err)
	}

	found, err := nodegroups.Get(containerInfraClient, rs.Primary.ID).Extract()
	if err != nil {
		return nil, nil, err
	}
	return rs, found, nil
}

func testAccKubernetesNodeGroupBasic(clusterName, clusterResource string, fixture *nodegroups.CreateOpts) string {
	return fmt.Sprintf(
		nodeGroupResourceFixture,
		clusterResource,
		fixture.Name,
		clusterName,
		fixture.FlavorID,
		fixture.NodeCount,
		fixture.MaxNodes,
		fixture.MinNodes,
		fixture.Autoscaling,
	)
}
