package ovc

import (
	"strconv"

	"github.com/gig-tech/ovc-sdk-go/v2/ovc"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceOvcSizes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOvcSizesRead,

		Schema: map[string]*schema.Schema{
			"sizes_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cloudspace_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vcpus": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"memory": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceOvcSizesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*ovc.Client)
	sid, err := client.Sizes.GetByVcpusAndMemory(d.Get("vcpus").(int), d.Get("memory").(int), d.Get("cloudspace_id").(string))
	if err != nil {
		return err
	}
	id := strconv.Itoa(sid.ID)
	d.SetId(id)
	d.Set("sizes_id", sid.ID)
	d.Set("vcpus", sid.Vcpus)
	d.Set("name", sid.Name)
	d.Set("description", sid.Description)
	d.Set("memory", sid.Memory)
	return nil
}
