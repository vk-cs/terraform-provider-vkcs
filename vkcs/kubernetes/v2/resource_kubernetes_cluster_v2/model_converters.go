package resource_kubernetes_cluster_v2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/clusters"
	dskubeclusterv2 "github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes/v2/datasource_kubernetes_cluster_v2"
)

// ToClusterModel конвертирует clusters.Cluster в KubernetesClusterV2Model
func ToClusterModel(ctx context.Context, cluster *clusters.Cluster) (KubernetesClusterV2Model, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model KubernetesClusterV2Model

	// Set basic fields
	model.Id = types.StringValue(cluster.ID)
	model.Uuid = types.StringValue(cluster.UUID)
	model.Name = types.StringValue(cluster.Name)
	model.Version = types.StringValue(cluster.Version)
	model.Status = types.StringValue(cluster.Status)
	model.Description = types.StringValue(cluster.Description)
	model.CreatedAt = types.StringValue(cluster.CreatedAt)
	model.ProjectId = types.StringValue(cluster.ProjectID)
	model.ApiLbFip = types.StringValue(cluster.ExternalIP)
	model.ApiLbVip = types.StringValue(cluster.InternalIP)
	model.ApiAddress = types.StringValue(cluster.ApiAddress)

	// Set labels
	model.Labels, diags = dskubeclusterv2.FlattenLabels(cluster.Labels)
	if diags.HasError() {
		return model, diags
	}

	// Set insecure registries
	model.InsecureRegistries, diags = dskubeclusterv2.FlattenInsecureRegistries(cluster.InsecureRegistries)
	if diags.HasError() {
		return model, diags
	}

	// Set master spec fields
	model.MasterFlavor = types.StringValue(cluster.MasterSpec.Engine.NovaEngine.FlavorID)
	model.MasterCount = types.Int64Value(int64(cluster.MasterSpec.Replicas))

	// Set master disks
	model.MasterDisks, diags = dskubeclusterv2.FlattenMasterDisks(ctx, cluster.MasterSpec.Disks)
	if diags.HasError() {
		return model, diags
	}

	// Set deployment type fields
	model.ClusterType, model.AvailabilityZones, diags = dskubeclusterv2.FlattenAvailabilityZones(cluster.DeploymentType)
	if diags.HasError() {
		return model, diags
	}

	// Set network config fields
	model.NetworkPlugin, model.PodsIpv4Cidr, diags = dskubeclusterv2.FlattenNetworkPlugin(cluster.NetworkConfig.Plugin)
	if diags.HasError() {
		return model, diags
	}

	// Set network engine fields
	model.NetworkId = types.StringValue(cluster.NetworkConfig.Engine.SprutEngine.NetworkID)
	model.SubnetId = types.StringValue(cluster.NetworkConfig.Engine.SprutEngine.SubnetID)
	model.ExternalNetworkId = types.StringValue(cluster.NetworkConfig.Engine.SprutEngine.ExternalNetworkID)

	// Set load balancer config fields
	model.EnablePublicIp = types.BoolValue(cluster.LoadBalancerConfig.OctaviaEngine.EnablePublicIP)
	model.LoadbalancerSubnetId = types.StringValue(cluster.LoadBalancerConfig.OctaviaEngine.LoadbalancerSubnetID)

	// Set loadbalancer allowed CIDRs
	model.LoadbalancerAllowedCidrs, diags = dskubeclusterv2.FlattenLoadbalancerAllowedCIDRs(cluster.LoadBalancerConfig.OctaviaEngine.AllowedCIDRs)
	if diags.HasError() {
		return model, diags
	}

	// Set node groups
	model.NodeGroups, diags = dskubeclusterv2.FlattenNodeGroups(ctx, cluster.NodeGroups)
	if diags.HasError() {
		return model, diags
	}

	return model, diags
}
