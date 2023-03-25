package device42

import (
	"context"
	"encoding/json"
	"strconv"

	resty "github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
)

var (
	// Errors
	HTTPConnectionError = diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "An error occurred",
		Detail:   "Something went wrong with the HTTP request, for example a connection timeout, or TLS handshake error",
	}
	InvalidJSONError = diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Invalid JSON response",
		Detail:   "A response was received, but could not be parsed as JSON; is your hostname correctly configured?",
	}
	Device42PasswordEncryptionPassphraseError = diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "No password encryption secret is configured",
		Detail:   "Go to Tools > Settings > Password Security in your Device42 instance to configure one",
	}
	Device42NoSecretsError = diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "No passwords are stored",
		Detail:   "Your Device42 instance seems to be correctly configured, but does not seem to contain any Secrets",
	}
	Device42NoMatchingSecretsError = diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "No results were returned",
		Detail:   "Potential causes: the secret does not exist, or you do not have View permission to it",
	}
	Device42APIPermissionError = diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "REST API Access Denied",
		Detail:   "Your user does not have permissions to access the Password APIs. The user may need to be added to one of the groups listed in Tools > Admins & Permissions > Admin Groups which has the 'Password | Can view password' permission",
	}
	Device42MultiplePasswordsWarning = diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "multiple passwords were returned",
		Detail:   "multiple entries were returned; the first is being used. this should not be able to happen...",
	}
)

type (
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

func dataSourcePasswordRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client := meta.(*resty.Client)
	var genericJSON interface{}
	passwordAPIResponse, err := client.R().
		SetQueryParam("id", strconv.Itoa(d.Get("id").(int))).
		SetQueryParam("plain_text", "yes").
		SetHeader("Accept", "application/json").
		Get("/api/1.0/passwords/")

	// return up any errors which are probably on the connection level at this point (i.e. timeout, handshake failure)
	if err != nil {
		diags = append(diags, HTTPConnectionError)
	}

	// API permission errors come out as a non-JSON-wrapped quoted string
	if passwordAPIResponse.String() == `"You don't have permissions to access this resource"` {
		diags = append(diags, Device42APIPermissionError)
	}

	err = json.Unmarshal(passwordAPIResponse.Body(), &genericJSON)
	if err != nil {
		diags = append(diags, InvalidJSONError)
	}

	if diags.HasError() {
		return
	}

	genericJSONAsMap := genericJSON.(map[string]interface{})

	msg, msgPresent := genericJSONAsMap["msg"].(string)
	code, codePresent := genericJSONAsMap["code"].(float64)
	totalCount, totalCountPresent := genericJSONAsMap["total_count"].(float64)

	if (codePresent && code == 2) && (msgPresent && msg == "Please enter passphrase first. Go to Tools > Settings > Password Security") {
		// code=2 and this message means the password encryption passphrase is unset
		diags = append(diags, Device42PasswordEncryptionPassphraseError)
	}

	if (codePresent && code == 0) && (msgPresent && msg == "") {
		// password secret is configured, but no passwords are stored
		diags = append(diags, Device42NoSecretsError)
	}

	if totalCountPresent && totalCount == 0 {
		// API limits are enabled, no match
		diags = append(diags, Device42NoMatchingSecretsError)
	}

	passwordList, passwordListPresent := genericJSONAsMap["Passwords"].([]interface{})
	if passwordListPresent && len(passwordList) == 0 {
		// API limits are disabled, no match
		diags = append(diags, Device42NoMatchingSecretsError)
	}

	if passwordListPresent && len(passwordList) > 1 {
		// somehow there was more than one match
		diags = append(diags, Device42MultiplePasswordsWarning)
	}

	if diags.HasError() {
		// a fatal error has already occurred so the real process can't succeed
		return
	}

	if passwordListPresent && len(passwordList) >= 1 {
		// read entry 0 (which should be the only entry since ID is unique) and set the correct fields on the returned Terraform object
		passwordObject := new(Device42Password)
		mapstructure.Decode(passwordList[0], &passwordObject)
		d.Set("username", passwordObject.Username)
		d.Set("password", passwordObject.Password)
		d.Set("label", passwordObject.Label)
		d.SetId(strconv.Itoa(passwordObject.ID))
	}

	// return including any warnings
	return
}
