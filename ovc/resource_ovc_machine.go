package ovc

import (
	"log"

	"github.com/gig-tech/ovc-sdk-go/ovc"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOvcMachine() *schema.Resource {
	return &schema.Resource{
		Create: resourceOvcMachineCreate,
		Read:   resourceOvcMachineRead,
		Update: resourceOvcMachineUpdate,
		Delete: resourceOvcMachineDelete,
		Exists: resourceOvcMachineExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"cloudspace_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"size_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"memory", "vcpus"},
				Computed:      true,
			},
			"memory": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"size_id"},
				Computed:      true,
			},
			"vcpus": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"size_id"},
				Computed:      true,
			},
			"image_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"disksize": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeFloat,
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
			"username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"password": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"size_max": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"userdata": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceOvcMachineExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*ovc.OvcClient)
	_, err := client.Machines.Get(d.Id())
	if err != nil {
		return false, nil
	}
	return true, nil
}

func resourceOvcMachineRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*ovc.OvcClient)
	machineInfo, err := client.Machines.Get(d.Id())
	if err != nil {
		log.Println("machine not found error")
		d.SetId("")
		log.Println("machine id na read:" + d.Id())
		return nil
	}
	d.Set("hostname", machineInfo.Hostname)
	if len(machineInfo.Accounts) > 0 {
		d.Set("username", machineInfo.Accounts[0].Login)
		d.Set("password", machineInfo.Accounts[0].Password)
	}
	if len(machineInfo.Interfaces) > 0 {
		d.Set("ip_address", machineInfo.Interfaces[0].IPAddress)
	}
	d.Set("memory", machineInfo.Memory)
	d.Set("name", machineInfo.Name)
	d.Set("description", machineInfo.Description)
	d.Set("cloudspace_id", machineInfo.Cloudspaceid)
	d.Set("size_id", machineInfo.Sizeid)
	d.Set("vcpus", machineInfo.Vcpus)
	d.Set("disks", flattenDisks(machineInfo))

	return nil
}

func resourceOvcMachineCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*ovc.OvcClient)
	machineConfig := ovc.MachineConfig{}
	machineConfig.CloudspaceID = d.Get("cloudspace_id").(int)
	machineConfig.Name = d.Get("name").(string)
	machineConfig.Description = d.Get("description").(string)
	machineConfig.SizeID = d.Get("size_id").(int)
	machineConfig.ImageID = d.Get("image_id").(int)
	machineConfig.Disksize = d.Get("disksize").(int)
	machineConfig.Memory = d.Get("memory").(int)
	machineConfig.Vcpus = d.Get("vcpus").(int)
	machineConfig.Userdata = d.Get("userdata").(string)
	machineID, err := client.Machines.Create(&machineConfig)
	if err != nil {
		return err
	}
	log.Println("NEW MACHINE ID: " + machineID)
	d.SetId(machineID)
	log.Println(d.Id())
	return resourceOvcMachineRead(d, m)

}

func resourceOvcMachineUpdate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ovc.OvcClient)
	machineConfig := ovc.MachineConfig{}
	machineConfig.MachineID = d.Id()
	if d.HasChange("name") {
		machineConfig.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		machineConfig.Description = d.Get("description").(string)
		_, err = client.Machines.Update(&machineConfig)
		if err != nil {
			return err
		}
	}
	if d.HasChange("memory") || d.HasChange("vcpus") {
		machineConfig.Memory = d.Get("memory").(int)
		machineConfig.Vcpus = d.Get("vcpus").(int)
		_, err = client.Machines.Resize(&machineConfig)
		if err != nil {
			return err
		}
	}

	return resourceOvcMachineRead(d, m)
}

func resourceOvcMachineDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*ovc.OvcClient)
	machineConfig := ovc.MachineConfig{}
	machineConfig.MachineID = d.Id()
	machineConfig.Permanently = true
	err := client.Machines.Delete(&machineConfig)
	if err != nil {
		return err
	}
	return nil
}

func flattenDisks(machineInfo *ovc.MachineInfo) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)

	if machineInfo != nil {
		for _, disk := range machineInfo.Disks {
			diskinfo := make(map[string]interface{})
			diskinfo["status"] = disk.Status
			diskinfo["size_max"] = disk.SizeMax
			diskinfo["name"] = disk.Name
			diskinfo["description"] = disk.Descr
			diskinfo["type"] = disk.Type
			diskinfo["id"] = disk.ID

			result = append(result, diskinfo)
		}
		log.Printf("disks in map: %v", result)
	}
	return result
}
