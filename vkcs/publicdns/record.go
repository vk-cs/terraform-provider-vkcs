package publicdns

import (
	"fmt"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/publicdns/v2/records"
)

type recordResult interface {
	ExtractA() (*records.RecordA, error)
	ExtractAAAA() (*records.RecordAAAA, error)
	ExtractCNAME() (*records.RecordCNAME, error)
	ExtractMX() (*records.RecordMX, error)
	ExtractNS() (*records.RecordNS, error)
	ExtractSRV() (*records.RecordSRV, error)
	ExtractTXT() (*records.RecordTXT, error)
}

func PublicDNSRecordExtract(res recordResult, recordType string) (interface{}, error) {
	var (
		r   interface{}
		err error
	)

	switch recordType {
	case recordTypeA:
		r, err = res.ExtractA()
	case recordTypeAAAA:
		r, err = res.ExtractAAAA()
	case recordTypeCNAME:
		r, err = res.ExtractCNAME()
	case recordTypeMX:
		r, err = res.ExtractMX()
	case recordTypeNS:
		r, err = res.ExtractNS()
	case recordTypeSRV:
		r, err = res.ExtractSRV()
	case recordTypeTXT:
		r, err = res.ExtractTXT()
	}

	return r, err
}

func publicDNSRecordStateRefreshFunc(client *gophercloud.ServiceClient, zoneID string, id string, recordType string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res := records.Get(client, zoneID, id, recordType)
		record, err := PublicDNSRecordExtract(res, recordType)

		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return record, recordStatusDeleted, nil
			}
			return nil, "", err
		}

		return record, recordStatusActive, nil
	}
}

func publicDNSRecordParseZoneID(dns string) (string, error) {
	dnsParts := strings.Split(dns, "/") // /v2/dns/<zone-uuid>
	if len(dnsParts) != 4 {
		return "", fmt.Errorf("unable to determine vkcs_publicdns_record zone ID from raw DNS: %s", dns)
	}
	zoneID := dnsParts[3]
	return zoneID, nil
}

func publicDNSRecordParseID(id string) (string, string, string, error) {
	parts := strings.Split(id, "/") // <zone-uuid>/<record-type>/<record-uuid>
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("unable to determine vkcs_publicdns_record ID from %s", id)
	}

	zoneID := parts[0]
	recordType := parts[1]
	recordID := parts[2]

	return zoneID, recordType, recordID, nil
}
