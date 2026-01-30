package datasource_kubernetes_cluster_v2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/clusters"
	mshared "github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes/models_shared"
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
	m.MasterDisks, diags = FlattenMasterDisks(ctx, cluster.MasterSpec.Disks)
	if diags.HasError() {
		return diags
	}

	// Set deployment type fields
	m.ClusterType, m.AvailabilityZones, diags = FlattenDeploymentType(cluster.DeploymentType)
	if diags.HasError() {
		return diags
	}

	// Set network config fields
	m.NetworkPlugin, m.PodsIpv4Cidr, diags = FlattenNetworkPlugin(cluster.NetworkConfig.Plugin)
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
	m.NodeGroups, diags = FlattenNodeGroups(ctx, cluster.NodeGroups)
	if diags.HasError() {
		return diags
	}

	return
}

func FlattenMasterDisks(ctx context.Context, disks []clusters.MasterSpecDisk) (types.Set, diag.Diagnostics) {
	if len(disks) == 0 {
		return types.SetNull(MasterDisksValue{}.Type(ctx)), nil
	}

	resSet := make([]attr.Value, len(disks))
	for i, disk := range disks {
		resSet[i] = MasterDisksValue{
			Size:            types.Int64Value(int64(disk.Size)),
			MasterDisksType: types.StringValue(disk.Type),
			state:           attr.ValueStateKnown,
		}
	}

	return types.SetValue(MasterDisksValue{}.Type(ctx), resSet)
}

func FlattenDeploymentType(deploymentType clusters.DeploymentType) (
	clusterType types.String,
	availabilityZones types.Set,
	diags diag.Diagnostics,
) {
	switch {
	case deploymentType.ZonalDeployment != nil:
		clusterType = types.StringValue("standard")
		zones := []attr.Value{types.StringValue(deploymentType.ZonalDeployment.Zone)}
		availabilityZones, diags = types.SetValue(types.StringType, zones)
	case deploymentType.MultiZonalDeployment != nil:
		clusterType = types.StringValue("regional")
		availabilityZones, diags = mshared.FlattenStringSet(deploymentType.MultiZonalDeployment.Zones)
	default:
		clusterType = types.StringNull()
		availabilityZones = types.SetNull(types.StringType)
	}

	return
}

func FlattenNetworkPlugin(plugin clusters.NetworkPlugin) (
	networkPlugin types.String,
	podsIpv4Cidr types.String,
	_ diag.Diagnostics,
) {
	switch {
	case plugin.Calico != nil:
		networkPlugin = types.StringValue("calico")
		podsIpv4Cidr = types.StringValue(plugin.Calico.PodsIPv4CIDR)
	case plugin.Cilium != nil:
		networkPlugin = types.StringValue("cilium")
		podsIpv4Cidr = types.StringValue(plugin.Cilium.PodsIPv4CIDR)
	default:
		networkPlugin = types.StringNull()
		podsIpv4Cidr = types.StringNull()
	}

	return
}

func FlattenNodeGroups(ctx context.Context, nodeGroups []clusters.NodeGroup) (types.Set, diag.Diagnostics) {
	if len(nodeGroups) == 0 {
		return types.SetValue(NodeGroupsValue{}.Type(ctx), nil)
	}

	resSet := make([]attr.Value, len(nodeGroups))
	for i, ng := range nodeGroups {
		ngZone := ""
		if len(ng.Zones) > 0 {
			ngZone = ng.Zones[0]
		}

		resSet[i] = NodeGroupsValue{
			Id:               types.StringValue(ng.ID),
			Name:             types.StringValue(ng.Name),
			Flavor:           types.StringValue(ng.VMEngine.NovaEngine.FlavorID),
			NodeCount:        types.Int64Value(int64(ng.GetActualSize())),
			AvailabilityZone: types.StringValue(ngZone),
			state:            attr.ValueStateKnown,
		}
	}

	return types.SetValue(NodeGroupsValue{}.Type(ctx), resSet)
}
