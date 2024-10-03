package resource_ssl_certificate

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/cdn/v1/ssldata"
)

func (m *SslCertificateModel) UpdateFromSslCertificate(ctx context.Context, sslCertificate *ssldata.SSLCertificate) diag.Diagnostics {
	var diags diag.Diagnostics

	if sslCertificate == nil {
		return diags
	}

	m.Id = types.Int64Value(int64(sslCertificate.ID))
	m.Issuer = types.StringValue(sslCertificate.CertIssuer)
	m.Name = types.StringValue(sslCertificate.Name)
	m.SubjectCn = types.StringValue(sslCertificate.CertSubjectCN)
	m.ValidityNotAfter = types.StringValue(sslCertificate.ValidityNotAfter)
	m.ValidityNotBefore = types.StringValue(sslCertificate.ValidityNotBefore)

	return diags
}
