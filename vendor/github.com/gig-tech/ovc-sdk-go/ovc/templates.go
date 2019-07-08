package ovc

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// TemplateList is a list of templates
// Returned when using the List method
type TemplateList []struct {
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
	List(int) (*TemplateList, error)
}

// TemplateServiceOp handles communication with the image related methods of the
// OVC API
type TemplateServiceOp struct {
	client *Client
}

// List all images
func (s *TemplateServiceOp) List(accountID int) (*TemplateList, error) {
	templateMap := make(map[string]interface{})
	templateMap["accountId"] = 4
	templateJSON, err := json.Marshal(templateMap)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/images/list", bytes.NewBuffer(templateJSON))
	if err != nil {
		return nil, err
	}
	body, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	templates := new(TemplateList)
	err = json.Unmarshal(body, &templates)
	if err != nil {
		return nil, err
	}

	return templates, nil
}
