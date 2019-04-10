package ovc

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gig-tech/ovc-sdk-go/ovc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOvcCloudSpace() *schema.Resource {
	return &schema.Resource{
		Create: resourceOvcCloudSpaceCreate,
		Read:   resourceOvcCloudSpaceRead,
		Update: resourceOvcCloudSpaceUpdate,
		Delete: resourceOvcCloudSpaceDelete,
		Exists: resourceOvcCloudspaceExists,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"account": {
				Type:     schema.TypeString,
				Required: true,
			},
			"external_network_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"allowed_vm_sizes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"private_network": {
				Type:     schema.TypeString,
				Optional: true,
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
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_memory_capacity": {
							Type:     schema.TypeFloat,
							Optional: true,
							Default:  -1.0,
						},
						"max_disk_capacity": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  -1,
						},
						"max_cpu_capacity": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  -1,
						},
						"max_network_peer_transfer": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  -1,
						},
						"max_num_public_ip": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  -1,
						},
					},
				},
			},
		},
	}
}

func resourceOvcCloudspaceExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*ovc.OvcClient)
	cloudspace, err := client.CloudSpaces.Get(d.Id())
	if err != nil || cloudspace.Status == "DESTROYED" {
		return false, nil
	}
	return true, nil
}

func resourceOvcCloudSpaceRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*ovc.OvcClient)
	cloudSpaceID := d.Id()
	cloudSpace, err := client.CloudSpaces.Get(cloudSpaceID)
	if err != nil {
		return err
	}
	if cloudSpace.Status == "DESTROYED" {
		log.Println("cloudspace destroyed baby")
		d.SetId("")
		return nil
	}
	d.Set("status", cloudSpace.Status)
	rl := make(map[string]interface{})
	rl["max_memory_capacity"] = strconv.FormatFloat(cloudSpace.ResourceLimits.CUM, 'f', -1, 64)
	rl["max_disk_capacity"] = strconv.Itoa(cloudSpace.ResourceLimits.CUD)
	rl["max_cpu_capacity"] = strconv.Itoa(cloudSpace.ResourceLimits.CUC)
	rl["max_network_peer_transfer"] = strconv.Itoa(cloudSpace.ResourceLimits.CUNP)
	rl["max_num_public_ip"] = strconv.Itoa(cloudSpace.ResourceLimits.CUI)
	d.Set("resource_limits", rl)
	d.Set("description", cloudSpace.Description)
	d.Set("external_network_ip", cloudSpace.Externalnetworkip)
	d.Set("location", cloudSpace.Location)
	return nil

}

func resourceOvcCloudSpaceCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*ovc.OvcClient)
	account := d.Get("account").(string)
	accountID, err := client.Accounts.GetIDByName(account)
	if err != nil {
		return err
	}
	location := client.GetLocation()
	cloudSpaceConfig := ovc.CloudSpaceConfig{
		Access:                 client.Access,
		AccountID:              accountID,
		Location:               location,
		Name:                   d.Get("name").(string),
		MaxMemoryCapacity:      -1,
		MaxCPUCapacity:         -1,
		MaxDiskCapacity:        -1,
		MaxNetworkPeerTransfer: -1,
		MaxNumPublicIP:         -1,
	}

	if v, ok := d.GetOk("resource_limits"); ok {
		rL := v.(map[string]interface{})
		if rL["max_memory_capacity"] != nil {
			v2, err := strconv.ParseFloat(rL["max_memory_capacity"].(string), 64)
			if err != nil {
				return err
			}
			cloudSpaceConfig.MaxMemoryCapacity = v2
		}
		if rL["max_disk_capacity"] != nil {
			v2, err := strconv.Atoi(rL["max_disk_capacity"].(string))
			if err != nil {
				return err
			}
			cloudSpaceConfig.MaxDiskCapacity = v2
		}
		if rL["max_cpu_capacity"] != nil {
			v2, err := strconv.Atoi(rL["max_cpu_capacity"].(string))
			if err != nil {
				return err
			}
			cloudSpaceConfig.MaxCPUCapacity = v2
		}
		if rL["max_network_peer_transfer"] != nil {
			v2, err := strconv.Atoi(rL["max_network_peer_transfer"].(string))
			if err != nil {
				return err
			}
			cloudSpaceConfig.MaxNetworkPeerTransfer = v2
		}
		if rL["max_num_public_ip"] != nil {
			v2, err := strconv.Atoi(rL["max_num_public_ip"].(string))
			if err != nil {
				return err
			}
			cloudSpaceConfig.MaxNumPublicIP = v2
		}
	}
	cloudspaceID, err := client.CloudSpaces.Create(&cloudSpaceConfig)
	d.SetId(cloudspaceID)
	if err != nil {
		return err
	}
	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		cloudspace, _ := client.CloudSpaces.Get(cloudspaceID)

		if cloudspace.Status != "DEPLOYED" {
			log.Print("[DEBUG] Cloudspace is still deploying")
			return resource.RetryableError(fmt.Errorf("Cloudspace is in state: %s", cloudspace.Status))
		}
		return resource.NonRetryableError(resourceOvcCloudSpaceRead(d, m))
	})
}

func resourceOvcCloudSpaceUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*ovc.OvcClient)
	if d.HasChange("resource_limits") {
		cloudSpaceID, _ := strconv.Atoi(d.Id())
		cloudSpaceConfig := ovc.CloudSpaceConfig{
			MaxMemoryCapacity:      -1,
			MaxCPUCapacity:         -1,
			MaxDiskCapacity:        -1,
			MaxNetworkPeerTransfer: -1,
			MaxNumPublicIP:         -1,
		}
		cloudSpaceConfig.CloudSpaceID = cloudSpaceID
		cloudSpaceConfig.Name = d.Get("name").(string)
		if v, ok := d.GetOk("resource_limits"); ok {
			rl := v.(map[string]interface{})
			if rl["max_memory_capacity"] != nil {
				val, err := strconv.ParseFloat(rl["max_memory_capacity"].(string), 64)
				if err != nil {
					return err
				}
				cloudSpaceConfig.MaxMemoryCapacity = val
			}
			if rl["max_cpu_capacity"] != nil {
				val, err := strconv.Atoi(rl["max_cpu_capacity"].(string))
				if err != nil {
					return err
				}
				cloudSpaceConfig.MaxCPUCapacity = val
			}
			if rl["max_disk_capacity"] != nil {
				log.Println("has change")
				val, err := strconv.Atoi(rl["max_disk_capacity"].(string))
				if err != nil {
					return err
				}
				cloudSpaceConfig.MaxDiskCapacity = val
			}
			if rl["max_network_peer_transfer"] != nil {
				val, err := strconv.Atoi(rl["max_network_peer_transfer"].(string))
				if err != nil {
					return err
				}
				cloudSpaceConfig.MaxNetworkPeerTransfer = val
			}
			if rl["max_num_public_ip"] != nil {
				val, err := strconv.Atoi(rl["max_num_public_ip"].(string))
				if err != nil {
					return err
				}
				cloudSpaceConfig.MaxNumPublicIP = val
			}
		}
		err := client.CloudSpaces.Update(&cloudSpaceConfig)

		if err != nil {
			return err
		}
	}
	return resourceOvcCloudSpaceRead(d, m)
}

func resourceOvcCloudSpaceDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*ovc.OvcClient)
	cloudSpaceConfig := ovc.CloudSpaceDeleteConfig{}
	cloudSpaceID, err := strconv.Atoi(d.Id())
	cloudSpaceConfig.CloudSpaceID = cloudSpaceID
	if err != nil {
		return err
	}
	cloudSpaceConfig.Permanently = true
	err = client.CloudSpaces.Delete(&cloudSpaceConfig)
	if err != nil {
		return err
	}
	return nil
}
