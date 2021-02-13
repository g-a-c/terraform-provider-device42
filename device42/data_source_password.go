package device42

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	resty "github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// Device42NoPassphraseResponse -
	Device42GenericResponse struct {
		Message string `json:"msg,omitempty"`
		Code    string `json:"code,omitempty"`
	}

	// Device42PasswordList is a list containing one or more Device42Password objects
	Device42PasswordList struct {
		Passwords []Device42Password `json:"Passwords"`
	}
	// Device42Password is a struct containing all potential documented fields which may be returned from the Device42 password API
	Device42Password struct {
		Username string `json:"username,omitempty"`
		// Category       string   `json:"category,omitempty"`
		// DeviceIds      []string `json:"device_ids,omitempty"`
		// ViewUsers      string   `json:"view_users,omitempty"`
		// ViewGroups     string   `json:"view_groups,omitempty"`
		// LastPwChange   string   `json:"last_pw_change,omitempty"`
		// Notes          string   `json:"notes,omitempty"`
		// Storage        string   `json:"storage,omitempty"`
		// UseOnlyUsers   string   `json:"use_only_users,omitempty"`
		Label string `json:"label,omitempty"`
		// ViewEditGroups string   `json:"view_edit_groups,omitempty"`
		// FirstAdded     string   `json:"first_added,omitempty"`
		// UseOnlyGroups  string   `json:"use_only_groups,omitempty"`
		// StorageID      int      `json:"storage_id,omitempty"`
		// ViewEditUsers  string   `json:"view_edit_users,omitempty"`
		Password string `json:"password"`
		ID       int    `json:"id"`
		// CustomFields   []string `json:"custom_fields,omitempty"`
	}
)

func dataSourcePassword() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePasswordRead,
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

func dataSourcePasswordRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	client := m.(*resty.Client)
	passwordAPIResponse, err := client.R().
		SetQueryParam("id", strconv.Itoa(d.Get("id").(int))).
		SetQueryParam("plain_text", "yes").
		SetHeader("Accept", "application/json").
		Get("/api/1.0/passwords/")

	// return up any errors which are probably on the connection level at this point (i.e. timeout, handshake failure)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Something bad happened",
				Detail:   "Something went wrong with the HTTP request",
			},
		}
	}

	// permission errors come out as a non-JSON-wrapped quoted string because why not?
	if string(passwordAPIResponse.Body()) == "\"You don't have permissions to access this resource\"" {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Access Denied",
				Detail:   "This secret appears to exist, but you do not have view permission to it",
			},
		}
	}

	// if we got here, try to Unmarshal into a list-of-passwords which should contain one element
	passwordAPIList := new(Device42PasswordList)
	err = json.Unmarshal(passwordAPIResponse.Body(), passwordAPIList)
	if err == nil {
		if len(passwordAPIList.Passwords) == 1 {
			// read entry 0 (which will be the only entry since ID is unique) and set the correct fields on the returned Terraform object
			d.Set("username", passwordAPIList.Passwords[0].Username)
			d.Set("password", passwordAPIList.Passwords[0].Password)
			d.Set("label", passwordAPIList.Passwords[0].Label)
			d.SetId(fmt.Sprintf("%b", passwordAPIList.Passwords[0].ID))
			return nil
		}

		// If there are no entries in the list then Device42 is configured with an encryption passphrase, but no secret exists
		// so return an error saying such
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "No secret was found",
				Detail:   "No secret exists in your Device42 instance with that ID",
			},
		}
	}

	// if we couldn't unmarshal to the expected list-of-password struct, try unmarshalling to a generic struct with msg/code attributes
	genericResponse := new(Device42GenericResponse)
	err = json.Unmarshal(passwordAPIResponse.Body(), genericResponse)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Device42 returned an error, there may be more information below",
				Detail:   fmt.Sprintf("%v", genericResponse.Message),
			},
		}
	}

	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "There was an unhandled error",
		},
	}
}
