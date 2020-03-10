package ovc

import (
	"fmt"
	"strconv"

	"github.com/gig-tech/ovc-sdk-go/v2/ovc"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceOvcExternalNetwork() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOvcExternalNetworkRead,

		Schema: map[string]*schema.Schema{
			"network_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"network_id"},
			},
			"account": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gateway": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnetmask": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOvcExternalNetworkRead(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	networkID := d.Get("network_id").(int)

	if networkID == 0 && name == "" {
		return fmt.Errorf("Either 'name' or 'network_id' should be given to define external network datasource")
	}

	client := m.(*ovc.Client)

	var err error
	var accountID int
	account := d.Get("account")
	if account != "" {
		accountID, err = client.Accounts.GetIDByName(account.(string))
		if err != nil {
			return err
		}
	}

	externalNetworks, err := client.ExternalNetworks.List(accountID)
	if err != nil {
		return err
	}

	filteredExternalNetworks := make([]ovc.ExternalNetworkInfo, 0)

	for _, externalNetwork := range *externalNetworks {
		// select images by name
		if (name == externalNetwork.Name || networkID == externalNetwork.ID) && (externalNetwork.AccountID == accountID || externalNetwork.AccountID == 0) {
			filteredExternalNetworks = append(filteredExternalNetworks, externalNetwork)
		}
	}

	if len(filteredExternalNetworks) < 1 {
		return fmt.Errorf("No external network with name '%s' is accessible for account '%s'", name, account)
	}

	if len(filteredExternalNetworks) > 1 {
		return fmt.Errorf("More than one external network was found with name '%s'", name)
	}

	d.SetId(strconv.Itoa(filteredExternalNetworks[0].ID))
	d.Set("network_id", filteredExternalNetworks[0].ID)
	d.Set("name", filteredExternalNetworks[0].Name)
	d.Set("network", filteredExternalNetworks[0].Network)
	d.Set("gateway", filteredExternalNetworks[0].Gateway)
	d.Set("subnetmask", filteredExternalNetworks[0].Subnetmask)
	return nil
}
