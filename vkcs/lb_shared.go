package vkcs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud"
	octavialisteners "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/listeners"
	octavialoadbalancers "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/loadbalancers"
	octaviamonitors "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/monitors"
	octaviapools "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/pools"
	neutronl7policies "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/l7policies"
	neutronlisteners "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/listeners"
	neutronloadbalancers "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/loadbalancers"
	neutronmonitors "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/monitors"
	neutronpools "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/pools"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
)

const octaviaLBClientType = "load-balancer"

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

// chooseLBClient will determine which load balacing client to use:
// either the Octavia/LBaaS client or the Neutron/Networking v2 client.
func chooseLBClient(d *schema.ResourceData, config *config) (*gophercloud.ServiceClient, error) {
	if config.UseOctavia {
		return config.LoadBalancerV2Client(getRegion(d, config))
	}
	return config.NetworkingV2Client(getRegion(d, config), getSDN(d))
}

// chooseLBAccTestClient will determine which load balacing client to use:
// either the Octavia/LBaaS client or the Neutron/Networking v2 client.
// This is similar to the chooseLBClient function but specific for acceptance
// tests.
func chooseLBAccTestClient(config *config, region string) (*gophercloud.ServiceClient, error) {
	if config.UseOctavia {
		return config.LoadBalancerV2Client(region)
	}
	return config.NetworkingV2Client(region, defaultSDN)
}

// chooseLBLoadbalancerUpdateOpts will determine which load balancer update options to use:
// either the Octavia/LBaaS or the Neutron/Networking v2.
func chooseLBLoadbalancerUpdateOpts(d *schema.ResourceData, config *config) (neutronloadbalancers.UpdateOptsBuilder, error) {
	var hasChange bool

	if config.UseOctavia {
		// Use Octavia.
		var updateOpts octavialoadbalancers.UpdateOpts

		if d.HasChange("name") {
			hasChange = true
			name := d.Get("name").(string)
			updateOpts.Name = &name
		}
		if d.HasChange("description") {
			hasChange = true
			description := d.Get("description").(string)
			updateOpts.Description = &description
		}
		if d.HasChange("admin_state_up") {
			hasChange = true
			asu := d.Get("admin_state_up").(bool)
			updateOpts.AdminStateUp = &asu
		}

		if d.HasChange("tags") {
			hasChange = true
			if v, ok := d.GetOk("tags"); ok {
				tags := v.(*schema.Set).List()
				tagsToUpdate := expandToStringSlice(tags)
				updateOpts.Tags = &tagsToUpdate
			} else {
				updateOpts.Tags = &[]string{}
			}
		}

		if hasChange {
			return updateOpts, nil
		}
	}

	// Use Neutron.
	var updateOpts neutronloadbalancers.UpdateOpts

	if d.HasChange("name") {
		hasChange = true
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}
	if d.HasChange("description") {
		hasChange = true
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}
	if d.HasChange("admin_state_up") {
		hasChange = true
		asu := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &asu
	}

	if hasChange {
		return updateOpts, nil
	}

	return nil, nil
}

// chooseLBListenerCreateOpts will determine which load balancer listener Create options to use:
// either the Octavia/LBaaS or the Neutron/Networking v2.
func chooseLBListenerCreateOpts(d *schema.ResourceData, config *config) (neutronlisteners.CreateOptsBuilder, error) {
	adminStateUp := d.Get("admin_state_up").(bool)

	var sniContainerRefs []string
	if raw, ok := d.GetOk("sni_container_refs"); ok {
		for _, v := range raw.([]interface{}) {
			sniContainerRefs = append(sniContainerRefs, v.(string))
		}
	}

	var createOpts neutronlisteners.CreateOptsBuilder

	if config.UseOctavia {
		// Use Octavia.
		opts := octavialisteners.CreateOpts{
			// Protocol SCTP requires octavia minor version 2.23
			Protocol:               octavialisteners.Protocol(d.Get("protocol").(string)),
			ProtocolPort:           d.Get("protocol_port").(int),
			LoadbalancerID:         d.Get("loadbalancer_id").(string),
			Name:                   d.Get("name").(string),
			DefaultPoolID:          d.Get("default_pool_id").(string),
			Description:            d.Get("description").(string),
			DefaultTlsContainerRef: d.Get("default_tls_container_ref").(string),
			SniContainerRefs:       sniContainerRefs,
			AdminStateUp:           &adminStateUp,
		}

		if v, ok := d.GetOk("connection_limit"); ok {
			connectionLimit := v.(int)
			opts.ConnLimit = &connectionLimit
		}

		if v, ok := d.GetOk("timeout_client_data"); ok {
			timeoutClientData := v.(int)
			opts.TimeoutClientData = &timeoutClientData
		}

		if v, ok := d.GetOk("timeout_member_connect"); ok {
			timeoutMemberConnect := v.(int)
			opts.TimeoutMemberConnect = &timeoutMemberConnect
		}

		if v, ok := d.GetOk("timeout_member_data"); ok {
			timeoutMemberData := v.(int)
			opts.TimeoutMemberData = &timeoutMemberData
		}

		if v, ok := d.GetOk("timeout_tcp_inspect"); ok {
			timeoutTCPInspect := v.(int)
			opts.TimeoutTCPInspect = &timeoutTCPInspect
		}

		// Get and check insert  headers map.
		rawHeaders := d.Get("insert_headers").(map[string]interface{})
		headers, err := expandLBListenerHeadersMap(rawHeaders)
		if err != nil {
			return nil, fmt.Errorf("unable to parse insert_headers argument: %s", err)
		}

		opts.InsertHeaders = headers

		if raw, ok := d.GetOk("allowed_cidrs"); ok {
			allowedCidrs := make([]string, len(raw.([]interface{})))
			for i, v := range raw.([]interface{}) {
				allowedCidrs[i] = v.(string)
			}
			opts.AllowedCIDRs = allowedCidrs
		}

		createOpts = opts

		return createOpts, nil
	}

	// Use Neutron.
	opts := neutronlisteners.CreateOpts{
		Protocol:               neutronlisteners.Protocol(d.Get("protocol").(string)),
		ProtocolPort:           d.Get("protocol_port").(int),
		LoadbalancerID:         d.Get("loadbalancer_id").(string),
		Name:                   d.Get("name").(string),
		DefaultPoolID:          d.Get("default_pool_id").(string),
		Description:            d.Get("description").(string),
		DefaultTlsContainerRef: d.Get("default_tls_container_ref").(string),
		SniContainerRefs:       sniContainerRefs,
		AdminStateUp:           &adminStateUp,
	}

	if v, ok := d.GetOk("connection_limit"); ok {
		connectionLimit := v.(int)
		opts.ConnLimit = &connectionLimit
	}

	createOpts = opts

	return createOpts, nil
}

// chooseLBListenerUpdateOpts will determine which load balancer listener Update options to use:
// either the Octavia/LBaaS or the Neutron/Networking v2.
func chooseLBListenerUpdateOpts(d *schema.ResourceData, config *config) (neutronlisteners.UpdateOptsBuilder, error) {
	var hasChange bool

	if config.UseOctavia {
		// Use Octavia.
		var opts octavialisteners.UpdateOpts
		if d.HasChange("name") {
			hasChange = true
			name := d.Get("name").(string)
			opts.Name = &name
		}

		if d.HasChange("description") {
			hasChange = true
			description := d.Get("description").(string)
			opts.Description = &description
		}

		if d.HasChange("connection_limit") {
			hasChange = true
			connLimit := d.Get("connection_limit").(int)
			opts.ConnLimit = &connLimit
		}

		if d.HasChange("timeout_client_data") {
			hasChange = true
			timeoutClientData := d.Get("timeout_client_data").(int)
			opts.TimeoutClientData = &timeoutClientData
		}

		if d.HasChange("timeout_member_connect") {
			hasChange = true
			timeoutMemberConnect := d.Get("timeout_member_connect").(int)
			opts.TimeoutMemberConnect = &timeoutMemberConnect
		}

		if d.HasChange("timeout_member_data") {
			hasChange = true
			timeoutMemberData := d.Get("timeout_member_data").(int)
			opts.TimeoutMemberData = &timeoutMemberData
		}

		if d.HasChange("timeout_tcp_inspect") {
			hasChange = true
			timeoutTCPInspect := d.Get("timeout_tcp_inspect").(int)
			opts.TimeoutTCPInspect = &timeoutTCPInspect
		}

		if d.HasChange("default_pool_id") {
			hasChange = true
			defaultPoolID := d.Get("default_pool_id").(string)
			opts.DefaultPoolID = &defaultPoolID
		}

		if d.HasChange("default_tls_container_ref") {
			hasChange = true
			defaultTLSContainerRef := d.Get("default_tls_container_ref").(string)
			opts.DefaultTlsContainerRef = &defaultTLSContainerRef
		}

		if d.HasChange("sni_container_refs") {
			hasChange = true
			var sniContainerRefs []string
			if raw, ok := d.GetOk("sni_container_refs"); ok {
				for _, v := range raw.([]interface{}) {
					sniContainerRefs = append(sniContainerRefs, v.(string))
				}
			}
			opts.SniContainerRefs = &sniContainerRefs
		}

		if d.HasChange("admin_state_up") {
			hasChange = true
			asu := d.Get("admin_state_up").(bool)
			opts.AdminStateUp = &asu
		}

		if d.HasChange("insert_headers") {
			hasChange = true

			// Get and check insert headers map.
			rawHeaders := d.Get("insert_headers").(map[string]interface{})
			headers, err := expandLBListenerHeadersMap(rawHeaders)
			if err != nil {
				return nil, fmt.Errorf("unable to parse insert_headers argument: %s", err)
			}

			opts.InsertHeaders = &headers
		}

		if d.HasChange("allowed_cidrs") {
			hasChange = true
			var allowedCidrs []string
			if raw, ok := d.GetOk("allowed_cidrs"); ok {
				for _, v := range raw.([]interface{}) {
					allowedCidrs = append(allowedCidrs, v.(string))
				}
			}
			opts.AllowedCIDRs = &allowedCidrs
		}

		if hasChange {
			return opts, nil
		}
	}

	// Use Neutron.
	var opts neutronlisteners.UpdateOpts
	if d.HasChange("name") {
		hasChange = true
		name := d.Get("name").(string)
		opts.Name = &name
	}

	if d.HasChange("description") {
		hasChange = true
		description := d.Get("description").(string)
		opts.Description = &description
	}

	if d.HasChange("connection_limit") {
		hasChange = true
		connLimit := d.Get("connection_limit").(int)
		opts.ConnLimit = &connLimit
	}

	if d.HasChange("default_pool_id") {
		hasChange = true
		defaultPoolID := d.Get("default_pool_id").(string)
		opts.DefaultPoolID = &defaultPoolID
	}

	if d.HasChange("default_tls_container_ref") {
		hasChange = true
		defaultTLSContainerRef := d.Get("default_tls_container_ref").(string)
		opts.DefaultTlsContainerRef = &defaultTLSContainerRef
	}

	if d.HasChange("sni_container_refs") {
		hasChange = true
		var sniContainerRefs []string
		if raw, ok := d.GetOk("sni_container_refs"); ok {
			for _, v := range raw.([]interface{}) {
				sniContainerRefs = append(sniContainerRefs, v.(string))
			}
		}
		opts.SniContainerRefs = &sniContainerRefs
	}

	if d.HasChange("admin_state_up") {
		hasChange = true
		asu := d.Get("admin_state_up").(bool)
		opts.AdminStateUp = &asu
	}

	if hasChange {
		return opts, nil
	}

	return nil, nil
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

func waitForLBListener(ctx context.Context, lbClient *gophercloud.ServiceClient, listener *neutronlisteners.Listener, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for vkcs_lb_listener %s to become %s.", listener.ID, target)

	if len(listener.Loadbalancers) == 0 {
		return fmt.Errorf("Failed to detect a vkcs_lb_loadbalancer for the %s vkcs_lb_listener", listener.ID)
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

		return fmt.Errorf("Error waiting for vkcs_lb_listener %s to become %s: %s", listener.ID, target, err)
	}

	return nil
}

func resourceLBListenerRefreshFunc(lbClient *gophercloud.ServiceClient, lbID string, listener *neutronlisteners.Listener) resource.StateRefreshFunc {
	if listener.ProvisioningStatus != "" {
		return func() (interface{}, string, error) {
			lb, status, err := resourceLBLoadBalancerRefreshFunc(lbClient, lbID)()
			if err != nil {
				return lb, status, err
			}
			if !strSliceContains(getLbSkipStatuses(), status) {
				return lb, status, nil
			}

			listener, err := neutronlisteners.Get(lbClient, listener.ID).Extract()
			if err != nil {
				return nil, "", err
			}

			return listener, listener.ProvisioningStatus, nil
		}
	}

	return resourceLBLoadBalancerStatusRefreshFuncNeutron(lbClient, lbID, "listener", listener.ID, "")
}

// chooseLBMonitorCreateOpts will determine which load balancer monitor Create options to use:
// either the Octavia/LBaaS or the Neutron/Networking v2.
func chooseLBMonitorCreateOpts(d *schema.ResourceData, config *config) neutronmonitors.CreateOptsBuilder {
	adminStateUp := d.Get("admin_state_up").(bool)

	var createOpts neutronmonitors.CreateOptsBuilder

	if config.UseOctavia {
		// Use Octavia.
		opts := octaviamonitors.CreateOpts{
			PoolID:         d.Get("pool_id").(string),
			Type:           d.Get("type").(string),
			Delay:          d.Get("delay").(int),
			Timeout:        d.Get("timeout").(int),
			MaxRetries:     d.Get("max_retries").(int),
			MaxRetriesDown: d.Get("max_retries_down").(int),
			URLPath:        d.Get("url_path").(string),
			HTTPMethod:     d.Get("http_method").(string),
			ExpectedCodes:  d.Get("expected_codes").(string),
			Name:           d.Get("name").(string),
			AdminStateUp:   &adminStateUp,
		}

		createOpts = opts
	} else {
		// Use Neutron.
		opts := neutronmonitors.CreateOpts{
			PoolID:        d.Get("pool_id").(string),
			Type:          d.Get("type").(string),
			Delay:         d.Get("delay").(int),
			Timeout:       d.Get("timeout").(int),
			MaxRetries:    d.Get("max_retries").(int),
			URLPath:       d.Get("url_path").(string),
			HTTPMethod:    d.Get("http_method").(string),
			ExpectedCodes: d.Get("expected_codes").(string),
			Name:          d.Get("name").(string),
			AdminStateUp:  &adminStateUp,
		}

		createOpts = opts
	}

	return createOpts
}

// chooseLBMonitorUpdateOpts will determine which load balancer monitor Update options to use:
// either the Octavia/LBaaS or the Neutron/Networking v2.
func chooseLBMonitorUpdateOpts(d *schema.ResourceData, config *config) neutronmonitors.UpdateOptsBuilder {
	var hasChange bool

	if config.UseOctavia {
		// Use Octavia.
		var opts octaviamonitors.UpdateOpts

		if d.HasChange("url_path") {
			hasChange = true
			opts.URLPath = d.Get("url_path").(string)
		}
		if d.HasChange("expected_codes") {
			hasChange = true
			opts.ExpectedCodes = d.Get("expected_codes").(string)
		}
		if d.HasChange("delay") {
			hasChange = true
			opts.Delay = d.Get("delay").(int)
		}
		if d.HasChange("timeout") {
			hasChange = true
			opts.Timeout = d.Get("timeout").(int)
		}
		if d.HasChange("max_retries") {
			hasChange = true
			opts.MaxRetries = d.Get("max_retries").(int)
		}
		if d.HasChange("max_retries_down") {
			hasChange = true
			opts.MaxRetriesDown = d.Get("max_retries_down").(int)
		}
		if d.HasChange("admin_state_up") {
			hasChange = true
			asu := d.Get("admin_state_up").(bool)
			opts.AdminStateUp = &asu
		}
		if d.HasChange("name") {
			hasChange = true
			name := d.Get("name").(string)
			opts.Name = &name
		}
		if d.HasChange("http_method") {
			hasChange = true
			opts.HTTPMethod = d.Get("http_method").(string)
		}

		if hasChange {
			return opts
		}
	} else {
		// Use Neutron.
		var opts neutronmonitors.UpdateOpts

		if d.HasChange("url_path") {
			hasChange = true
			opts.URLPath = d.Get("url_path").(string)
		}
		if d.HasChange("expected_codes") {
			hasChange = true
			opts.ExpectedCodes = d.Get("expected_codes").(string)
		}
		if d.HasChange("delay") {
			hasChange = true
			opts.Delay = d.Get("delay").(int)
		}
		if d.HasChange("timeout") {
			hasChange = true
			opts.Timeout = d.Get("timeout").(int)
		}
		if d.HasChange("max_retries") {
			hasChange = true
			opts.MaxRetries = d.Get("max_retries").(int)
		}
		if d.HasChange("admin_state_up") {
			hasChange = true
			asu := d.Get("admin_state_up").(bool)
			opts.AdminStateUp = &asu
		}
		if d.HasChange("name") {
			hasChange = true
			name := d.Get("name").(string)
			opts.Name = &name
		}
		if d.HasChange("http_method") {
			hasChange = true
			opts.HTTPMethod = d.Get("http_method").(string)
		}

		if hasChange {
			return opts
		}
	}

	return nil
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
				return fmt.Errorf("Error: loadbalancer %s not found: %s", lbID, err)
			}
		}
		return fmt.Errorf("Error waiting for loadbalancer %s to become %s: %s", lbID, target, err)
	}

	return nil
}

func resourceLBLoadBalancerRefreshFunc(lbClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		lb, err := neutronloadbalancers.Get(lbClient, id).Extract()
		if err != nil {
			return nil, "", err
		}

		return lb, lb.ProvisioningStatus, nil
	}
}

func waitForLBMember(ctx context.Context, lbClient *gophercloud.ServiceClient, parentPool *neutronpools.Pool, member *neutronpools.Member, target string, pending []string, timeout time.Duration) error {
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

		return fmt.Errorf("Error waiting for member %s to become %s: %s", member.ID, target, err)
	}

	return nil
}

func resourceLBMemberRefreshFunc(lbClient *gophercloud.ServiceClient, lbID string, poolID string, member *neutronpools.Member) resource.StateRefreshFunc {
	if member.ProvisioningStatus != "" {
		return func() (interface{}, string, error) {
			lb, status, err := resourceLBLoadBalancerRefreshFunc(lbClient, lbID)()
			if err != nil {
				return lb, status, err
			}
			if !strSliceContains(getLbSkipStatuses(), status) {
				return lb, status, nil
			}

			member, err := neutronpools.GetMember(lbClient, poolID, member.ID).Extract()
			if err != nil {
				return nil, "", err
			}

			return member, member.ProvisioningStatus, nil
		}
	}

	return resourceLBLoadBalancerStatusRefreshFuncNeutron(lbClient, lbID, "member", member.ID, poolID)
}

func waitForLBMonitor(ctx context.Context, lbClient *gophercloud.ServiceClient, parentPool *neutronpools.Pool, monitor *neutronmonitors.Monitor, target string, pending []string, timeout time.Duration) error {
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
		return fmt.Errorf("Error waiting for vkcs_lb_monitor %s to become %s: %s", monitor.ID, target, err)
	}

	return nil
}

func resourceLBMonitorRefreshFunc(lbClient *gophercloud.ServiceClient, lbID string, monitor *neutronmonitors.Monitor) resource.StateRefreshFunc {
	if monitor.ProvisioningStatus != "" {
		return func() (interface{}, string, error) {
			lb, status, err := resourceLBLoadBalancerRefreshFunc(lbClient, lbID)()
			if err != nil {
				return lb, status, err
			}
			if !strSliceContains(getLbSkipStatuses(), status) {
				return lb, status, nil
			}

			monitor, err := neutronmonitors.Get(lbClient, monitor.ID).Extract()
			if err != nil {
				return nil, "", err
			}

			return monitor, monitor.ProvisioningStatus, nil
		}
	}

	return resourceLBLoadBalancerStatusRefreshFuncNeutron(lbClient, lbID, "monitor", monitor.ID, "")
}

func waitForLBPool(ctx context.Context, lbClient *gophercloud.ServiceClient, pool *neutronpools.Pool, target string, pending []string, timeout time.Duration) error {
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

		return fmt.Errorf("Error waiting for pool %s to become %s: %s", pool.ID, target, err)
	}

	return nil
}

func resourceLBPoolRefreshFunc(lbClient *gophercloud.ServiceClient, lbID string, pool *neutronpools.Pool) resource.StateRefreshFunc {
	if pool.ProvisioningStatus != "" {
		return func() (interface{}, string, error) {
			lb, status, err := resourceLBLoadBalancerRefreshFunc(lbClient, lbID)()
			if err != nil {
				return lb, status, err
			}
			if !strSliceContains(getLbSkipStatuses(), status) {
				return lb, status, nil
			}

			pool, err := neutronpools.Get(lbClient, pool.ID).Extract()
			if err != nil {
				return nil, "", err
			}

			return pool, pool.ProvisioningStatus, nil
		}
	}

	return resourceLBLoadBalancerStatusRefreshFuncNeutron(lbClient, lbID, "pool", pool.ID, "")
}

func lbFindLBIDviaPool(lbClient *gophercloud.ServiceClient, pool *neutronpools.Pool) (string, error) {
	if len(pool.Loadbalancers) > 0 {
		return pool.Loadbalancers[0].ID, nil
	}

	if len(pool.Listeners) > 0 {
		listenerID := pool.Listeners[0].ID
		listener, err := neutronlisteners.Get(lbClient, listenerID).Extract()
		if err != nil {
			return "", err
		}

		if len(listener.Loadbalancers) > 0 {
			return listener.Loadbalancers[0].ID, nil
		}
	}

	return "", fmt.Errorf("Unable to determine loadbalancer ID from pool %s", pool.ID)
}

func resourceLBLoadBalancerStatusRefreshFuncNeutron(lbClient *gophercloud.ServiceClient, lbID, resourceType, resourceID string, parentID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		statuses, err := neutronloadbalancers.GetStatuses(lbClient, lbID).Extract()
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
			return nil, "", fmt.Errorf("Unable to get statuses from the Load Balancer %s statuses tree: %s", lbID, err)
		}

		// Don't fail, when statuses returns "null"
		if statuses == nil || statuses.Loadbalancer == nil {
			statuses = new(neutronloadbalancers.StatusTree)
			statuses.Loadbalancer = new(neutronloadbalancers.LoadBalancer)
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
			listener, err := neutronlisteners.Get(lbClient, resourceID).Extract()
			return listener, "ACTIVE", err

		case "pool":
			for _, pool := range statuses.Loadbalancer.Pools {
				if pool.ID == resourceID {
					if pool.ProvisioningStatus != "" {
						return pool, pool.ProvisioningStatus, nil
					}
				}
			}
			pool, err := neutronpools.Get(lbClient, resourceID).Extract()
			return pool, "ACTIVE", err

		case "monitor":
			for _, pool := range statuses.Loadbalancer.Pools {
				if pool.Monitor.ID == resourceID {
					if pool.Monitor.ProvisioningStatus != "" {
						return pool.Monitor, pool.Monitor.ProvisioningStatus, nil
					}
				}
			}
			monitor, err := neutronmonitors.Get(lbClient, resourceID).Extract()
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
			member, err := neutronpools.GetMember(lbClient, parentID, resourceID).Extract()
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
			l7policy, err := neutronl7policies.Get(lbClient, resourceID).Extract()
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
			l7Rule, err := neutronl7policies.GetRule(lbClient, parentID, resourceID).Extract()
			return l7Rule, "ACTIVE", err
		}

		return nil, "", fmt.Errorf("An unexpected error occurred querying the status of %s %s by loadbalancer %s", resourceType, resourceID, lbID)
	}
}

func resourceLBL7PolicyRefreshFunc(lbClient *gophercloud.ServiceClient, lbID string, l7policy *neutronl7policies.L7Policy) resource.StateRefreshFunc {
	if l7policy.ProvisioningStatus != "" {
		return func() (interface{}, string, error) {
			lb, status, err := resourceLBLoadBalancerRefreshFunc(lbClient, lbID)()
			if err != nil {
				return lb, status, err
			}
			if !strSliceContains(getLbSkipStatuses(), status) {
				return lb, status, nil
			}

			l7policy, err := neutronl7policies.Get(lbClient, l7policy.ID).Extract()
			if err != nil {
				return nil, "", err
			}

			return l7policy, l7policy.ProvisioningStatus, nil
		}
	}

	return resourceLBLoadBalancerStatusRefreshFuncNeutron(lbClient, lbID, "l7policy", l7policy.ID, "")
}

func waitForLBL7Policy(ctx context.Context, lbClient *gophercloud.ServiceClient, parentListener *neutronlisteners.Listener, l7policy *neutronl7policies.L7Policy, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for l7policy %s to become %s.", l7policy.ID, target)

	if len(parentListener.Loadbalancers) == 0 {
		return fmt.Errorf("Unable to determine loadbalancer ID from listener %s", parentListener.ID)
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

		return fmt.Errorf("Error waiting for l7policy %s to become %s: %s", l7policy.ID, target, err)
	}

	return nil
}

func getListenerIDForL7Policy(lbClient *gophercloud.ServiceClient, id string) (string, error) {
	log.Printf("[DEBUG] Trying to get Listener ID associated with the %s L7 Policy ID", id)
	lbsPages, err := neutronloadbalancers.List(lbClient, neutronloadbalancers.ListOpts{}).AllPages()
	if err != nil {
		return "", fmt.Errorf("No Load Balancers were found: %s", err)
	}

	lbs, err := neutronloadbalancers.ExtractLoadBalancers(lbsPages)
	if err != nil {
		return "", fmt.Errorf("Unable to extract Load Balancers list: %s", err)
	}

	for _, lb := range lbs {
		statuses, err := neutronloadbalancers.GetStatuses(lbClient, lb.ID).Extract()
		if err != nil {
			return "", fmt.Errorf("Failed to get Load Balancer statuses: %s", err)
		}
		for _, listener := range statuses.Loadbalancer.Listeners {
			for _, l7policy := range listener.L7Policies {
				if l7policy.ID == id {
					return listener.ID, nil
				}
			}
		}
	}

	return "", fmt.Errorf("Unable to find Listener ID associated with the %s L7 Policy ID", id)
}

func resourceLBL7RuleRefreshFunc(lbClient *gophercloud.ServiceClient, lbID string, l7policyID string, l7rule *neutronl7policies.Rule) resource.StateRefreshFunc {
	if l7rule.ProvisioningStatus != "" {
		return func() (interface{}, string, error) {
			lb, status, err := resourceLBLoadBalancerRefreshFunc(lbClient, lbID)()
			if err != nil {
				return lb, status, err
			}
			if !strSliceContains(getLbSkipStatuses(), status) {
				return lb, status, nil
			}

			l7rule, err := neutronl7policies.GetRule(lbClient, l7policyID, l7rule.ID).Extract()
			if err != nil {
				return nil, "", err
			}

			return l7rule, l7rule.ProvisioningStatus, nil
		}
	}

	return resourceLBLoadBalancerStatusRefreshFuncNeutron(lbClient, lbID, "l7rule", l7rule.ID, l7policyID)
}

func waitForLBL7Rule(ctx context.Context, lbClient *gophercloud.ServiceClient, parentListener *neutronlisteners.Listener, parentL7policy *neutronl7policies.L7Policy, l7rule *neutronl7policies.Rule, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for l7rule %s to become %s.", l7rule.ID, target)

	if len(parentListener.Loadbalancers) == 0 {
		return fmt.Errorf("Unable to determine loadbalancer ID from listener %s", parentListener.ID)
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

		return fmt.Errorf("Error waiting for l7rule %s to become %s: %s", l7rule.ID, target, err)
	}

	return nil
}

func flattenLBPoolPersistence(p neutronpools.SessionPersistence) []map[string]interface{} {
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

func resourceLoadBalancerSetSecurityGroups(networkingClient *gophercloud.ServiceClient, vipPortID string, d *schema.ResourceData) error {
	if vipPortID != "" {
		if v, ok := d.GetOk("security_group_ids"); ok {
			securityGroups := expandToStringSlice(v.(*schema.Set).List())
			updateOpts := ports.UpdateOpts{
				SecurityGroups: &securityGroups,
			}

			log.Printf("[DEBUG] Adding security groups to vkcs_lb_loadbalancer "+
				"VIP port %s: %#v", vipPortID, updateOpts)

			_, err := ports.Update(networkingClient, vipPortID, updateOpts).Extract()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func resourceLoadBalancerGetSecurityGroups(networkingClient *gophercloud.ServiceClient, vipPortID string, d *schema.ResourceData) error {
	port, err := ports.Get(networkingClient, vipPortID).Extract()
	if err != nil {
		return err
	}

	d.Set("security_group_ids", port.SecurityGroups)

	return nil
}
