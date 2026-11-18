package planmodifiers

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestResourceOptionWithStatePreservesKnownPriorState(t *testing.T) {
	ctx := context.Background()
	attributeTypes := map[string]attr.Type{
		"enabled": types.BoolType,
	}
	state := types.ObjectValueMust(attributeTypes, map[string]attr.Value{
		"enabled": types.BoolValue(true),
	})
	request := planmodifier.ObjectRequest{
		ConfigValue: types.ObjectNull(attributeTypes),
		PlanValue:   types.ObjectUnknown(attributeTypes),
		StateValue:  state,
	}
	response := planmodifier.ObjectResponse{PlanValue: request.PlanValue}

	ResourceOptionWithState().PlanModifyObject(ctx, request, &response)

	if !response.PlanValue.Equal(state) {
		t.Fatalf("planned value = %s, want prior state %s", response.PlanValue, state)
	}
}

func TestResourceOptionWithStatePreservesNullPriorState(t *testing.T) {
	ctx := context.Background()
	attributeTypes := map[string]attr.Type{
		"enabled": types.BoolType,
	}
	state := types.ObjectNull(attributeTypes)
	request := planmodifier.ObjectRequest{
		ConfigValue: types.ObjectNull(attributeTypes),
		PlanValue:   types.ObjectUnknown(attributeTypes),
		StateValue:  state,
	}
	response := planmodifier.ObjectResponse{PlanValue: request.PlanValue}

	ResourceOptionWithState().PlanModifyObject(ctx, request, &response)

	if !response.PlanValue.Equal(state) {
		t.Fatalf("planned value = %s, want null prior state %s", response.PlanValue, state)
	}
}

func TestResourceOptionWithStateDoesNotReplaceUnknownConfiguration(t *testing.T) {
	ctx := context.Background()
	attributeTypes := map[string]attr.Type{
		"enabled": types.BoolType,
	}
	configuration := types.ObjectUnknown(attributeTypes)
	plan := types.ObjectUnknown(attributeTypes)
	state := types.ObjectValueMust(attributeTypes, map[string]attr.Value{
		"enabled": types.BoolValue(true),
	})
	request := planmodifier.ObjectRequest{
		ConfigValue: configuration,
		PlanValue:   plan,
		StateValue:  state,
	}
	response := planmodifier.ObjectResponse{PlanValue: plan}

	ResourceOptionWithState().PlanModifyObject(ctx, request, &response)

	if !response.PlanValue.Equal(plan) {
		t.Fatalf("planned value = %s, want unresolved plan %s", response.PlanValue, plan)
	}
}
