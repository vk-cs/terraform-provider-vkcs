package vkcs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud"
	octavial7policies "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/l7policies"
	octavialisteners "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/listeners"
	octavialoadbalancers "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/loadbalancers"
	octaviamonitors "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/monitors"
	octaviapools "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/pools"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
)

const (
	lbPendingCreate = "PENDING_CREATE"
	lbPendingUpdate = "PENDING_UPDATE"
	lbPendingDelete = "PENDING_DELETE"
	lbActive        = "ACTIVE"
	lbError         = "ERROR"
)

// lbPendingStatuses are the valid statuses a LoadBalancer will be in while
// it's updating.
func getLbPendingStatuses() []string {
	return []string{lbPendingCreate, lbPendingUpdate}
}

// lbPendingDeleteStatuses are the valid statuses a LoadBalancer will be before delete.
func getLbPendingDeleteStatuses() []string {
	return []string{lbError, lbPendingUpdate, lbPendingDelete, lbActive}
}

func getLbSkipStatuses() []string {
	return []string{lbError, lbActive}
}

func expandLBListenerHeadersMap(raw map[string]interface{}) (map[string]string, error) {
	m := make(map[string]string, len(raw))
	for key, val := range raw {
		labelValue, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("label %s value should be string", key)
		}

		m[key] = labelValue
	}

	return m, nil
}

func waitForLBListener(ctx context.Context, lbClient *gophercloud.ServiceClient, listener *octavialisteners.Listener, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for vkcs_lb_listener %s to become %s.", listener.ID, target)

	if len(listener.Loadbalancers) == 0 {
		return fmt.Errorf("failed to detect a vkcs_lb_loadbalancer for the %s vkcs_lb_listener", listener.ID)
	}

	lbID := listener.Loadbalancers[0].ID

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceLBListenerRefreshFunc(lbClient, lbID, listener),
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			if target == "DELETED" {
				return nil
			}
		}

		return fmt.Errorf("error waiting for vkcs_lb_listener %s to become %s: %s", listener.ID, target, err)
	}

	return nil
}

func resourceLBListenerRefreshFunc(lbClient *gophercloud.ServiceClient, lbID string, listener *octavialisteners.Listener) resource.StateRefreshFunc {
	if listener.ProvisioningStatus != "" {
		return func() (interface{}, string, error) {
			lb, status, err := resourceLBLoadBalancerRefreshFunc(lbClient, lbID)()
			if err != nil {
				return lb, status, err
			}
			if !strSliceContains(getLbSkipStatuses(), status) {
				return lb, status, nil
			}

			listener, err := octavialisteners.Get(lbClient, listener.ID).Extract()
			if err != nil {
				return nil, "", err
			}

			return listener, listener.ProvisioningStatus, nil
		}
	}

	return resourceLBLoadBalancerStatusRefreshFuncNeutron(lbClient, lbID, "listener", listener.ID, "")
}

func waitForLBLoadBalancer(ctx context.Context, lbClient *gophercloud.ServiceClient, lbID string, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for loadbalancer %s to become %s.", lbID, target)

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceLBLoadBalancerRefreshFunc(lbClient, lbID),
		Timeout:    timeout,
		Delay:      0,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			switch target {
			case "DELETED":
				return nil
			default:
				return fmt.Errorf("error: loadbalancer %s not found: %s", lbID, err)
			}
		}
		return fmt.Errorf("error waiting for loadbalancer %s to become %s: %s", lbID, target, err)
	}

	return nil
}

func resourceLBLoadBalancerRefreshFunc(lbClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		lb, err := octavialoadbalancers.Get(lbClient, id).Extract()
		if err != nil {
			return nil, "", err
		}

		return lb, lb.ProvisioningStatus, nil
	}
}

func waitForLBMember(ctx context.Context, lbClient *gophercloud.ServiceClient, parentPool *octaviapools.Pool, member *octaviapools.Member, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for member %s to become %s.", member.ID, target)

	lbID, err := lbFindLBIDviaPool(lbClient, parentPool)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceLBMemberRefreshFunc(lbClient, lbID, parentPool.ID, member),
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			if target == "DELETED" {
				return nil
			}
		}

		return fmt.Errorf("error waiting for member %s to become %s: %s", member.ID, target, err)
	}

	return nil
}

func resourceLBMemberRefreshFunc(lbClient *gophercloud.ServiceClient, lbID string, poolID string, member *octaviapools.Member) resource.StateRefreshFunc {
	if member.ProvisioningStatus != "" {
		return func() (interface{}, string, error) {
			lb, status, err := resourceLBLoadBalancerRefreshFunc(lbClient, lbID)()
			if err != nil {
				return lb, status, err
			}
			if !strSliceContains(getLbSkipStatuses(), status) {
				return lb, status, nil
			}

			member, err := octaviapools.GetMember(lbClient, poolID, member.ID).Extract()
			if err != nil {
				return nil, "", err
			}

			return member, member.ProvisioningStatus, nil
		}
	}

	return resourceLBLoadBalancerStatusRefreshFuncNeutron(lbClient, lbID, "member", member.ID, poolID)
}

func waitForLBMonitor(ctx context.Context, lbClient *gophercloud.ServiceClient, parentPool *octaviapools.Pool, monitor *octaviamonitors.Monitor, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for vkcs_lb_monitor %s to become %s.", monitor.ID, target)

	lbID, err := lbFindLBIDviaPool(lbClient, parentPool)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceLBMonitorRefreshFunc(lbClient, lbID, monitor),
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			if target == "DELETED" {
				return nil
			}
		}
		return fmt.Errorf("error waiting for vkcs_lb_monitor %s to become %s: %s", monitor.ID, target, err)
	}

	return nil
}

func resourceLBMonitorRefreshFunc(lbClient *gophercloud.ServiceClient, lbID string, monitor *octaviamonitors.Monitor) resource.StateRefreshFunc {
	if monitor.ProvisioningStatus != "" {
		return func() (interface{}, string, error) {
			lb, status, err := resourceLBLoadBalancerRefreshFunc(lbClient, lbID)()
			if err != nil {
				return lb, status, err
			}
			if !strSliceContains(getLbSkipStatuses(), status) {
				return lb, status, nil
			}

			monitor, err := octaviamonitors.Get(lbClient, monitor.ID).Extract()
			if err != nil {
				return nil, "", err
			}

			return monitor, monitor.ProvisioningStatus, nil
		}
	}

	return resourceLBLoadBalancerStatusRefreshFuncNeutron(lbClient, lbID, "monitor", monitor.ID, "")
}

func waitForLBPool(ctx context.Context, lbClient *gophercloud.ServiceClient, pool *octaviapools.Pool, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for pool %s to become %s.", pool.ID, target)

	lbID, err := lbFindLBIDviaPool(lbClient, pool)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceLBPoolRefreshFunc(lbClient, lbID, pool),
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			if target == "DELETED" {
				return nil
			}
		}

		return fmt.Errorf("error waiting for pool %s to become %s: %s", pool.ID, target, err)
	}

	return nil
}

func resourceLBPoolRefreshFunc(lbClient *gophercloud.ServiceClient, lbID string, pool *octaviapools.Pool) resource.StateRefreshFunc {
	if pool.ProvisioningStatus != "" {
		return func() (interface{}, string, error) {
			lb, status, err := resourceLBLoadBalancerRefreshFunc(lbClient, lbID)()
			if err != nil {
				return lb, status, err
			}
			if !strSliceContains(getLbSkipStatuses(), status) {
				return lb, status, nil
			}

			pool, err := octaviapools.Get(lbClient, pool.ID).Extract()
			if err != nil {
				return nil, "", err
			}

			return pool, pool.ProvisioningStatus, nil
		}
	}

	return resourceLBLoadBalancerStatusRefreshFuncNeutron(lbClient, lbID, "pool", pool.ID, "")
}

func lbFindLBIDviaPool(lbClient *gophercloud.ServiceClient, pool *octaviapools.Pool) (string, error) {
	if len(pool.Loadbalancers) > 0 {
		return pool.Loadbalancers[0].ID, nil
	}

	if len(pool.Listeners) > 0 {
		listenerID := pool.Listeners[0].ID
		listener, err := octavialisteners.Get(lbClient, listenerID).Extract()
		if err != nil {
			return "", err
		}

		if len(listener.Loadbalancers) > 0 {
			return listener.Loadbalancers[0].ID, nil
		}
	}

	return "", fmt.Errorf("unable to determine loadbalancer ID from pool %s", pool.ID)
}

func resourceLBLoadBalancerStatusRefreshFuncNeutron(lbClient *gophercloud.ServiceClient, lbID, resourceType, resourceID string, parentID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		statuses, err := octavialoadbalancers.GetStatuses(lbClient, lbID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return nil, "", gophercloud.ErrDefault404{
					ErrUnexpectedResponseCode: gophercloud.ErrUnexpectedResponseCode{
						BaseError: gophercloud.BaseError{
							DefaultErrString: fmt.Sprintf("Unable to get statuses from the Load Balancer %s statuses tree: %s", lbID, err),
						},
					},
				}
			}
			return nil, "", fmt.Errorf("unable to get statuses from the Load Balancer %s statuses tree: %s", lbID, err)
		}

		// Don't fail, when statuses returns "null"
		if statuses == nil || statuses.Loadbalancer == nil {
			statuses = new(octavialoadbalancers.StatusTree)
			statuses.Loadbalancer = new(octavialoadbalancers.LoadBalancer)
		} else if !strSliceContains(getLbSkipStatuses(), statuses.Loadbalancer.ProvisioningStatus) {
			return statuses.Loadbalancer, statuses.Loadbalancer.ProvisioningStatus, nil
		}

		switch resourceType {
		case "listener":
			for _, listener := range statuses.Loadbalancer.Listeners {
				if listener.ID == resourceID {
					if listener.ProvisioningStatus != "" {
						return listener, listener.ProvisioningStatus, nil
					}
				}
			}
			listener, err := octavialisteners.Get(lbClient, resourceID).Extract()
			return listener, "ACTIVE", err

		case "pool":
			for _, pool := range statuses.Loadbalancer.Pools {
				if pool.ID == resourceID {
					if pool.ProvisioningStatus != "" {
						return pool, pool.ProvisioningStatus, nil
					}
				}
			}
			pool, err := octaviapools.Get(lbClient, resourceID).Extract()
			return pool, "ACTIVE", err

		case "monitor":
			for _, pool := range statuses.Loadbalancer.Pools {
				if pool.Monitor.ID == resourceID {
					if pool.Monitor.ProvisioningStatus != "" {
						return pool.Monitor, pool.Monitor.ProvisioningStatus, nil
					}
				}
			}
			monitor, err := octaviamonitors.Get(lbClient, resourceID).Extract()
			return monitor, "ACTIVE", err

		case "member":
			for _, pool := range statuses.Loadbalancer.Pools {
				for _, member := range pool.Members {
					if member.ID == resourceID {
						if member.ProvisioningStatus != "" {
							return member, member.ProvisioningStatus, nil
						}
					}
				}
			}
			member, err := octaviapools.GetMember(lbClient, parentID, resourceID).Extract()
			return member, "ACTIVE", err

		case "l7policy":
			for _, listener := range statuses.Loadbalancer.Listeners {
				for _, l7policy := range listener.L7Policies {
					if l7policy.ID == resourceID {
						if l7policy.ProvisioningStatus != "" {
							return l7policy, l7policy.ProvisioningStatus, nil
						}
					}
				}
			}
			l7policy, err := octavial7policies.Get(lbClient, resourceID).Extract()
			return l7policy, "ACTIVE", err

		case "l7rule":
			for _, listener := range statuses.Loadbalancer.Listeners {
				for _, l7policy := range listener.L7Policies {
					for _, l7rule := range l7policy.Rules {
						if l7rule.ID == resourceID {
							if l7rule.ProvisioningStatus != "" {
								return l7rule, l7rule.ProvisioningStatus, nil
							}
						}
					}
				}
			}
			l7Rule, err := octavial7policies.GetRule(lbClient, parentID, resourceID).Extract()
			return l7Rule, "ACTIVE", err
		}

		return nil, "", fmt.Errorf("an unexpected error occurred querying the status of %s %s by loadbalancer %s", resourceType, resourceID, lbID)
	}
}

func resourceLBL7PolicyRefreshFunc(lbClient *gophercloud.ServiceClient, lbID string, l7policy *octavial7policies.L7Policy) resource.StateRefreshFunc {
	if l7policy.ProvisioningStatus != "" {
		return func() (interface{}, string, error) {
			lb, status, err := resourceLBLoadBalancerRefreshFunc(lbClient, lbID)()
			if err != nil {
				return lb, status, err
			}
			if !strSliceContains(getLbSkipStatuses(), status) {
				return lb, status, nil
			}

			l7policy, err := octavial7policies.Get(lbClient, l7policy.ID).Extract()
			if err != nil {
				return nil, "", err
			}

			return l7policy, l7policy.ProvisioningStatus, nil
		}
	}

	return resourceLBLoadBalancerStatusRefreshFuncNeutron(lbClient, lbID, "l7policy", l7policy.ID, "")
}

func waitForLBL7Policy(ctx context.Context, lbClient *gophercloud.ServiceClient, parentListener *octavialisteners.Listener, l7policy *octavial7policies.L7Policy, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for l7policy %s to become %s.", l7policy.ID, target)

	if len(parentListener.Loadbalancers) == 0 {
		return fmt.Errorf("unable to determine loadbalancer ID from listener %s", parentListener.ID)
	}

	lbID := parentListener.Loadbalancers[0].ID

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceLBL7PolicyRefreshFunc(lbClient, lbID, l7policy),
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			if target == "DELETED" {
				return nil
			}
		}

		return fmt.Errorf("error waiting for l7policy %s to become %s: %s", l7policy.ID, target, err)
	}

	return nil
}

func getListenerIDForL7Policy(lbClient *gophercloud.ServiceClient, id string) (string, error) {
	log.Printf("[DEBUG] Trying to get Listener ID associated with the %s L7 Policy ID", id)
	lbsPages, err := octavialoadbalancers.List(lbClient, octavialoadbalancers.ListOpts{}).AllPages()
	if err != nil {
		return "", fmt.Errorf("no Load Balancers were found: %s", err)
	}

	lbs, err := octavialoadbalancers.ExtractLoadBalancers(lbsPages)
	if err != nil {
		return "", fmt.Errorf("unable to extract Load Balancers list: %s", err)
	}

	for _, lb := range lbs {
		statuses, err := octavialoadbalancers.GetStatuses(lbClient, lb.ID).Extract()
		if err != nil {
			return "", fmt.Errorf("failed to get Load Balancer statuses: %s", err)
		}
		for _, listener := range statuses.Loadbalancer.Listeners {
			for _, l7policy := range listener.L7Policies {
				if l7policy.ID == id {
					return listener.ID, nil
				}
			}
		}
	}

	return "", fmt.Errorf("unable to find Listener ID associated with the %s L7 Policy ID", id)
}

func resourceLBL7RuleRefreshFunc(lbClient *gophercloud.ServiceClient, lbID string, l7policyID string, l7rule *octavial7policies.Rule) resource.StateRefreshFunc {
	if l7rule.ProvisioningStatus != "" {
		return func() (interface{}, string, error) {
			lb, status, err := resourceLBLoadBalancerRefreshFunc(lbClient, lbID)()
			if err != nil {
				return lb, status, err
			}
			if !strSliceContains(getLbSkipStatuses(), status) {
				return lb, status, nil
			}

			l7rule, err := octavial7policies.GetRule(lbClient, l7policyID, l7rule.ID).Extract()
			if err != nil {
				return nil, "", err
			}

			return l7rule, l7rule.ProvisioningStatus, nil
		}
	}

	return resourceLBLoadBalancerStatusRefreshFuncNeutron(lbClient, lbID, "l7rule", l7rule.ID, l7policyID)
}

func waitForLBL7Rule(ctx context.Context, lbClient *gophercloud.ServiceClient, parentListener *octavialisteners.Listener, parentL7policy *octavial7policies.L7Policy, l7rule *octavial7policies.Rule, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for l7rule %s to become %s.", l7rule.ID, target)

	if len(parentListener.Loadbalancers) == 0 {
		return fmt.Errorf("unable to determine loadbalancer ID from listener %s", parentListener.ID)
	}

	lbID := parentListener.Loadbalancers[0].ID

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceLBL7RuleRefreshFunc(lbClient, lbID, parentL7policy.ID, l7rule),
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			if target == "DELETED" {
				return nil
			}
		}

		return fmt.Errorf("error waiting for l7rule %s to become %s: %s", l7rule.ID, target, err)
	}

	return nil
}

func flattenLBPoolPersistence(p octaviapools.SessionPersistence) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"type":        p.Type,
			"cookie_name": p.CookieName,
		},
	}
}

func flattenLBMembers(members []octaviapools.Member) []map[string]interface{} {
	m := make([]map[string]interface{}, len(members))

	for i, member := range members {
		m[i] = map[string]interface{}{
			"name":           member.Name,
			"weight":         member.Weight,
			"admin_state_up": member.AdminStateUp,
			"subnet_id":      member.SubnetID,
			"address":        member.Address,
			"protocol_port":  member.ProtocolPort,
			"id":             member.ID,
			"backup":         member.Backup,
		}
	}

	return m
}

func expandLBMembers(members *schema.Set, lbClient *gophercloud.ServiceClient) []octaviapools.BatchUpdateMemberOpts {
	var m []octaviapools.BatchUpdateMemberOpts

	if members != nil {
		for _, raw := range members.List() {
			rawMap := raw.(map[string]interface{})
			name := rawMap["name"].(string)
			subnetID := rawMap["subnet_id"].(string)
			weight := rawMap["weight"].(int)
			adminStateUp := rawMap["admin_state_up"].(bool)

			member := octaviapools.BatchUpdateMemberOpts{
				Address:      rawMap["address"].(string),
				ProtocolPort: rawMap["protocol_port"].(int),
				Name:         &name,
				SubnetID:     &subnetID,
				Weight:       &weight,
				AdminStateUp: &adminStateUp,
			}

			// backup requires octavia minor version 2.1. Only set when specified
			if val, ok := rawMap["backup"]; ok {
				backup := val.(bool)
				member.Backup = &backup
			}

			m = append(m, member)
		}
	}

	return m
}

func resourceLoadBalancerGetSecurityGroups(networkingClient *gophercloud.ServiceClient, vipPortID string, d *schema.ResourceData) error {
	port, err := ports.Get(networkingClient, vipPortID).Extract()
	if err != nil {
		return err
	}

	d.Set("security_group_ids", port.SecurityGroups)

	return nil
}
