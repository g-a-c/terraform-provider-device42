package device42

import (
	"encoding/json"
	"fmt"

	resty "github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// Device42PasswordPostResponseWrapper describes the expected response from making a POST request to store a password
	// {"msg": ["Password added with username testcli (and label )", 2, "testcli", true, true], "code": 0}
	Device42PasswordPostResponseWrapper struct {
		Message []interface{} `json:"msg"`
		Code    int
	}
)

func resourcePassword() *schema.Resource {
	return &schema.Resource{
		Create: resourcePasswordCreate,
		Update: resourcePasswordUpdate,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
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

func resourcePasswordCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	passwordAPIResponse, err := client.R().
		SetQueryParam("id", fmt.Sprintf("%v", d.Get("id"))).
		SetQueryParam("plain_text", "yes").
		Get(client.HostURL)

	if err != nil {
		return nil
	}

	// grab the full response into a list-of-passwords object
	passwordAPIList := new(Device42PasswordList)
	json.Unmarshal([]byte(fmt.Sprintf("%v", passwordAPIResponse)), &passwordAPIList)
	// read entry 0 (which will be the only entry since we retrieve by ID, which is unique) and set the correct fields on the returned Terraform object
	d.Set("username", passwordAPIList.Passwords[0].Username)
	d.Set("password", passwordAPIList.Passwords[0].Password)
	d.Set("label", passwordAPIList.Passwords[0].Label)
	d.SetId(fmt.Sprintf("%b", passwordAPIList.Passwords[0].ID))
	return nil
}

func resourcePasswordUpdate(d *schema.ResourceData, m interface{}) error {
	// stuff
	return nil
}

func resourcePasswordDelete(d *schema.ResourceData, m interface{}) error {
	// stuff
	return nil
}
