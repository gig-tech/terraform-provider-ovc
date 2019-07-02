package ovc

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"

	"github.com/gig-tech/ovc-sdk-go/ovc"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceOvcImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOvcImageRead,

		Schema: map[string]*schema.Schema{
			"account": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"most_recent": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			// Computed values
			"image_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOvcImageRead(d *schema.ResourceData, m interface{}) error {
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

	filteredImages := make([]ovc.ImageInfo, 0)

	for _, image := range *images {
		// select images by name
		if nameRegex == "" || re.FindString(image.Name) != "" {
			filteredImages = append(filteredImages, image)
		}
	}

	if len(filteredImages) < 1 {
		return fmt.Errorf("No images were found with given criteria")
	}
	if len(filteredImages) > 1 {
		if !d.Get("most_recent").(bool) {
			return fmt.Errorf("More than one image matches given criteria. " +
				"Try more specific criteria or set `most_recent` attribute to true")
		}
		sort.Slice(filteredImages, func(i, j int) bool {
			return filteredImages[i].ID > filteredImages[j].ID
		})
	}

	d.Set("image_id", filteredImages[0].ID)
	if filteredImages[0].AccountID != 0 {
		d.Set("account", account)
	}
	d.Set("name", filteredImages[0].Name)
	d.SetId(strconv.Itoa(filteredImages[0].ID))
	return nil
}
