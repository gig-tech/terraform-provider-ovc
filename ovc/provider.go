package ovc

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/nuberabe/ovc-sdk-go/ovc"
)

// Provider method to define all user inputs
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"server_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OPENVCLOUD_SERVER_URL", nil),
				Description: "OpenvCloud URL",
			},
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ITSYOU_ONLINE_CLIENT_ID", nil),
				Description: "Client Id",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ITSYOU_ONLINE_CLIENT_SECRET", nil),
				Description: "Client Secret",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"ovc_machine":         dataSourceOvcMachine(),
			"ovc_cloudspace":      dataSourceOvcCloudSpace(),
			"ovc_sizes":           dataSourceOvcSizes(),
			"ovc_disk":            dataSourceOvcDisk(),
			"ovc_cloudspaces":     dataSourceOvcCloudSpaces(),
			"ovc_machines":        dataSourceOvcMachines(),
			"ovc_port_forwarding": dataSourceOvcPortForwarding(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"ovc_machine":         resourceOvcMachine(),
			"ovc_port_forwarding": resourcePortForwarding(),
			"ovc_disk":            resourceOvcDisk(),
			"ovc_cloudspace":      resourceOvcCloudSpace(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := ovc.Config{
		Hostname:     d.Get("server_url").(string) + "/restmachine",
		ClientID:     d.Get("client_id").(string),
		ClientSecret: d.Get("client_secret").(string),
	}

	client := ovc.NewClient(&config, d.Get("server_url").(string))

	return client, nil
}
