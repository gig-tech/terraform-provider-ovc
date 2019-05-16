package ovc

import (
	"fmt"

	"github.com/gig-tech/ovc-sdk-go/ovc"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceIpsec() *schema.Resource {
	return &schema.Resource{
		Create: resourceIpsecCreate,
		Read:   resourceIpsecRead,
		Delete: resourceIpsecDelete,

		Schema: map[string]*schema.Schema{
			"cloudspace_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"remote_public_ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"remote_private_network": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"psk": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceIpsecRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*ovc.Client)
	ipsecConfig := ovc.IpsecConfig{}
	ipsecConfig.CloudspaceID = d.Get("cloudspace_id").(int)
	tunnelList, err := client.Ipsec.List(&ipsecConfig)
	if err != nil {
		return err
	}
	for _, tunnel := range *tunnelList {
		remotePublicIP := d.Get("remote_public_ip").(string)
		remotePrivateNetwork := d.Get("remote_private_network").(string)
		if tunnel.RemoteAddr == remotePublicIP && tunnel.RemotePrivateNetwork == remotePrivateNetwork {
			tunnelID := fmt.Sprintf("%s:%s", remotePublicIP, remotePrivateNetwork)
			d.SetId(tunnelID)
		}
	}
	return nil
}

func resourceIpsecCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*ovc.Client)
	ipsecConfig := ovc.IpsecConfig{}
	ipsecConfig.CloudspaceID = d.Get("cloudspace_id").(int)
	ipsecConfig.RemotePublicAddr = d.Get("remote_public_ip").(string)
	ipsecConfig.RemotePrivateNetwork = d.Get("remote_private_network").(string)
	ipsecConfig.PskSecret = d.Get("psk").(string)
	PskSecret, err := client.Ipsec.Create(&ipsecConfig)
	if err != nil {
		return err
	}
	d.Set("psk", PskSecret)
	return resourceIpsecRead(d, m)

}

func resourceIpsecDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*ovc.Client)
	ipsecConfig := ovc.IpsecConfig{}
	ipsecConfig.CloudspaceID = d.Get("cloudspace_id").(int)
	ipsecConfig.RemotePublicAddr = d.Get("remote_public_ip").(string)
	ipsecConfig.RemotePrivateNetwork = d.Get("remote_private_network").(string)
	err := client.Ipsec.Delete(&ipsecConfig)
	if err != nil {
		return err
	}
	return nil
}
