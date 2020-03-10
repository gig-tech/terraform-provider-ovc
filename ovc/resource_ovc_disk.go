package ovc

import (
	"strconv"

	"github.com/gig-tech/ovc-sdk-go/v2/ovc"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOvcDisk() *schema.Resource {
	return &schema.Resource{
		Create: resourceOvcDiskCreate,
		Read:   resourceOvcDiskRead,
		Update: resourceOvcDiskUpdate,
		Delete: resourceOvcDiskDelete,
		Exists: resourceOvcDiskExists,

		Schema: map[string]*schema.Schema{
			"machine_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"disk_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ssd_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"iops": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceOvcDiskExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*ovc.Client)
	disk, err := client.Disks.Get(d.Id())
	if err != nil || disk.Status == "DESTROYED" {
		return false, nil
	}

	return true, nil
}

func resourceOvcDiskRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*ovc.Client)
	diskID := d.Id()
	_, err := client.Disks.Get(diskID)
	if err != nil {
		d.SetId("")
		return nil
	}

	return nil
}

func resourceOvcDiskCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*ovc.Client)
	diskConfig := ovc.DiskConfig{}
	diskConfig.MachineID = d.Get("machine_id").(int)
	diskConfig.DiskName = d.Get("disk_name").(string)
	diskConfig.Description = d.Get("description").(string)
	diskConfig.Size = d.Get("size").(int)
	diskConfig.Type = d.Get("type").(string)
	diskConfig.SSDSize = d.Get("ssd_size").(int)
	diskConfig.IOPS = d.Get("iops").(int)
	diskID, err := client.Disks.CreateAndAttach(&diskConfig)
	if err != nil {
		return err
	}
	d.SetId(diskID)

	return resourceOvcDiskRead(d, m)
}

func resourceOvcDiskUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*ovc.Client)
	diskConfig := ovc.DiskConfig{}
	update := false
	diskID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	diskConfig.DiskID = diskID

	if d.HasChange("size") {
		diskConfig.Size = d.Get("size").(int)
		update = true
	}

	if d.HasChange("iops") {
		diskConfig.IOPS = d.Get("iops").(int)
		update = true
	}

	if update {
		err = client.Disks.Update(&diskConfig)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceOvcDiskDelete(d *schema.ResourceData, m interface{}) error {
	defer ovc.ReleaseLock(d.Get("machine_id").(int))
	ovc.GetLock(d.Get("machine_id").(int))
	client := m.(*ovc.Client)
	diskConfig := ovc.DiskDeleteConfig{}
	diskID, err := strconv.Atoi(d.Id())
	diskConfig.DiskID = diskID
	if err != nil {
		return err
	}
	diskConfig.Detach = true
	diskConfig.Permanently = true
	err = client.Disks.Delete(&diskConfig)
	return err
}
