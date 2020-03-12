package ovc

import (
	"encoding/json"
)

// Template is a list of templates
// Returned when using the List method
type Template struct {
	Username    interface{} `json:"username"`
	Status      string      `json:"status"`
	Description string      `json:"description"`
	Name        string      `json:"name"`
	Size        int         `json:"size"`
	Type        string      `json:"type"`
	ID          int         `json:"id"`
	AccountID   int         `json:"accountId"`
}

// TemplateService is an interface for interfacing with the Images
// endpoints of the OVC API
type TemplateService interface {
	List(int) (*[]Template, error)
}

// TemplateServiceOp handles communication with the image related methods of the
// OVC API
type TemplateServiceOp struct {
	client *Client
}

// List all images
func (s *TemplateServiceOp) List(accountID int) (*[]Template, error) {
	templateMap := make(map[string]interface{})
	templateMap["accountId"] = 4

	body, err := s.client.Post("/cloudapi/images/list", templateMap, ModelActionTimeout)
	if err != nil {
		return nil, err
	}

	templates := new([]Template)
	err = json.Unmarshal(body, &templates)
	if err != nil {
		return nil, err
	}

	return templates, nil
}
