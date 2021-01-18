package device42

import (
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
				Computed: true,
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
	// stuff
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
