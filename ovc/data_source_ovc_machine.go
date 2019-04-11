package ovc

import (
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/gig-tech/ovc-sdk-go/ovc"
)

func dataSourceOvcMachine() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOvcMachineRead,

		Schema: map[string]*schema.Schema{
			"machine_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloudspace_id": {
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
			"size_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"image_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"disksize": {
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
			"os_image": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"locked": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"accounts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"guid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"login": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"password": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"interfaces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mac_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"guid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_id": {
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

func dataSourceOvcMachineRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*ovc.Client)
	var machine *ovc.MachineInfo
	var err error
	if v, ok := d.GetOk("machine_id"); ok {
		machine, err = client.Machines.Get(v.(string))
		if err != nil {
			return err
		}
	} else {
		machine, err = client.Machines.GetByName(d.Get("name").(string), d.Get("cloudspace_id").(string))

		if err != nil {
			return err
		}
	}
	d.SetId(strconv.Itoa(machine.ID))
	d.Set("status", machine.Status)
	d.Set("name", machine.Name)
	d.Set("size_id", machine.Sizeid)
	d.Set("description", machine.Description.(string))
	d.Set("update_time", machine.UpdateTime)
	d.Set("cloudspace_id", strconv.Itoa(machine.Cloudspaceid))
	d.Set("machine_id", machine.ID)
	d.Set("image_id", machine.Imageid)
	d.Set("hostname", machine.Hostname)
	d.Set("creationtime", machine.CreationTime)
	d.Set("os_image", machine.OsImage)
	d.Set("storage", machine.Storage)
	d.Set("locked", machine.Locked)
	interfaces := make([]map[string]interface{}, len(machine.Interfaces))
	for i := range machine.Interfaces {
		machineInterface := make(map[string]interface{})
		machineInterface["status"] = machine.Interfaces[i].Status
		machineInterface["mac_address"] = machine.Interfaces[i].MacAddress
		machineInterface["ip_address"] = machine.Interfaces[i].IPAddress
		machineInterface["guid"] = machine.Interfaces[i].GUID
		machineInterface["type"] = machine.Interfaces[i].Type
		machineInterface["network_id"] = machine.Interfaces[i].NetworkID
		interfaces[i] = machineInterface
	}
	d.Set("interfaces", interfaces)
	accounts := make([]map[string]interface{}, len(machine.Accounts))
	for i := range machine.Accounts {
		account := make(map[string]interface{})
		account["guid"] = machine.Accounts[i].GUID
		account["login"] = machine.Accounts[i].Login
		account["password"] = machine.Accounts[i].Password
		accounts[i] = account
	}
	d.Set("accounts", accounts)
	return nil

}
