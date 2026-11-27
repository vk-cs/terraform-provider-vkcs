package cdn

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/cdn/resource_resource"
)

type resourceResourceStateV0 struct {
	Active             types.Bool                            `tfsdk:"active"`
	Cname              types.String                          `tfsdk:"cname"`
	Id                 types.Int64                           `tfsdk:"id"`
	Options            types.Object                          `tfsdk:"options"`
	OriginGroup        types.Int64                           `tfsdk:"origin_group"`
	OriginProtocol     types.String                          `tfsdk:"origin_protocol"`
	PresetApplied      types.Bool                            `tfsdk:"preset_applied"`
	Region             types.String                          `tfsdk:"region"`
	SecondaryHostnames types.List                            `tfsdk:"secondary_hostnames"`
	Shielding          resource_resource.ShieldingValue      `tfsdk:"shielding"`
	SslCertificate     resource_resource.SslCertificateValue `tfsdk:"ssl_certificate"`
	Status             types.String                          `tfsdk:"status"`
	Timeouts           timeouts.Value                        `tfsdk:"timeouts"`
	VpEnabled          types.Bool                            `tfsdk:"vp_enabled"`
}

func (r *resourceResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	priorSchema := resourceResourceSchemaV0(ctx)

	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema: &priorSchema,
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var prior resourceResourceStateV0
				resp.Diagnostics.Append(req.State.Get(ctx, &prior)...)
				if resp.Diagnostics.HasError() {
					return
				}

				options, diags := upgradeResourceOptionsV0(ctx, prior.Options)
				resp.Diagnostics.Append(diags...)
				if resp.Diagnostics.HasError() {
					return
				}

				upgraded := resource_resource.ResourceModel{
					Active:             prior.Active,
					Cname:              prior.Cname,
					Id:                 prior.Id,
					Options:            options,
					OriginGroup:        prior.OriginGroup,
					OriginProtocol:     prior.OriginProtocol,
					PresetApplied:      prior.PresetApplied,
					Region:             prior.Region,
					SecondaryHostnames: prior.SecondaryHostnames,
					Shielding:          prior.Shielding,
					SslCertificate:     prior.SslCertificate,
					Status:             prior.Status,
					Timeouts:           prior.Timeouts,
					VpEnabled:          prior.VpEnabled,
				}

				resp.Diagnostics.Append(resp.State.Set(ctx, &upgraded)...)
			},
		},
	}
}

func resourceResourceSchemaV0(ctx context.Context) schema.Schema {
	priorSchema := resource_resource.ResourceResourceSchema(ctx)
	options := priorSchema.Attributes["options"].(schema.SingleNestedAttribute)
	options.CustomType = nil
	delete(options.Attributes, "static_response_headers")

	tlsVersions := options.Attributes["tls_versions"].(schema.SingleNestedAttribute)
	tlsVersions.CustomType = nil
	tlsVersions.Attributes["value"] = schema.ListAttribute{
		ElementType: types.StringType,
		Optional:    true,
		Computed:    true,
	}
	options.Attributes["tls_versions"] = tlsVersions
	priorSchema.Attributes["options"] = options
	priorSchema.Version = 0

	return priorSchema
}

func upgradeResourceOptionsV0(ctx context.Context, prior types.Object) (resource_resource.OptionsValue, diag.Diagnostics) {
	if prior.IsNull() {
		return resource_resource.NewOptionsValueNull(), nil
	}
	if prior.IsUnknown() {
		return resource_resource.NewOptionsValueUnknown(), nil
	}

	attributes := prior.Attributes()
	attributes["static_response_headers"] = types.ObjectNull(resource_resource.StaticResponseHeadersValue{}.AttributeTypes(ctx))
	priorTls, ok := attributes["tls_versions"].(types.Object)
	if !ok {
		var diags diag.Diagnostics
		diags.AddError("Invalid Prior State", "The prior CDN options state is missing the `tls_versions` object.")
		return resource_resource.NewOptionsValueUnknown(), diags
	}

	var tlsVersionsValue types.Object
	switch {
	case priorTls.IsNull():
		tlsVersionsValue = types.ObjectNull(resource_resource.TlsVersionsValue{}.AttributeTypes(ctx))
	case priorTls.IsUnknown():
		tlsVersionsValue = types.ObjectUnknown(resource_resource.TlsVersionsValue{}.AttributeTypes(ctx))
	default:
		tlsAttributes := priorTls.Attributes()
		priorValue, ok := tlsAttributes["value"].(types.List)
		if !ok {
			var diags diag.Diagnostics
			diags.AddError("Invalid Prior State", "The prior CDN TLS versions state does not contain a list value.")
			return resource_resource.NewOptionsValueUnknown(), diags
		}

		value, diags := upgradeTLSVersionsListToSet(priorValue)
		if diags.HasError() {
			return resource_resource.NewOptionsValueUnknown(), diags
		}

		tlsVersionsValue, diags = types.ObjectValue(resource_resource.TlsVersionsValue{}.AttributeTypes(ctx), map[string]attr.Value{
			"enabled": tlsAttributes["enabled"],
			"value":   value,
		})
		if diags.HasError() {
			return resource_resource.NewOptionsValueUnknown(), diags
		}
	}

	currentAttributeTypes := resource_resource.OptionsValue{}.AttributeTypes(ctx)
	attributes["tls_versions"] = tlsVersionsValue
	for name, value := range attributes {
		if name == "tls_versions" {
			continue
		}

		terraformValue, err := value.ToTerraformValue(ctx)
		if err != nil {
			var diags diag.Diagnostics
			diags.AddError("Invalid Prior State", fmt.Sprintf("Unable to read prior CDN option %q: %s", name, err))
			return resource_resource.NewOptionsValueUnknown(), diags
		}

		upgradedValue, err := currentAttributeTypes[name].ValueFromTerraform(ctx, terraformValue)
		if err != nil {
			var diags diag.Diagnostics
			diags.AddError("Invalid Prior State", fmt.Sprintf("Unable to upgrade CDN option %q: %s", name, err))
			return resource_resource.NewOptionsValueUnknown(), diags
		}
		attributes[name] = upgradedValue
	}

	upgradedObject, diags := types.ObjectValue(currentAttributeTypes, attributes)
	if diags.HasError() {
		return resource_resource.NewOptionsValueUnknown(), diags
	}

	upgraded, diags := resource_resource.OptionsType{}.ValueFromObject(ctx, upgradedObject)
	if diags.HasError() {
		return resource_resource.NewOptionsValueUnknown(), diags
	}
	return upgraded.(resource_resource.OptionsValue), diags
}

func upgradeTLSVersionsListToSet(prior types.List) (types.Set, diag.Diagnostics) {
	switch {
	case prior.IsNull():
		return types.SetNull(types.StringType), nil
	case prior.IsUnknown():
		return types.SetUnknown(types.StringType), nil
	default:
		return types.SetValue(types.StringType, prior.Elements())
	}
}
