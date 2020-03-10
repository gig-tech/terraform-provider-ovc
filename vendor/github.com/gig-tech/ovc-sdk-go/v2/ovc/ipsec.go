package ovc

import (
	"encoding/json"
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
	body, err := s.client.Post("/cloudapi/ipsec/addTunnelToCloudspace", *ipsecConfig, OperationalActionTimeout)
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
	_, err := s.client.Post("/cloudapi/ipsec/removeTunnelFromCloudspace", *ipsecConfig, OperationalActionTimeout)
	return err
}

// List all ipsec of a cloudspace
func (s *IpsecServiceOp) List(ipsecConfig *IpsecConfig) (*IpsecList, error) {
	body, err := s.client.Post("/cloudapi/ipsec/listTunnels", *ipsecConfig, ModelActionTimeout)
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
