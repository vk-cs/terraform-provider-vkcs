package vkcs

import (
	"net/http"
	"time"

	"github.com/gophercloud/gophercloud"
)

// monitoringClient performs request to cloud monitoring api
type monitoringClient interface {
	Get(url string, JSONResponse interface{}, opts *gophercloud.RequestOpts) (*http.Response, error)
	Post(url string, JSONBody interface{}, JSONResponse interface{}, opts *gophercloud.RequestOpts) (*http.Response, error)
	Patch(url string, JSONBody interface{}, JSONResponse interface{}, opts *gophercloud.RequestOpts) (*http.Response, error)
	Delete(url string, opts *gophercloud.RequestOpts) (*http.Response, error)
	Head(url string, opts *gophercloud.RequestOpts) (*http.Response, error)
	Put(url string, JSONBody interface{}, JSONResponse interface{}, opts *gophercloud.RequestOpts) (*http.Response, error)
	ServiceURL(parts ...string) string
}

type CreateTrigger struct {
	Name                 string   `json:"name"`
	Status               string   `json:"status"`
	Namespace            string   `json:"namespace"`
	Query                string   `json:"query"`
	Interval             int      `json:"interval"`
	NotificationTitle    string   `json:"notification_title"`
	NotificationChannels []string `json:"notification_channels"`
}

type TriggerIn struct {
	Trigger CreateTrigger `json:"trigger"`
}

// Map converts opts to a map (for a request body)
func (opts TriggerIn) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

type TriggerData struct {
	CreatedAt            time.Time `json:"created_at"`
	Id                   string    `json:"id"`
	Namespace            string    `json:"namespace"`
	Name                 string    `json:"name"`
	Interval             int       `json:"interval"`
	NotificationChannels []string  `json:"notification_channels"`
	NotificationTitle    string    `json:"notification_title"`
	Query                string    `json:"query"`
	Status               string    `json:"status"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type TriggerOut struct {
	Trigger TriggerData `json:"trigger"`
}

type TriggersList struct {
	Triggers []TriggerData `json:"triggers"`
}

type ChannelIn struct {
	Name        string `json:"name" required:"true"`
	ChannelType string `json:"channel_type" required:"true"`
	Address     string `json:"address" required:"true"`
}

// Map converts opts to a map (for a request body)
func (opts ChannelIn) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

type ChannelOut struct {
	Channel struct {
		Address     string    `json:"address"`
		ChannelType string    `json:"channel_type"`
		CreatedAt   time.Time `json:"created_at"`
		ID          string    `json:"id"`
		Name        string    `json:"name"`
		UpdatedAt   time.Time `json:"updated_at"`
	} `json:"channel"`
}

type ChannelList struct {
	Channels []struct {
		Address     string    `json:"address"`
		ChannelType string    `json:"channel_type"`
		CreatedAt   time.Time `json:"created_at"`
		ID          string    `json:"id"`
		InUse       bool      `json:"in_use"`
		Name        string    `json:"name"`
		UpdatedAt   time.Time `json:"updated_at"`
	} `json:"channels"`
}

type TemplateIn struct {
	InstanceID   string   `json:"instance_id"`
	Capabilities []string `json:"capabilities"`
}

// Map converts opts to a map (for a request body)
func (opts TemplateIn) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

type TemplateOut struct {
	Script string `json:"script"`
	LinkId string `json:"link_id"`
}
