package utils

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/net/context"
)

func IsKnown(v attr.Value) bool {
	return !v.IsNull() && !v.IsUnknown()
}

func GetFirstNotEmptyValue(values ...types.String) string {
	for _, value := range values {
		if len(value.ValueString()) != 0 {
			return value.ValueString()
		}
	}

	return ""
}

func ImportStatePassthroughInt64ID(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	rawID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("The resource ID must be a valid integer", err.Error())
		return
	}
	id := types.Int64Value(int64(rawID))
	resp.State.SetAttribute(ctx, path.Root("id"), id)
}
