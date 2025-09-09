package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func RequiresReplaceIfWasPresent() planmodifier.String {
	return stringplanmodifier.RequiresReplaceIf(
		func(ctx context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
			if req.StateValue.IsNull() {
				return
			}

			parentPath := req.Path.ParentPath()
			var oldUsername, newUsername string
			req.State.GetAttribute(ctx, parentPath.AtName("username"), &oldUsername)
			req.Config.GetAttribute(ctx, parentPath.AtName("username"), &newUsername)
			if oldUsername != newUsername {
				return
			}

			resp.RequiresReplace = true
		},
		"If the value of this attribute is configured and changes, Terraform will destroy and recreate the resource.",
		"If the value of this attribute is configured and changes, Terraform will destroy and recreate the resource.",
	)
}
