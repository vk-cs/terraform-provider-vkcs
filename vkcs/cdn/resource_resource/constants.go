package resource_resource

type ResourceStatus string

var (
	ResourceStatusActive    ResourceStatus = "active"
	ResourceStatusProcessed ResourceStatus = "processed"
	ResourceStatusSuspended ResourceStatus = "suspended"
)

type SslCertificateProviderType string

var (
	SslCertificateProviderTypeNotUsed     SslCertificateProviderType = "not_used"
	SslCertificateProviderTypeLetsEncrypt SslCertificateProviderType = "lets_encrypt"
	SslCertificateProviderTypeOwn         SslCertificateProviderType = "own"
)

type SslCertificateStatus string

var (
	SslCertificateStatusBeingIssued SslCertificateStatus = "being_issued"
	SslCertificateStatusReady       SslCertificateStatus = "ready"
)
