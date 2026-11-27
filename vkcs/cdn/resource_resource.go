package cdn

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/cdn/resource_resource"
	resource_validators "github.com/vk-cs/terraform-provider-vkcs/vkcs/cdn/resource_resource/validators"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	fwutils "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/framework/utils"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/cdn/v1/resources"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

var (
	_ resource.Resource                     = (*resourceResource)(nil)
	_ resource.ResourceWithConfigure        = (*resourceResource)(nil)
	_ resource.ResourceWithConfigValidators = (*resourceResource)(nil)
	_ resource.ResourceWithImportState      = (*resourceResource)(nil)
	_ resource.ResourceWithValidateConfig   = (*resourceResource)(nil)
	_ resource.ResourceWithModifyPlan       = (*resourceResource)(nil)
	_ resource.ResourceWithUpgradeState     = (*resourceResource)(nil)
)

func NewResourceResource() resource.Resource {
	return &resourceResource{}
}

type resourceResource struct {
	config clients.Config
}

func (r *resourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cdn_resource"
}

func (r *resourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_resource.ResourceResourceSchema(ctx)
	resp.Schema.Version = 1
}

func (r *resourceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.config = req.ProviderData.(clients.Config)
}

func (r *resourceResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resource_validators.ConflictingEnabled(
			path.MatchRoot("options").AtName("fetch_compressed"),
			path.MatchRoot("options").AtName("gzip_on"),
		),
		resource_validators.ConflictingEnabled(
			path.MatchRoot("options").AtName("slice"),
			path.MatchRoot("options").AtName("gzip_on"),
		),
		resource_validators.ConflictingEnabled(
			path.MatchRoot("options").AtName("fetch_compressed"),
			path.MatchRoot("options").AtName("gzip_compression"),
		),
		resource_validators.ConflictingEnabled(
			path.MatchRoot("options").AtName("slice"),
			path.MatchRoot("options").AtName("gzip_compression"),
		),
		resource_validators.ConflictingEnabled(
			path.MatchRoot("options").AtName("forward_host_header"),
			path.MatchRoot("options").AtName("host_header"),
		),
		resource_validators.ConflictingEnabled(
			path.MatchRoot("options").AtName("ignore_query_string"),
			path.MatchRoot("options").AtName("query_params_blacklist"),
			path.MatchRoot("options").AtName("query_params_whitelist"),
		),
	}
}

func (r *resourceResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config resource_resource.ResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.Options.IsNull() && !config.Options.IsUnknown() {
		resp.Diagnostics.Append(validateStaticResponseHeaders(ctx, config.Options.StaticResponseHeaders)...)
	}

	sslCertType := config.SslCertificate.SslCertificateType.ValueString()
	if sslCertType == string(resource_resource.SslCertificateProviderTypeOwn) && config.SslCertificate.Id.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("ssl_certificate"),
			"Missing Attribute Value",
			"`ssl_certificate.id` must be configured when `ssl_certificate.type` is \"own\"",
		)
	}

	shieldingEnabled := config.Shielding.Enabled.ValueBool()
	if shieldingEnabled && config.Shielding.PopId.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("shielding"),
			"Missing Attribute Value",
			"`shielding.pop_id` must be configured when `shielding.enabled` is \"true\"",
		)
	}
}

func selectConfiguredResponseHeaderOption(
	ctx context.Context,
	plan *resource_resource.OptionsValue,
	config resource_resource.OptionsValue,
	previous *resource_resource.OptionsValue,
) {
	if plan == nil || plan.IsNull() || plan.IsUnknown() || config.IsNull() {
		return
	}
	if config.IsUnknown() {
		return
	}

	legacyConfigured := !config.StaticHeaders.IsNull()
	responseConfigured := !config.StaticResponseHeaders.IsNull()

	if responseConfigured {
		if previous == nil || !isKnownOption(previous.StaticResponseHeaders) ||
			!isKnownOption(plan.StaticResponseHeaders) || !plan.StaticResponseHeaders.Equal(previous.StaticResponseHeaders) {
			// The CDN API can return the deprecated staticHeaders alias after
			// receiving static_response_headers. Keep that alias unknown when
			// the canonical option changes so the API read-back is consistent.
			plan.StaticHeaders = types.ObjectUnknown(resource_resource.StaticHeadersValue{}.AttributeTypes(ctx))
		}
	}

	if legacyConfigured {
		if previous == nil || !isKnownOption(previous.StaticHeaders) ||
			!isKnownOption(plan.StaticHeaders) || !plan.StaticHeaders.Equal(previous.StaticHeaders) {
			// The CDN API can return the canonical static_response_headers form
			// after receiving the deprecated staticHeaders field. Keep the plan
			// unknown when that legacy setting changes so the API's read-back is
			// consistent with the plan, while preserving known state on no-op plans.
			plan.StaticResponseHeaders = types.ObjectUnknown(resource_resource.StaticResponseHeadersValue{}.AttributeTypes(ctx))
		}
	}
}

// preserveUnconfiguredOptions restores prior API values for omitted fields
// during no-op planning, preventing computed option defaults from producing a
// perpetual update.
func preserveUnconfiguredOptions(
	ctx context.Context,
	plan *resource_resource.OptionsValue,
	config resource_resource.OptionsValue,
	previous resource_resource.OptionsValue,
) diag.Diagnostics {
	var diags diag.Diagnostics
	if plan == nil || config.IsUnknown() {
		return diags
	}
	if config.IsNull() {
		if plan.IsNull() || plan.IsUnknown() {
			*plan = previous
		}
		return diags
	}
	if plan.IsNull() || plan.IsUnknown() || previous.IsNull() || previous.IsUnknown() {
		return diags
	}

	configValue, d := config.ToObjectValue(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}
	planValue, d := plan.ToObjectValue(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}
	previousValue, d := previous.ToObjectValue(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	mergedValue, changed, d := preserveUnconfiguredObjectAttributes(ctx, configValue, planValue, previousValue)
	diags.Append(d...)
	if diags.HasError() || !changed {
		return diags
	}

	mergedOptions, d := resource_resource.OptionsType{}.ValueFromObject(ctx, mergedValue)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}
	*plan = mergedOptions.(resource_resource.OptionsValue)
	return diags
}

func preserveUnconfiguredShielding(ctx context.Context, plan *resource_resource.ShieldingValue, config, previous resource_resource.ShieldingValue) diag.Diagnostics {
	if config.IsNull() && (plan.IsNull() || plan.IsUnknown()) {
		*plan = previous
		return nil
	}
	if config.IsNull() || config.IsUnknown() || plan.IsNull() || plan.IsUnknown() || previous.IsNull() || previous.IsUnknown() {
		return nil
	}

	configObject, diags := config.ToObjectValue(ctx)
	if diags.HasError() {
		return diags
	}
	planObject, diags := plan.ToObjectValue(ctx)
	if diags.HasError() {
		return diags
	}
	previousObject, diags := previous.ToObjectValue(ctx)
	if diags.HasError() {
		return diags
	}
	mergedObject, changed, diags := preserveUnconfiguredObjectAttributes(ctx, configObject, planObject, previousObject)
	if diags.HasError() || !changed {
		return diags
	}
	mergedValue, diags := resource_resource.ShieldingType{}.ValueFromObject(ctx, mergedObject)
	if diags.HasError() {
		return diags
	}
	*plan = mergedValue.(resource_resource.ShieldingValue)
	return diags
}

func preserveUnconfiguredSslCertificate(ctx context.Context, plan *resource_resource.SslCertificateValue, config, previous resource_resource.SslCertificateValue) diag.Diagnostics {
	if config.IsNull() && (plan.IsNull() || plan.IsUnknown()) {
		*plan = previous
		return nil
	}
	if config.IsNull() || config.IsUnknown() || plan.IsNull() || plan.IsUnknown() || previous.IsNull() || previous.IsUnknown() {
		return nil
	}

	configObject, diags := config.ToObjectValue(ctx)
	if diags.HasError() {
		return diags
	}
	planObject, diags := plan.ToObjectValue(ctx)
	if diags.HasError() {
		return diags
	}
	previousObject, diags := previous.ToObjectValue(ctx)
	if diags.HasError() {
		return diags
	}
	mergedObject, changed, diags := preserveUnconfiguredObjectAttributes(ctx, configObject, planObject, previousObject)
	if diags.HasError() {
		return diags
	}
	attributes := mergedObject.Attributes()
	previousAttributes := previousObject.Attributes()
	status, statusExists := attributes["status"]
	previousStatus, previousStatusExists := previousAttributes["status"]
	if statusExists && previousStatusExists && (status.IsNull() || status.IsUnknown()) && !previousStatus.IsNull() && !previousStatus.IsUnknown() {
		attributes["status"] = previousStatus
		changed = true
	}
	if !changed {
		return diags
	}
	mergedObject, diags = types.ObjectValue(mergedObject.AttributeTypes(ctx), attributes)
	if diags.HasError() {
		return diags
	}
	mergedValue, diags := resource_resource.SslCertificateType{}.ValueFromObject(ctx, mergedObject)
	if diags.HasError() {
		return diags
	}
	*plan = mergedValue.(resource_resource.SslCertificateValue)
	return diags
}

func preserveUnconfiguredObjectAttributes(
	ctx context.Context,
	config types.Object,
	plan types.Object,
	previous types.Object,
) (types.Object, bool, diag.Diagnostics) {
	var diags diag.Diagnostics
	if config.IsNull() || config.IsUnknown() || plan.IsNull() || plan.IsUnknown() || previous.IsNull() || previous.IsUnknown() {
		return plan, false, diags
	}

	configAttributes := config.Attributes()
	planAttributes := plan.Attributes()
	previousAttributes := previous.Attributes()
	changed := false

	for name, planValue := range planAttributes {
		configValue, ok := configAttributes[name]
		if !ok {
			continue
		}
		previousValue, ok := previousAttributes[name]
		if !ok {
			continue
		}

		preserve := preserveUnconfiguredValue
		if name == "edge_cache_settings" {
			preserve = preserveEdgeCacheSettingsValue
		}
		preservedValue, valueChanged, d := preserve(ctx, configValue, planValue, previousValue)
		diags.Append(d...)
		if diags.HasError() {
			return plan, false, diags
		}
		if valueChanged {
			planAttributes[name] = preservedValue
			changed = true
		}
	}

	if !changed {
		return plan, false, diags
	}

	merged, d := types.ObjectValue(plan.AttributeTypes(ctx), planAttributes)
	diags.Append(d...)
	return merged, true, diags
}

func preserveEdgeCacheSettingsValue(ctx context.Context, config, plan, previous attr.Value) (attr.Value, bool, diag.Diagnostics) {
	if config.IsNull() && (plan.IsNull() || plan.IsUnknown()) {
		return previous, true, nil
	}
	if config.IsUnknown() || plan.IsNull() || plan.IsUnknown() || previous.IsNull() || previous.IsUnknown() {
		return plan, false, nil
	}

	configObject, configIsObject := config.(types.Object)
	planObject, planIsObject := plan.(types.Object)
	previousObject, previousIsObject := previous.(types.Object)
	if !configIsObject || !planIsObject || !previousIsObject {
		return plan, false, nil
	}

	configAttributes := configObject.Attributes()
	planAttributes := planObject.Attributes()
	previousAttributes := previousObject.Attributes()
	defaultValue, defaultExists := configAttributes["default"]
	valueValue, valueExists := configAttributes["value"]
	defaultConfigured := defaultExists && !defaultValue.IsNull() && !defaultValue.IsUnknown()
	valueConfigured := valueExists && !valueValue.IsNull() && !valueValue.IsUnknown()
	changed := false

	for name, planValue := range planAttributes {
		configValue, ok := configAttributes[name]
		if !ok {
			continue
		}
		previousValue, ok := previousAttributes[name]
		if !ok {
			continue
		}
		// Edge cache API input rejects a non-empty `value` alongside a configured
		// `default` (and vice versa). Do not restore that stale counterpart.
		if name == "value" && defaultConfigured && edgeCacheStringIsSet(previousValue) {
			continue
		}
		if name == "default" && valueConfigured && edgeCacheStringIsSet(previousValue) {
			continue
		}
		preservedValue, valueChanged, diags := preserveUnconfiguredValue(ctx, configValue, planValue, previousValue)
		if diags.HasError() {
			return plan, false, diags
		}
		if valueChanged {
			planAttributes[name] = preservedValue
			changed = true
		}
	}

	if !changed {
		return plan, false, nil
	}
	value, diags := types.ObjectValue(planObject.AttributeTypes(ctx), planAttributes)
	return value, true, diags
}

func edgeCacheStringIsSet(value attr.Value) bool {
	stringValue, ok := value.(types.String)
	return ok && !stringValue.IsNull() && !stringValue.IsUnknown() && stringValue.ValueString() != ""
}

func preserveUnconfiguredValue(ctx context.Context, config, plan, previous attr.Value) (attr.Value, bool, diag.Diagnostics) {
	if config.IsNull() && (plan.IsNull() || plan.IsUnknown()) {
		return previous, true, nil
	}
	if config.IsUnknown() || plan.IsNull() || plan.IsUnknown() || previous.IsNull() || previous.IsUnknown() {
		return plan, false, nil
	}

	configObject, configIsObject := config.(types.Object)
	planObject, planIsObject := plan.(types.Object)
	previousObject, previousIsObject := previous.(types.Object)
	if !configIsObject || !planIsObject || !previousIsObject {
		return plan, false, nil
	}

	merged, changed, diags := preserveUnconfiguredObjectAttributes(ctx, configObject, planObject, previousObject)
	if diags.HasError() || !changed {
		return plan, changed, diags
	}
	return merged, true, diags
}

func isKnownOption(value types.Object) bool {
	return !value.IsNull() && !value.IsUnknown()
}

func responseHeaderOptionsForAPIRequest(
	ctx context.Context,
	options resource_resource.OptionsValue,
	config resource_resource.OptionsValue,
) resource_resource.OptionsValue {
	if options.IsNull() || options.IsUnknown() || config.IsNull() || config.IsUnknown() {
		return options
	}

	legacyConfigured := !config.StaticHeaders.IsNull()
	responseConfigured := !config.StaticResponseHeaders.IsNull()

	if legacyConfigured {
		options.StaticResponseHeaders = types.ObjectNull(resource_resource.StaticResponseHeadersValue{}.AttributeTypes(ctx))
	} else if responseConfigured || isKnownOption(options.StaticResponseHeaders) {
		options.StaticHeaders = types.ObjectNull(resource_resource.StaticHeadersValue{}.AttributeTypes(ctx))
	}

	return options
}

func mergeResourceOptions(planned, prior *resources.ResourceOptions) *resources.ResourceOptions {
	if planned == nil {
		return prior
	}
	if prior == nil {
		return planned
	}

	return &resources.ResourceOptions{
		AllowedHttpMethods:    preferResourceOption(planned.AllowedHttpMethods, prior.AllowedHttpMethods),
		BrotliCompression:     preferResourceOption(planned.BrotliCompression, prior.BrotliCompression),
		BrowserCacheSettings:  preferResourceOption(planned.BrowserCacheSettings, prior.BrowserCacheSettings),
		CORS:                  preferResourceOption(planned.CORS, prior.CORS),
		EdgeCacheSettings:     preferResourceOption(planned.EdgeCacheSettings, prior.EdgeCacheSettings),
		FetchCompressed:       preferResourceOption(planned.FetchCompressed, prior.FetchCompressed),
		ForceReturn:           preferResourceOption(planned.ForceReturn, prior.ForceReturn),
		ForwardHostHeader:     preferResourceOption(planned.ForwardHostHeader, prior.ForwardHostHeader),
		GzipOn:                preferResourceOption(planned.GzipOn, prior.GzipOn),
		GzipCompression:       preferResourceOption(planned.GzipCompression, prior.GzipCompression),
		HostHeader:            preferResourceOption(planned.HostHeader, prior.HostHeader),
		IgnoreQueryString:     preferResourceOption(planned.IgnoreQueryString, prior.IgnoreQueryString),
		IgnoreCookie:          preferResourceOption(planned.IgnoreCookie, prior.IgnoreCookie),
		QueryParamsBlacklist:  preferResourceOption(planned.QueryParamsBlacklist, prior.QueryParamsBlacklist),
		QueryParamsWhitelist:  preferResourceOption(planned.QueryParamsWhitelist, prior.QueryParamsWhitelist),
		CountryACL:            preferResourceOption(planned.CountryACL, prior.CountryACL),
		ReferrerACL:           preferResourceOption(planned.ReferrerACL, prior.ReferrerACL),
		IpAddressACL:          preferResourceOption(planned.IpAddressACL, prior.IpAddressACL),
		UserAgentACL:          preferResourceOption(planned.UserAgentACL, prior.UserAgentACL),
		Slice:                 preferResourceOption(planned.Slice, prior.Slice),
		Stale:                 preferResourceOption(planned.Stale, prior.Stale),
		StaticHeaders:         preferResourceOption(planned.StaticHeaders, prior.StaticHeaders),
		StaticRequestHeaders:  preferResourceOption(planned.StaticRequestHeaders, prior.StaticRequestHeaders),
		StaticResponseHeaders: preferResourceOption(planned.StaticResponseHeaders, prior.StaticResponseHeaders),
		TlsVersions:           preferResourceOption(planned.TlsVersions, prior.TlsVersions),
		Websockets:            preferResourceOption(planned.Websockets, prior.Websockets),
		SecureKey:             preferResourceOption(planned.SecureKey, prior.SecureKey),
	}
}

func preferResourceOption[T any](planned, prior *T) *T {
	if planned != nil {
		return planned
	}
	return prior
}

func validateStaticResponseHeaders(ctx context.Context, value types.Object) diag.Diagnostics {
	var diags diag.Diagnostics
	if value.IsNull() || value.IsUnknown() {
		return diags
	}

	optionValuable, d := resource_resource.StaticResponseHeadersType{}.ValueFromObject(ctx, value)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}
	option := optionValuable.(resource_resource.StaticResponseHeadersValue)
	if option.Value.IsUnknown() {
		return diags
	}
	if option.Value.IsNull() {
		if option.Enabled.IsUnknown() {
			return diags
		}
		if !option.Enabled.IsNull() && !option.Enabled.IsUnknown() && !option.Enabled.ValueBool() {
			return diags
		}
		diags.AddAttributeError(
			path.Root("options").AtName("static_response_headers").AtName("value"),
			"Missing Attribute Value",
			"Set `value` (an empty list is allowed) or explicitly disable `static_response_headers`.",
		)
		return diags
	}

	return diags
}

func (r *resourceResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan *resource_resource.ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config resource_resource.ResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	originalOptions := plan.Options

	var previousOptions *resource_resource.OptionsValue
	var previousState *resource_resource.ResourceModel
	if !req.State.Raw.IsNull() {
		var state resource_resource.ResourceModel
		resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
		if resp.Diagnostics.HasError() {
			return
		}
		previousOptions = &state.Options
		previousState = &state
	}

	if previousState != nil {
		configuredChanges, diags := hasConfiguredResourceChanges(ctx, &config, previousState)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		// Reuse API-computed state only for a no-op plan. On a real update,
		// fields such as resource status and SSL status may change after apply.
		if !configuredChanges {
			resp.Diagnostics.Append(preserveUnconfiguredRootValues(ctx, resp, plan, &config, previousState)...)
			previousPlanShielding := plan.Shielding
			resp.Diagnostics.Append(preserveUnconfiguredShielding(ctx, &plan.Shielding, config.Shielding, previousState.Shielding)...)
			if !plan.Shielding.Equal(previousPlanShielding) {
				resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("shielding"), plan.Shielding)...)
			}
			previousPlanSslCertificate := plan.SslCertificate
			resp.Diagnostics.Append(preserveUnconfiguredSslCertificate(ctx, &plan.SslCertificate, config.SslCertificate, previousState.SslCertificate)...)
			if !plan.SslCertificate.Equal(previousPlanSslCertificate) {
				resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("ssl_certificate"), plan.SslCertificate)...)
			}
			resp.Diagnostics.Append(preserveUnconfiguredOptions(ctx, &plan.Options, config.Options, previousState.Options)...)
		}
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if planDisabledShieldingPopID(&plan.Shielding, config.Shielding) {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("shielding"), plan.Shielding)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	selectConfiguredResponseHeaderOption(ctx, &plan.Options, config.Options, previousOptions)

	gzipCompression := plan.Options.GzipCompression
	if !gzipCompression.IsNull() && !gzipCompression.IsUnknown() {
		gzipCompressionObjV, d := resource_resource.GzipCompressionType{}.ValueFromObject(ctx, gzipCompression)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		gzipCompression := gzipCompressionObjV.(resource_resource.GzipCompressionValue)

		plan.Options.GzipOn = types.BoolValue(gzipCompression.Enabled.ValueBool())
	}
	if !plan.Options.Equal(originalOptions) {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("options"), plan.Options)...)
	}
}

func hasConfiguredResourceChanges(ctx context.Context, config, previous *resource_resource.ResourceModel) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics
	if config == nil || previous == nil {
		return false, diags
	}

	for _, changed := range map[string]bool{
		"cname":               configuredValueChanged(config.Cname, previous.Cname),
		"origin_group":        configuredValueChanged(config.OriginGroup, previous.OriginGroup),
		"active":              configuredValueChanged(config.Active, previous.Active),
		"origin_protocol":     configuredValueChanged(config.OriginProtocol, previous.OriginProtocol),
		"region":              configuredValueChanged(config.Region, previous.Region),
		"secondary_hostnames": configuredValueChanged(config.SecondaryHostnames, previous.SecondaryHostnames),
	} {
		if changed {
			return true, diags
		}
	}

	if !config.Options.IsNull() && !config.Options.IsUnknown() {
		configOptions, d := config.Options.ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return false, diags
		}
		previousOptions, d := previous.Options.ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return false, diags
		}
		for name, configValue := range configOptions.Attributes() {
			previousValue, ok := previousOptions.Attributes()[name]
			if ok && configuredValueChanged(configValue, previousValue) {
				return true, diags
			}
		}
	}

	if !config.Shielding.IsNull() && !config.Shielding.IsUnknown() {
		configShielding, d := config.Shielding.ToObjectValue(ctx)
		diags.Append(d...)
		previousShielding, d := previous.Shielding.ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return false, diags
		}
		if configuredValueChanged(configShielding, previousShielding) {
			return true, diags
		}
	}

	if !config.SslCertificate.IsNull() && !config.SslCertificate.IsUnknown() {
		configSslCertificate, d := config.SslCertificate.ToObjectValue(ctx)
		diags.Append(d...)
		previousSslCertificate, d := previous.SslCertificate.ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return false, diags
		}
		if configuredValueChanged(configSslCertificate, previousSslCertificate) {
			return true, diags
		}
	}

	return false, diags
}

func configuredValueChanged(config, previous attr.Value) bool {
	if config.IsNull() || config.IsUnknown() {
		return false
	}
	if previous.IsNull() || previous.IsUnknown() {
		return true
	}

	configObject, configIsObject := config.(types.Object)
	previousObject, previousIsObject := previous.(types.Object)
	if !configIsObject || !previousIsObject {
		return !config.Equal(previous)
	}

	configAttributes := configObject.Attributes()
	previousAttributes := previousObject.Attributes()
	for name, configValue := range configAttributes {
		previousValue, ok := previousAttributes[name]
		if !ok || configuredValueChanged(configValue, previousValue) {
			return true
		}
	}

	return false
}

func planDisabledShieldingPopID(plan *resource_resource.ShieldingValue, config resource_resource.ShieldingValue) bool {
	if plan == nil || config.IsNull() || config.IsUnknown() || config.Enabled.IsNull() || config.Enabled.IsUnknown() || config.Enabled.ValueBool() {
		return false
	}
	if plan.IsNull() || plan.IsUnknown() {
		*plan = config
	}
	plan.Enabled = types.BoolValue(false)
	plan.PopId = types.Int64Value(0)
	return true
}

func preserveUnconfiguredRootValues(
	ctx context.Context,
	resp *resource.ModifyPlanResponse,
	plan, config, previous *resource_resource.ResourceModel,
) diag.Diagnostics {
	var diags diag.Diagnostics
	if plan == nil || config == nil || previous == nil {
		return diags
	}

	if preserveUnconfiguredRootValue(&plan.Active, config.Active, previous.Active) {
		diags.Append(resp.Plan.SetAttribute(ctx, path.Root("active"), plan.Active)...)
	}
	if preserveUnknownRootValue(&plan.Id, previous.Id) {
		diags.Append(resp.Plan.SetAttribute(ctx, path.Root("id"), plan.Id)...)
	}
	if preserveUnconfiguredRootValue(&plan.OriginProtocol, config.OriginProtocol, previous.OriginProtocol) {
		diags.Append(resp.Plan.SetAttribute(ctx, path.Root("origin_protocol"), plan.OriginProtocol)...)
	}
	if preserveUnknownRootValue(&plan.PresetApplied, previous.PresetApplied) {
		diags.Append(resp.Plan.SetAttribute(ctx, path.Root("preset_applied"), plan.PresetApplied)...)
	}
	if preserveUnconfiguredRootValue(&plan.Region, config.Region, previous.Region) {
		diags.Append(resp.Plan.SetAttribute(ctx, path.Root("region"), plan.Region)...)
	}
	if preserveUnconfiguredRootValue(&plan.SecondaryHostnames, config.SecondaryHostnames, previous.SecondaryHostnames) {
		diags.Append(resp.Plan.SetAttribute(ctx, path.Root("secondary_hostnames"), plan.SecondaryHostnames)...)
	}
	if preserveUnconfiguredRootValue(&plan.Shielding, config.Shielding, previous.Shielding) {
		diags.Append(resp.Plan.SetAttribute(ctx, path.Root("shielding"), plan.Shielding)...)
	}
	if preserveUnconfiguredRootValue(&plan.SslCertificate, config.SslCertificate, previous.SslCertificate) {
		diags.Append(resp.Plan.SetAttribute(ctx, path.Root("ssl_certificate"), plan.SslCertificate)...)
	}
	if preserveUnknownRootValue(&plan.Status, previous.Status) {
		diags.Append(resp.Plan.SetAttribute(ctx, path.Root("status"), plan.Status)...)
	}
	if preserveUnknownRootValue(&plan.VpEnabled, previous.VpEnabled) {
		diags.Append(resp.Plan.SetAttribute(ctx, path.Root("vp_enabled"), plan.VpEnabled)...)
	}

	return diags
}

type nullableUnknownValue interface {
	IsNull() bool
	IsUnknown() bool
}

func preserveUnconfiguredRootValue[T nullableUnknownValue](plan *T, config, previous T) bool {
	if config.IsNull() && ((*plan).IsNull() || (*plan).IsUnknown()) {
		*plan = previous
		return true
	}
	return false
}

func preserveUnknownRootValue[T nullableUnknownValue](plan *T, previous T) bool {
	if (*plan).IsNull() || (*plan).IsUnknown() {
		*plan = previous
		return true
	}
	return false
}

func (r *resourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_resource.ResourceModel
	var config resource_resource.ResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Create(ctx, resource_resource.ResourceReadyTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}
	data.Region = types.StringValue(region)

	client, err := r.config.CDNV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating CDN API client", err.Error())
		return
	}

	secondaryHostnames := make([]string, 0, len(data.SecondaryHostnames.Elements()))
	if !data.SecondaryHostnames.IsNull() && !data.SecondaryHostnames.IsUnknown() {
		resp.Diagnostics.Append(data.SecondaryHostnames.ElementsAs(ctx, &secondaryHostnames, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	requestOptions := responseHeaderOptionsForAPIRequest(ctx, data.Options, config.Options)
	opts, diags := requestOptions.ToResourceOptions(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOpts := resources.CreateOpts{
		Active:             data.Active.ValueBool(),
		CNAME:              data.Cname.ValueString(),
		Options:            opts,
		OriginGroup:        int(data.OriginGroup.ValueInt64()),
		OriginProtocol:     resources.ResourceOriginProtocol(data.OriginProtocol.ValueString()),
		SecondaryHostnames: secondaryHostnames,
	}

	if sslOpts := data.SslCertificate.ToSslOpts(); sslOpts != nil {
		createOpts.SSLEnabled = sslOpts.Enabled
		createOpts.SSLData = sslOpts.Data
	}

	projectID := r.config.GetProjectID()

	tflog.Trace(ctx, "Calling CDN API to create CDN resource", map[string]interface{}{"opts": fmt.Sprintf("%#v", createOpts)})

	resource, err := resources.Create(client, projectID, &createOpts).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling CDN API to create CDN resource", err.Error())
		return
	}

	tflog.Trace(ctx, "Called CDN API to create CDN resource", map[string]interface{}{"resource": fmt.Sprintf("%#v", resource)})

	id := resource.ID
	resp.State.SetAttribute(ctx, path.Root("id"), id)

	resp.Diagnostics.Append(resource_resource.WaitForResourceReady(ctx, client, projectID, id, timeout)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sslCertType := data.SslCertificate.SslCertificateType.ValueString()
	if sslCertType == string(resource_resource.SslCertificateProviderTypeLetsEncrypt) {
		resp.Diagnostics.Append(issueLetsEncrypt(ctx, client, projectID, id)...)
		if resp.Diagnostics.HasError() {
			return
		}

		resp.Diagnostics.Append(resource_resource.WaitForResourceReady(ctx, client, projectID, id, timeout)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	shieldingOpts := data.Shielding.ToUpdateShieldingOpts()
	if shieldingOpts != nil {
		resourceShielding, diags := updateShielding(ctx, client, projectID, id, shieldingOpts)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		data.Shielding = resource_resource.ShieldingValue{}.FromResourceShielding(resourceShielding)

		resp.Diagnostics.Append(resource_resource.WaitForResourceReady(ctx, client, projectID, id, timeout)...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		data.Shielding = resource_resource.NewShieldingValueNull()
	}

	resource, diags = retrieveCDNResource(ctx, client, projectID, id)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(data.UpdateFromResource(ctx, resource)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(updateFromLetsEncryptStatus(ctx, client, projectID, id, &data.SslCertificate)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_resource.ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}
	data.Region = types.StringValue(region)

	client, err := r.config.CDNV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating CDN API client", err.Error())
		return
	}

	id := int(data.Id.ValueInt64())
	ctx = tflog.SetField(ctx, "resource_id", id)
	projectID := r.config.GetProjectID()

	tflog.Trace(ctx, "Calling CDN API to retrieve CDN resource")

	resource, err := resources.Get(client, projectID, id).Extract()
	if errutil.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error calling CDN API to retrieve CDN resource", err.Error())
		return
	}

	tflog.Trace(ctx, "Called CDN API to retrieve CDN resource", map[string]interface{}{"resource": fmt.Sprintf("%#v", resource)})

	resp.Diagnostics.Append(data.UpdateFromResource(ctx, resource)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(updateFromLetsEncryptStatus(ctx, client, projectID, id, &data.SslCertificate)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Calling CDN API to retrieve origin shielding settings")

	resourceShielding, err := resources.GetShielding(client, projectID, id).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling CDN API to retrieve origin shielding settings", err.Error())
		return
	}

	tflog.Trace(ctx, "Called CDN API to retrieve origin shielding settings", map[string]interface{}{"resource_shielding": fmt.Sprintf("%#v", resourceShielding)})

	data.Shielding = resource_resource.ShieldingValue{}.FromResourceShielding(resourceShielding)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan resource_resource.ResourceModel
	var data resource_resource.ResourceModel
	var config resource_resource.ResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := plan.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}
	data.Region = types.StringValue(region)

	timeout, diags := data.Timeouts.Update(ctx, resource_resource.ResourceReadyTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	client, err := r.config.CDNV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating CDN API client", err.Error())
		return
	}

	id := int(data.Id.ValueInt64())
	ctx = tflog.SetField(ctx, "resource_id", id)

	secondaryHostnames := make([]string, 0, len(plan.SecondaryHostnames.Elements()))
	if !plan.SecondaryHostnames.IsNull() && !plan.SecondaryHostnames.IsUnknown() {
		resp.Diagnostics.Append(plan.SecondaryHostnames.ElementsAs(ctx, &secondaryHostnames, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	requestOptions := responseHeaderOptionsForAPIRequest(ctx, plan.Options, config.Options)
	plannedOptions, diags := requestOptions.ToResourceOptions(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The CDN update API requires unchanged options too. Optional+computed
	// attributes can be unknown in the plan (for example beside a nested list),
	// so fill only missing planned options from prior state.
	priorRequestOptions := responseHeaderOptionsForAPIRequest(ctx, data.Options, config.Options)
	priorOptions, diags := priorRequestOptions.ToResourceOptions(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	opts := mergeResourceOptions(plannedOptions, priorOptions)

	updateOpts := resources.UpdateOpts{
		Active:             plan.Active.ValueBoolPointer(),
		Options:            opts,
		OriginGroup:        int(plan.OriginGroup.ValueInt64()),
		OriginProtocol:     resources.ResourceOriginProtocol(plan.OriginProtocol.ValueString()),
		SecondaryHostnames: secondaryHostnames,
	}

	if sslOpts := plan.SslCertificate.ToSslOpts(); sslOpts != nil {
		updateOpts.SSLEnabled = &sslOpts.Enabled
		updateOpts.SSLData = sslOpts.Data
	}

	projectID := r.config.GetProjectID()

	tflog.Trace(ctx, "Calling CDN API to update CDN resource", map[string]interface{}{"opts": fmt.Sprintf("%#v", updateOpts)})

	resource, err := resources.Update(client, projectID, id, &updateOpts).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling CDN API to update CDN resource", err.Error())
		return
	}

	tflog.Trace(ctx, "Called CDN API to update CDN resource", map[string]interface{}{"resource": fmt.Sprintf("%#v", resource)})

	resp.Diagnostics.Append(resource_resource.WaitForResourceReady(ctx, client, projectID, id, timeout)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sslCertType := plan.SslCertificate.SslCertificateType.ValueString()
	if !plan.SslCertificate.SslCertificateType.Equal(data.SslCertificate.SslCertificateType) && sslCertType == string(resource_resource.SslCertificateProviderTypeLetsEncrypt) {
		resp.Diagnostics.Append(issueLetsEncrypt(ctx, client, projectID, id)...)
		if resp.Diagnostics.HasError() {
			return
		}

		resp.Diagnostics.Append(resource_resource.WaitForResourceReady(ctx, client, projectID, id, timeout)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !plan.Shielding.Equal(data.Shielding) {
		shieldingOpts := plan.Shielding.ToUpdateShieldingOpts()
		if shieldingOpts != nil {
			resourceShielding, diags := updateShielding(ctx, client, projectID, id, shieldingOpts)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}

			data.Shielding = resource_resource.ShieldingValue{}.FromResourceShielding(resourceShielding)

			resp.Diagnostics.Append(resource_resource.WaitForResourceReady(ctx, client, projectID, id, timeout)...)
			if resp.Diagnostics.HasError() {
				return
			}
		}
	}

	resource, diags = retrieveCDNResource(ctx, client, projectID, id)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(data.UpdateFromResource(ctx, resource)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(updateFromLetsEncryptStatus(ctx, client, projectID, id, &data.SslCertificate)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_resource.ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	client, err := r.config.CDNV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating CDN API client", err.Error())
		return
	}

	id := data.Id.ValueInt64()
	ctx = tflog.SetField(ctx, "resource_id", id)

	tflog.Trace(ctx, "Calling CDN API to delete CDN resource")

	err = resources.Delete(client, r.config.GetProjectID(), int(id)).ExtractErr()
	if errutil.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error calling CDN API to delete CDN resource", err.Error())
		return
	}

	tflog.Trace(ctx, "Called CDN API to delete CDN resource")
}

func (r *resourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	fwutils.ImportStatePassthroughInt64ID(ctx, req, resp)
}

func retrieveCDNResource(ctx context.Context, client *gophercloud.ServiceClient, projectID string, resourceID int) (*resources.Resource, diag.Diagnostics) {
	var diags diag.Diagnostics
	tflog.Trace(ctx, "Calling CDN API to retrieve CDN resource")

	resource, err := resources.Get(client, projectID, resourceID).Extract()
	if err != nil {
		diags.AddError("Error calling CDN API to retrieve CDN resource", err.Error())
		return nil, diags
	}

	tflog.Trace(ctx, "Called CDN API to retrieve CDN resource", map[string]interface{}{"resource": fmt.Sprintf("%#v", resource)})
	return resource, diags
}

func issueLetsEncrypt(ctx context.Context, client *gophercloud.ServiceClient, projectID string, resourceID int) diag.Diagnostics {
	var diags diag.Diagnostics
	tflog.Trace(ctx, "Calling CDN API to issue a Let's Encrypt certificate")

	err := resources.IssueLetsEncrypt(client, projectID, resourceID).ExtractErr()
	if err != nil {
		diags.AddError("Error calling CDN API to issue a Let's Encrypt certificate", err.Error())
		return diags
	}

	tflog.Trace(ctx, "Called CDN API to issue a Let's Encrypt certificate")
	return diags
}

func updateFromLetsEncryptStatus(ctx context.Context, client *gophercloud.ServiceClient, projectID string, resourceID int, sslCertificateData *resource_resource.SslCertificateValue) diag.Diagnostics {
	var diags diag.Diagnostics

	tflog.Trace(ctx, "Calling CDN API to get Let's Encrypt certificate issuing details")

	leStatus, err := resources.GetLetsEncryptStatus(client, projectID, resourceID).Extract()
	if err != nil {
		if !errutil.IsNotFound(err) {
			diags.AddError("Error calling CDN API to get Let's Encrypt certificate issuing details", err.Error())
		}
		return diags
	}

	if leStatus.Active {
		sslCertificateData.SslCertificateType = types.StringValue(string(resource_resource.SslCertificateProviderTypeLetsEncrypt))
		sslCertificateData.Status = types.StringValue(string(resource_resource.SslCertificateStatusBeingIssued))
	}

	tflog.Trace(ctx, "Called CDN API to get Let's Encrypt certificate issuing details", map[string]any{"status": fmt.Sprintf("%#v", leStatus)})
	return diags
}

func updateShielding(ctx context.Context, client *gophercloud.ServiceClient, projectID string, resourceID int, opts *resources.UpdateShieldingOpts) (*resources.ResourceShielding, diag.Diagnostics) {
	var diags diag.Diagnostics
	tflog.Trace(ctx, "Calling CDN API to update origin shielding settings")

	resourceShielding, err := resources.UpdateShielding(client, projectID, resourceID, opts).Extract()
	if err != nil {
		diags.AddError("Error calling CDN API to update origin shielding settings", err.Error())
		return nil, diags
	}

	tflog.Trace(ctx, "Called CDN API to update origin shielding settings")
	return resourceShielding, diags
}
