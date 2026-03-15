package datasource_kubernetes_cluster_v2

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/clusters"
)

// FlattenLabels конвертирует map[string]string в types.Map
func FlattenLabels(labels map[string]string) (types.Map, diag.Diagnostics) {
	if labels == nil {
		return types.MapNull(types.StringType), nil
	}

	labelsMap := make(map[string]attr.Value, len(labels))
	for k, v := range labels {
		labelsMap[k] = types.StringValue(v)
	}

	return types.MapValue(types.StringType, labelsMap)
}

// FlattenInsecureRegistries конвертирует []string в types.Set
func FlattenInsecureRegistries(registries []string) (types.Set, diag.Diagnostics) {
	if len(registries) == 0 {
		return types.SetNull(types.StringType), nil
	}

	registriesList := make([]attr.Value, len(registries))
	for i, reg := range registries {
		registriesList[i] = types.StringValue(reg)
	}

	return types.SetValue(types.StringType, registriesList)
}

// FlattenMasterDisks конвертирует []clusters.MasterSpecDisk в types.Set
func FlattenMasterDisks(ctx context.Context, disks []clusters.MasterSpecDisk) (types.Set, diag.Diagnostics) {
	if len(disks) == 0 {
		return types.SetNull(MasterDisksValue{}.Type(ctx)), nil
	}

	masterDisks := make([]attr.Value, len(disks))
	for i, disk := range disks {
		masterDisks[i] = MasterDisksValue{
			Size:            types.StringValue(strconv.Itoa(disk.Size)),
			MasterDisksType: types.StringValue(disk.Type),
			state:           attr.ValueStateKnown,
		}
	}

	return types.SetValue(MasterDisksValue{}.Type(ctx), masterDisks)
}

// FlattenAvailabilityZones конвертирует DeploymentType в cluster_type и availability_zones
func FlattenAvailabilityZones(deploymentType clusters.DeploymentType) (types.String, types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics
	var clusterType types.String
	var availabilityZones types.Set

	switch {
	case deploymentType.ZonalDeployment != nil:
		clusterType = types.StringValue("standard")
		zones := []attr.Value{types.StringValue(deploymentType.ZonalDeployment.Zone)}
		availabilityZones, diags = types.SetValue(types.StringType, zones)
	case deploymentType.MultiZonalDeployment != nil:
		clusterType = types.StringValue("regional")
		zones := make([]attr.Value, len(deploymentType.MultiZonalDeployment.Zones))
		for i, zone := range deploymentType.MultiZonalDeployment.Zones {
			zones[i] = types.StringValue(zone)
		}
		availabilityZones, diags = types.SetValue(types.StringType, zones)
	default:
		clusterType = types.StringNull()
		availabilityZones = types.SetNull(types.StringType)
	}

	return clusterType, availabilityZones, diags
}

// FlattenNetworkPlugin конвертирует NetworkConfig.Plugin в network_plugin и pods_ipv4_cidr
func FlattenNetworkPlugin(plugin clusters.NetworkPlugin) (types.String, types.String, diag.Diagnostics) {
	var networkPlugin types.String
	var podsIpv4Cidr types.String

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

	return networkPlugin, podsIpv4Cidr, nil
}

// FlattenLoadbalancerAllowedCIDRs конвертирует []string в types.Set
func FlattenLoadbalancerAllowedCIDRs(allowedCIDRs []string) (types.Set, diag.Diagnostics) {
	if len(allowedCIDRs) == 0 {
		return types.SetNull(types.StringType), nil
	}

	cidrs := make([]attr.Value, len(allowedCIDRs))
	for i, cidr := range allowedCIDRs {
		cidrs[i] = types.StringValue(cidr)
	}

	return types.SetValue(types.StringType, cidrs)
}

// FlattenNodeGroups конвертирует []clusters.NodeGroup в types.Set
func FlattenNodeGroups(ctx context.Context, nodeGroups []clusters.NodeGroup) (types.Set, diag.Diagnostics) {
	if len(nodeGroups) == 0 {
		return types.SetNull(NodeGroupsValue{}.Type(ctx)), nil
	}

	ngList := make([]attr.Value, len(nodeGroups))
	for i, ng := range nodeGroups {
		ngZone := ""
		if len(ng.Zones) > 0 {
			ngZone = ng.Zones[0]
		}

		ngList[i] = NodeGroupsValue{
			Id:               types.StringValue(ng.ID),
			Name:             types.StringValue(ng.Name),
			Flavor:           types.StringValue(ng.VMEngine.NovaEngine.FlavorID),
			NodeCount:        types.Int64Value(int64(ng.GetActualSize())),
			AvailabilityZone: types.StringValue(ngZone),
			state:            attr.ValueStateKnown,
		}
	}

	return types.SetValue(NodeGroupsValue{}.Type(ctx), ngList)
}

// FlattenCluster конвертирует clusters.Cluster в KubernetesClusterV2Model
func FlattenCluster(ctx context.Context, cluster *clusters.Cluster) (KubernetesClusterV2Model, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model KubernetesClusterV2Model

	// Set basic fields
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
	model.Id = types.StringValue(cluster.ID)

	// Set labels
	model.Labels, diags = FlattenLabels(cluster.Labels)
	if diags.HasError() {
		return model, diags
	}

	// Set insecure registries
	model.InsecureRegistries, diags = FlattenInsecureRegistries(cluster.InsecureRegistries)
	if diags.HasError() {
		return model, diags
	}

	// Set master spec fields
	model.MasterFlavor = types.StringValue(cluster.MasterSpec.Engine.NovaEngine.FlavorID)
	model.MasterCount = types.Int64Value(int64(cluster.MasterSpec.Replicas))

	// Set master disks
	model.MasterDisks, diags = FlattenMasterDisks(ctx, cluster.MasterSpec.Disks)
	if diags.HasError() {
		return model, diags
	}

	// Set deployment type fields
	model.ClusterType, model.AvailabilityZones, diags = FlattenAvailabilityZones(cluster.DeploymentType)
	if diags.HasError() {
		return model, diags
	}

	// Set network config fields
	model.NetworkPlugin, model.PodsIpv4Cidr, diags = FlattenNetworkPlugin(cluster.NetworkConfig.Plugin)
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
	model.LoadbalancerAllowedCidrs, diags = FlattenLoadbalancerAllowedCIDRs(cluster.LoadBalancerConfig.OctaviaEngine.AllowedCIDRs)
	if diags.HasError() {
		return model, diags
	}

	// Set node groups
	model.NodeGroups, diags = FlattenNodeGroups(ctx, cluster.NodeGroups)
	if diags.HasError() {
		return model, diags
	}

	return model, diags
}
