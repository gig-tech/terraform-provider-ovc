package ovc

import (
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/nuberabe/ovc-sdk-go/ovc"
)

func dataSourceOvcMachines() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOvcMachinesRead,

		Schema: map[string]*schema.Schema{
			"cloudspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"machine_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"size_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"image_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"hostname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creationtime": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceOvcMachinesRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*ovc.OvcClient)
	cloudspaceID := d.Get("cloudspace_id")
	cid, err := strconv.Atoi(cloudspaceID.(string))
	if err != nil {
		return err
	}
	machines, err := c.Machines.List(cid)
	if err != nil {
		return err
	}
	entities := make([]map[string]interface{}, len(*machines))
	for i, mc := range *machines {
		entity := make(map[string]interface{})
		entity["machine_id"] = strconv.Itoa(mc.ID)
		entity["name"] = mc.Name
		entity["size_id"] = mc.SizeID
		entity["image_id"] = mc.ImageID
		entity["status"] = mc.Status
		entity["update_time"] = mc.UpdateTime
		entity["creationtime"] = mc.CreationTime
		entities[i] = entity
	}
	if err := d.Set("entities", entities); err != nil {
		return err
	}
	d.SetId(cloudspaceID.(string))
	return nil
}
