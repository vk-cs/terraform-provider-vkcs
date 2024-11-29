package l7policies

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/l7policies"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func Create(c *gophercloud.ServiceClient, opts l7policies.CreateOptsBuilder) l7policies.CreateResult {
	r := l7policies.Create(c, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func Get(c *gophercloud.ServiceClient, id string) l7policies.GetResult {
	r := l7policies.Get(c, id)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func Delete(c *gophercloud.ServiceClient, id string) l7policies.DeleteResult {
	r := l7policies.Delete(c, id)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func Update(c *gophercloud.ServiceClient, id string, opts l7policies.UpdateOptsBuilder) l7policies.UpdateResult {
	r := l7policies.Update(c, id, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func CreateRule(c *gophercloud.ServiceClient, policyID string, opts l7policies.CreateRuleOpts) l7policies.CreateRuleResult {
	r := l7policies.CreateRule(c, policyID, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func GetRule(c *gophercloud.ServiceClient, policyID string, ruleID string) l7policies.GetRuleResult {
	r := l7policies.GetRule(c, policyID, ruleID)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func DeleteRule(c *gophercloud.ServiceClient, policyID string, ruleID string) l7policies.DeleteRuleResult {
	r := l7policies.DeleteRule(c, policyID, ruleID)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func UpdateRule(c *gophercloud.ServiceClient, policyID string, ruleID string, opts l7policies.UpdateRuleOptsBuilder) l7policies.UpdateRuleResult {
	r := l7policies.UpdateRule(c, policyID, ruleID, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}
