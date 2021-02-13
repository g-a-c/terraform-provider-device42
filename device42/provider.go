package device42

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"

	resty "github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"hostname": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "hostname of the Device42 instance",
				DefaultFunc: func() (interface{}, error) {
					if v := os.Getenv("D42_HOSTNAME"); v != "" {
						return v, nil
					}

					return "swaggerdemo.device42.com", nil
				},
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "username with access to the Device42 instance",
				DefaultFunc: func() (interface{}, error) {
					if v := os.Getenv("D42_USERNAME"); v != "" {
						return v, nil
					}

					return nil, errors.New("no username was provided via the D42_USERNAME environment variable")
				},
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "password for the username with access",
				DefaultFunc: func() (interface{}, error) {
					if v := os.Getenv("D42_PASSWORD"); v != "" {
						return v, nil
					}

					return nil, errors.New("no password was provided via the D42_PASSWORD environment variable")
				},
			},
		},
		ResourcesMap: map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"device42_password": dataSourcePassword(),
		},
		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	client := resty.New()
	if os.Getenv("D42_INSECURE") == "true" {
		insecureTransport := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		client.SetTransport(insecureTransport)
	}
	client.SetHostURL(fmt.Sprintf("https://%s", d.Get("hostname").(string)))
	client.SetBasicAuth(d.Get("username").(string), d.Get("password").(string))

	return client, nil
}
