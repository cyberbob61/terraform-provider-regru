package regru

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRegruDNSRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceRegruDNSRecordCreate,
		Read:   resourceRegruDNSRecordRead,
		Delete: resourceRegruDNSRecordDelete,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"record": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"zone": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
		},
	}
}

func resourceRegruDNSRecordCreate(d *schema.ResourceData, m interface{}) error {

	recordType := d.Get("type").(string)
	recordName := d.Get("name").(string)
	value := d.Get("record").(string)
	zone := d.Get("zone").(string)

	c := m.(*Client)

	baseRequest := CreateRecordRequest{
		Username:          c.username,
		Password:          c.password,
		Domains:           []Domain{{DName: zone}},
		SubDomain:         recordName,
		OutputContentType: "plain",
	}

	var request interface{}
	switch strings.ToUpper(recordType) {
	case "A":
		request = CreateARecordRequest{
			CreateRecordRequest: baseRequest,
			IPAddr:              value,
		}

	case "AAAA":
		request = CreateAAAARecordRequest{
			CreateRecordRequest: baseRequest,
			IPAddr:              value,
		}

	case "CNAME":
		request = CreateCnameRecordRequest{
			CreateRecordRequest: baseRequest,
			CanonicalName:       value,
		}

	case "MX":
		fields := strings.Fields(value)
		if len(fields) != 2 {
			return fmt.Errorf("invalid MX record format, expected 'priority mailserver'")
		}
		request = CreateMxRecordRequest{
			CreateRecordRequest: baseRequest,
			MailServer:          fields[1],
			Priority:            fields[0],
		}

	case "TXT":
		request = CreateTxtRecordRequest{
			CreateRecordRequest: baseRequest,
			Text:                value,
		}

	default:
		return fmt.Errorf("invalid record type '%s'", recordType)
	}

	resp, err := c.doRequest(request, "zone", "add_alias") // Явно указываем путь
	if err != nil {
		return fmt.Errorf("failed to create DNS record: %w", err)
	}

	if resp.HasError() != nil {
		return fmt.Errorf("API error: %w", resp.HasError())
	}

	recordID := generateRecordID(recordName, zone)
	d.SetId(recordID)

	return nil
}

func generateRecordID(recordName, zone string) string {
	return fmt.Sprintf("%s.%s", recordName, zone)
}

func resourceRegruDNSRecordRead(_ *schema.ResourceData, _ interface{}) error {
	// Placeholder: function not implemented yet
	// func resourceRegruDNSRecordRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceRegruDNSRecordDelete(d *schema.ResourceData, m interface{}) error {
	recordName := d.Get("name").(string)
	zone := d.Get("zone").(string)
	recordType := d.Get("type").(string)
	recordValue := d.Get("record").(string)

	c := m.(*Client)

	request := DeleteRecordRequest{
		Username:          c.username,
		Password:          c.password,
		Domains:           []Domain{{DName: zone}},
		SubDomain:         recordName,
		Content:           recordValue,
		RecordType:        recordType,
		OutputContentType: "plain",
	}

	resp, err := c.doRequest(request, "zone", "remove_record")
	if err != nil {
		return fmt.Errorf("failed to delete DNS record: %w", err)
	}

	if resp.HasError() != nil {
		return fmt.Errorf("API error: %w", resp.HasError())
	}

	d.SetId("")

	return nil
}
