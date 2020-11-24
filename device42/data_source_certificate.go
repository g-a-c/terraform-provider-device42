package device42

import (
	"encoding/json"
	"fmt"

	resty "github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	device42CertificateResponse struct {
		TotalCount         string                       `json:"total_count"`
		CertificateDetails []device42CertificateDetails `json:"certificate_details"`
	}

	device42CertificateDetails struct {
		ContentCommitmentUsage bool `json:"content_commitment_usage"`
		CrlSignUsage           bool `json:"crl_sign_usage"`
		// "custom_fields": [],
		// "data_encipherment_usage": false,
		// "days_to_expiry": 1077,
		// "decipher_only_usage": false,
		// "digital_signature_usage": true,
		// "encipher_only_usage": false,
		// "extended_key_usage": "SERVERAUTH(1.3.6.1.5.5.7.3.1)\nCLIENTAUTH(1.3.6.1.5.5.7.3.2)\n",
		// "id": "3",
		// "issued_by": "",
		// "issued_to": "registration.device42.com",
		// "key_agreement_usage": false,
		// "key_cert_sign_usage": false,
		// "key_encipherment_usage": true,
		// "parent_cert": "",
		// "serial_number": "77eb9b55e9228635f2157fd374b8da8",
		// "signature_algorithm": "sha256WithRSAEncryption",
		// "signature_hash": "708489795",
		// "subject": "/OU=Domain Control Validated/OU=PositiveSSL/CN=registration.device42.com",
		// "valid_from": "2014-08-10",
		// "valid_to": "2019-08-09",
		// "vendor": "",
		// "version": 2
	}
)

func dataSourceCertificate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCertificateRead,
		Schema: map[string]*schema.Schema{
			"certificate_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"label": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCertificateRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	passwordAPIResponse, err := client.R().
		SetQueryParam("id", fmt.Sprintf("%v", d.Get("id"))).
		SetQueryParam("plain_text", "yes").
		Get(client.HostURL)

	if err != nil {
		return nil
	}

	passwordAPIList := new(Device42PasswordList)
	json.Unmarshal([]byte(fmt.Sprintf("%v", passwordAPIResponse)), &passwordAPIList)
	d.Set("username", passwordAPIList.Passwords[0].Username)
	d.Set("password", passwordAPIList.Passwords[0].Password)
	d.Set("label", passwordAPIList.Passwords[0].Label)
	d.SetId(fmt.Sprintf("%b", passwordAPIList.Passwords[0].ID))
	return nil
}
