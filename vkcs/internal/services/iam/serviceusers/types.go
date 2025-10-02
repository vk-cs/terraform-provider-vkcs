package serviceusers

type OrderByType string

const (
	OrderByTypeName OrderByType = "NAME"
)

type OrderByDirection string

const (
	OrderByDirectionASC  OrderByDirection = "ASC"
	OrderByDirectionDESC OrderByDirection = "DESC"
)

type ServiceUser struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	TenantID    string   `json:"tenant_id"`
	Description string   `json:"description"`
	CreatorName string   `json:"creator_name"`
	RoleNames   []string `json:"role_names"`
	CreatedAt   string   `json:"created_at"`
}

type CreateServiceUserResponse struct {
	ServiceUser
	Login    string `json:"login"`
	Password string `json:"password"`
}
