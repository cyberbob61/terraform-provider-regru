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
	var action string

	switch strings.ToUpper(record_type) {
	case "A":
		request = CreateARecordRequest{
			CreateRecordRequest: baseRequest,
			IPAddr:              value,
		}
		action = "add_alias"

	case "AAAA":
		request = CreateAAAARecordRequest{
			CreateRecordRequest: baseRequest,
			IPAddr:              value,
		}
		action = "add_alias_ipv6"

	case "CNAME":
		request = CreateCnameRecordRequest{
			CreateRecordRequest: baseRequest,
			CanonicalName:       value,
		}
		action = "add_cname"

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
		action = "add_mx"
	case "TXT":
		request = CreateTxtRecordRequest{
			CreateRecordRequest: baseRequest,
			Text:                value,
		}
		action = "add_txt"
	default:
		return fmt.Errorf("invalid record type '%s'", record_type)
	}

	resp, err := c.doRequest(request, "zone", action)
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
