package datasource_service_user

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/iam/serviceusers"
)

func (m *ServiceUserModel) UpdateFromServiceUser(ctx context.Context, serviceUser *serviceusers.ServiceUser) diag.Diagnostics {
	if serviceUser == nil {
		return nil
	}

	m.Id = types.StringValue(serviceUser.ID)
	m.Name = types.StringValue(serviceUser.Name)
	m.CreatedAt = types.StringValue(serviceUser.CreatedAt)
	m.CreatorName = types.StringValue(serviceUser.CreatorName)
	m.Description = types.StringValue(serviceUser.Description)

	var d diag.Diagnostics
	m.RoleNames, d = types.ListValueFrom(ctx, types.StringType, serviceUser.RoleNames)
	if d.HasError() {
		return d
	}

	return nil
}
