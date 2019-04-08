package ovc

import (
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/nuberabe/ovc-sdk-go/ovc"
)

func dataSourceOvcCloudSpace() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOvcCloudSpaceRead,

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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"external_network_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_network": {
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
			"resource_limits": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_memory_capacity": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"max_disk_capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_cpu_capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_network_peer_transfer": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_num_public_ip": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceOvcCloudSpaceRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*ovc.OvcClient)
	var cloudSpace *ovc.CloudSpace
	var err error
	if v, ok := d.GetOk("cloudspace_id"); ok {
		cloudSpace, err = client.CloudSpaces.Get(v.(string))
		if err != nil {
			return err
		}
	} else {
		cloudSpace, err = client.CloudSpaces.GetByNameAndAccount(d.Get("name").(string), d.Get("account").(string))
		if err != nil {
			return err
		}

	}
	d.Set("status", cloudSpace.Status)
	rl := make(map[string]interface{})
	rl["max_memory_capacity"] = strconv.FormatFloat(cloudSpace.ResourceLimits.CUM, 'f', -1, 64)
	rl["max_disk_capacity"] = strconv.Itoa(cloudSpace.ResourceLimits.CUD)
	rl["max_cpu_capacity"] = strconv.Itoa(cloudSpace.ResourceLimits.CUC)
	rl["max_network_peer_transfer"] = strconv.Itoa(cloudSpace.ResourceLimits.CUNP)
	rl["max_num_public_ip"] = strconv.Itoa(cloudSpace.ResourceLimits.CUI)
	d.Set("name", cloudSpace.Name)
	d.Set("account_id", cloudSpace.AccountID)
	err = d.Set("resource_limits", rl)
	if err != nil {
		return err
	}
	d.Set("description", cloudSpace.Description)
	d.Set("external_network_ip", cloudSpace.Externalnetworkip)
	d.Set("location", cloudSpace.Location)
	d.SetId(strconv.Itoa(cloudSpace.ID))
	return nil

}
