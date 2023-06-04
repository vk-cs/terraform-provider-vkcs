package util

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
)

var DecoderConfig = &mapstructure.DecoderConfig{
	TagName: "json",
}

// MapStructureDecoder ...
func MapStructureDecoder(strct interface{}, v *map[string]interface{}, config *mapstructure.DecoderConfig) error {
	config.Result = strct
	decoder, _ := mapstructure.NewDecoder(config)
	return decoder.Decode(*v)
}

// StructToMap converts a structure to map with keys set to lowered
// structure field names
// NOTE: This function does not implement mapping nested structures
func StructToMap(s interface{}) (map[string]interface{}, error) {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("got %T instead of struct", s)
	}

	m := make(map[string]interface{})
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := t.Field(i)
		if f.IsExported() {
			m[strings.ToLower(f.Name)] = v.Field(i).Interface()
		}
	}

	return m, nil
}

// GetTimestamp ...
func GetTimestamp(t *time.Time) string {
	if t != nil {
		return t.Format(time.RFC3339)
	}
	return ""
}

// BuildRequest takes an opts struct and builds a request body for
// Gophercloud to execute.
func BuildRequest(opts interface{}, parent string) (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	b = AddValueSpecs(b)

	return map[string]interface{}{parent: b}, nil
}

// CheckDeleted checks the error to see if it's a 404 (Not Found) and, if so,
// sets the resource ID to the empty string instead of throwing an error.
func CheckDeleted(d *schema.ResourceData, err error, msg string) error {
	if _, ok := err.(gophercloud.ErrDefault404); ok {
		d.SetId("")
		return nil
	}

	return fmt.Errorf("%s %s: %s", msg, d.Id(), err)
}

// checkDeletedResource checks the error to see if it's a 404 (Not Found) and, if so,
// sets the resource ID to the empty string instead of throwing an error.
func CheckDeletedResource(ctx context.Context, r *resource.ReadResponse, err error) error {
	if _, ok := err.(gophercloud.ErrDefault404); ok {
		r.State.RemoveResource(ctx)
		return nil
	}
	var id string
	r.State.GetAttribute(ctx, path.Root("id"), &id)

	return fmt.Errorf("%s: %s", id, err)
}

// checkDeletedDatasource checks the error to see if it's a 404 (Not Found) and, if so,
// sets the resource ID to the empty string instead of throwing an error.
func CheckDeletedDatasource(ctx context.Context, r *datasource.ReadResponse, err error) error {
	if _, ok := err.(gophercloud.ErrDefault404); ok {
		r.State.SetAttribute(ctx, path.Root("id"), "")
		return nil
	}
	var id string
	r.State.GetAttribute(ctx, path.Root("id"), &id)

	return fmt.Errorf("%s: %s", id, err)
}

func CheckAlreadyExists(err error, msg, resourceName, conflict string) error {
	if _, ok := err.(gophercloud.ErrDefault409); ok {
		return fmt.Errorf("%s: %s with %s already exists", msg, resourceName, conflict)
	}

	return fmt.Errorf("%s: %s", msg, err)
}

// GetRegion returns the region that was specified in the resource. If a
// region was not set, the provider-level region is checked. The provider-level
// region can either be set by the region argument or by OS_REGION_NAME.
func GetRegion(d *schema.ResourceData, config clients.Config) string {
	if v, ok := d.GetOk("region"); ok {
		return v.(string)
	}

	return config.GetRegion()
}

// AddValueSpecs expands the 'value_specs' object and removes 'value_specs'
// from the reqeust body.
func AddValueSpecs(body map[string]interface{}) map[string]interface{} {
	if body["value_specs"] != nil {
		for k, v := range body["value_specs"].(map[string]interface{}) {
			body[k] = v
		}
		delete(body, "value_specs")
	}

	return body
}

// MapValueSpecs converts ResourceData into a map.
func MapValueSpecs(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("value_specs").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

func CheckForRetryableError(err error) *retry.RetryError {
	switch e := err.(type) {
	case gophercloud.ErrDefault500:
		return retry.RetryableError(err)
	case gophercloud.ErrDefault409:
		return retry.RetryableError(err)
	case gophercloud.ErrDefault503:
		return retry.RetryableError(err)
	case gophercloud.ErrUnexpectedResponseCode:
		if e.GetStatusCode() == 504 || e.GetStatusCode() == 502 {
			return retry.RetryableError(err)
		} else {
			return retry.NonRetryableError(err)
		}
	default:
		return retry.NonRetryableError(err)
	}
}

func ExpandVendorOptions(vendOptsRaw []interface{}) map[string]interface{} {
	vendorOptions := make(map[string]interface{})

	for _, option := range vendOptsRaw {
		for optKey, optValue := range option.(map[string]interface{}) {
			vendorOptions[optKey] = optValue
		}
	}

	return vendorOptions
}

func ExpandObjectReadTags(d *schema.ResourceData, tags []string) {
	d.Set("all_tags", tags)

	allTags := d.Get("all_tags").(*schema.Set)
	desiredTags := d.Get("tags").(*schema.Set)
	actualTags := allTags.Intersection(desiredTags)
	if !actualTags.Equal(desiredTags) {
		d.Set("tags", ExpandToStringSlice(actualTags.List()))
	}
}

func ExpandObjectUpdateTags(d *schema.ResourceData) []string {
	allTags := d.Get("all_tags").(*schema.Set)
	oldTagsRaw, newTagsRaw := d.GetChange("tags")
	oldTags, newTags := oldTagsRaw.(*schema.Set), newTagsRaw.(*schema.Set)

	allTagsWithoutOld := allTags.Difference(oldTags)

	return ExpandToStringSlice(allTagsWithoutOld.Union(newTags).List())
}

func ExpandObjectTags(d *schema.ResourceData) []string {
	rawTags := d.Get("tags").(*schema.Set).List()
	tags := make([]string, len(rawTags))

	for i, raw := range rawTags {
		tags[i] = raw.(string)
	}

	return tags
}

func ExpandToMapStringString(v map[string]interface{}) map[string]string {
	m := make(map[string]string)
	for key, val := range v {
		if strVal, ok := val.(string); ok {
			m[key] = strVal
		}
	}

	return m
}

func ExpandToStringSlice(v []interface{}) []string {
	s := make([]string, len(v))
	for i, val := range v {
		if strVal, ok := val.(string); ok {
			s[i] = strVal
		}
	}

	return s
}

// util.StrSliceContains checks if a given string is contained in a slice
// When anybody asks why Go needs generics, here you go.
func StrSliceContains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

func SliceUnion(a, b []string) []string {
	var res []string
	for _, i := range a {
		if !StrSliceContains(res, i) {
			res = append(res, i)
		}
	}
	for _, k := range b {
		if !StrSliceContains(res, k) {
			res = append(res, k)
		}
	}
	return res
}

func IsOperationNotSupported(d string, types ...string) bool {
	for _, t := range types {
		if d == t {
			return true
		}
	}
	return false
}

func EnsureOnlyOnePresented(d *schema.ResourceData, keys ...string) (string, error) {
	var isPresented bool
	var keyPresented string
	for _, key := range keys {
		_, ok := d.GetOk(key)

		if ok {
			if isPresented {
				return "", fmt.Errorf("only one of %v keys can be presented", keys)
			}

			isPresented = true
			keyPresented = key
		}
	}

	if !isPresented {
		return "", fmt.Errorf("no one of %v keys are presented", keys)
	}

	return keyPresented, nil
}

func CopyToMap(dst, src *map[string]string) {
	for k, v := range *src {
		(*dst)[k] = v
	}
}
