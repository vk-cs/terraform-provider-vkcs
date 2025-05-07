package datasource_products

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func ProductsDataSourceSchema(ctx context.Context) schema.Schema {
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
								"connections": schema.SetNestedAttribute{
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"desc_i18n_key": schema.StringAttribute{
												Computed:            true,
												Description:         "Connection description i18n key.",
												MarkdownDescription: "Connection description i18n key.",
											},
											"is_require": schema.BoolAttribute{
												Computed:            true,
												Description:         "Is connection required.",
												MarkdownDescription: "Is connection required.",
											},
											"name_i18n_key": schema.StringAttribute{
												Computed:            true,
												Description:         "Connection name i18n key.",
												MarkdownDescription: "Connection name i18n key.",
											},
											"plug": schema.StringAttribute{
												Computed:            true,
												Description:         "Connection plug.",
												MarkdownDescription: "Connection plug.",
											},
											"position": schema.Int64Attribute{
												Computed:            true,
												Description:         "Connection position.",
												MarkdownDescription: "Connection position.",
											},
											"required_group": schema.StringAttribute{
												Computed:            true,
												Description:         "Connection required group.",
												MarkdownDescription: "Connection required group.",
											},
											"settings": schema.SetNestedAttribute{
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"alias": schema.StringAttribute{
															Computed:            true,
															Description:         "Setting alias.",
															MarkdownDescription: "Setting alias.",
														},
														"backref": schema.StringAttribute{
															Computed:            true,
															Description:         "Setting backref.",
															MarkdownDescription: "Setting backref.",
														},
														"default_value": schema.StringAttribute{
															Computed:            true,
															Description:         "Setting default value.",
															MarkdownDescription: "Setting default value.",
														},
														"dependencies": schema.SetNestedAttribute{
															NestedObject: schema.NestedAttributeObject{
																Attributes: map[string]schema.Attribute{
																	"anchors": schema.SetAttribute{
																		ElementType:         types.StringType,
																		Computed:            true,
																		Description:         "Dependency anchors.",
																		MarkdownDescription: "Dependency anchors.",
																	},
																	"kind": schema.StringAttribute{
																		Computed:            true,
																		Description:         "Dependency kind.",
																		MarkdownDescription: "Dependency kind.",
																	},
																},
															},
															Computed:            true,
															Description:         "Setting dependencies.",
															MarkdownDescription: "Setting dependencies.",
														},
														"is_require": schema.BoolAttribute{
															Computed:            true,
															Description:         "Is setting required.",
															MarkdownDescription: "Is setting required.",
														},
														"is_sensitive": schema.StringAttribute{
															Computed:            true,
															Description:         "Is setting sensitive.",
															MarkdownDescription: "Is setting sensitive.",
														},
														"max_value": schema.StringAttribute{
															Computed:            true,
															Description:         "Setting max value.",
															MarkdownDescription: "Setting max value.",
														},
														"min_value": schema.StringAttribute{
															Computed:            true,
															Description:         "Setting min value.",
															MarkdownDescription: "Setting min value.",
														},
														"policy": schema.SingleNestedAttribute{
															Attributes: map[string]schema.Attribute{
																"add_on_create": schema.BoolAttribute{
																	Computed:            true,
																	Description:         "Add on create policy.",
																	MarkdownDescription: "Add on create policy.",
																},
																"live_update": schema.BoolAttribute{
																	Computed:            true,
																	Description:         "Live update policy.",
																	MarkdownDescription: "Live update policy.",
																},
															},
															Computed:            true,
															Description:         "Setting policy.",
															MarkdownDescription: "Setting policy.",
														},
														"regexp": schema.StringAttribute{
															Computed:            true,
															Description:         "Setting regexp.",
															MarkdownDescription: "Setting regexp.",
														},
														"string_variation": schema.SetAttribute{
															ElementType:         types.StringType,
															Computed:            true,
															Description:         "Setting string variations.",
															MarkdownDescription: "Setting string variations.",
														},
														"validation": schema.StringAttribute{
															Computed:            true,
															Description:         "Setting validation.",
															MarkdownDescription: "Setting validation.",
														},
													},
												},
												Optional:            true,
												Computed:            true,
												Description:         "Additional connection settings.",
												MarkdownDescription: "Additional connection settings.",
											},
										},
									},
									Optional:            true,
									Computed:            true,
									Description:         "Connections settings.",
									MarkdownDescription: "Connections settings.",
								},
								"cron_tabs": schema.SetNestedAttribute{
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"name": schema.StringAttribute{
												Computed:            true,
												Description:         "Cron tab name.",
												MarkdownDescription: "Cron tab name.",
											},
											"required": schema.BoolAttribute{
												Computed:            true,
												Description:         "Is cron required.",
												MarkdownDescription: "Is cron required.",
											},
											"settings": schema.SetNestedAttribute{
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"alias": schema.StringAttribute{
															Computed:            true,
															Description:         "Setting alias.",
															MarkdownDescription: "Setting alias.",
														},
														"backref": schema.StringAttribute{
															Computed:            true,
															Description:         "Setting backref.",
															MarkdownDescription: "Setting backref.",
														},
														"default_value": schema.StringAttribute{
															Computed:            true,
															Description:         "Setting default value.",
															MarkdownDescription: "Setting default value.",
														},
														"dependencies": schema.SetNestedAttribute{
															NestedObject: schema.NestedAttributeObject{
																Attributes: map[string]schema.Attribute{
																	"anchors": schema.SetAttribute{
																		ElementType:         types.StringType,
																		Computed:            true,
																		Description:         "Dependency anchors.",
																		MarkdownDescription: "Dependency anchors.",
																	},
																	"kind": schema.StringAttribute{
																		Computed:            true,
																		Description:         "Dependency kind.",
																		MarkdownDescription: "Dependency kind.",
																	},
																},
															},
															Computed:            true,
															Description:         "Setting dependencies.",
															MarkdownDescription: "Setting dependencies.",
														},
														"is_require": schema.BoolAttribute{
															Computed:            true,
															Description:         "Is setting required.",
															MarkdownDescription: "Is setting required.",
														},
														"is_sensitive": schema.StringAttribute{
															Computed:            true,
															Description:         "Is setting sensitive.",
															MarkdownDescription: "Is setting sensitive.",
														},
														"max_value": schema.StringAttribute{
															Computed:            true,
															Description:         "Setting max value.",
															MarkdownDescription: "Setting max value.",
														},
														"min_value": schema.StringAttribute{
															Computed:            true,
															Description:         "Setting min value.",
															MarkdownDescription: "Setting min value.",
														},
														"policy": schema.SingleNestedAttribute{
															Attributes: map[string]schema.Attribute{
																"add_on_create": schema.BoolAttribute{
																	Computed:            true,
																	Description:         "Add on create policy.",
																	MarkdownDescription: "Add on create policy.",
																},
																"live_update": schema.BoolAttribute{
																	Computed:            true,
																	Description:         "Live update policy.",
																	MarkdownDescription: "Live update policy.",
																},
															},
															Computed:            true,
															Description:         "Setting policy.",
															MarkdownDescription: "Setting policy.",
														},
														"regexp": schema.StringAttribute{
															Computed:            true,
															Description:         "Setting regexp.",
															MarkdownDescription: "Setting regexp.",
														},
														"string_variation": schema.SetAttribute{
															ElementType:         types.StringType,
															Computed:            true,
															Description:         "Setting string variations.",
															MarkdownDescription: "Setting string variations.",
														},
														"validation": schema.StringAttribute{
															Computed:            true,
															Description:         "Setting validation.",
															MarkdownDescription: "Setting validation.",
														},
													},
												},
												Optional:            true,
												Computed:            true,
												Description:         "Additional cron settings.",
												MarkdownDescription: "Additional cron settings.",
											},
											"start": schema.StringAttribute{
												Computed:            true,
												Description:         "Cron schedule.",
												MarkdownDescription: "Cron schedule.",
											},
										},
									},
									Optional:            true,
									Computed:            true,
									Description:         "Cron tabs settings.",
									MarkdownDescription: "Cron tabs settings.",
								},
								"extensions": schema.SetNestedAttribute{
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"access_control": schema.StringAttribute{
												Computed:            true,
												Description:         "Extension access control mode.",
												MarkdownDescription: "Extension access control mode.",
											},
											"desc_i18n_key": schema.StringAttribute{
												Computed:            true,
												Description:         "Extension description i18n key.",
												MarkdownDescription: "Extension description i18n key.",
											},
											"name_i18n_key": schema.StringAttribute{
												Computed:            true,
												Description:         "Extension name i18n key.",
												MarkdownDescription: "Extension name i18n key.",
											},
											"persistent": schema.BoolAttribute{
												Computed:            true,
												Description:         "Is extensions persistent.",
												MarkdownDescription: "Is extensions persistent.",
											},
											"settings": schema.SetNestedAttribute{
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"alias": schema.StringAttribute{
															Computed:            true,
															Description:         "Setting alias.",
															MarkdownDescription: "Setting alias.",
														},
														"backref": schema.StringAttribute{
															Computed:            true,
															Description:         "Setting backref.",
															MarkdownDescription: "Setting backref.",
														},
														"default_value": schema.StringAttribute{
															Computed:            true,
															Description:         "Setting default value.",
															MarkdownDescription: "Setting default value.",
														},
														"dependencies": schema.SetNestedAttribute{
															NestedObject: schema.NestedAttributeObject{
																Attributes: map[string]schema.Attribute{
																	"anchors": schema.SetAttribute{
																		ElementType:         types.StringType,
																		Computed:            true,
																		Description:         "Dependency anchors.",
																		MarkdownDescription: "Dependency anchors.",
																	},
																	"kind": schema.StringAttribute{
																		Computed:            true,
																		Description:         "Dependency kind.",
																		MarkdownDescription: "Dependency kind.",
																	},
																},
															},
															Computed:            true,
															Description:         "Setting dependencies.",
															MarkdownDescription: "Setting dependencies.",
														},
														"is_require": schema.BoolAttribute{
															Computed:            true,
															Description:         "Is setting required.",
															MarkdownDescription: "Is setting required.",
														},
														"is_sensitive": schema.StringAttribute{
															Computed:            true,
															Description:         "Is setting sensitive.",
															MarkdownDescription: "Is setting sensitive.",
														},
														"max_value": schema.StringAttribute{
															Computed:            true,
															Description:         "Setting max value.",
															MarkdownDescription: "Setting max value.",
														},
														"min_value": schema.StringAttribute{
															Computed:            true,
															Description:         "Setting min value.",
															MarkdownDescription: "Setting min value.",
														},
														"policy": schema.SingleNestedAttribute{
															Attributes: map[string]schema.Attribute{
																"add_on_create": schema.BoolAttribute{
																	Computed:            true,
																	Description:         "Add on create policy.",
																	MarkdownDescription: "Add on create policy.",
																},
																"live_update": schema.BoolAttribute{
																	Computed:            true,
																	Description:         "Live update policy.",
																	MarkdownDescription: "Live update policy.",
																},
															},
															Computed:            true,
															Description:         "Setting policy.",
															MarkdownDescription: "Setting policy.",
														},
														"regexp": schema.StringAttribute{
															Computed:            true,
															Description:         "Setting regexp.",
															MarkdownDescription: "Setting regexp.",
														},
														"string_variation": schema.SetAttribute{
															ElementType:         types.StringType,
															Computed:            true,
															Description:         "Setting string variations.",
															MarkdownDescription: "Setting string variations.",
														},
														"validation": schema.StringAttribute{
															Computed:            true,
															Description:         "Setting validation.",
															MarkdownDescription: "Setting validation.",
														},
													},
												},
												Optional:            true,
												Computed:            true,
												Description:         "Additional extension settings.",
												MarkdownDescription: "Additional extension settings.",
											},
											"type": schema.StringAttribute{
												Computed:            true,
												Description:         "Extension type.",
												MarkdownDescription: "Extension type.",
											},
											"version": schema.StringAttribute{
												Computed:            true,
												Description:         "Extension version.",
												MarkdownDescription: "Extension version.",
											},
										},
									},
									Optional:            true,
									Computed:            true,
									Description:         "Extensions settings.",
									MarkdownDescription: "Extensions settings.",
								},
								"settings": schema.SetNestedAttribute{
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"alias": schema.StringAttribute{
												Computed:            true,
												Description:         "Setting alias.",
												MarkdownDescription: "Setting alias.",
											},
											"backref": schema.StringAttribute{
												Computed:            true,
												Description:         "Setting backref.",
												MarkdownDescription: "Setting backref.",
											},
											"default_value": schema.StringAttribute{
												Computed:            true,
												Description:         "Setting default value.",
												MarkdownDescription: "Setting default value.",
											},
											"dependencies": schema.SetNestedAttribute{
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"anchors": schema.SetAttribute{
															ElementType:         types.StringType,
															Computed:            true,
															Description:         "Dependency anchors.",
															MarkdownDescription: "Dependency anchors.",
														},
														"kind": schema.StringAttribute{
															Computed:            true,
															Description:         "Dependency kind.",
															MarkdownDescription: "Dependency kind.",
														},
													},
												},
												Computed:            true,
												Description:         "Setting dependencies.",
												MarkdownDescription: "Setting dependencies.",
											},
											"is_require": schema.BoolAttribute{
												Computed:            true,
												Description:         "Is setting required.",
												MarkdownDescription: "Is setting required.",
											},
											"is_sensitive": schema.StringAttribute{
												Computed:            true,
												Description:         "Is setting sensitive.",
												MarkdownDescription: "Is setting sensitive.",
											},
											"max_value": schema.StringAttribute{
												Computed:            true,
												Description:         "Setting max value.",
												MarkdownDescription: "Setting max value.",
											},
											"min_value": schema.StringAttribute{
												Computed:            true,
												Description:         "Setting min value.",
												MarkdownDescription: "Setting min value.",
											},
											"policy": schema.SingleNestedAttribute{
												Attributes: map[string]schema.Attribute{
													"add_on_create": schema.BoolAttribute{
														Computed:            true,
														Description:         "Add on create policy.",
														MarkdownDescription: "Add on create policy.",
													},
													"live_update": schema.BoolAttribute{
														Computed:            true,
														Description:         "Live update policy.",
														MarkdownDescription: "Live update policy.",
													},
												},
												Computed:            true,
												Description:         "Setting policy.",
												MarkdownDescription: "Setting policy.",
											},
											"regexp": schema.StringAttribute{
												Computed:            true,
												Description:         "Setting regexp.",
												MarkdownDescription: "Setting regexp.",
											},
											"string_variation": schema.SetAttribute{
												ElementType:         types.StringType,
												Computed:            true,
												Description:         "Setting string variations.",
												MarkdownDescription: "Setting string variations.",
											},
											"validation": schema.StringAttribute{
												Computed:            true,
												Description:         "Setting validation.",
												MarkdownDescription: "Setting validation.",
											},
										},
									},
									Optional:            true,
									Computed:            true,
									Description:         "Additional settings.",
									MarkdownDescription: "Additional settings.",
								},
								"user_accesses": schema.SetNestedAttribute{
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"settings": schema.SetNestedAttribute{
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"alias": schema.StringAttribute{
															Computed:            true,
															Description:         "Setting alias.",
															MarkdownDescription: "Setting alias.",
														},
														"backref": schema.StringAttribute{
															Computed:            true,
															Description:         "Setting backref.",
															MarkdownDescription: "Setting backref.",
														},
														"default_value": schema.StringAttribute{
															Computed:            true,
															Description:         "Setting default value.",
															MarkdownDescription: "Setting default value.",
														},
														"dependencies": schema.SetNestedAttribute{
															NestedObject: schema.NestedAttributeObject{
																Attributes: map[string]schema.Attribute{
																	"anchors": schema.SetAttribute{
																		ElementType:         types.StringType,
																		Computed:            true,
																		Description:         "Dependency anchors.",
																		MarkdownDescription: "Dependency anchors.",
																	},
																	"kind": schema.StringAttribute{
																		Computed:            true,
																		Description:         "Dependency kind.",
																		MarkdownDescription: "Dependency kind.",
																	},
																},
															},
															Computed:            true,
															Description:         "Setting dependencies.",
															MarkdownDescription: "Setting dependencies.",
														},
														"is_require": schema.BoolAttribute{
															Computed:            true,
															Description:         "Is setting required.",
															MarkdownDescription: "Is setting required.",
														},
														"is_sensitive": schema.StringAttribute{
															Computed:            true,
															Description:         "Is setting sensitive.",
															MarkdownDescription: "Is setting sensitive.",
														},
														"max_value": schema.StringAttribute{
															Computed:            true,
															Description:         "Setting max value.",
															MarkdownDescription: "Setting max value.",
														},
														"min_value": schema.StringAttribute{
															Computed:            true,
															Description:         "Setting min value.",
															MarkdownDescription: "Setting min value.",
														},
														"policy": schema.SingleNestedAttribute{
															Attributes: map[string]schema.Attribute{
																"add_on_create": schema.BoolAttribute{
																	Computed:            true,
																	Description:         "Add on create policy.",
																	MarkdownDescription: "Add on create policy.",
																},
																"live_update": schema.BoolAttribute{
																	Computed:            true,
																	Description:         "Live update policy.",
																	MarkdownDescription: "Live update policy.",
																},
															},
															Computed:            true,
															Description:         "Setting policy.",
															MarkdownDescription: "Setting policy.",
														},
														"regexp": schema.StringAttribute{
															Computed:            true,
															Description:         "Setting regexp.",
															MarkdownDescription: "Setting regexp.",
														},
														"string_variation": schema.SetAttribute{
															ElementType:         types.StringType,
															Computed:            true,
															Description:         "Setting string variations.",
															MarkdownDescription: "Setting string variations.",
														},
														"validation": schema.StringAttribute{
															Computed:            true,
															Description:         "Setting validation.",
															MarkdownDescription: "Setting validation.",
														},
													},
												},
												Optional:            true,
												Computed:            true,
												Description:         "Additional user access settings.",
												MarkdownDescription: "Additional user access settings.",
											},
										},
									},
									Optional:            true,
									Computed:            true,
									Description:         "User accesses settings.",
									MarkdownDescription: "User accesses settings.",
								},
							},
							Optional:            true,
							Computed:            true,
							Description:         "Product configuration.",
							MarkdownDescription: "Product configuration.",
						},
						"product_name": schema.StringAttribute{
							Computed:            true,
							Description:         "Product name.",
							MarkdownDescription: "Product name.",
						},
						"product_version": schema.StringAttribute{
							Computed:            true,
							Description:         "Product name.",
							MarkdownDescription: "Product name.",
						},
					},
				},
				Optional:            true,
				Computed:            true,
				Description:         "Products info.",
				MarkdownDescription: "Products info.",
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

type ProductsModel struct {
	Id       types.String `tfsdk:"id"`
	Products types.Set    `tfsdk:"products"`
	Region   types.String `tfsdk:"region"`
}

type ProductsValue struct {
	Configs        basetypes.ObjectValue `tfsdk:"configs"`
	ProductName    basetypes.StringValue `tfsdk:"product_name"`
	ProductVersion basetypes.StringValue `tfsdk:"product_version"`
}

type ConfigsValue struct {
	Connections  basetypes.SetValue `tfsdk:"connections"`
	CronTabs     basetypes.SetValue `tfsdk:"cron_tabs"`
	Extensions   basetypes.SetValue `tfsdk:"extensions"`
	Settings     basetypes.SetValue `tfsdk:"settings"`
	UserAccesses basetypes.SetValue `tfsdk:"user_accesses"`
}

type ConnectionsValue struct {
	DescI18nKey   basetypes.StringValue `tfsdk:"desc_i18n_key"`
	IsRequire     basetypes.BoolValue   `tfsdk:"is_require"`
	NameI18nKey   basetypes.StringValue `tfsdk:"name_i18n_key"`
	Plug          basetypes.StringValue `tfsdk:"plug"`
	Position      basetypes.Int64Value  `tfsdk:"position"`
	RequiredGroup basetypes.StringValue `tfsdk:"required_group"`
	Settings      basetypes.SetValue    `tfsdk:"settings"`
}

type SettingsValue struct {
	Alias           basetypes.StringValue `tfsdk:"alias"`
	Backref         basetypes.StringValue `tfsdk:"backref"`
	DefaultValue    basetypes.StringValue `tfsdk:"default_value"`
	Dependencies    basetypes.SetValue    `tfsdk:"dependencies"`
	IsRequire       basetypes.BoolValue   `tfsdk:"is_require"`
	IsSensitive     basetypes.StringValue `tfsdk:"is_sensitive"`
	MaxValue        basetypes.StringValue `tfsdk:"max_value"`
	MinValue        basetypes.StringValue `tfsdk:"min_value"`
	Policy          basetypes.ObjectValue `tfsdk:"policy"`
	Regexp          basetypes.StringValue `tfsdk:"regexp"`
	StringVariation basetypes.SetValue    `tfsdk:"string_variation"`
	Validation      basetypes.StringValue `tfsdk:"validation"`
}

type DependenciesValue struct {
	Anchors basetypes.SetValue    `tfsdk:"anchors"`
	Kind    basetypes.StringValue `tfsdk:"kind"`
}

type PolicyValue struct {
	AddOnCreate basetypes.BoolValue `tfsdk:"add_on_create"`
	LiveUpdate  basetypes.BoolValue `tfsdk:"live_update"`
}

type CronTabsValue struct {
	Name     basetypes.StringValue `tfsdk:"name"`
	Required basetypes.BoolValue   `tfsdk:"required"`
	Settings basetypes.SetValue    `tfsdk:"settings"`
	Start    basetypes.StringValue `tfsdk:"start"`
}

type ExtensionsValue struct {
	AccessControl  basetypes.StringValue `tfsdk:"access_control"`
	DescI18nKey    basetypes.StringValue `tfsdk:"desc_i18n_key"`
	NameI18nKey    basetypes.StringValue `tfsdk:"name_i18n_key"`
	Persistent     basetypes.BoolValue   `tfsdk:"persistent"`
	Settings       basetypes.SetValue    `tfsdk:"settings"`
	ExtensionsType basetypes.StringValue `tfsdk:"type"`
	Version        basetypes.StringValue `tfsdk:"version"`
}

type UserAccessesValue struct {
	Settings basetypes.SetValue `tfsdk:"settings"`
}
