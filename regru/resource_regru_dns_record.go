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
	record_type := d.Get("type").(string)
	record_name := d.Get("name").(string)
	value := d.Get("record").(string)
	zone := d.Get("zone").(string)

	c := m.(*Client)
	baseRequest := CreateRecordRequest{
		Username:          c.username,
		Password:          c.password,
		Domains:           []Domain{{DName: zone}},
		SubDomain:         record_name,
		OutputContentType: "plain",
	}

	var request interface{}

	switch strings.ToUpper(record_type) {
	case "A":
		request = CreateARecordRequest{
			CreateRecordRequest: baseRequest,
			IPAddr:              value,
		}
		d.Set("api_endpoint", "https://api.reg.ru/api/regru2/zone/add_alias")

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
		return fmt.Errorf("invalid record type '%s'", record_type)
	}

	resp, err := c.doRequest(request, "zone")
	if err != nil {
		return err
	}
	if resp.HasError() != nil {
		return resp.HasError()
	}
	d.SetId(strings.Join([]string{record_name, zone}, "."))
	return nil
}

func resourceRegruDNSRecordRead(_ *schema.ResourceData, _ interface{}) error {
	// Placeholder: function not implemented yet
	// func resourceRegruDNSRecordRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceRegruDNSRecordDelete(d *schema.ResourceData, m interface{}) error {
	record_type := d.Get("type").(string)
	record_name := d.Get("name").(string)
	value := d.Get("record").(string)
	zone := d.Get("zone").(string)

	c := m.(*Client)

	request := DeleteRecordRequest{
		Username:          c.username,
		Password:          c.password,
		Domains:           []Domain{{DName: zone}},
		SubDomain:         record_name,
		Content:           value,
		RecordType:        strings.ToUpper(record_type),
		OutputContentType: "plain",
	}

	resp, err := c.doRequest(request, "zone", "remove_record")
	if err != nil {
		return err
	}

	return resp.HasError()
}
