package ovc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// ExternalNetworkConfig is used when getting an external network
type ExternalNetworkConfig struct {
	Name      string `json:"name"`
	AccountID int    `json:"accountId"`
}

// ExternalNetworkList is a list of external networks
// Returned when using the List method
type ExternalNetworkList []struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	AccountID  int    `json:"accountId"`
	Network    string `json:"network"`
	Gateway    string `json:"gateway"`
	Subnetmask string `json:"subnetmask"`
	DHCP       bool   `json:"dhcp"`
}

// ExternalNetworkInfo contains information about the external network returned by API
type ExternalNetworkInfo struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	AccountID  int    `json:"accountId"`
	Network    string `json:"network"`
	Gateway    string `json:"gateway"`
	Subnetmask string `json:"subnetmask"`
	DHCP       bool   `json:"dhcp"`
}

// ExternalNetworkService is an interface for interfacing with the external networks
// of the OVC API
type ExternalNetworkService interface {
	Get(string) (*ExternalNetworkInfo, error)
	List(int) (*ExternalNetworkList, error)
}

// ExternalNetworkServiceOp handles communication with the external network related methods of the
// OVC API
type ExternalNetworkServiceOp struct {
	client *Client
}

// Get external network
func (s *ExternalNetworkServiceOp) Get(id string) (*ExternalNetworkInfo, error) {
	externalNetworkIDMap := make(map[string]interface{})
	externalNetworkIDMap["id"], _ = strconv.Atoi(id)
	externalNetworkJSON, err := json.Marshal(externalNetworkIDMap)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/externalnetwork/get", bytes.NewBuffer(externalNetworkJSON))
	if err != nil {
		return nil, err
	}
	body, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	externalNetworkInfo := new(ExternalNetworkInfo)
	err = json.Unmarshal(body, &externalNetworkInfo)
	if err != nil {
		return nil, err
	}

	return externalNetworkInfo, nil
}

// GetByName gets an individual external network from its name
func (s *ExternalNetworkServiceOp) GetByName(name string, accountID string) (*ExternalNetworkInfo, error) {
	aid, err := strconv.Atoi(accountID)
	if err != nil {
		return nil, err
	}
	externalNetworks, err := s.List(aid)
	if err != nil {
		return nil, err
	}
	for _, externalNetwork := range *externalNetworks {
		if externalNetwork.Name == name {
			return s.Get(strconv.Itoa(externalNetwork.ID))
		}
	}

	return nil, fmt.Errorf("External Network %s not found", name)
}

// List all external networks
func (s *ExternalNetworkServiceOp) List(accountID int) (*ExternalNetworkList, error) {
	accountIDMap := make(map[string]interface{})
	accountIDMap["accountId"] = accountID
	accountIDJson, err := json.Marshal(accountIDMap)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/externalnetwork/list", bytes.NewBuffer(accountIDJson))
	if err != nil {
		return nil, err
	}
	body, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	externalNetworks := new(ExternalNetworkList)
	err = json.Unmarshal(body, &externalNetworks)
	if err != nil {
		return nil, err
	}

	return externalNetworks, nil
}
