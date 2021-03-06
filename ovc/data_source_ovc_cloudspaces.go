package ovc

import (
	"github.com/gig-tech/ovc-sdk-go/v3/ovc"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceOvcCloudSpaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOvcCloudSpacesRead,

		Schema: map[string]*schema.Schema{
			"entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"account": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"account_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cloudspace_id": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"external_network_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"location": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceOvcCloudSpacesRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*ovc.Client)
	cloudSpaces, err := c.CloudSpaces.List()
	if err != nil {
		return err
	}
	entities := make([]map[string]interface{}, len(*cloudSpaces))
	for i, cp := range *cloudSpaces {
		entity := make(map[string]interface{})
		entity["status"] = cp.Status
		entity["cloudspace_id"] = cp.ID
		entity["name"] = cp.Name
		entity["account_id"] = cp.AccountID
		entity["external_network_ip"] = cp.Externalnetworkip
		entity["location"] = cp.Location
		entity["type"] = cp.Type
		entity["mode"] = cp.Mode
		entity["description"] = cp.Descr
		entity["account"] = cp.AccountName
		entities[i] = entity
	}
	if err := d.Set("entities", entities); err != nil {
		return err
	}
	d.SetId("1")
	return nil
}
