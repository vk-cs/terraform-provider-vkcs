package s3accounts

type S3Account struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TenantID    string `json:"tenant_id"`
	AccountID   string `json:"account_id"`
	AccountName string `json:"account_name"`
	AccessKey   string `json:"access_key"`
	CreatedAt   string `json:"created_at"`
}

type CreateS3AccountResponse struct {
	S3Account
	SecretKey string `json:"secret_key"`
}
