package ovc

import (
	"strconv"

	"github.com/gig-tech/ovc-sdk-go/ovc"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceOvcExternalNetworks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOvcExternalNetworksRead,

		Schema: map[string]*schema.Schema{
			"account": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"account_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
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
						"dhcp": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceOvcExternalNetworksRead(d *schema.ResourceData, m interface{}) error {
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

	name := d.Get("name").(string)

	entities := make([]map[string]interface{}, len(*externalNetworks))

	for i, externalNetwork := range *externalNetworks {
		// select externalNetworks by name of network
		if (name == "" || name == externalNetwork.Name) && (externalNetwork.AccountID == accountID || externalNetwork.AccountID == 0) {
			entity := make(map[string]interface{})
			entity["id"] = strconv.Itoa(externalNetwork.ID)
			entity["name"] = externalNetwork.Name
			entity["network"] = externalNetwork.Network
			entity["gateway"] = externalNetwork.Gateway
			entity["subnetmask"] = externalNetwork.Subnetmask
			entity["dhcp"] = externalNetwork.DHCP
			entity["account_id"] = externalNetwork.AccountID
			entities[i] = entity
		}
	}

	if err = d.Set("entities", entities); err != nil {
		return err
	}

	d.SetId("1")
	return nil
}
