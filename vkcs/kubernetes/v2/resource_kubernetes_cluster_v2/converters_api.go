package resource_kubernetes_cluster_v2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/clusters"
)

// ToCreateOpts конвертирует KubernetesClusterV2Model в clusters.CreateOpts
func ToCreateOpts(ctx context.Context, model *KubernetesClusterV2Model) (*clusters.CreateOpts, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Build labels map
	labels := make(map[string]string)
	if !model.Labels.IsNull() {
		elements := make(map[string]types.String, len(model.Labels.Elements()))
		diags.Append(model.Labels.ElementsAs(ctx, &elements, false)...)
		if diags.HasError() {
			return nil, diags
		}
		for k, v := range elements {
			labels[k] = v.ValueString()
		}
	}

	// Build master spec
	masterSpec := clusters.MasterSpecOpts{
		Engine: clusters.MasterEngineOpts{
			NovaEngine: clusters.NovaEngineOpts{
				FlavorID: model.MasterFlavor.ValueString(),
			},
		},
		Replicas: int(model.MasterCount.ValueInt64()),
	}

	// Build insecure registries
	var insecureRegistries []string
	if !model.InsecureRegistries.IsNull() {
		elements := make([]types.String, 0, len(model.InsecureRegistries.Elements()))
		diags.Append(model.InsecureRegistries.ElementsAs(ctx, &elements, false)...)
		if diags.HasError() {
			return nil, diags
		}
		for _, reg := range elements {
			insecureRegistries = append(insecureRegistries, reg.ValueString())
		}
	}

	// Build deployment type
	var deploymentType clusters.DeploymentTypeOpts
	clusterType := model.ClusterType.ValueString()

	if !model.AvailabilityZones.IsNull() {
		var zones []types.String
		diags.Append(model.AvailabilityZones.ElementsAs(ctx, &zones, false)...)
		if diags.HasError() {
			return nil, diags
		}

		switch clusterType {
		case "standard":
			deploymentType = clusters.DeploymentTypeOpts{
				ZonalDeployment: &clusters.ZonalDeploymentOpts{
					Zone: zones[0].ValueString(),
				},
			}
		case "regional":
			zoneStrings := make([]string, len(zones))
			for i, zone := range zones {
				zoneStrings[i] = zone.ValueString()
			}
			deploymentType = clusters.DeploymentTypeOpts{
				MultiZonalDeployment: &clusters.MultiZonalDeploymentOpts{
					Zones: zoneStrings,
				},
			}
		}
	}

	// Build network plugin
	var networkPlugin clusters.NetworkPluginOpts
	networkPluginType := model.NetworkPlugin.ValueString()
	podsIPv4CIDR := model.PodsIpv4Cidr.ValueString()

	switch networkPluginType {
	case "calico":
		networkPlugin = clusters.NetworkPluginOpts{
			Calico: &clusters.CalicoPluginOpts{
				PodsIPv4CIDR: podsIPv4CIDR,
			},
		}
	case "cillium":
		networkPlugin = clusters.NetworkPluginOpts{
			Cilium: &clusters.CiliumPluginOpts{
				PodsIPv4CIDR: podsIPv4CIDR,
			},
		}
	}

	// Build load balancer config
	var allowedCIDRs []string
	if !model.LoadbalancerAllowedCidrs.IsNull() {
		var cidrs []types.String
		diags.Append(model.LoadbalancerAllowedCidrs.ElementsAs(ctx, &cidrs, false)...)
		if diags.HasError() {
			return nil, diags
		}
		for _, cidr := range cidrs {
			allowedCIDRs = append(allowedCIDRs, cidr.ValueString())
		}
	}

	loadBalancerConfig := clusters.LoadBalancerConfigOpts{
		OctaviaEngine: clusters.OctaviaEngineOpts{
			EnablePublicIP:       model.EnablePublicIp.ValueBool(),
			LoadbalancerSubnetID: model.LoadbalancerSubnetId.ValueString(),
			AllowedCIDRs:         allowedCIDRs,
		},
	}

	// Build create options
	createOpts := clusters.CreateOpts{
		UUID:               model.Uuid.ValueString(),
		Name:               model.Name.ValueString(),
		Version:            model.Version.ValueString(),
		Description:        model.Description.ValueString(),
		Labels:             labels,
		InsecureRegistries: insecureRegistries,
		MasterSpec:         masterSpec,
		DeploymentType:     deploymentType,
		NetworkConfig: clusters.NetworkConfigOpts{
			Plugin: networkPlugin,
			Engine: clusters.NetworkEngineOpts{
				SprutEngine: clusters.SprutEngineOpts{
					NetworkID:         model.NetworkId.ValueString(),
					SubnetID:          model.SubnetId.ValueString(),
					ExternalNetworkID: model.ExternalNetworkId.ValueString(),
				},
			},
		},
		LoadBalancerConfig: loadBalancerConfig,
	}

	return &createOpts, diags
}
