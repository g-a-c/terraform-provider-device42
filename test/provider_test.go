package test

import (
	"context"
	"os"
	"testing"

	"github.com/g-a-c/terraform-provider-device42/device42"
	resty "github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jarcoal/httpmock"
)

var providerFactories = map[string]func() (*schema.Provider, error){
	"device42": func() (*schema.Provider, error) {
		os.Setenv("D42_HOSTNAME", "example.com")
		os.Setenv("D42_USERNAME", "username")
		os.Setenv("D42_PASSWORD", "password")
		os.Setenv("D42_INSECURE", "true")
		provider := device42.Provider()
		provider.Configure(context.TODO(), &terraform.ResourceConfig{
			Config: map[string]interface{}{
				"hostname": "example.com",
				"username": "username",
				"password": "password",
				"insecure": "true",
			},
		})
		// pass the httpmock client with its fake responses
		provider.SetMeta(getMockClient())
		httpmock.RegisterResponderWithQuery(
			"GET",
			"https://example.com/api/1.0/passwords/",
			map[string]string{
				"id":         "1",
				"plain_text": "yes",
			},
			httpmock.NewStringResponder(200, `{"Passwords": [{"id": 1, "username": "validUsername", "password": "validPassword", "label": "validLabel"}]}`),
		)
		return provider, nil
	},
}

func getMockClient() (client *resty.Client) {
	client = resty.New()
	client.SetHostURL("https://example.com")
	client.SetBasicAuth("mockUsername", "mockPassword")
	httpmock.ActivateNonDefault(client.GetClient())
	return client
}

func TestProvider(t *testing.T) {
	if err := device42.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
