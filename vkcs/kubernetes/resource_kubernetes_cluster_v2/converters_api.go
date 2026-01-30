package resource_kubernetes_cluster_v2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/clusters"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func ToCreateOpts(ctx context.Context, model *KubernetesClusterV2Model) (*clusters.CreateOpts, diag.Diagnostics) {
	if model == nil {
		return nil, nil
	}

	var diags diag.Diagnostics

	// Build labels map
	labels := make(map[string]string, len(model.Labels.Elements()))
	if !util.IsNullOrUnknown(model.Labels) {
		diags.Append(model.Labels.ElementsAs(ctx, &labels, false)...)
		if diags.HasError() {
			return nil, diags
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
	insecureRegistries := make([]string, 0, len(model.InsecureRegistries.Elements()))
	if !util.IsNullOrUnknown(model.InsecureRegistries) {
		diags.Append(model.InsecureRegistries.ElementsAs(ctx, &insecureRegistries, false)...)
		if diags.HasError() {
			return nil, diags
		}
	}

	// Build deployment type
	var deploymentType clusters.DeploymentTypeOpts
	clusterType := model.ClusterType.ValueString()

	zones := make([]string, 0, len(model.AvailabilityZones.Elements()))
	diags.Append(model.AvailabilityZones.ElementsAs(ctx, &zones, false)...)
	if diags.HasError() {
		return nil, diags
	}

	switch clusterType {
	case "standard":
		deploymentType = clusters.DeploymentTypeOpts{
			ZonalDeployment: &clusters.ZonalDeploymentOpts{
				Zone: zones[0],
			},
		}
	case "regional":
		deploymentType = clusters.DeploymentTypeOpts{
			MultiZonalDeployment: &clusters.MultiZonalDeploymentOpts{
				Zones: zones,
			},
		}
	default:
		diags.AddAttributeError(
			path.Root("cluster_type"),
			"Unknown cluster type",
			"Cluster type must be either standard or regional, got "+clusterType,
		)
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
	case "cilium":
		networkPlugin = clusters.NetworkPluginOpts{
			Cilium: &clusters.CiliumPluginOpts{
				PodsIPv4CIDR: podsIPv4CIDR,
			},
		}
	default:
		diags.AddAttributeError(
			path.Root("network_plugin"),
			"Unknown network plugin",
			"Network plugin must be either calico or cilium, got "+networkPluginType,
		)
	}

	// Build load balancer config
	allowedCIDRs := make([]string, 0, len(model.LoadbalancerAllowedCidrs.Elements()))
	if !util.IsNullOrUnknown(model.LoadbalancerAllowedCidrs) {
		diags.Append(model.LoadbalancerAllowedCidrs.ElementsAs(ctx, &allowedCIDRs, false)...)
		if diags.HasError() {
			return nil, diags
		}
	}

	loadBalancerConfig := clusters.LoadBalancerConfigOpts{
		OctaviaEngine: clusters.OctaviaEngineOpts{
			EnablePublicIP:       model.PublicIp.ValueBool(),
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
