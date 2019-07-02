package ovc

import (
	"fmt"
	"log"
	"strconv"

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

		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			if diff.Id() != "" && diff.HasChange("image_id") {
				return fmt.Errorf("Cannot change Image ID on existing machine")
			}

			return nil
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
			"iops": {
				Type:     schema.TypeInt,
				Optional: true,
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
			"interfaces": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_id": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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
	client := m.(*ovc.Client)
	_, err := client.Machines.Get(d.Id())
	if err != nil {
		return false, nil
	}
	return true, nil
}

func resourceOvcMachineRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*ovc.Client)
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
	d.Set("cloudspace_id", machineInfo.CloudspaceID)
	d.Set("size_id", machineInfo.SizeID)
	d.Set("vcpus", machineInfo.Vcpus)
	d.Set("disks", flattenDisks(machineInfo))
	d.Set("interfaces", flattenNics(machineInfo))

	return nil
}

func resourceOvcMachineCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*ovc.Client)
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
	log.Println("[DEBUG] New machine ID: " + machineID)
	d.SetId(machineID)
	log.Println("[DEBUG] Resource machine ID: " + d.Id())
	machineIDInt, err := strconv.Atoi(machineID)
	if err != nil {
		return err
	}
	// Set IOPS boot disk
	iops := d.Get("iops")
	if iops != nil {
		bootDiskID, err := GetBootDiskID(client, machineID)
		if err != nil {
			return err
		}
		diskConfig := &ovc.DiskConfig{
			DiskID: bootDiskID,
			IOPS:   iops.(int),
		}
		err = client.Disks.Update(diskConfig)
		if err != nil {
			return err
		}
	}
	// attach to external networks
	if v, ok := d.GetOk("interfaces"); ok {
		nics := v.([]interface{})
		for _, nici := range nics {
			var networkID int
			if nici != nil {
				nic := nici.(map[string]interface{})
				if nic["network_id"] != nil {
					// if network ID is given, attach to this network
					networkID = nic["network_id"].(int)
				}
			}
			if err := client.Machines.AddExternalIP(machineIDInt, networkID); err != nil {
				return err
			}
		}
	}
	return resourceOvcMachineRead(d, m)
}

func resourceOvcMachineUpdate(d *schema.ResourceData, m interface{}) error {

	var err error
	client := m.(*ovc.Client)
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

	if d.HasChange("iops") || d.HasChange("disksize") {
		bootDiskID, err := GetBootDiskID(client, d.Id())
		if err != nil {
			return err
		}

		log.Println("[DEBUG] Updating machine boot disk")
		diskConfig := &ovc.DiskConfig{
			DiskID: bootDiskID,
		}
		if d.HasChange("iops") {
			diskConfig.IOPS = d.Get("iops").(int)
		}
		if d.HasChange("disksize") {
			diskConfig.Size = d.Get("disksize").(int)
		}
		err = client.Disks.Update(diskConfig)
		if err != nil {
			return err
		}
	}

	if d.HasChange("interfaces") {
		machineIDInt, err := strconv.Atoi(machineConfig.MachineID)
		if err != nil {
			return err
		}
		if _, ok := d.GetOk("interfaces"); ok {
			old, new := d.GetChange("interfaces")
			oldNics := old.([]interface{})
			newNics := new.([]interface{})
			oldNetworks := countAttachedNetworks(oldNics)
			newNetworks := countAttachedNetworks(newNics)

			for networkID := range oldNetworks {
				if len(newNetworks[networkID]) == 0 {
					log.Println("[DEBUG] Detaching from external network")
					if err = client.Machines.DeleteExternalIP(machineIDInt, networkID, ""); err != nil {
						return err
					}
				}
			}

			for networkID, ips := range newNetworks {
				if err != nil {
					return err
				}
				switch diffNet := len(newNetworks[networkID]) - len(oldNetworks[networkID]); {
				case diffNet > 0:
					for i := 0; i < diffNet; i++ {
						// add missing external IPs
						log.Println("[DEBUG] Attaching to external network")
						if err := client.Machines.AddExternalIP(machineIDInt, networkID); err != nil {
							return err
						}
					}
				case diffNet < 0:
					for i := 0; i > diffNet; i-- {
						log.Println("[DEBUG] Deleting external IP")
						if err := client.Machines.DeleteExternalIP(machineIDInt, networkID, ips[0]); err != nil {
							return err
						}
						ips = ips[1:len(ips)]
					}
				}
			}
		} else {
			// delete all interfaces
			log.Println("[DEBUG] Detaching from all external networks")
			if err := client.Machines.DeleteExternalIP(machineIDInt, 0, ""); err != nil {
				return err
			}
		}
	}
	return resourceOvcMachineRead(d, m)
}

func countAttachedNetworks(nics []interface{}) map[int][]string {
	attachedNetworks := make(map[int][]string)
	for _, nicInterface := range nics {
		var networkID int
		var networkIP string
		if nicInterface != nil {
			nic := nicInterface.(map[string]interface{})
			networkID = nic["network_id"].(int)
			networkIP = nic["ip_address"].(string)
		}
		attachedNetworks[networkID] = append(attachedNetworks[networkID], networkIP)
	}
	return attachedNetworks
}

func resourceOvcMachineDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*ovc.Client)
	machineConfig := ovc.MachineConfig{}
	machineConfig.MachineID = d.Id()
	machineConfig.Permanently = true
	err := client.Machines.Delete(&machineConfig)
	return err
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

func flattenNics(machineInfo *ovc.MachineInfo) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)
	if machineInfo != nil {
		for _, cp := range machineInfo.Interfaces {
			if cp.Type == "PUBLIC" {
				// show only public interfaces
				nic := make(map[string]interface{})
				nic["network_id"] = cp.NetworkID
				nic["ip_address"] = cp.IPAddress
				result = append(result, nic)
			}
		}
		log.Printf("nics in map: %v", result)
	}
	return result
}

// GetBootDiskID gets ID of the boot disk of the machine
func GetBootDiskID(client *ovc.Client, id string) (int, error) {
	machineInfo, err := client.Machines.Get(id)
	if err != nil {
		return 0, err
	}
	for _, disk := range machineInfo.Disks {
		if disk.Type == "B" {
			return disk.ID, nil
		}
	}
	return 0, fmt.Errorf("Machine %s has no boot disk", machineInfo.Name)
}
