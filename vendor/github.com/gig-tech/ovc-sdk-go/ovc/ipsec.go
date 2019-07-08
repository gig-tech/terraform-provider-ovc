package ovc

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// IpsecConfig is used when creating/deleting/listing ipsec
type IpsecConfig struct {
	CloudspaceID         int    `json:"cloudspaceId"`
	RemotePublicAddr     string `json:"remotePublicAddr,omitempty"`
	RemotePrivateNetwork string `json:"remotePrivateNetwork,omitempty"`
	PskSecret            string `json:"pskSecret,omitempty"`
}

// IpsecList is a list of ipsec of a cloudspace
// Returned when using the List method
type IpsecList []struct {
	RemoteAddr           string `json:"remoteAddr"`
	RemotePrivateNetwork string `json:"remoteprivatenetwork"`
	PSK                  string `json:"psk"`
}

// IpsecService is an interface for interfacing with ipsec
// endpoints of the OVC API
type IpsecService interface {
	Create(*IpsecConfig) (string, error)
	List(*IpsecConfig) (*IpsecList, error)
	Delete(*IpsecConfig) error
}

// IpsecServiceOp handles communication with the ipsec related methods of the
// OVC API
type IpsecServiceOp struct {
	client *Client
}

// Create a new ipsec tunnel
func (s *IpsecServiceOp) Create(ipsecConfig *IpsecConfig) (string, error) {
	ipsecJSON, err := json.Marshal(*ipsecConfig)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/ipsec/addTunnelToCloudspace", bytes.NewBuffer(ipsecJSON))
	if err != nil {
		return "", err
	}
	body, err := s.client.Do(req)
	if err != nil {
		return "", err
	}

	result := ""
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	return result, nil
}

// Delete an existing ipsec tunnel
func (s *IpsecServiceOp) Delete(ipsecConfig *IpsecConfig) error {
	ipsecJSON, err := json.Marshal(*ipsecConfig)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/ipsec/removeTunnelFromCloudspace", bytes.NewBuffer(ipsecJSON))
	if err != nil {
		return err
	}
	_, err = s.client.Do(req)

	return err
}

// List all ipsec of a cloudspace
func (s *IpsecServiceOp) List(ipsecConfig *IpsecConfig) (*IpsecList, error) {
	ipsecJSON, err := json.Marshal(*ipsecConfig)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/ipsec/listTunnels", bytes.NewBuffer(ipsecJSON))
	if err != nil {
		return nil, err
	}
	body, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	ipsecList := new(IpsecList)
	err = json.Unmarshal(body, &ipsecList)
	if err != nil {
		return nil, err
	}

	return ipsecList, nil
}
