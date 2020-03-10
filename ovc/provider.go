package ovc

import (
	"os"
	"strings"

	"github.com/gig-tech/ovc-sdk-go/v2/ovc"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sirupsen/logrus"
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
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"client_jwt"},
				DefaultFunc:   schema.EnvDefaultFunc("ITSYOU_ONLINE_CLIENT_ID", nil),
				Description:   "Client Id",
			},
			"client_secret": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"client_jwt"},
				DefaultFunc:   schema.EnvDefaultFunc("ITSYOU_ONLINE_CLIENT_SECRET", nil),
				Description:   "Client Secret",
			},
			"client_jwt": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"client_id", "client_secret"},
				DefaultFunc:   schema.EnvDefaultFunc("ITSYOU_ONLINE_CLIENT_JWT", nil),
				Description:   "Client JWT",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"ovc_machine":           dataSourceOvcMachine(),
			"ovc_machines":          dataSourceOvcMachines(),
			"ovc_cloudspace":        dataSourceOvcCloudSpace(),
			"ovc_cloudspaces":       dataSourceOvcCloudSpaces(),
			"ovc_sizes":             dataSourceOvcSizes(),
			"ovc_disk":              dataSourceOvcDisk(),
			"ovc_port_forwarding":   dataSourceOvcPortForwarding(),
			"ovc_image":             dataSourceOvcImage(),
			"ovc_images":            dataSourceOvcImages(),
			"ovc_external_network":  dataSourceOvcExternalNetwork(),
			"ovc_external_networks": dataSourceOvcExternalNetworks(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"ovc_machine":         resourceOvcMachine(),
			"ovc_port_forwarding": resourcePortForwarding(),
			"ovc_disk":            resourceOvcDisk(),
			"ovc_cloudspace":      resourceOvcCloudSpace(),
			"ovc_ipsec":           resourceIpsec(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var ovcLogger ovc.Logger = nil
	g8LogFile, found := os.LookupEnv("G8_API_ACCESS_LOG_FILE")
	if found {
		logger := logrus.New()
		switch g8LogLevel, _ := os.LookupEnv("G8_API_ACCESS_LOG_LEVEL"); strings.ToLower(g8LogLevel) {
		case "debug":
			logger.SetLevel(logrus.DebugLevel)
		default:
			logger.SetLevel(logrus.InfoLevel)
		}
		f, err := os.OpenFile(g8LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			logger.Fatalf("error opening file: %v", err)
			logger.SetOutput(os.Stdout)
		} else {
			logger.SetOutput(f)
		}
		ovcLogger = ovc.LogrusAdapter{FieldLogger: logger}
	}
	config := ovc.Config{
		URL:          d.Get("server_url").(string),
		ClientID:     d.Get("client_id").(string),
		ClientSecret: d.Get("client_secret").(string),
		JWT:          d.Get("client_jwt").(string),
		Logger:       ovcLogger,
	}
	return ovc.NewClient(&config)
}
