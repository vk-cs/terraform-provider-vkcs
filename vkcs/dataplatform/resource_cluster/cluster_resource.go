package resource_cluster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func ClusterResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"availability_zone": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Availability zone to create cluster in.",
				MarkdownDescription: "Availability zone to create cluster in.",
			},
			"cluster_template_id": schema.StringAttribute{
				Required:            true,
				Description:         "ID of the cluster template.",
				MarkdownDescription: "ID of the cluster template.",
			},
			"configs": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"features": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"volume_autoresize": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"data": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"enabled": schema.BoolAttribute{
												Optional:            true,
												Computed:            true,
												Description:         "Enables option.",
												MarkdownDescription: "Enables option.",
											},
											"max_scale_size": schema.Int64Attribute{
												Optional:            true,
												Computed:            true,
												Description:         "Maximum scale size.",
												MarkdownDescription: "Maximum scale size.",
											},
											"scale_step_size": schema.Int64Attribute{
												Optional:            true,
												Computed:            true,
												Description:         "Scale step size.",
												MarkdownDescription: "Scale step size.",
											},
											"size_scale_threshold": schema.Int64Attribute{
												Optional:            true,
												Computed:            true,
												Description:         "Size scale threshold.",
												MarkdownDescription: "Size scale threshold.",
											},
										},
										Optional:            true,
										Computed:            true,
										Description:         "Data volume options.",
										MarkdownDescription: "Data volume options.",
									},
									"wal": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"enabled": schema.BoolAttribute{
												Optional:            true,
												Computed:            true,
												Description:         "Enables option.",
												MarkdownDescription: "Enables option.",
											},
											"max_scale_size": schema.Int64Attribute{
												Optional:            true,
												Computed:            true,
												Description:         "Maximum scale size.",
												MarkdownDescription: "Maximum scale size.",
											},
											"scale_step_size": schema.Int64Attribute{
												Optional:            true,
												Computed:            true,
												Description:         "Scale step size.",
												MarkdownDescription: "Scale step size.",
											},
											"size_scale_threshold": schema.Int64Attribute{
												Optional:            true,
												Computed:            true,
												Description:         "Size scale threshold.",
												MarkdownDescription: "Size scale threshold.",
											},
										},
										Optional:            true,
										Computed:            true,
										Description:         "Wal volume options.",
										MarkdownDescription: "Wal volume options.",
									},
								},
								Optional:            true,
								Computed:            true,
								Description:         "Volume autoresize options.",
								MarkdownDescription: "Volume autoresize options.",
							},
						},
						Optional:            true,
						Computed:            true,
						Description:         "Product features.",
						MarkdownDescription: "Product features.",
					},
					"maintenance": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"backup": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"differential": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"enabled": schema.BoolAttribute{
												Computed:            true,
												Description:         "Whether full backup is enabled.",
												MarkdownDescription: "Whether full backup is enabled.",
											},
											"keep_count": schema.Int64Attribute{
												Optional:            true,
												Computed:            true,
												Description:         "Differential backup keep count.",
												MarkdownDescription: "Differential backup keep count.",
											},
											"keep_time": schema.Int64Attribute{
												Optional:            true,
												Computed:            true,
												Description:         "Differential backup keep time.",
												MarkdownDescription: "Differential backup keep time.",
											},
											"start": schema.StringAttribute{
												Required:            true,
												Description:         "Differential backup schedule.",
												MarkdownDescription: "Differential backup schedule.",
											},
										},
										Optional:            true,
										Computed:            true,
										Description:         "Differential backup settings.",
										MarkdownDescription: "Differential backup settings.",
									},
									"full": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"enabled": schema.BoolAttribute{
												Computed:            true,
												Description:         "Whether full backup is enabled.",
												MarkdownDescription: "Whether full backup is enabled.",
											},
											"keep_count": schema.Int64Attribute{
												Optional:            true,
												Computed:            true,
												Description:         "Full backup keep count.",
												MarkdownDescription: "Full backup keep count.",
											},
											"keep_time": schema.Int64Attribute{
												Optional:            true,
												Computed:            true,
												Description:         "Full backup keep time.",
												MarkdownDescription: "Full backup keep time.",
											},
											"start": schema.StringAttribute{
												Required:            true,
												Description:         "Full backup schedule.",
												MarkdownDescription: "Full backup schedule.",
											},
										},
										Optional:            true,
										Computed:            true,
										Description:         "Full backup settings.",
										MarkdownDescription: "Full backup settings.",
									},
									"incremental": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"enabled": schema.BoolAttribute{
												Computed:            true,
												Description:         "Whether full backup is enabled.",
												MarkdownDescription: "Whether full backup is enabled.",
											},
											"keep_count": schema.Int64Attribute{
												Optional:            true,
												Computed:            true,
												Description:         "Incremental backup keep count.",
												MarkdownDescription: "Incremental backup keep count.",
											},
											"keep_time": schema.Int64Attribute{
												Optional:            true,
												Computed:            true,
												Description:         "Incremental backup keep time.",
												MarkdownDescription: "Incremental backup keep time.",
											},
											"start": schema.StringAttribute{
												Required:            true,
												Description:         "Incremental backup schedule.",
												MarkdownDescription: "Incremental backup schedule.",
											},
										},
										Optional:            true,
										Computed:            true,
										Description:         "Incremental backup settings.",
										MarkdownDescription: "Incremental backup settings.",
									},
								},
								Optional:            true,
								Computed:            true,
								Description:         "Backup settings.",
								MarkdownDescription: "Backup settings.",
							},
							"cron_tabs": schema.SetNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"name": schema.StringAttribute{
											Required:            true,
											Description:         "Cron tab name.",
											MarkdownDescription: "Cron tab name.",
										},
										"required": schema.BoolAttribute{
											Computed:            true,
											Description:         "Whether cron tab is required.",
											MarkdownDescription: "Whether cron tab is required.",
										},
										"settings": schema.SetNestedAttribute{
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"alias": schema.StringAttribute{
														Required:            true,
														Description:         "Setting alias.",
														MarkdownDescription: "Setting alias.",
													},
													"value": schema.StringAttribute{
														Required:            true,
														Description:         "Setting value.",
														MarkdownDescription: "Setting value.",
													},
												},
											},
											Optional:            true,
											Computed:            true,
											Description:         "Additional cron settings.",
											MarkdownDescription: "Additional cron settings.",
										},
										"start": schema.StringAttribute{
											Optional:            true,
											Computed:            true,
											Description:         "Cron tab schedule.",
											MarkdownDescription: "Cron tab schedule.",
										},
									},
								},
								Optional:            true,
								Computed:            true,
								Description:         "Cron tabs settings.",
								MarkdownDescription: "Cron tabs settings.",
							},
							"start": schema.StringAttribute{
								Optional:            true,
								Computed:            true,
								Description:         "Maintenance cron schedule.",
								MarkdownDescription: "Maintenance cron schedule.",
							},
						},
						Required:            true,
						Description:         "Maintenance settings.",
						MarkdownDescription: "Maintenance settings.",
					},
					"settings": schema.SetNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"alias": schema.StringAttribute{
									Required:            true,
									Description:         "Setting alias.",
									MarkdownDescription: "Setting alias.",
								},
								"value": schema.StringAttribute{
									Required:            true,
									Description:         "Setting value.",
									MarkdownDescription: "Setting value.",
								},
							},
						},
						Optional:            true,
						Computed:            true,
						Description:         "Additional common settings.",
						MarkdownDescription: "Additional common settings.",
					},
					"users": schema.SetNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"access": schema.SingleNestedAttribute{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Computed:            true,
											Description:         "Access ID.",
											MarkdownDescription: "Access ID.",
										},
										"settings": schema.SetNestedAttribute{
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"alias": schema.StringAttribute{
														Required:            true,
														Description:         "Setting alias.",
														MarkdownDescription: "Setting alias.",
													},
													"value": schema.StringAttribute{
														Required:            true,
														Description:         "Setting value.",
														MarkdownDescription: "Setting value.",
													},
												},
											},
											Optional:            true,
											Computed:            true,
											Description:         "Access users settings.",
											MarkdownDescription: "Access users settings.",
										},
									},
									Optional:            true,
									Computed:            true,
									Description:         "Access settings.",
									MarkdownDescription: "Access settings.",
								},
								"created_at": schema.StringAttribute{
									Computed:            true,
									Description:         "User creation timestamp.",
									MarkdownDescription: "User creation timestamp.",
								},
								"password": schema.StringAttribute{
									Required:            true,
									Description:         "Password.",
									MarkdownDescription: "Password.",
								},
								"user": schema.StringAttribute{
									Required:            true,
									Description:         "Username.",
									MarkdownDescription: "Username.",
								},
								"user_id": schema.StringAttribute{
									Computed:            true,
									Description:         "User ID.",
									MarkdownDescription: "User ID.",
								},
							},
						},
						Optional:            true,
						Computed:            true,
						Description:         "Users settings.",
						MarkdownDescription: "Users settings.",
					},
					"warehouses": schema.SetNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"connections": schema.SetNestedAttribute{
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"created_at": schema.StringAttribute{
												Computed:            true,
												Description:         "Connection creation timestamp.",
												MarkdownDescription: "Connection creation timestamp.",
											},
											"id": schema.StringAttribute{
												Computed:            true,
												Description:         "Connection ID.",
												MarkdownDescription: "Connection ID.",
											},
											"name": schema.StringAttribute{
												Required:            true,
												Description:         "Connection name.",
												MarkdownDescription: "Connection name.",
											},
											"plug": schema.StringAttribute{
												Required:            true,
												Description:         "Connection plug.",
												MarkdownDescription: "Connection plug.",
											},
											"settings": schema.SetNestedAttribute{
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"alias": schema.StringAttribute{
															Required:            true,
															Description:         "Setting alias.",
															MarkdownDescription: "Setting alias.",
														},
														"value": schema.StringAttribute{
															Required:            true,
															Description:         "Setting value.",
															MarkdownDescription: "Setting value.",
														},
													},
												},
												Optional:            true,
												Computed:            true,
												Description:         "Additional warehouse settings.",
												MarkdownDescription: "Additional warehouse settings.",
											},
										},
									},
									Optional:            true,
									Computed:            true,
									Description:         "Warehouse connections.",
									MarkdownDescription: "Warehouse connections.",
								},
								"extensions": schema.SetNestedAttribute{
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"created_at": schema.StringAttribute{
												Computed:            true,
												Description:         "Extension creation timestamp.",
												MarkdownDescription: "Extension creation timestamp.",
											},
											"id": schema.StringAttribute{
												Computed:            true,
												Description:         "Extension ID.",
												MarkdownDescription: "Extension ID.",
											},
											"settings": schema.SetNestedAttribute{
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"alias": schema.StringAttribute{
															Required:            true,
															Description:         "Setting alias.",
															MarkdownDescription: "Setting alias.",
														},
														"value": schema.StringAttribute{
															Required:            true,
															Description:         "Setting value.",
															MarkdownDescription: "Setting value.",
														},
													},
												},
												Optional:            true,
												Computed:            true,
												Description:         "Additional extension settings.",
												MarkdownDescription: "Additional extension settings.",
											},
											"type": schema.StringAttribute{
												Required:            true,
												Description:         "Extension type.",
												MarkdownDescription: "Extension type.",
											},
											"version": schema.StringAttribute{
												Optional:            true,
												Computed:            true,
												Description:         "Extension version.",
												MarkdownDescription: "Extension version.",
											},
										},
									},
									Optional:            true,
									Computed:            true,
									Description:         "Warehouse extensions.",
									MarkdownDescription: "Warehouse extensions.",
								},
								"id": schema.StringAttribute{
									Computed:            true,
									Description:         "Warehouse ID.",
									MarkdownDescription: "Warehouse ID.",
								},
								"name": schema.StringAttribute{
									Optional:            true,
									Computed:            true,
									Description:         "Warehouse name.",
									MarkdownDescription: "Warehouse name.",
								},
								"users": schema.SetAttribute{
									ElementType:         types.StringType,
									Optional:            true,
									Computed:            true,
									Description:         "Warehouse users.",
									MarkdownDescription: "Warehouse users.",
								},
							},
						},
						Optional:            true,
						Computed:            true,
						Description:         "Warehouses settings.",
						MarkdownDescription: "Warehouses settings.",
					},
				},
				Optional:            true,
				Computed:            true,
				Description:         "Product configuration.",
				MarkdownDescription: "Product configuration.",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				Description:         "Cluster creation timestamp.",
				MarkdownDescription: "Cluster creation timestamp.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Cluster description.",
				MarkdownDescription: "Cluster description.",
			},
			"floating_ip_pool": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Floating IP pool ID.",
				MarkdownDescription: "Floating IP pool ID.",
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				Description:         "ID of the cluster.",
				MarkdownDescription: "ID of the cluster.",
			},
			"multi_az": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Enables multi az support.",
				MarkdownDescription: "Enables multi az support.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "Name of the cluster.",
				MarkdownDescription: "Name of the cluster.",
			},
			"network_id": schema.StringAttribute{
				Required:            true,
				Description:         "ID of the cluster network.",
				MarkdownDescription: "ID of the cluster network.",
			},
			"pod_groups": schema.SetNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"count": schema.Int64Attribute{
							Optional:            true,
							Computed:            true,
							Description:         "Pod count.",
							MarkdownDescription: "Pod count.",
						},
						"floating_ip_pool": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "Floating IP pool ID.",
							MarkdownDescription: "Floating IP pool ID.",
						},
						"id": schema.StringAttribute{
							Computed:            true,
							Description:         "Pod group ID.",
							MarkdownDescription: "Pod group ID.",
						},
						"node_processes": schema.SetAttribute{
							ElementType:         types.StringType,
							Optional:            true,
							Computed:            true,
							Description:         "Node processes.",
							MarkdownDescription: "Node processes.",
						},
						"pod_group_template_id": schema.StringAttribute{
							Required:            true,
							Description:         "Pod group template ID.",
							MarkdownDescription: "Pod group template ID.",
						},
						"resource": schema.SingleNestedAttribute{
							Attributes: map[string]schema.Attribute{
								"cpu_limit": schema.StringAttribute{
									Computed:            true,
									Description:         "CPU limit.",
									MarkdownDescription: "CPU limit.",
								},
								"cpu_request": schema.StringAttribute{
									Optional:            true,
									Computed:            true,
									Description:         "CPU request settings.",
									MarkdownDescription: "CPU request settings.",
								},
								"ram_limit": schema.StringAttribute{
									Computed:            true,
									Description:         "RAM limit settings.",
									MarkdownDescription: "RAM limit settings.",
								},
								"ram_request": schema.StringAttribute{
									Optional:            true,
									Computed:            true,
									Description:         "RAM request settings.",
									MarkdownDescription: "RAM request settings.",
								},
							},
							Optional:            true,
							Computed:            true,
							Description:         "Resource request settings.",
							MarkdownDescription: "Resource request settings.",
						},
						"volumes": schema.SetNestedAttribute{
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"count": schema.Int64Attribute{
										Required:            true,
										Description:         "Volume count.",
										MarkdownDescription: "Volume count.",
									},
									"storage": schema.StringAttribute{
										Required:            true,
										Description:         "Storage size.",
										MarkdownDescription: "Storage size.",
									},
									"storage_class_name": schema.StringAttribute{
										Required:            true,
										Description:         "Storage class name.",
										MarkdownDescription: "Storage class name.",
									},
									"type": schema.StringAttribute{
										Required:            true,
										Description:         "Volume type.",
										MarkdownDescription: "Volume type.",
									},
								},
							},
							Optional:            true,
							Computed:            true,
							Description:         "Volumes settings.",
							MarkdownDescription: "Volumes settings.",
						},
					},
				},
				Optional:            true,
				Computed:            true,
				Description:         "Cluster pod groups.",
				MarkdownDescription: "Cluster pod groups.",
			},
			"product_name": schema.StringAttribute{
				Required:            true,
				Description:         "Name of the product.",
				MarkdownDescription: "Name of the product.",
			},
			"product_version": schema.StringAttribute{
				Required:            true,
				Description:         "Version of the product.",
				MarkdownDescription: "Version of the product.",
			},
			"region": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The region in which to obtain the Data Platform client. If omitted, the `region` argument of the provider is used. Changing this creates a new resource.",
				MarkdownDescription: "The region in which to obtain the Data Platform client. If omitted, the `region` argument of the provider is used. Changing this creates a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"stack_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "ID of the cluster stack.",
				MarkdownDescription: "ID of the cluster stack.",
			},
			"subnet_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "ID of the cluster subnet.",
				MarkdownDescription: "ID of the cluster subnet.",
			},
		},
	}
}

type ClusterModel struct {
	AvailabilityZone  types.String `tfsdk:"availability_zone"`
	ClusterTemplateId types.String `tfsdk:"cluster_template_id"`
	Configs           ConfigsValue `tfsdk:"configs"`
	CreatedAt         types.String `tfsdk:"created_at"`
	Description       types.String `tfsdk:"description"`
	FloatingIpPool    types.String `tfsdk:"floating_ip_pool"`
	Id                types.Int64  `tfsdk:"id"`
	MultiAz           types.Bool   `tfsdk:"multi_az"`
	Name              types.String `tfsdk:"name"`
	NetworkId         types.String `tfsdk:"network_id"`
	PodGroups         types.Set    `tfsdk:"pod_groups"`
	ProductName       types.String `tfsdk:"product_name"`
	ProductVersion    types.String `tfsdk:"product_version"`
	Region            types.String `tfsdk:"region"`
	StackId           types.String `tfsdk:"stack_id"`
	SubnetId          types.String `tfsdk:"subnet_id"`
}

type ConfigsValue struct {
	Features    basetypes.ObjectValue `tfsdk:"features"`
	Maintenance basetypes.ObjectValue `tfsdk:"maintenance"`
	Settings    basetypes.SetValue    `tfsdk:"settings"`
	Users       basetypes.SetValue    `tfsdk:"users"`
	Warehouses  basetypes.SetValue    `tfsdk:"warehouses"`
}

type FeaturesValue struct {
	VolumeAutoresize basetypes.ObjectValue `tfsdk:"volume_autoresize"`
}

type VolumeAutoresizeValue struct {
	Data basetypes.ObjectValue `tfsdk:"data"`
	Wal  basetypes.ObjectValue `tfsdk:"wal"`
}

type DataValue struct {
	Enabled            basetypes.BoolValue  `tfsdk:"enabled"`
	MaxScaleSize       basetypes.Int64Value `tfsdk:"max_scale_size"`
	ScaleStepSize      basetypes.Int64Value `tfsdk:"scale_step_size"`
	SizeScaleThreshold basetypes.Int64Value `tfsdk:"size_scale_threshold"`
}

type WalValue struct {
	Enabled            basetypes.BoolValue  `tfsdk:"enabled"`
	MaxScaleSize       basetypes.Int64Value `tfsdk:"max_scale_size"`
	ScaleStepSize      basetypes.Int64Value `tfsdk:"scale_step_size"`
	SizeScaleThreshold basetypes.Int64Value `tfsdk:"size_scale_threshold"`
}

type MaintenanceValue struct {
	Backup   basetypes.ObjectValue `tfsdk:"backup"`
	CronTabs basetypes.SetValue    `tfsdk:"cron_tabs"`
	Start    basetypes.StringValue `tfsdk:"start"`
}

type BackupValue struct {
	Differential basetypes.ObjectValue `tfsdk:"differential"`
	Full         basetypes.ObjectValue `tfsdk:"full"`
	Incremental  basetypes.ObjectValue `tfsdk:"incremental"`
}

type DifferentialValue struct {
	Enabled   basetypes.BoolValue   `tfsdk:"enabled"`
	KeepCount basetypes.Int64Value  `tfsdk:"keep_count"`
	KeepTime  basetypes.Int64Value  `tfsdk:"keep_time"`
	Start     basetypes.StringValue `tfsdk:"start"`
}

type FullValue struct {
	Enabled   basetypes.BoolValue   `tfsdk:"enabled"`
	KeepCount basetypes.Int64Value  `tfsdk:"keep_count"`
	KeepTime  basetypes.Int64Value  `tfsdk:"keep_time"`
	Start     basetypes.StringValue `tfsdk:"start"`
}

type IncrementalValue struct {
	Enabled   basetypes.BoolValue   `tfsdk:"enabled"`
	KeepCount basetypes.Int64Value  `tfsdk:"keep_count"`
	KeepTime  basetypes.Int64Value  `tfsdk:"keep_time"`
	Start     basetypes.StringValue `tfsdk:"start"`
}

type CronTabsValue struct {
	Name     basetypes.StringValue `tfsdk:"name"`
	Required basetypes.BoolValue   `tfsdk:"required"`
	Settings basetypes.SetValue    `tfsdk:"settings"`
	Start    basetypes.StringValue `tfsdk:"start"`
}

type SettingsValue struct {
	Alias basetypes.StringValue `tfsdk:"alias"`
	Value basetypes.StringValue `tfsdk:"value"`
}

type UsersValue struct {
	Access    basetypes.ObjectValue `tfsdk:"access"`
	CreatedAt basetypes.StringValue `tfsdk:"created_at"`
	Password  basetypes.StringValue `tfsdk:"password"`
	User      basetypes.StringValue `tfsdk:"user"`
	UserId    basetypes.StringValue `tfsdk:"user_id"`
}

type AccessValue struct {
	Id       basetypes.StringValue `tfsdk:"id"`
	Settings basetypes.SetValue    `tfsdk:"settings"`
}

type WarehousesValue struct {
	Connections basetypes.SetValue    `tfsdk:"connections"`
	Extensions  basetypes.SetValue    `tfsdk:"extensions"`
	Id          basetypes.StringValue `tfsdk:"id"`
	Name        basetypes.StringValue `tfsdk:"name"`
	Users       basetypes.SetValue    `tfsdk:"users"`
}

type ConnectionsValue struct {
	CreatedAt basetypes.StringValue `tfsdk:"created_at"`
	Id        basetypes.StringValue `tfsdk:"id"`
	Name      basetypes.StringValue `tfsdk:"name"`
	Plug      basetypes.StringValue `tfsdk:"plug"`
	Settings  basetypes.SetValue    `tfsdk:"settings"`
}

type ExtensionsValue struct {
	CreatedAt      basetypes.StringValue `tfsdk:"created_at"`
	Id             basetypes.StringValue `tfsdk:"id"`
	Settings       basetypes.SetValue    `tfsdk:"settings"`
	ExtensionsType basetypes.StringValue `tfsdk:"type"`
	Version        basetypes.StringValue `tfsdk:"version"`
}

type PodGroupsValue struct {
	Count              basetypes.Int64Value  `tfsdk:"count"`
	FloatingIpPool     basetypes.StringValue `tfsdk:"floating_ip_pool"`
	Id                 basetypes.StringValue `tfsdk:"id"`
	NodeProcesses      basetypes.SetValue    `tfsdk:"node_processes"`
	PodGroupTemplateId basetypes.StringValue `tfsdk:"pod_group_template_id"`
	Resource           basetypes.ObjectValue `tfsdk:"resource"`
	Volumes            basetypes.SetValue    `tfsdk:"volumes"`
}

type ResourceValue struct {
	CpuLimit   basetypes.StringValue `tfsdk:"cpu_limit"`
	CpuRequest basetypes.StringValue `tfsdk:"cpu_request"`
	RamLimit   basetypes.StringValue `tfsdk:"ram_limit"`
	RamRequest basetypes.StringValue `tfsdk:"ram_request"`
}

type VolumesValue struct {
	Count            basetypes.Int64Value  `tfsdk:"count"`
	Storage          basetypes.StringValue `tfsdk:"storage"`
	StorageClassName basetypes.StringValue `tfsdk:"storage_class_name"`
	VolumesType      basetypes.StringValue `tfsdk:"type"`
}
