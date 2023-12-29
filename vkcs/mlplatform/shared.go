package mlplatform

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/framework/planmodifiers"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/mlplatform/v1/backups"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/mlplatform/v1/instances"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

const (
	jupyterHubInstanceType   = "JUPYTERHUB"
	mlFlowInstanceType       = "MLFLOW"
	mlFlowDeployInstanceType = "DEPLOY"
)

const (
	instanceDelay         = 10 * time.Second
	instanceMinTimeout    = 10 * time.Second
	instanceCreateTimeout = 30 * time.Minute
	instanceUpdateTimeout = 30 * time.Minute
	instanceDeleteTimeout = 30 * time.Minute
)

const (
	instanceStatusPrepareDBAAS   = "PREPARE_DBAAS"
	instanceStatusCreating       = "CREATING"
	instanceStatusInstallScripts = "INSTALL_SCRIPTS"
	instanceStatusStarting       = "STARTING"
	instanceStatusRunning        = "RUNNING"
	instanceStatusDeleting       = "DELETING"
	instanceStatusDeleted        = "DELETED"
	instanceStatusCreateFailed   = "CREATE_FAILED"
)

const (
	updateFlavorActionType = "FLAVOR"
	resizeVolumeActionType = "VOLUME"
)

type MLPlatformVolumeModel struct {
	Name       types.String `tfsdk:"name"`
	Size       types.Int64  `tfsdk:"size"`
	VolumeType types.String `tfsdk:"volume_type"`
	VolumeID   types.String `tfsdk:"volume_id"`
}

type MLPlatformNetworkModel struct {
	IPPool    types.String `tfsdk:"ip_pool"`
	NetworkID types.String `tfsdk:"network_id"`
}

func getCommonInstanceSchema(ctx context.Context, resp *resource.SchemaResponse) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "ID of the resource",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},

		"name": schema.StringAttribute{
			Required:    true,
			Description: "Instance name. Changing this creates a new resource",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},

		"flavor_id": schema.StringAttribute{
			Required:    true,
			Description: "Flavor ID",
		},

		"availability_zone": schema.StringAttribute{
			Required:    true,
			Description: "The availability zone in which to create the resource. Changing this creates a new resource",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},

		"boot_volume": schema.SingleNestedAttribute{
			Attributes: map[string]schema.Attribute{
				"size": schema.Int64Attribute{
					Optional:    true,
					Computed:    true,
					Description: "Size of the volume",
				},
				"volume_type": schema.StringAttribute{
					Required:    true,
					Description: "Type of the volume",
				},
				"name": schema.StringAttribute{
					Computed:    true,
					Description: "Name of the volume",
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"volume_id": schema.StringAttribute{
					Computed:    true,
					Description: "ID of the volume",
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			Required:    true,
			Description: "Instance's boot volume configuration",
		},

		"networks": schema.ListNestedAttribute{
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"ip_pool": schema.StringAttribute{
						Optional:    true,
						Description: "ID of the ip pool",
					},
					"network_id": schema.StringAttribute{
						Required:    true,
						Description: "ID of the network",
					},
				},
			},
			Required:    true,
			Description: "Network configuration",
			Validators: []validator.List{
				listvalidator.SizeAtMost(1),
			},
		},

		// Computed fields

		"created_at": schema.StringAttribute{
			Computed:    true,
			Description: "Creation timestamp",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},

		"private_ip": schema.StringAttribute{
			Computed:    true,
			Description: "Private IP address",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},

		"dns_name": schema.StringAttribute{
			Computed:    true,
			Description: "DNS name",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},

		"region": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The `region` to fetch availability zones from, defaults to the provider's `region`.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplaceIf(planmodifiers.GetRegionPlanModifier(resp),
					"require replacement if configuration value changes", "require replacement if configuration value changes"),
				stringplanmodifier.UseStateForUnknown(),
			},
		},

		"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
			Create: true,
			Update: true,
			Delete: true,
		}),
	}
}

func expandVolumeOpts(bootVolume *MLPlatformVolumeModel, bootVolumeSize int64, dataVolumes []*MLPlatformVolumeModel, availabilityZone string) []instances.VolumeCreateOpts {
	var volumeOpts []instances.VolumeCreateOpts

	bootVolumeOpts := instances.VolumeCreateOpts{
		Name:             bootVolume.Name.ValueString(),
		Size:             int(bootVolumeSize),
		VolumeType:       bootVolume.VolumeType.ValueString(),
		AvailabilityZone: availabilityZone,
	}

	volumeOpts = append(volumeOpts, bootVolumeOpts)

	for _, dataVolume := range dataVolumes {
		dataVolumeOpts := instances.VolumeCreateOpts{
			Name:             dataVolume.Name.ValueString(),
			Size:             int(dataVolume.Size.ValueInt64()),
			VolumeType:       dataVolume.VolumeType.ValueString(),
			AvailabilityZone: availabilityZone,
		}
		volumeOpts = append(volumeOpts, dataVolumeOpts)
	}
	return volumeOpts
}

func expandNetworkOpts(networks []*MLPlatformNetworkModel) instances.NetworkCreateOpts {
	networkOpts := instances.NetworkCreateOpts{
		IPPool:    networks[0].IPPool.ValueString(),
		NetworkID: networks[0].NetworkID.ValueString(),
	}
	return networkOpts
}

func flattenVolumeOpts(volumes []instances.VolumeResponse) (*MLPlatformVolumeModel, []*MLPlatformVolumeModel, string) {
	var bootVolume *MLPlatformVolumeModel
	var dataVolumes []*MLPlatformVolumeModel
	for _, volume := range volumes {
		volumeFlattened := MLPlatformVolumeModel{
			Name:       types.StringValue(volume.Name),
			Size:       types.Int64Value(int64(volume.Size)),
			VolumeType: types.StringValue(volume.VolumeType),
			VolumeID:   types.StringValue(volume.CinderID),
		}
		if strings.Contains(volume.Name, "boot") {
			bootVolume = &volumeFlattened
		} else {
			dataVolumes = append(dataVolumes, &volumeFlattened)
		}
	}
	sort.Slice(dataVolumes, func(i, j int) bool {
		return dataVolumes[i].Name.ValueString() < dataVolumes[j].Name.ValueString()
	})
	return bootVolume, dataVolumes, volumes[0].AvailabilityZone
}

func instanceStateRefreshFunc(client *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		i, err := instances.Get(client, id).Extract()
		if err != nil {
			if errutil.IsNotFound(err) {
				return i, instanceStatusDeleted, nil
			}
			return nil, "", err
		}

		if i.Status == instanceStatusCreateFailed {
			return i, i.Status, fmt.Errorf("instance is in failed status, retry the operation or contact support")
		}

		return i, i.Status, nil
	}
}

func instanceUpdateFlavor(ctx context.Context, client *gophercloud.ServiceClient, id string, newFlavorID string, timeout time.Duration) error {
	updateFlavorOpts := instances.ActionOpts{
		Action: instances.ResizeAction{
			Resize: instances.ResizeActionOpts{
				Type: updateFlavorActionType,
				Params: instances.ResizeActionParams{
					Flavor: newFlavorID,
				},
			},
		},
	}

	_, err := instances.Action(client, id, &updateFlavorOpts).Extract()
	if err != nil {
		return err
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{instanceStatusPrepareDBAAS, instanceStatusCreating, instanceStatusInstallScripts, instanceStatusStarting},
		Target:     []string{instanceStatusRunning},
		Refresh:    instanceStateRefreshFunc(client, id),
		Timeout:    timeout,
		Delay:      instanceDelay,
		MinTimeout: instanceMinTimeout,
	}

	tflog.Debug(ctx, "Waiting for the instance to update flavor", map[string]interface{}{"timeout": timeout})
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func instanceUpdateVolumes(ctx context.Context, client *gophercloud.ServiceClient, id string, volumes []instances.ResizeVolumeParams, timeout time.Duration) error {
	for _, volume := range volumes {
		resizeVolumesOpts := instances.ActionOpts{
			Action: instances.ResizeAction{
				Resize: instances.ResizeActionOpts{
					Type: resizeVolumeActionType,
					Params: instances.ResizeActionParams{
						Volumes: []instances.ResizeVolumeParams{volume},
					},
				},
			},
		}

		_, err := instances.Action(client, id, &resizeVolumesOpts).Extract()
		if err != nil {
			return err
		}

		jupyterHubStateConf := &retry.StateChangeConf{
			Pending:    []string{instanceStatusPrepareDBAAS, instanceStatusCreating, instanceStatusInstallScripts, instanceStatusStarting},
			Target:     []string{instanceStatusRunning},
			Refresh:    instanceStateRefreshFunc(client, id),
			Timeout:    timeout,
			Delay:      instanceDelay,
			MinTimeout: instanceMinTimeout,
		}

		tflog.Debug(ctx, "Waiting for the instance to resize volume", map[string]interface{}{"timeout": timeout})
		_, err = jupyterHubStateConf.WaitForStateContext(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func flattenBackupOpts(backupsRaw []*backups.Response) []*MLPlatformBackupModel {
	backups := make([]*MLPlatformBackupModel, len(backupsRaw))

	for i, backup := range backupsRaw {
		backupFlattened := MLPlatformBackupModel{
			VolumeID:  types.StringValue(backup.CinderID),
			CreatedAt: types.StringValue(backup.CreatedAt),
			BackupID:  types.StringValue(backup.BackupID),
			Comment:   types.StringValue(backup.Comment),
		}
		backups[i] = &backupFlattened
	}

	return backups
}
