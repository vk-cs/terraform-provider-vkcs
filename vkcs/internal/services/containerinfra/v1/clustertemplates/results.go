package clustertemplates

import "github.com/gophercloud/gophercloud"

type clusterTemplateResult struct {
	gophercloud.Result
}

type clusterTemplatesResult struct {
	gophercloud.Result
}

// Extract parses result into params for cluster template.
func (r clusterTemplateResult) Extract() (*ClusterTemplate, error) {
	var s *ClusterTemplate
	err := r.ExtractInto(&s)
	return s, err
}

// Extract parses result into params for cluster templates.
func (r clusterTemplatesResult) Extract() ([]ClusterTemplate, error) {
	var s *ClusterTemplates
	err := r.ExtractInto(&s)
	if err != nil {
		return nil, err
	}

	return s.Templates, err
}
