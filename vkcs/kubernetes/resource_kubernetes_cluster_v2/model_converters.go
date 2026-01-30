package resource_kubernetes_cluster_v2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/clusters"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	dskubeclusterv2 "github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes/datasource_kubernetes_cluster_v2"
	mshared "github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes/models_shared"
)

const (
	createClusterDefaultTimeoutV2 = "60m"
	deleteClusterDefaultTimeoutV2 = "30m"
	updateClusterDefaultTimeoutV2 = "60m"
)

func (m *KubernetesClusterV2Model) UpdateFromCluster(ctx context.Context, cluster *clusters.Cluster, kubeconfig *string) (diags diag.Diagnostics) {
	if cluster != nil {
		diags.Append(m.updateClusterAttributes(ctx, cluster)...)
		if diags.HasError() {
			return
		}
	}

	if kubeconfig != nil {
		m.K8sConfig = types.StringValue(*kubeconfig)
	}

	// Good for terraform import state
	if util.IsNullOrUnknown(m.Timeouts) {
		m.Timeouts = GetDefaultClusterV2Timeouts(ctx)
	}

	return
}

func (m *KubernetesClusterV2Model) updateClusterAttributes(ctx context.Context, cluster *clusters.Cluster) (diags diag.Diagnostics) {
	// Set basic fields
	m.Id = types.StringValue(cluster.ID)
	m.Uuid = types.StringValue(cluster.UUID)
	m.Name = types.StringValue(cluster.Name)
	m.Version = types.StringValue(cluster.Version)
	m.Status = types.StringValue(cluster.Status)
	m.Description = types.StringValue(cluster.Description)
	m.CreatedAt = types.StringValue(cluster.CreatedAt)
	m.ProjectId = types.StringValue(cluster.ProjectID)
	m.ApiLbFip = types.StringValue(cluster.ExternalIP)
	m.ApiLbVip = types.StringValue(cluster.InternalIP)
	m.ApiAddress = types.StringValue(cluster.ApiAddress)

	// Set labels
	m.Labels, diags = mshared.FlattenStringMap(cluster.Labels)
	if diags.HasError() {
		return diags
	}

	// Set insecure registries. Sort registries in ascending order.
	m.InsecureRegistries, diags = mshared.FlattenStringSet(cluster.InsecureRegistries)
	if diags.HasError() {
		return diags
	}

	// Set master spec fields
	m.MasterFlavor = types.StringValue(cluster.MasterSpec.Engine.NovaEngine.FlavorID)
	m.MasterCount = types.Int64Value(int64(cluster.MasterSpec.Replicas))

	// Set master disks
	m.MasterDisks, diags = dskubeclusterv2.FlattenMasterDisks(ctx, cluster.MasterSpec.Disks)
	if diags.HasError() {
		return diags
	}

	// Set deployment type fields
	m.ClusterType, m.AvailabilityZones, diags = dskubeclusterv2.FlattenDeploymentType(cluster.DeploymentType)
	if diags.HasError() {
		return diags
	}

	// Set network config fields
	m.NetworkPlugin, m.PodsIpv4Cidr, diags = dskubeclusterv2.FlattenNetworkPlugin(cluster.NetworkConfig.Plugin)
	if diags.HasError() {
		return diags
	}

	// Set network engine fields
	m.NetworkId = types.StringValue(cluster.NetworkConfig.Engine.SprutEngine.NetworkID)
	m.SubnetId = types.StringValue(cluster.NetworkConfig.Engine.SprutEngine.SubnetID)
	m.ExternalNetworkId = types.StringValue(cluster.NetworkConfig.Engine.SprutEngine.ExternalNetworkID)

	// Set load balancer config fields
	m.PublicIp = types.BoolValue(cluster.LoadBalancerConfig.OctaviaEngine.EnablePublicIP)
	m.LoadbalancerSubnetId = types.StringValue(cluster.LoadBalancerConfig.OctaviaEngine.LoadbalancerSubnetID)

	// Set loadbalancer allowed CIDRs
	m.LoadbalancerAllowedCidrs, diags = mshared.FlattenStringSet(cluster.LoadBalancerConfig.OctaviaEngine.AllowedCIDRs)
	if diags.HasError() {
		return diags
	}

	// Set node groups
	m.NodeGroups, diags = dskubeclusterv2.FlattenNodeGroups(ctx, cluster.NodeGroups)
	if diags.HasError() {
		return diags
	}

	return
}

func GetDefaultClusterV2Timeouts(ctx context.Context) TimeoutsValue {
	return NewTimeoutsValueMust(
		TimeoutsValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"create": types.StringValue(createClusterDefaultTimeoutV2),
			"delete": types.StringValue(deleteClusterDefaultTimeoutV2),
			"update": types.StringValue(updateClusterDefaultTimeoutV2),
		},
	)
}

func GetDefaultClusterV2CreateTimeout() basetypes.StringValue {
	return types.StringValue(createClusterDefaultTimeoutV2)
}

func GetDefaultClusterV2DeleteTimeout() basetypes.StringValue {
	return types.StringValue(deleteClusterDefaultTimeoutV2)
}

func GetDefaultClusterV2UpdateTimeout() basetypes.StringValue {
	return types.StringValue(updateClusterDefaultTimeoutV2)
}
