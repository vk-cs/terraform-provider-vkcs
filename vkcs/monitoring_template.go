package vkcs

import (
	"github.com/gophercloud/gophercloud"
	"net/http"
)

func createTemplateURL(c ContainerClient, pid string) string {
	return c.ServiceURL("project", pid, "link")
}

type commonTemplateResult struct {
	gophercloud.Result
}

// extract is used to extract result into short response struct
func (r commonTemplateResult) extract() (*TemplateOut, error) {
	var t *TemplateOut
	if err := r.ExtractInto(&t); err != nil {
		return nil, err
	}
	return t, nil
}

func templateCreate(client monitoringClient, pid string, opts createOptsBuilder) commonTemplateResult {
	b, err := opts.Map()
	var r commonTemplateResult
	if err != nil {
		r.Err = err
		return r
	}
	var result *http.Response
	reqOpts := getMonRequestOpts(http.StatusCreated)
	result, r.Err = client.Post(createTemplateURL(client, pid), b, &r.Body, reqOpts)
	if r.Err == nil {
		r.Header = result.Header
	}
	return r
}
