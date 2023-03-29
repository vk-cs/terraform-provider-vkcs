package vkcs

import (
	"github.com/gophercloud/gophercloud"
	"net/http"
)

type commonTriggerResult struct {
	gophercloud.Result
}

// extract is used to extract result into short response struct
func (r commonTriggerResult) extract() (*TriggerOut, error) {
	var c *TriggerOut
	if err := r.ExtractInto(&c); err != nil {
		return nil, err
	}
	return c, nil
}

type triggerDeleteResult struct {
	gophercloud.Result
}

func (r triggerDeleteResult) extractErr() error {
	return r.Err
}

func getMonRequestOpts(codes ...int) *gophercloud.RequestOpts {
	reqOpts := &gophercloud.RequestOpts{
		OkCodes: codes,
	}
	if len(codes) != 0 {
		reqOpts.OkCodes = codes
	}
	return reqOpts
}

func instanceTriggerURL(c ContainerClient, pid string, id string) string {
	return c.ServiceURL(pid, "triggers", id)
}

func createTriggerURL(c ContainerClient, pid string) string {
	return c.ServiceURL(pid, "triggers")
}

func triggerCreate(client monitoringClient, pid string, opts createOptsBuilder) commonTriggerResult {
	b, err := opts.Map()
	var r commonTriggerResult
	if err != nil {
		r.Err = err
		return r
	}
	var result *http.Response
	reqOpts := getMonRequestOpts(http.StatusCreated)
	result, r.Err = client.Post(createTriggerURL(client, pid), b, &r.Body, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return r
}

func triggerUpdate(client monitoringClient, pid string, id string, opts createOptsBuilder) commonTriggerResult {
	b, err := opts.Map()
	var r commonTriggerResult
	if err != nil {
		r.Err = err
		return r
	}
	var result *http.Response
	reqOpts := getMonRequestOpts(http.StatusOK)
	result, r.Err = client.Patch(instanceTriggerURL(client, pid, id), b, &r.Body, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return r
}

// triggerGet performs request to get trigger instance
func triggerGet(client monitoringClient, pid string, id string) (r commonTriggerResult) {
	reqOpts := getMonRequestOpts(http.StatusOK)
	var result *http.Response
	result, r.Err = client.Get(instanceTriggerURL(client, pid, id), &r.Body, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

// triggerDelete performs request to delete trigger
func triggerDelete(client monitoringClient, pid string, id string) (r triggerDeleteResult) {
	reqOpts := getMonRequestOpts(http.StatusNoContent)
	var result *http.Response
	result, r.Err = client.Delete(instanceTriggerURL(client, pid, id), reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}
