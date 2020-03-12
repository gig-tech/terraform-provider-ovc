package ovc

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// ExternalNetworkConfig is used when getting an external network
type ExternalNetworkConfig struct {
	Name      string `json:"name"`
	AccountID int    `json:"accountId"`
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
	Get(int) (*ExternalNetworkInfo, error)
	List(int) (*[]ExternalNetworkInfo, error)
}

// ExternalNetworkServiceOp handles communication with the external network related methods of the
// OVC API
type ExternalNetworkServiceOp struct {
	client *Client
}

// Get external network
func (s *ExternalNetworkServiceOp) Get(id int) (*ExternalNetworkInfo, error) {
	externalNetworkIDMap := make(map[string]interface{})
	externalNetworkIDMap["id"] = id
	body, err := s.client.Post("/cloudapi/externalnetwork/get", externalNetworkIDMap, ModelActionTimeout)
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
			return s.Get(externalNetwork.ID)
		}
	}

	return nil, fmt.Errorf("External Network %s not found", name)
}

// List all external networks
func (s *ExternalNetworkServiceOp) List(accountID int) (*[]ExternalNetworkInfo, error) {
	accountIDMap := make(map[string]interface{})
	accountIDMap["accountId"] = accountID

	body, err := s.client.Post("/cloudapi/externalnetwork/list", accountIDMap, ModelActionTimeout)
	if err != nil {
		return nil, err
	}

	externalNetworks := new([]ExternalNetworkInfo)
	err = json.Unmarshal(body, &externalNetworks)
	if err != nil {
		return nil, err
	}

	return externalNetworks, nil
}
