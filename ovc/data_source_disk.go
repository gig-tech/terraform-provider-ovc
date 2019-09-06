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
			"account": {
				Type:     schema.TypeString,
				Optional: true,
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
				Optional: true,
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
	var err error
	var accountID int
	account := d.Get("account")
	if account != "" {
		accountID, err = client.Accounts.GetIDByName(account.(string))
		if err != nil {
			return err
		}
	}
	var disk *ovc.DiskInfo
	if v, ok := d.GetOk("disk_id"); ok {
		disk, err = client.Disks.Get(v.(string))
	} else {
		disk, err = client.Disks.GetByName(d.Get("name").(string), accountID, d.Get("type").(string))
	}
	if err != nil {
		return err
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
