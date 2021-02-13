package test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/g-a-c/terraform-provider-device42/device42"
)

func TestDevice42_Password_ValidResponse(t *testing.T) {
	t.Parallel()
	jsonResponse := `{"Passwords": [{"id": 1, "username": "validUsername", "password": "validPassword", "label": "validLabel"}]}`
	processedResponse := new(device42.Device42PasswordList)
	json.Unmarshal([]byte(jsonResponse), processedResponse)
	expectedResponse := device42.Device42PasswordList{
		Passwords: []device42.Device42Password{
			{
				ID:       1,
				Username: "validUsername",
				Password: "validPassword",
				Label:    "validLabel",
			},
		},
	}
	if !reflect.DeepEqual(*processedResponse, expectedResponse) {
		t.Fatal("value was not as expected")
	}
}

func TestDevice42_Password_ValidResponseWithAPILimitEnforced(t *testing.T) {
	t.Parallel()
	jsonResponse := `{"total_count": 1, "Passwords": [{"username": "validUsername", "category": null, "device_ids": [], "view_users": "", "view_groups": "", "last_pw_change": "2020-09-08T11:35:08.808Z", "notes": "", "storage": "Normal", "use_only_users": "", "label": "validLabel", "view_edit_groups": "", "first_added": "2020-09-08T11:35:08.810Z", "use_only_groups": "", "storage_id": 1, "view_edit_users": "admin", "password": "validPassword", "id": 1, "custom_fields": []}], "limit": 1000, "offset": 0}`
	processedResponse := new(device42.Device42PasswordList)
	json.Unmarshal([]byte(jsonResponse), processedResponse)
	expectedResponse := device42.Device42PasswordList{
		Passwords: []device42.Device42Password{
			{
				ID:       1,
				Username: "validUsername",
				Password: "validPassword",
				Label:    "validLabel",
			},
		},
	}
	if !reflect.DeepEqual(*processedResponse, expectedResponse) {
		t.Fatal("value was not as expected")
	}
}

// special cases to test for in future:
// {"msg": "Please enter the passphrase first. Go to Tools > Settings > Password Security", "code": 2}
// if the password encryption secret is not configured
//
// {"msg": "", "code": 0}
// if the password secret is configured, but there are no passwords, this happens
//
// {"Passwords": []}
// if the password secret is configured, there are passwords stored, but no matches for the ID, this happens
