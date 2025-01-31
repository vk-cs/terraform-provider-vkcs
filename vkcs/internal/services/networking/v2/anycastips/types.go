package anycastips

// AnycastIPAssociationType association type
type AnycastIPAssociationType string

const (
	AnycastIPAssociationTypePort        AnycastIPAssociationType = "port"
	AnycastIPAssociationTypeDCInterface AnycastIPAssociationType = "dc_interface"
	AnycastIPAssociationTypeOctavia     AnycastIPAssociationType = "octavia"
)

// AnycastIPAssociationTypeValues returns list of all possible values for
// a AnycastIPAssociationType enumeration.
func AnycastIPAssociationTypeValues() []string {
	return []string{
		string(AnycastIPAssociationTypePort),
		string(AnycastIPAssociationTypeDCInterface),
		string(AnycastIPAssociationTypeOctavia),
	}
}

type AnycastIPAssociation struct {
	ID   string                   `json:"id" required:"true"`
	Type AnycastIPAssociationType `json:"type" required:"true"`
}

type AnycastIPHealthCheckType string

const (
	AnycastIPHealthCheckTypeTCP  AnycastIPHealthCheckType = "TCP"
	AnycastIPHealthCheckTypeICMP AnycastIPHealthCheckType = "ICMP"
)

// AnycastIPHealthCheckTypeValues returns list of all possible values for
// a AnycastIPHealthCheckType enumeration.
func AnycastIPHealthCheckTypeValues() []string {
	return []string{
		string(AnycastIPHealthCheckTypeTCP),
		string(AnycastIPHealthCheckTypeICMP),
	}
}

type AnycastIPHealthCheck struct {
	Type AnycastIPHealthCheckType `json:"type"`
	Port int                      `json:"port"`
}
