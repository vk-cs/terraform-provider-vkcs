package vkcs

import (
	"github.com/gophercloud/gophercloud"
	"net/http"
)

func instanceChannelURL(c ContainerClient, pid string, id string) string {
	return c.ServiceURL(pid, "notification_channels", id)
}

func createChannelURL(c ContainerClient, pid string) string {
	return c.ServiceURL(pid, "notification_channels")
}

func listChannelURL(c ContainerClient, pid string) string {
	return c.ServiceURL(pid, "notification_channels")
}

type commonChannelResult struct {
	gophercloud.Result
}

// extract is used to extract result into short response struct
func (r commonChannelResult) extract() (*ChannelOut, error) {
	var c *ChannelOut
	if err := r.ExtractInto(&c); err != nil {
		return nil, err
	}
	return c, nil
}

type deleteChannelResult struct {
	gophercloud.Result
}

func (r deleteChannelResult) extractErr() error {
	return r.Err
}

func channelCreate(client monitoringClient, pid string, opts createOptsBuilder) commonChannelResult {
	b, err := opts.Map()
	var r commonChannelResult
	if err != nil {
		r.Err = err
		return r
	}
	var result *http.Response
	reqOpts := getMonRequestOpts(http.StatusCreated)
	result, r.Err = client.Post(createChannelURL(client, pid), b, &r.Body, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return r
}

func channelUpdate(client monitoringClient, pid string, id string, opts createOptsBuilder) commonChannelResult {
	b, err := opts.Map()
	var r commonChannelResult
	if err != nil {
		r.Err = err
		return r
	}
	var result *http.Response
	reqOpts := getMonRequestOpts(http.StatusCreated)
	result, r.Err = client.Patch(instanceChannelURL(client, pid, id), b, &r.Body, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return r
}

// triggerGet performs request to get trigger instance
func channelGet(client monitoringClient, pid string, id string) (r commonChannelResult) {
	reqOpts := getMonRequestOpts(http.StatusOK)
	var result *http.Response
	result, r.Err = client.Get(instanceChannelURL(client, pid, id), &r.Body, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

// triggerDelete performs request to delete trigger
func channelDelete(client monitoringClient, pid string, id string) (r deleteChannelResult) {
	reqOpts := getMonRequestOpts(http.StatusNoContent)
	var result *http.Response
	result, r.Err = client.Delete(instanceChannelURL(client, pid, id), reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}
