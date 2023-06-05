package publicdns

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"text/template"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mitchellh/mapstructure"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/publicdns/v2/records"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

const (
	recordDelay         = 10 * time.Second
	recordMinTimeout    = 3 * time.Second
	recordCreateTimeout = 10 * time.Minute
	recordDeleteTimeout = 10 * time.Minute
)

const (
	recordStatusActive  = "active"
	recordStatusDeleted = "deleted"
)

const (
	recordTypeA     = "A"
	recordTypeAAAA  = "AAAA"
	recordTypeCNAME = "CNAME"
	recordTypeMX    = "MX"
	recordTypeNS    = "NS"
	recordTypeSRV   = "SRV"
	recordTypeTXT   = "TXT"
)

var (
	recordTypeSharedArgs = [5]string{"region", "zone_id", "type", "name", "ttl"}
	recordTypeAArgs      = [1]string{"ip"}
	recordTypeAAAAArgs   = [1]string{"ip"}
	recordTypeCNAMEArgs  = [2]string{"name", "content"}
	recordTypeMXArgs     = [2]string{"priority", "content"}
	recordTypeNSArgs     = [1]string{"content"}
	recordTypeSRVArgs    = [6]string{"service", "proto", "priority", "weight", "host", "port"}
	recordTypeTXTArgs    = [1]string{"content"}
)

func ResourcePublicDNSRecord() *schema.Resource {
	return &schema.Resource{
		CustomizeDiff: resourcePublicDNSRecordCustomizeDiff,

		CreateContext: resourcePublicDNSRecordCreate,
		ReadContext:   resourcePublicDNSRecordRead,
		UpdateContext: resourcePublicDNSRecordUpdate,
		DeleteContext: resourcePublicDNSRecordDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zoneID, recordType, recordID, err := publicDNSRecordParseID(d.Id())
				if err != nil {
					return nil, err
				}

				recordType = strings.ToUpper(recordType)
				d.Set("zone_id", zoneID)
				d.Set("type", recordType)
				d.SetId(recordID)

				return []*schema.ResourceData{d}, nil
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(recordCreateTimeout),
			Delete: schema.DefaultTimeout(recordDeleteTimeout),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region in which to obtain the V2 Public DNS client. If omitted, the `region` argument of the provider is used. Changing this creates a new record.",
			},

			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the zone to attach the record to.",
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					recordTypeA,
					recordTypeAAAA,
					recordTypeCNAME,
					recordTypeMX,
					recordTypeNS,
					recordTypeSRV,
					recordTypeTXT,
				}, false),
				Description: "The type of the record. Must be one of following: \"A\", \"AAAA\", \"CNAME\", \"MX\", \"NS\", \"SRV\", \"TXT\".",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the record.",
			},

			"ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The IP address of the record. It should be IPv4 for record of type \"A\" and IPv6 for record of type \"AAAA\".",
			},

			"content": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The content of the record.",
			},

			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The priority of the record's server.",
			},

			"weight": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The relative weight of the record's server.",
			},

			"service": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the desired service.",
			},

			"proto": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the desired protocol.",
			},

			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The domain name of the target host.",
			},

			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The port on the target host of the service.",
			},

			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The time to live of the record.",
			},

			// Computed values
			"full_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The full name of the SRV record.",
			},
		},
		Description: resourcePublicDNSRecordDescription(),
	}
}

func resourcePublicDNSRecordDescription() string {
	var templ = `Manages a public DNS zone record resource within VKCS.<br>
**Note:** Although some arguments are marked as optional, it is actually required to set values for them depending on record \"type\". Use this map to get information about which arguments you have to set:

| Record type | Required arguments |
| ----------- | ------------------ |
{{range $k, $v := .}}| {{$k}} | {{join $v ", "}} |
{{end}}

`
	t := template.Must(template.New("").Funcs(template.FuncMap{"join": strings.Join}).Parse(templ))
	args := getPublicDNSRecordArgsMap()
	buf := &bytes.Buffer{}
	_ = t.Execute(buf, args)
	return buf.String()
}

func resourcePublicDNSRecordCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	client, err := config.PublicDNSV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS public DNS client: %s", err)
	}

	zoneID := d.Get("zone_id").(string)
	recordType := d.Get("type").(string)

	var createOpts records.CreateOptsBuilder
	switch recordType {
	case recordTypeA:
		createOpts = records.RecordACreateOpts{}
	case recordTypeAAAA:
		createOpts = records.RecordAAAACreateOpts{}
	case recordTypeCNAME:
		createOpts = records.RecordCNAMECreateOpts{}
	case recordTypeMX:
		createOpts = records.RecordMXCreateOpts{}
	case recordTypeNS:
		createOpts = records.RecordNSCreateOpts{}
	case recordTypeSRV:
		createOpts = records.RecordSRVCreateOpts{}
	case recordTypeTXT:
		createOpts = records.RecordTXTCreateOpts{}
	}

	m := publicDNSRecordResourceDataMap(d)
	if err = mapstructure.Decode(m, &createOpts); err != nil {
		return diag.Errorf("Error retrieving VKCS public DNS record create options: %s", err)
	}

	log.Printf("[DEBUG] vkcs_publicdns_record create options: zone_id: %s, type: %s, opts: %#v",
		zoneID, recordType, createOpts)

	res := records.Create(client, zoneID, createOpts, recordType)
	r, err := PublicDNSRecordExtract(res, recordType)
	if err != nil {
		return diag.Errorf("Error creating vkcs_publicdns_record: %s", err)
	}

	record, err := util.StructToMap(r)
	if err != nil {
		return diag.Errorf("Error creating vkcs_publicdns_record: %s", err)
	}

	d.SetId(record["uuid"].(string))
	log.Printf("[DEBUG] Created vkcs_publicdns_record %#v", record)

	return resourcePublicDNSRecordRead(ctx, d, meta)
}

func resourcePublicDNSRecordRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	client, err := config.PublicDNSV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS public DNS client: %s", err)
	}

	zoneID := d.Get("zone_id").(string)
	recordType := d.Get("type").(string)

	res := records.Get(client, zoneID, d.Id(), recordType)
	r, err := PublicDNSRecordExtract(res, recordType)

	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error retrieving vkcs_publicdns_record"))
	}

	record, err := util.StructToMap(r)
	if err != nil {
		return diag.Errorf("Error retrieving vkcs_publicdns_record: %s", err)
	}

	log.Printf("[DEBUG] Retrieved vkcs_publicdns_record %s: %#v", d.Id(), record)

	zoneID, err = publicDNSRecordParseZoneID(record["dns"].(string))
	if err != nil {
		return diag.Errorf("Error retrieving vkcs_publicdns_record: %s", err)
	}

	d.Set("zone_id", zoneID)

	if recordType == recordTypeSRV {
		fullN := record["name"].(string)
		d.Set("full_name", fullN)

		n, err := extractPublicDNSRecordSRVName(fullN)
		if err != nil {
			return diag.Errorf("Error retrieving vkcs_publicdns_record: %s", err)
		}
		d.Set("name", n)
	} else {
		d.Set("name", record["name"].(string))
	}

	switch recordType {
	case recordTypeA:
		d.Set("ip", record["ipv4"].(string))
	case recordTypeAAAA:
		d.Set("ip", record["ipv6"].(string))
	case recordTypeCNAME, recordTypeNS, recordTypeTXT:
		d.Set("content", record["content"].(string))
	case recordTypeMX:
		d.Set("priority", record["priority"].(int))
		d.Set("content", record["content"].(string))
	case recordTypeSRV:
		d.Set("priority", record["priority"].(int))
		d.Set("weight", record["weight"].(int))
		d.Set("host", record["host"].(string))
		d.Set("port", record["port"].(int))
	}

	d.Set("ttl", record["ttl"].(int))

	return nil
}

func resourcePublicDNSRecordUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	client, err := config.PublicDNSV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS public DNS client: %s", err)
	}

	zoneID := d.Get("zone_id").(string)
	recordType := d.Get("type").(string)

	var updateOpts records.UpdateOptsBuilder
	switch recordType {
	case recordTypeA:
		updateOpts = records.RecordACreateOpts{}
	case recordTypeAAAA:
		updateOpts = records.RecordAAAACreateOpts{}
	case recordTypeCNAME:
		updateOpts = records.RecordCNAMECreateOpts{}
	case recordTypeMX:
		updateOpts = records.RecordMXCreateOpts{}
	case recordTypeNS:
		updateOpts = records.RecordNSCreateOpts{}
	case recordTypeSRV:
		updateOpts = records.RecordSRVCreateOpts{}
	case recordTypeTXT:
		updateOpts = records.RecordTXTCreateOpts{}
	}

	m := publicDNSRecordResourceDataMap(d)
	if err = mapstructure.Decode(m, &updateOpts); err != nil {
		return diag.Errorf("Error retrieving VKCS public DNS record update options: %s", err)
	}

	log.Printf("[DEBUG] vkcs_publicdns_record update options: zone_id: %s, type: %s, opts: %#v",
		zoneID, recordType, updateOpts)

	res := records.Update(client, zoneID, d.Id(), updateOpts, recordType)
	_, err = PublicDNSRecordExtract(res, recordType)

	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error updating vkcs_publicdns_record"))
	}

	return resourcePublicDNSRecordRead(ctx, d, meta)
}

func resourcePublicDNSRecordDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	client, err := config.PublicDNSV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS public DNS client: %s", err)
	}

	zoneID := d.Get("zone_id").(string)
	recordType := d.Get("type").(string)

	err = records.Delete(client, zoneID, d.Id(), recordType).ExtractErr()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err,
			fmt.Sprintf("Error deleting vkcs_publicdns_record: zone_id: %s, type: %s, id:", zoneID, recordType)))
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{recordStatusActive},
		Target:     []string{recordStatusDeleted},
		Refresh:    publicDNSRecordStateRefreshFunc(client, zoneID, d.Id(), recordType),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      recordDelay,
		MinTimeout: recordMinTimeout,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for vkcs_publicdns_record %s to become deleted: %s", d.Id(), err)
	}

	return nil
}

func resourcePublicDNSRecordCustomizeDiff(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
	recordType := diff.Get("type").(string)

	if recordType == recordTypeAAAA {
		if diff.HasChange("ip") {
			new := net.ParseIP(diff.Get("ip").(string))
			if new != nil {
				diff.SetNew("ip", new.String())
			}
		}
	}

	args := getPublicDNSRecordArgs(recordType)

	var changedArgs []string
	for k := range diff.GetRawPlan().AsValueMap() {
		if diff.HasChange(k) {
			changedArgs = append(changedArgs, k)
		}
	}

OuterLoop:
	for _, k := range changedArgs {
		for _, ak := range args {
			if k == ak {
				continue OuterLoop
			}
		}
		return fmt.Errorf("\"%s\" is not supported for records of type \"%s\", supported args are: %v", k, recordType, args)
	}

	if recordType == recordTypeSRV {
		diff.SetNew("full_name", extractPublicDNSRecordSRVFullName(diff))
	}

	return nil
}

func extractPublicDNSRecordSRVName(fullName string) (string, error) {
	parts := strings.SplitN(fullName, ".", 3)
	if len(parts) == 3 {
		return parts[2], nil
	}
	if len(parts) == 2 {
		return "", nil
	}
	return "", fmt.Errorf("error extracting SRV record name from %s", fullName)
}

func extractPublicDNSRecordSRVFullName(diff *schema.ResourceDiff) string {
	s := diff.Get("service").(string)
	p := diff.Get("proto").(string)
	n := fmt.Sprintf("%s.%s", s, p)

	if v, ok := diff.GetOk("name"); ok {
		n += fmt.Sprintf(".%s", v.(string))
	}

	return n
}

func getPublicDNSRecordArgs(recordType string) []string {
	argsMap := getPublicDNSRecordArgsMap()
	var args []string
	args = append(args, recordTypeSharedArgs[:]...)
	args = append(args, argsMap[recordType]...)
	return args
}

func getPublicDNSRecordArgsMap() map[string][]string {
	return map[string][]string{
		recordTypeA:     recordTypeAArgs[:],
		recordTypeAAAA:  recordTypeAAAAArgs[:],
		recordTypeCNAME: recordTypeCNAMEArgs[:],
		recordTypeMX:    recordTypeMXArgs[:],
		recordTypeNS:    recordTypeNSArgs[:],
		recordTypeSRV:   recordTypeSRVArgs[:],
		recordTypeTXT:   recordTypeTXTArgs[:],
	}
}

func publicDNSRecordResourceDataMap(d *schema.ResourceData) map[string]interface{} {
	recordType := d.Get("type").(string)

	m := map[string]interface{}{
		"content":  d.Get("content"),
		"priority": d.Get("priority"),
		"weight":   d.Get("weight"),
		"host":     d.Get("host"),
		"port":     d.Get("port"),
		"ttl":      d.Get("ttl"),
	}
	if recordType == recordTypeSRV {
		m["name"] = d.Get("full_name")
	} else {
		m["name"] = d.Get("name")
	}
	if recordType == recordTypeA {
		m["ipv4"] = d.Get("ip")
	} else if recordType == recordTypeAAAA {
		m["ipv6"] = d.Get("ip")
	}

	return m
}
