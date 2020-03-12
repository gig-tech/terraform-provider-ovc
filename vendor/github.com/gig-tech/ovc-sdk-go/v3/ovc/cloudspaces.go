package ovc

import (
	"encoding/json"
	"errors"
	"strconv"
)

// CloudSpaceConfig is used when creating a CloudSpace
type CloudSpaceConfig struct {
	CloudSpaceID           int     `json:"cloudspaceId,omitempty"`
	AccountID              int     `json:"accountId,omitempty"`
	Location               string  `json:"location,omitempty"`
	Name                   string  `json:"name,omitempty"`
	Access                 string  `json:"access,omitempty"`
	MaxMemoryCapacity      float64 `json:"maxMemoryCapacity,omitempty"`
	MaxCPUCapacity         int     `json:"maxCPUCapacity,omitempty"`
	MaxDiskCapacity        int     `json:"maxVDiskCapacity,omitempty"`
	MaxNetworkPeerTransfer int     `json:"maxNetworkPeerTransfer,omitempty"`
	MaxNumPublicIP         int     `json:"maxNumPublicIP,omitempty"`
	AllowedVMSizes         []int   `json:"allowedVMSizes,omitempty"`
	PrivateNetwork         string  `json:"privatenetwork,omitempty"`
	Mode                   string  `json:"mode,omitempty"`
	Type                   string  `json:"type,omitempty"`
	ExternalnetworkID      int     `json:"externalnetworkId,omitempty"`
}

// ResourceLimits contains all information related to resource limits
type ResourceLimits struct {
	CUM  float64 `json:"CU_M"`
	CUD  int     `json:"CU_D"`
	CUNP int     `json:"CU_NP"`
	CUI  int     `json:"CU_I"`
	CUC  int     `json:"CU_C"`
}

// CloudSpace contains all information related to a CloudSpace
type CloudSpace struct {
	Status            string         `json:"status"`
	UpdateTime        int            `json:"updateTime"`
	Externalnetworkip string         `json:"externalnetworkip"`
	Description       string         `json:"description"`
	ResourceLimits    ResourceLimits `json:"resourceLimits"`
	ID                int            `json:"id"`
	AccountID         int            `json:"accountId"`
	Name              string         `json:"name"`
	CreationTime      int            `json:"creationTime"`
	ACL               []ACL          `json:"acl"`
	Secret            string         `json:"secret"`
	GridID            int            `json:"gid"`
	Location          string         `json:"location"`
	Publicipaddress   string         `json:"publicipaddress"`
	PrivateNetwork    string         `json:"privatenetwork"`
	Type              string         `json:"type"`
	Mode              string         `json:"mode"`
}

// CloudSpaceInfo returns a list of CloudSpaces
type CloudSpaceInfo struct {
	Status            string     `json:"status"`
	UpdateTime        int        `json:"updateTime"`
	Externalnetworkip string     `json:"externalnetworkip"`
	Name              string     `json:"name"`
	Descr             string     `json:"descr"`
	CreationTime      int        `json:"creationTime"`
	ACL               []ACL      `json:"acl"`
	AccountACL        AccountACL `json:"accountAcl"`
	GridID            int        `json:"gid"`
	Location          string     `json:"location"`
	Mode              string     `json:"mode"`
	Type              string     `json:"type"`
	Publicipaddress   string     `json:"publicipaddress"`
	AccountName       string     `json:"accountName"`
	ID                int        `json:"id"`
	AccountID         int        `json:"accountId"`
}

// CloudSpaceDeleteConfig used to delete a CloudSpace
type CloudSpaceDeleteConfig struct {
	CloudSpaceID int  `json:"cloudspaceId"`
	Permanently  bool `json:"permanently"`
}

// CloudSpaceService is an interface for interfacing with the CloudSpace
// endpoints of the OVC API
type CloudSpaceService interface {
	List() (*[]CloudSpaceInfo, error)
	Get(int) (*CloudSpace, error)
	GetByNameAndAccount(string, string) (*CloudSpace, error)
	Create(*CloudSpaceConfig) (int, error)
	Update(*CloudSpaceConfig) error
	Delete(*CloudSpaceDeleteConfig) error
	SetDefaultGateway(int, string) error
}

// CloudSpaceServiceOp handles communication with the cloudspace related methods of the
// OVC API
type CloudSpaceServiceOp struct {
	client *Client
}

// List returns all cloudspaces
func (s *CloudSpaceServiceOp) List() (*[]CloudSpaceInfo, error) {
	cloudSpaceMap := make(map[string]interface{})
	cloudSpaceMap["includedeleted"] = false
	body, err := s.client.Post("/cloudapi/cloudspaces/list", cloudSpaceMap, ModelActionTimeout)
	if err != nil {
		return nil, err
	}
	cloudSpaces := new([]CloudSpaceInfo)
	err = json.Unmarshal(body, &cloudSpaces)
	if err != nil {
		return nil, err
	}

	return cloudSpaces, nil
}

// Get individual CloudSpace
func (s *CloudSpaceServiceOp) Get(id int) (*CloudSpace, error) {
	cloudSpaceIDMap := make(map[string]interface{})
	cloudSpaceIDMap["cloudspaceId"] = id

	body, err := s.client.Post("/cloudapi/cloudspaces/get", cloudSpaceIDMap, ModelActionTimeout)
	if err != nil {
		return nil, err
	}
	cloudSpace := new(CloudSpace)
	err = json.Unmarshal(body, &cloudSpace)
	if err != nil {
		s.client.logger.Debug("Unmarschalling result into cloudspace object failed")
		return nil, err
	}

	return cloudSpace, nil
}

// GetByNameAndAccount gets an individual cloudspace
func (s *CloudSpaceServiceOp) GetByNameAndAccount(cloudSpaceName string, account string) (*CloudSpace, error) {
	cloudspaces, err := s.client.CloudSpaces.List()
	if err != nil {
		return nil, err
	}
	for _, cp := range *cloudspaces {
		if cp.AccountName == account && cp.Name == cloudSpaceName {
			return s.client.CloudSpaces.Get(cp.ID)
		}
	}

	return nil, errors.New("Could not find cloudspace based on name")
}

// Create a new CloudSpace
func (s *CloudSpaceServiceOp) Create(cloudSpaceConfig *CloudSpaceConfig) (int, error) {
	body, err := s.client.Post("/cloudapi/cloudspaces/create", *cloudSpaceConfig, OperationalActionTimeout)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(body))
}

// Delete a CloudSpace
func (s *CloudSpaceServiceOp) Delete(cloudSpaceConfig *CloudSpaceDeleteConfig) error {
	_, err := s.client.Post("/cloudapi/cloudspaces/delete", *cloudSpaceConfig, OperationalActionTimeout)
	return err
}

// Update an existing CloudSpace
func (s *CloudSpaceServiceOp) Update(cloudSpaceConfig *CloudSpaceConfig) error {
	_, err := s.client.Post("/cloudapi/cloudspaces/update", *cloudSpaceConfig, ModelActionTimeout)
	return err
}

// SetDefaultGateway sets default gateway of the cloudspace to the given IP address
func (s *CloudSpaceServiceOp) SetDefaultGateway(id int, gateway string) error {
	csMap := make(map[string]interface{})
	csMap["cloudspaceId"] = id
	csMap["gateway"] = gateway

	_, err := s.client.Post("/cloudapi/cloudspaces/setDefaultGateway", csMap, OperationalActionTimeout)
	return err
}
