package test

import (
	"os"
	"testing"

	"github.com/g-a-c/terraform-provider-device42/device42"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccDevice42PasswordValidResponse(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		IsUnitTest: false,
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"device42": func() (*schema.Provider, error) {
				os.Setenv("D42_HOSTNAME", "192.168.128.5:7777")
				os.Setenv("D42_USERNAME", "tftest")
				os.Setenv("D42_PASSWORD", "tftest")
				os.Setenv("D42_INSECURE", "true")
				provider := device42.Provider()
				return provider, nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `data "device42_password" "test" { id = 2 }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.device42_password.test", "id", "2"),
					resource.TestCheckResourceAttr("data.device42_password.test", "username", "root"),
					resource.TestCheckResourceAttr("data.device42_password.test", "label", "test_label"),
					resource.TestCheckResourceAttr("data.device42_password.test", "password", "password"),
				),
			},
			{
				Config: `data "device42_password" "test" { id = 3 }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.device42_password.test", "id", "3"),
					resource.TestCheckResourceAttr("data.device42_password.test", "username", "root"),
					resource.TestCheckResourceAttr("data.device42_password.test", "label", ""),
					resource.TestCheckResourceAttr("data.device42_password.test", "password", "password"),
				),
			},
		},
	})
}
