package resource_anycastip

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking/v2/anycastips"
)

func (m *AnycastipModel) UpdateFromAnycastIP(ctx context.Context, anycastIP *anycastips.AnycastIP) diag.Diagnostics {
	var diags diag.Diagnostics

	if anycastIP == nil {
		return diags
	}

	m.Id = types.StringValue(anycastIP.ID)
	m.Name = types.StringValue(anycastIP.Name)
	m.Description = types.StringValue(anycastIP.Description)
	m.NetworkId = types.StringValue(anycastIP.NetworkID)
	m.SubnetId = types.StringValue(anycastIP.SubnetID)
	m.IpAddress = types.StringValue(anycastIP.IPAddress)

	associations, d := flattenAssociations(ctx, anycastIP.Associations)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	m.Associations = associations

	healthCheck := flattenHealthCheck(ctx, anycastIP.HealthCheck)
	m.HealthCheck = healthCheck

	return diags
}

func ExpandAssociations(ctx context.Context, associations types.Set) ([]anycastips.AnycastIPAssociation, diag.Diagnostics) {
	if associations.IsUnknown() || associations.IsNull() {
		return nil, nil
	}

	associationsV := make([]AssociationsValue, 0, len(associations.Elements()))
	diags := associations.ElementsAs(ctx, &associationsV, false)
	if diags.HasError() {
		return nil, diags
	}

	result := make([]anycastips.AnycastIPAssociation, len(associationsV))
	for i, o := range associationsV {
		result[i] = anycastips.AnycastIPAssociation{
			ID:   o.Id.ValueString(),
			Type: anycastips.AnycastIPAssociationType(o.AssociationsType.ValueString()),
		}
	}

	return result, nil
}

func ExpandHealthCheck(ctx context.Context, healthCheck HealthCheckValue) *anycastips.AnycastIPHealthCheck {
	if healthCheck.IsUnknown() || healthCheck.IsNull() {
		return nil
	}

	healthCheckType := anycastips.AnycastIPHealthCheckType(healthCheck.HealthCheckType.ValueString())
	port := int(healthCheck.Port.ValueInt64())

	result := anycastips.AnycastIPHealthCheck{
		Type: healthCheckType,
		Port: port,
	}

	return &result
}

func flattenAssociations(ctx context.Context, associations []anycastips.AnycastIPAssociation) (types.Set, diag.Diagnostics) {
	associationsVType := AssociationsValue{}.Type(ctx)

	if len(associations) == 0 {
		return types.SetNull(associationsVType), nil
	}

	associationsV := make([]attr.Value, len(associations))
	for i, a := range associations {
		associationsV[i] = AssociationsValue{
			Id:               types.StringValue(a.ID),
			AssociationsType: types.StringValue(string(a.Type)),
			state:            attr.ValueStateKnown,
		}
	}

	return types.SetValue(AssociationsValue{}.Type(ctx), associationsV)
}

func flattenHealthCheck(ctx context.Context, healthCheck *anycastips.AnycastIPHealthCheck) HealthCheckValue {
	if healthCheck.Type == "" {
		return NewHealthCheckValueNull()
	}

	healthCheckValue := HealthCheckValue{
		Port:            types.Int64Value(int64(healthCheck.Port)),
		HealthCheckType: types.StringValue(string(healthCheck.Type)),
		state:           attr.ValueStateKnown,
	}

	return healthCheckValue
}
