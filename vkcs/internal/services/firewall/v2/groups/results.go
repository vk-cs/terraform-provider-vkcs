package groups

import (
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/pagination"
)

func ExtractSecurityGroupInto(r groups.GetResult, v interface{}) error {
	return r.ExtractIntoStructPtr(v, "security_group")
}

func ExtractSecurityGroupsInto(r pagination.Page, v interface{}) error {
	return r.(groups.SecGroupPage).ExtractIntoSlicePtr(v, "security_groups")
}
