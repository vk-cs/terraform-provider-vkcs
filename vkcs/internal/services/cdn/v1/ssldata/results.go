package ssldata

import "github.com/gophercloud/gophercloud"

type commonResult struct {
	gophercloud.Result
}

// Extract interprets any SSL certificate result as a SSLCertificate, if possible.
func (r commonResult) Extract() (*SSLCertificate, error) {
	var res SSLCertificate
	err := r.ExtractInto(&res)
	return &res, err
}

// AddResult is the result of an add request. Call its Extract method
// to interpret a result as a SSLCertificate.
type AddResult struct {
	commonResult
}

// ListResult is th result of a list request. Call its Extract method
// to interpret the result as a slice of SSLCertificate.
type ListResult struct {
	commonResult
}

// Extract interprets a ListResult as a slice of SSLCertificate, if possible.
func (r ListResult) Extract() ([]SSLCertificate, error) {
	var sslCerts []SSLCertificate
	err := r.ExtractInto(&sslCerts)
	return sslCerts, err
}

// UpdateResult is the result of an update request. Call its Extract method
// to interpret a result as a SSLCertificate.
type UpdateResult struct {
	commonResult
}

// DeleteResult is the result of a delete request. Call its ExtractErr method
// to determine if a request succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}

// SSLCertificate represents a CDN SSL certificate.
type SSLCertificate struct {
	ID                  int    `json:"id"`
	Automated           bool   `json:"automated"`
	CertIssuer          string `json:"cert_issuer"`
	CertSubjectCN       string `json:"cert_subject_cn"`
	Deleted             bool   `json:"deleted"`
	HasRelatedResources bool   `json:"has_related_resources"`
	Name                string `json:"name"`
	ValidityNotAfter    string `json:"validity_not_after"`
	ValidityNotBefore   string `json:"validity_not_before"`
}
