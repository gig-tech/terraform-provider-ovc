package ovc

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
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
	PrivateNetwork         string  `json:"privatenetwork"`
	Mode                   string  `json:"mode"`
	Type                   string  `json:"type"`
	ExternalnetworkID      int     `json:"externalnetworkId"`
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
	ACL               []struct {
		Status       string `json:"status"`
		CanBeDeleted bool   `json:"canBeDeleted"`
		Right        string `json:"right"`
		Type         string `json:"type"`
		UserGroupID  string `json:"userGroupId"`
	} `json:"acl"`
	Secret          string `json:"secret"`
	GridID          int    `json:"gid"`
	Location        string `json:"location"`
	Publicipaddress string `json:"publicipaddress"`
	PrivateNetwork  string `json:"privatenetwork"`
	Type            string `json:"type"`
	Mode            string `json:"mode"`
}

// CloudSpaceList returns a list of CloudSpaces
type CloudSpaceList []struct {
	Status            string `json:"status"`
	UpdateTime        int    `json:"updateTime"`
	Externalnetworkip string `json:"externalnetworkip"`
	Name              string `json:"name"`
	Descr             string `json:"descr"`
	CreationTime      int    `json:"creationTime"`
	ACL               []struct {
		Status       string `json:"status"`
		CanBeDeleted bool   `json:"canBeDeleted"`
		Right        string `json:"right"`
		Type         string `json:"type"`
		UserGroupID  string `json:"userGroupId"`
	} `json:"acl"`
	AccountACL struct {
		Status      string `json:"status"`
		Right       string `json:"right"`
		Explicit    bool   `json:"explicit"`
		UserGroupID string `json:"userGroupId"`
		GUID        string `json:"guid"`
		Type        string `json:"type"`
	} `json:"accountAcl"`
	GridID          int    `json:"gid"`
	Location        string `json:"location"`
	Publicipaddress string `json:"publicipaddress"`
	AccountName     string `json:"accountName"`
	ID              int    `json:"id"`
	AccountID       int    `json:"accountId"`
}

// CloudSpaceDeleteConfig used to delete a CloudSpace
type CloudSpaceDeleteConfig struct {
	CloudSpaceID int  `json:"cloudspaceId"`
	Permanently  bool `json:"permanently"`
}

// CloudSpaceService is an interface for interfacing with the CloudSpace
// endpoints of the OVC API
type CloudSpaceService interface {
	List() (*CloudSpaceList, error)
	Get(string) (*CloudSpace, error)
	GetByNameAndAccount(string, string) (*CloudSpace, error)
	Create(*CloudSpaceConfig) (string, error)
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
func (s *CloudSpaceServiceOp) List() (*CloudSpaceList, error) {
	cloudSpaceMap := make(map[string]interface{})
	cloudSpaceMap["includedeleted"] = false
	cloudSpaceJSON, err := json.Marshal(cloudSpaceMap)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/cloudspaces/list", bytes.NewBuffer(cloudSpaceJSON))
	if err != nil {
		return nil, err
	}
	body, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	cloudSpaces := new(CloudSpaceList)
	err = json.Unmarshal(body, &cloudSpaces)
	if err != nil {
		return nil, err
	}

	return cloudSpaces, nil
}

// Get individual CloudSpace
func (s *CloudSpaceServiceOp) Get(cloudSpaceID string) (*CloudSpace, error) {
	cloudSpaceIDMap := make(map[string]interface{})

	cloudSpaceIDInt, err := strconv.Atoi(cloudSpaceID)
	if err != nil {
		return nil, err
	}
	cloudSpaceIDMap["cloudspaceId"] = cloudSpaceIDInt
	cloudSpaceIDJson, err := json.Marshal(cloudSpaceIDMap)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/cloudspaces/get", bytes.NewBuffer(cloudSpaceIDJson))
	if err != nil {
		return nil, err
	}
	body, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	cloudSpace := new(CloudSpace)
	err = json.Unmarshal(body, &cloudSpace)
	if err != nil {
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
			cid := strconv.Itoa(cp.ID)
			return s.client.CloudSpaces.Get(cid)
		}
	}

	return nil, errors.New("Could not find cloudspace based on name")
}

// Create a new CloudSpace
func (s *CloudSpaceServiceOp) Create(cloudSpaceConfig *CloudSpaceConfig) (string, error) {
	cloudSpaceJSON, err := json.Marshal(*cloudSpaceConfig)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/cloudspaces/create", bytes.NewBuffer(cloudSpaceJSON))
	if err != nil {
		return "", err
	}
	body, err := s.client.Do(req)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// Delete a CloudSpace
func (s *CloudSpaceServiceOp) Delete(cloudSpaceConfig *CloudSpaceDeleteConfig) error {
	cloudSpaceJSON, err := json.Marshal(*cloudSpaceConfig)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/cloudspaces/delete", bytes.NewBuffer(cloudSpaceJSON))
	if err != nil {
		return err
	}
	_, err = s.client.Do(req)

	return err
}

// Update an existing CloudSpace
func (s *CloudSpaceServiceOp) Update(cloudSpaceConfig *CloudSpaceConfig) error {
	cloudSpaceJSON, err := json.Marshal(*cloudSpaceConfig)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/cloudspaces/update", bytes.NewBuffer(cloudSpaceJSON))
	if err != nil {
		return err
	}
	_, err = s.client.Do(req)

	return err
}

// SetDefaultGateway sets default gateway of the cloudspace to the given IP address
func (s *CloudSpaceServiceOp) SetDefaultGateway(cloudspaceID int, gateway string) error {
	csMap := make(map[string]interface{})
	csMap["cloudspaceId"] = cloudspaceID
	csMap["gateway"] = gateway

	csMapJSON, err := json.Marshal(csMap)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/cloudspaces/setDefaultGateway", bytes.NewBuffer(csMapJSON))
	if err != nil {
		return err
	}
	_, err = s.client.Do(req)

	return err
}
