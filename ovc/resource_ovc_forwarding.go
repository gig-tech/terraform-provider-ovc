package ovc

import (
	"strconv"

	"github.com/gig-tech/ovc-sdk-go/ovc"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePortForwarding() *schema.Resource {
	return &schema.Resource{
		Create: resourcePortForwardingCreate,
		Read:   resourcePortForwardingRead,
		Update: resourcePortForwardingUpdate,
		Delete: resourcePortForwardingDelete,
		Exists: resourcePortForwardingExists,

		Schema: map[string]*schema.Schema{
			"cloudspace_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"public_port": {
				Type:     schema.TypeInt,
				Optional: true,
				// DiffSuppressFunc suppresses change when a ramdom public port was created
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == "0"
				},
			},
			"machine_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"local_port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
			},
			"machine_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePortForwardingExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*ovc.Client)
	publicPort := d.Get("public_port").(int)
	cloudspaceID := d.Get("cloudspace_id").(int)
	machineID := d.Get("machine_id").(int)
	portForwardingConfig := &ovc.PortForwardingConfig{
		MachineID:    machineID,
		CloudspaceID: cloudspaceID,
		PublicPort:   publicPort,
	}
	_, err := client.Portforwards.Get(portForwardingConfig)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func resourcePortForwardingRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*ovc.Client)
	portForwardingConfig := ovc.PortForwardingConfig{}
	portForwardingConfig.CloudspaceID = d.Get("cloudspace_id").(int)
	portForwardingConfig.MachineID = d.Get("machine_id").(int)
	portForwardingList, err := client.Portforwards.List(&portForwardingConfig)
	if err != nil {
		return err
	}
	for _, pf := range *portForwardingList {
		publicPort := strconv.Itoa(d.Get("public_port").(int))
		if pf.PublicPort == publicPort {
			d.SetId(strconv.Itoa(pf.ID))
		}
	}
	return nil
}

func resourcePortForwardingCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*ovc.Client)
	portForwardingConfig := ovc.PortForwardingConfig{}
	portForwardingConfig.CloudspaceID = d.Get("cloudspace_id").(int)
	portForwardingConfig.PublicIP = d.Get("public_ip").(string)
	portForwardingConfig.PublicPort = d.Get("public_port").(int)
	portForwardingConfig.MachineID = d.Get("machine_id").(int)
	portForwardingConfig.LocalPort = d.Get("local_port").(int)
	portForwardingConfig.Protocol = d.Get("protocol").(string)
	publicPort, err := client.Portforwards.Create(&portForwardingConfig)
	if err != nil {
		return err
	}
	d.Set("public_port", publicPort)
	return resourcePortForwardingRead(d, m)

}

func resourcePortForwardingUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*ovc.Client)
	portForwardingConfig := ovc.PortForwardingConfig{}
	portForwardingConfig.CloudspaceID = d.Get("cloudspace_id").(int)
	needForUpdate := false
	if d.HasChange("public_ip") {
		sourcePublicIP, publicIP := d.GetChange("public_ip")
		portForwardingConfig.SourcePublicIP = sourcePublicIP.(string)
		portForwardingConfig.PublicIP = publicIP.(string)
		needForUpdate = true
	} else {
		portForwardingConfig.SourcePublicIP = d.Get("public_ip").(string)
		portForwardingConfig.PublicIP = d.Get("public_ip").(string)
	}
	if d.HasChange("public_port") {
		sourcePublicPort, publicPort := d.GetChange("public_port")
		portForwardingConfig.SourcePublicPort = sourcePublicPort.(int)
		portForwardingConfig.PublicPort = publicPort.(int)
		needForUpdate = true
	} else {
		portForwardingConfig.SourcePublicPort = d.Get("public_port").(int)
		portForwardingConfig.PublicPort = d.Get("public_port").(int)
	}
	if d.HasChange("protocol") {
		sourceProtocol, protocol := d.GetChange("protocol")
		portForwardingConfig.SourceProtocol = sourceProtocol.(string)
		portForwardingConfig.Protocol = protocol.(string)
		needForUpdate = true
	} else {
		portForwardingConfig.SourceProtocol = d.Get("protocol").(string)
		portForwardingConfig.Protocol = d.Get("protocol").(string)
	}
	portForwardingConfig.MachineID = d.Get("machine_id").(int)
	portForwardingConfig.LocalPort = d.Get("local_port").(int)
	if d.HasChange("machine_id") || d.HasChange("local_port") {
		needForUpdate = true
	}
	if needForUpdate {
		err := client.Portforwards.Update(&portForwardingConfig)
		if err != nil {
			return err
		}
	}

	return resourcePortForwardingRead(d, m)
}

func resourcePortForwardingDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*ovc.Client)
	portForwardingConfig := ovc.PortForwardingConfig{}
	portForwardingConfig.CloudspaceID = d.Get("cloudspace_id").(int)
	portForwardingConfig.PublicIP = d.Get("public_ip").(string)
	portForwardingConfig.PublicPort = d.Get("public_port").(int)
	portForwardingConfig.Protocol = d.Get("protocol").(string)
	err := client.Portforwards.Delete(&portForwardingConfig)
	if err != nil {
		return err
	}
	return nil
}
