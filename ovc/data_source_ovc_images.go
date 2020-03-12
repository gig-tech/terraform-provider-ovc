package ovc

import (
	"regexp"
	"strconv"

	"github.com/gig-tech/ovc-sdk-go/v3/ovc"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceOvcImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOvcImagesRead,

		Schema: map[string]*schema.Schema{
			"account": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"account": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceOvcImagesRead(d *schema.ResourceData, m interface{}) error {
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

	images, err := client.Images.List(accountID)
	if err != nil {
		return err
	}

	nameRegex := d.Get("name_regex").(string)
	re := regexp.MustCompile(nameRegex)

	entities := make([]map[string]interface{}, len(*images))

	for i, image := range *images {
		// select images by name
		if nameRegex == "" || re.FindString(image.Name) != "" {
			entity := make(map[string]interface{})
			entity["id"] = strconv.Itoa(image.ID)
			entity["name"] = image.Name
			if accountID != 0 && image.AccountID == accountID {
				entity["account"] = account
			} else {
				entity["account"] = ""
			}
			entities[i] = entity
		}
	}

	if err = d.Set("entities", entities); err != nil {
		return err
	}

	d.SetId("1")
	return nil
}
