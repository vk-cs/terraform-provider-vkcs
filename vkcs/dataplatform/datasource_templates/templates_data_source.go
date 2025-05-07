package datasource_templates

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func TemplatesDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "ID of the data source.",
				MarkdownDescription: "ID of the data source.",
			},
			"products": schema.SetNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"configs": schema.SingleNestedAttribute{
							Attributes: map[string]schema.Attribute{
								"common": schema.SingleNestedAttribute{
									Attributes: map[string]schema.Attribute{
										"maintenance": schema.SingleNestedAttribute{
											Attributes: map[string]schema.Attribute{
												"backup": schema.SingleNestedAttribute{
													Attributes: map[string]schema.Attribute{
														"differential": schema.SingleNestedAttribute{
															Attributes: map[string]schema.Attribute{
																"backup_name_prefix": schema.StringAttribute{
																	Computed:            true,
																	Description:         "Backup name prefix.",
																	MarkdownDescription: "Backup name prefix.",
																},
																"backup_s3_bucket_name": schema.StringAttribute{
																	Computed:            true,
																	Description:         "Backup S3 bucket name.",
																	MarkdownDescription: "Backup S3 bucket name.",
																},
																"creation_timeout": schema.Int64Attribute{
																	Computed:            true,
																	Description:         "Backup creation timeout.",
																	MarkdownDescription: "Backup creation timeout.",
																},
																"enabled": schema.BoolAttribute{
																	Computed:            true,
																	Description:         "Whether differential backup is enabled.",
																	MarkdownDescription: "Whether differential backup is enabled.",
																},
																"keep_count": schema.Int64Attribute{
																	Computed:            true,
																	Description:         "Backup keep count.",
																	MarkdownDescription: "Backup keep count.",
																},
																"keep_time": schema.Int64Attribute{
																	Computed:            true,
																	Description:         "Backup keep time.",
																	MarkdownDescription: "Backup keep time.",
																},
																"start": schema.StringAttribute{
																	Computed:            true,
																	Description:         "Backup schedule.",
																	MarkdownDescription: "Backup schedule.",
																},
															},
															Computed:            true,
															Description:         "Differential backup settings.",
															MarkdownDescription: "Differential backup settings.",
														},
														"full": schema.SingleNestedAttribute{
															Attributes: map[string]schema.Attribute{
																"backup_name_prefix": schema.StringAttribute{
																	Computed:            true,
																	Description:         "Backup name prefix.",
																	MarkdownDescription: "Backup name prefix.",
																},
																"backup_s3_bucket_name": schema.StringAttribute{
																	Computed:            true,
																	Description:         "Backup S3 bucket name.",
																	MarkdownDescription: "Backup S3 bucket name.",
																},
																"creation_timeout": schema.Int64Attribute{
																	Computed:            true,
																	Description:         "Backup creation timeout.",
																	MarkdownDescription: "Backup creation timeout.",
																},
																"enabled": schema.BoolAttribute{
																	Computed:            true,
																	Description:         "Whether full backup is enabled.",
																	MarkdownDescription: "Whether full backup is enabled.",
																},
																"keep_count": schema.Int64Attribute{
																	Computed:            true,
																	Description:         "Backup keep count.",
																	MarkdownDescription: "Backup keep count.",
																},
																"keep_time": schema.Int64Attribute{
																	Computed:            true,
																	Description:         "Backup keep time.",
																	MarkdownDescription: "Backup keep time.",
																},
																"start": schema.StringAttribute{
																	Computed:            true,
																	Description:         "Backup schedule.",
																	MarkdownDescription: "Backup schedule.",
																},
															},
															Computed:            true,
															Description:         "Full backup settings.",
															MarkdownDescription: "Full backup settings.",
														},
														"incremental": schema.SingleNestedAttribute{
															Attributes: map[string]schema.Attribute{
																"backup_name_prefix": schema.StringAttribute{
																	Computed:            true,
																	Description:         "Backup name prefix.",
																	MarkdownDescription: "Backup name prefix.",
																},
																"backup_s3_bucket_name": schema.StringAttribute{
																	Computed:            true,
																	Description:         "Backup S3 bucket name.",
																	MarkdownDescription: "Backup S3 bucket name.",
																},
																"creation_timeout": schema.Int64Attribute{
																	Computed:            true,
																	Description:         "Backup creation timeout.",
																	MarkdownDescription: "Backup creation timeout.",
																},
																"enabled": schema.BoolAttribute{
																	Computed:            true,
																	Description:         "Whether incremental backup is enabled.",
																	MarkdownDescription: "Whether incremental backup is enabled.",
																},
																"keep_count": schema.Int64Attribute{
																	Computed:            true,
																	Description:         "Backup keep count.",
																	MarkdownDescription: "Backup keep count.",
																},
																"keep_time": schema.Int64Attribute{
																	Computed:            true,
																	Description:         "Backup keep time.",
																	MarkdownDescription: "Backup keep time.",
																},
																"start": schema.StringAttribute{
																	Computed:            true,
																	Description:         "Backup schedule.",
																	MarkdownDescription: "Backup schedule.",
																},
															},
															Computed:            true,
															Description:         "Incremental backup settings.",
															MarkdownDescription: "Incremental backup settings.",
														},
													},
													Computed:            true,
													Description:         "Backup settings.",
													MarkdownDescription: "Backup settings.",
												},
												"duration": schema.Int64Attribute{
													Computed:            true,
													Description:         "Maintenance duration.",
													MarkdownDescription: "Maintenance duration.",
												},
												"start": schema.StringAttribute{
													Computed:            true,
													Description:         "Maintenance cron schedule.",
													MarkdownDescription: "Maintenance cron schedule.",
												},
											},
											Computed:            true,
											Description:         "Maintenance settings.",
											MarkdownDescription: "Maintenance settings.",
										},
									},
									Computed:            true,
									Description:         "Common configs.",
									MarkdownDescription: "Common configs.",
								},
							},
							Computed:            true,
							Description:         "Cluster template configs.",
							MarkdownDescription: "Cluster template configs.",
						},
						"created_at": schema.StringAttribute{
							Computed:            true,
							Description:         "Cluster template creation timestamp.",
							MarkdownDescription: "Cluster template creation timestamp.",
						},
						"description": schema.StringAttribute{
							Computed:            true,
							Description:         "Cluster template name.",
							MarkdownDescription: "Cluster template name.",
						},
						"id": schema.StringAttribute{
							Computed:            true,
							Description:         "Cluster template id.",
							MarkdownDescription: "Cluster template id.",
						},
						"multiaz": schema.BoolAttribute{
							Computed:            true,
							Description:         "Is multiple available zones mode enabled.",
							MarkdownDescription: "Is multiple available zones mode enabled.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							Description:         "Cluster template name.",
							MarkdownDescription: "Cluster template name.",
						},
						"pod_groups": schema.SetNestedAttribute{
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"alias": schema.StringAttribute{
										Computed:            true,
										Description:         "Alias.",
										MarkdownDescription: "Alias.",
									},
									"backref": schema.StringAttribute{
										Computed:            true,
										Description:         "Backref.",
										MarkdownDescription: "Backref.",
									},
									"cluster_template_id": schema.StringAttribute{
										Computed:            true,
										Description:         "Cluster template id.",
										MarkdownDescription: "Cluster template id.",
									},
									"count": schema.Int64Attribute{
										Computed:            true,
										Description:         "Pod count.",
										MarkdownDescription: "Pod count.",
									},
									"created_at": schema.StringAttribute{
										Computed:            true,
										Description:         "Pod group creation timestamp.",
										MarkdownDescription: "Pod group creation timestamp.",
									},
									"description": schema.StringAttribute{
										Computed:            true,
										Description:         "Pod group name.",
										MarkdownDescription: "Pod group name.",
									},
									"id": schema.StringAttribute{
										Computed:            true,
										Description:         "Pod group id.",
										MarkdownDescription: "Pod group id.",
									},
									"name": schema.StringAttribute{
										Computed:            true,
										Description:         "Pod group name.",
										MarkdownDescription: "Pod group name.",
									},
									"node_processes": schema.SetAttribute{
										ElementType:         types.StringType,
										Computed:            true,
										Description:         "Node processes.",
										MarkdownDescription: "Node processes.",
									},
									"resource": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"cpu_margin": schema.Int64Attribute{
												Computed:            true,
												Description:         "CPU margin settings.",
												MarkdownDescription: "CPU margin settings.",
											},
											"cpu_request": schema.StringAttribute{
												Computed:            true,
												Description:         "CPU request settings.",
												MarkdownDescription: "CPU request settings.",
											},
											"ram_margin": schema.Int64Attribute{
												Computed:            true,
												Description:         "RAM margin settings.",
												MarkdownDescription: "RAM margin settings.",
											},
											"ram_request": schema.StringAttribute{
												Computed:            true,
												Description:         "RAM request settings.",
												MarkdownDescription: "RAM request settings.",
											},
										},
										Computed:            true,
										Description:         "Resource settings.",
										MarkdownDescription: "Resource settings.",
									},
									"template_type": schema.StringAttribute{
										Computed:            true,
										Description:         "Template type.",
										MarkdownDescription: "Template type.",
									},
									"volumes": schema.SetNestedAttribute{
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"count": schema.Int64Attribute{
													Computed:            true,
													Description:         "Volume count.",
													MarkdownDescription: "Volume count.",
												},
												"storage": schema.StringAttribute{
													Computed:            true,
													Description:         "Storage size.",
													MarkdownDescription: "Storage size.",
												},
												"storage_class_name": schema.StringAttribute{
													Computed:            true,
													Description:         "Storage class name.",
													MarkdownDescription: "Storage class name.",
												},
											},
										},
										Computed:            true,
										Description:         "Volumes settings.",
										MarkdownDescription: "Volumes settings.",
									},
								},
							},
							Computed:            true,
							Description:         "Cluster pod groups.",
							MarkdownDescription: "Cluster pod groups.",
						},
						"presets": schema.SetNestedAttribute{
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Computed:            true,
										Description:         "Preset name.",
										MarkdownDescription: "Preset name.",
									},
									"pod_groups": schema.SetNestedAttribute{
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"count": schema.Int64Attribute{
													Computed:            true,
													Description:         "Pod count.",
													MarkdownDescription: "Pod count.",
												},
												"name": schema.StringAttribute{
													Computed:            true,
													Description:         "Pod group name.",
													MarkdownDescription: "Pod group name.",
												},
												"resource": schema.SingleNestedAttribute{
													Attributes: map[string]schema.Attribute{
														"cpu_request": schema.StringAttribute{
															Computed:            true,
															Description:         "CPU request settings.",
															MarkdownDescription: "CPU request settings.",
														},
														"ram_request": schema.StringAttribute{
															Computed:            true,
															Description:         "RAM request settings.",
															MarkdownDescription: "RAM request settings.",
														},
													},
													Computed:            true,
													Description:         "Resource settings.",
													MarkdownDescription: "Resource settings.",
												},
												"volumes": schema.SetNestedAttribute{
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"count": schema.Int64Attribute{
																Computed:            true,
																Description:         "Volume count.",
																MarkdownDescription: "Volume count.",
															},
															"storage": schema.StringAttribute{
																Computed:            true,
																Description:         "Storage size.",
																MarkdownDescription: "Storage size.",
															},
															"storage_class_name": schema.StringAttribute{
																Computed:            true,
																Description:         "Storage class name.",
																MarkdownDescription: "Storage class name.",
															},
														},
													},
													Computed:            true,
													Description:         "Volumes settings.",
													MarkdownDescription: "Volumes settings.",
												},
											},
										},
										Computed:            true,
										Description:         "Preset pod groups.",
										MarkdownDescription: "Preset pod groups.",
									},
								},
							},
							Computed:            true,
							Description:         "Presets info.",
							MarkdownDescription: "Presets info.",
						},
						"product_name": schema.StringAttribute{
							Computed:            true,
							Description:         "Product name.",
							MarkdownDescription: "Product name.",
						},
						"product_type": schema.StringAttribute{
							Computed:            true,
							Description:         "Product type.",
							MarkdownDescription: "Product type.",
						},
						"product_version": schema.StringAttribute{
							Computed:            true,
							Description:         "Product version.",
							MarkdownDescription: "Product version.",
						},
						"template_type": schema.StringAttribute{
							Computed:            true,
							Description:         "Template type.",
							MarkdownDescription: "Template type.",
						},
					},
				},
				Optional:            true,
				Computed:            true,
				Description:         "Cluster templates info.",
				MarkdownDescription: "Cluster templates info.",
			},
			"region": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The region in which to obtain the Data Platform client. If omitted, the `region` argument of the provider is used.",
				MarkdownDescription: "The region in which to obtain the Data Platform client. If omitted, the `region` argument of the provider is used.",
			},
		},
	}
}

type TemplatesModel struct {
	Id       types.String `tfsdk:"id"`
	Products types.Set    `tfsdk:"products"`
	Region   types.String `tfsdk:"region"`
}

type ProductsValue struct {
	Configs        basetypes.ObjectValue `tfsdk:"configs"`
	CreatedAt      basetypes.StringValue `tfsdk:"created_at"`
	Description    basetypes.StringValue `tfsdk:"description"`
	Id             basetypes.StringValue `tfsdk:"id"`
	Multiaz        basetypes.BoolValue   `tfsdk:"multiaz"`
	Name           basetypes.StringValue `tfsdk:"name"`
	PodGroups      basetypes.SetValue    `tfsdk:"pod_groups"`
	Presets        basetypes.SetValue    `tfsdk:"presets"`
	ProductName    basetypes.StringValue `tfsdk:"product_name"`
	ProductType    basetypes.StringValue `tfsdk:"product_type"`
	ProductVersion basetypes.StringValue `tfsdk:"product_version"`
	TemplateType   basetypes.StringValue `tfsdk:"template_type"`
}

type ConfigsValue struct {
	Common basetypes.ObjectValue `tfsdk:"common"`
}

type CommonValue struct {
	Maintenance basetypes.ObjectValue `tfsdk:"maintenance"`
}

type MaintenanceValue struct {
	Backup   basetypes.ObjectValue `tfsdk:"backup"`
	Duration basetypes.Int64Value  `tfsdk:"duration"`
	Start    basetypes.StringValue `tfsdk:"start"`
}

type BackupValue struct {
	Differential basetypes.ObjectValue `tfsdk:"differential"`
	Full         basetypes.ObjectValue `tfsdk:"full"`
	Incremental  basetypes.ObjectValue `tfsdk:"incremental"`
}

type DifferentialValue struct {
	BackupNamePrefix   basetypes.StringValue `tfsdk:"backup_name_prefix"`
	BackupS3BucketName basetypes.StringValue `tfsdk:"backup_s3_bucket_name"`
	CreationTimeout    basetypes.Int64Value  `tfsdk:"creation_timeout"`
	Enabled            basetypes.BoolValue   `tfsdk:"enabled"`
	KeepCount          basetypes.Int64Value  `tfsdk:"keep_count"`
	KeepTime           basetypes.Int64Value  `tfsdk:"keep_time"`
	Start              basetypes.StringValue `tfsdk:"start"`
}

type FullValue struct {
	BackupNamePrefix   basetypes.StringValue `tfsdk:"backup_name_prefix"`
	BackupS3BucketName basetypes.StringValue `tfsdk:"backup_s3_bucket_name"`
	CreationTimeout    basetypes.Int64Value  `tfsdk:"creation_timeout"`
	Enabled            basetypes.BoolValue   `tfsdk:"enabled"`
	KeepCount          basetypes.Int64Value  `tfsdk:"keep_count"`
	KeepTime           basetypes.Int64Value  `tfsdk:"keep_time"`
	Start              basetypes.StringValue `tfsdk:"start"`
}

type IncrementalValue struct {
	BackupNamePrefix   basetypes.StringValue `tfsdk:"backup_name_prefix"`
	BackupS3BucketName basetypes.StringValue `tfsdk:"backup_s3_bucket_name"`
	CreationTimeout    basetypes.Int64Value  `tfsdk:"creation_timeout"`
	Enabled            basetypes.BoolValue   `tfsdk:"enabled"`
	KeepCount          basetypes.Int64Value  `tfsdk:"keep_count"`
	KeepTime           basetypes.Int64Value  `tfsdk:"keep_time"`
	Start              basetypes.StringValue `tfsdk:"start"`
}

type PodGroupsValue struct {
	Alias             basetypes.StringValue `tfsdk:"alias"`
	Backref           basetypes.StringValue `tfsdk:"backref"`
	ClusterTemplateId basetypes.StringValue `tfsdk:"cluster_template_id"`
	Count             basetypes.Int64Value  `tfsdk:"count"`
	CreatedAt         basetypes.StringValue `tfsdk:"created_at"`
	Description       basetypes.StringValue `tfsdk:"description"`
	Id                basetypes.StringValue `tfsdk:"id"`
	Name              basetypes.StringValue `tfsdk:"name"`
	NodeProcesses     basetypes.SetValue    `tfsdk:"node_processes"`
	Resource          basetypes.ObjectValue `tfsdk:"resource"`
	TemplateType      basetypes.StringValue `tfsdk:"template_type"`
	Volumes           basetypes.SetValue    `tfsdk:"volumes"`
}

type ResourceValue struct {
	CpuMargin  basetypes.Int64Value  `tfsdk:"cpu_margin"`
	CpuRequest basetypes.StringValue `tfsdk:"cpu_request"`
	RamMargin  basetypes.Int64Value  `tfsdk:"ram_margin"`
	RamRequest basetypes.StringValue `tfsdk:"ram_request"`
}

type VolumesValue struct {
	Count            basetypes.Int64Value  `tfsdk:"count"`
	Storage          basetypes.StringValue `tfsdk:"storage"`
	StorageClassName basetypes.StringValue `tfsdk:"storage_class_name"`
}

type PresetsValue struct {
	Name      basetypes.StringValue `tfsdk:"name"`
	PodGroups basetypes.SetValue    `tfsdk:"pod_groups"`
}

type PodGroupsPresetsValue struct {
	Count    basetypes.Int64Value  `tfsdk:"count"`
	Name     basetypes.StringValue `tfsdk:"name"`
	Resource basetypes.ObjectValue `tfsdk:"resource"`
	Volumes  basetypes.SetValue    `tfsdk:"volumes"`
}

type ResourcePresetsValue struct {
	CpuRequest basetypes.StringValue `tfsdk:"cpu_request"`
	RamRequest basetypes.StringValue `tfsdk:"ram_request"`
}

type VolumesPresetsValue struct {
	Count            basetypes.Int64Value  `tfsdk:"count"`
	Storage          basetypes.StringValue `tfsdk:"storage"`
	StorageClassName basetypes.StringValue `tfsdk:"storage_class_name"`
}
