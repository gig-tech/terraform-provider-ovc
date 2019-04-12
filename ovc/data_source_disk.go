package ovc

import (
	"strconv"

	"github.com/gig-tech/ovc-sdk-go/ovc"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceOvcDisk() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOvcDiskRead,

		Schema: map[string]*schema.Schema{
			"disk_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size_max": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size_used": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceOvcDiskRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*ovc.Client)
	var disk *ovc.DiskInfo
	var err error
	if v, ok := d.GetOk("disk_id"); ok {
		disk, err = client.Disks.Get(v.(string))
	} else {
		disk, err = client.Disks.GetByName(d.Get("name").(string), d.Get("account_id").(string))
		if err != nil {
			return err
		}
	}
	d.SetId(strconv.Itoa(disk.ID))
	d.Set("disk_id", disk.ID)
	d.Set("size_max", disk.SizeMax)
	d.Set("account_id", disk.AccountID)
	d.Set("name", disk.Name)
	d.Set("description", disk.Descr)
	d.Set("type", disk.Type)
	d.Set("size_used", disk.SizeUsed)
	return nil
}
