package ovc

import (
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/gig-tech/ovc-sdk-go/ovc"
)

func dataSourceOvcPortForwarding() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOvcPortForwardingRead,

		Schema: map[string]*schema.Schema{
			"machine_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cloudspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"local_port": {
				Type:     schema.TypeString,
				Required: true,
			},
			"machine_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"local_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port_forward_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOvcPortForwardingRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*ovc.Client)
	var list *ovc.PortForwardingList
	if machineID, ok := d.GetOk("machine_id"); ok {
		if cloudspaceID, ok := d.GetOk("cloudspace_id"); ok {
			cID, err := strconv.Atoi(cloudspaceID.(string))
			mID, err := strconv.Atoi(machineID.(string))
			portForwardingConfig := &ovc.PortForwardingConfig{
				CloudspaceID: cID,
				MachineID:    mID,
			}
			list, err = client.Portforwards.List(portForwardingConfig)
			if err != nil {
				return err
			}
		}
	}
	for _, port := range *list {
		if port.LocalPort == d.Get("local_port") {
			d.SetId(strconv.Itoa(port.ID))
			d.Set("protocol", port.Protocol)
			d.Set("machine_name", port.MachineName)
			d.Set("public_ip", port.PublicIP)
			d.Set("local_ip", port.LocalIP)
			d.Set("public_port", port.PublicPort)
			d.Set("port_forward_id", port.ID)
		}
	}
	return nil

}
