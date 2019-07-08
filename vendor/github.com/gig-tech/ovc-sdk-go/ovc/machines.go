package ovc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// MachineList is a list of machines
// Returned when using the List method
type MachineList []struct {
	Status      string `json:"status"`
	StackID     int    `json:"stackId"`
	UpdateTime  int    `json:"updateTime"`
	ReferenceID string `json:"referenceId"`
	Name        string `json:"name"`
	Nics        []struct {
		Status      string `json:"status"`
		MacAddress  string `json:"macAddress"`
		ReferenceID string `json:"referenceId"`
		DeviceName  string `json:"deviceName"`
		Type        string `json:"type"`
		Params      string `json:"params"`
		NetworkID   int    `json:"networkId"`
		GUID        string `json:"guid"`
		IPAddress   string `json:"ipAddress"`
	} `json:"nics"`
	SizeID       int   `json:"sizeId"`
	Disks        []int `json:"disks"`
	CreationTime int   `json:"creationTime"`
	ImageID      int   `json:"imageId"`
	Storage      int   `json:"storage"`
	Vcpus        int   `json:"vcpus"`
	Memory       int   `json:"memory"`
	ID           int   `json:"id"`
}

// MachineConfig is used when creating a machine
type MachineConfig struct {
	MachineID    string        `json:"machineId,omitempty"`
	CloudspaceID int           `json:"cloudspaceId,omitempty"`
	Name         string        `json:"name,omitempty"`
	Description  string        `json:"description,omitempty"`
	Memory       int           `json:"memory,omitempty"`
	Vcpus        int           `json:"vcpus,omitempty"`
	SizeID       int           `json:"sizeId,omitempty"`
	ImageID      int           `json:"imageId,omitempty"`
	Disksize     int           `json:"disksize,omitempty"`
	DataDisks    []interface{} `json:"datadisks,omitempty"`
	Permanently  bool          `json:"permanently,omitempty"`
	Userdata     string        `json:"userdata,omitempty"`
}

// MachineInfo contains all information related to a cloudspace
type MachineInfo struct {
	CloudspaceID int    `json:"cloudspaceid"`
	Status       string `json:"status"`
	UpdateTime   int    `json:"updateTime"`
	Hostname     string `json:"hostname"`
	Locked       bool   `json:"locked"`
	Name         string `json:"name"`
	CreationTime int    `json:"creationTime"`
	SizeID       int    `json:"sizeid"`
	Disks        []struct {
		Status  string `json:"status,omitempty"`
		SizeMax int    `json:"sizeMax,omitempty"`
		Name    string `json:"name,omitempty"`
		Descr   string `json:"descr,omitempty"`
		ACL     struct {
		} `json:"acl"`
		Type string `json:"type"`
		ID   int    `json:"id"`
	} `json:"disks"`
	Storage int `json:"storage"`
	ACL     []struct {
		Status       string `json:"status"`
		CanBeDeleted bool   `json:"canBeDeleted"`
		Right        string `json:"right"`
		Type         string `json:"type"`
		UserGroupID  string `json:"userGroupId"`
	} `json:"acl"`
	OsImage  string `json:"osImage"`
	Accounts []struct {
		GUID     string `json:"guid"`
		Login    string `json:"login"`
		Password string `json:"password"`
	} `json:"accounts"`
	Interfaces []struct {
		Status      string `json:"status"`
		MacAddress  string `json:"macAddress"`
		ReferenceID string `json:"referenceId"`
		DeviceName  string `json:"deviceName"`
		IPAddress   string `json:"ipAddress"`
		Params      string `json:"params"`
		NetworkID   int    `json:"networkId"`
		GUID        string `json:"guid"`
		Type        string `json:"type"`
	} `json:"interfaces"`
	ImageID     int         `json:"imageid"`
	ID          int         `json:"id"`
	Memory      int         `json:"memory"`
	Vcpus       int         `json:"vcpus"`
	Description interface{} `json:"description"`
}

// MachineService is an interface for interfacing with the Machine
// endpoints of the OVC API
type MachineService interface {
	List(int) (*MachineList, error)
	Get(string) (*MachineInfo, error)
	GetByName(string, string) (*MachineInfo, error)
	GetByReferenceID(string) (*MachineInfo, error)
	Create(*MachineConfig) (string, error)
	Update(*MachineConfig) (string, error)
	Resize(*MachineConfig) (string, error)
	Delete(*MachineConfig) error
	DeleteByID(int) error
	Template(int, string) error
	Shutdown(int) error
	AddExternalIP(int, int) error
	DeleteExternalIP(int, int, string) error
}

// MachineServiceOp handles communication with the machine related methods of the
// OVC API
type MachineServiceOp struct {
	client *Client
}

// List all machines
func (s *MachineServiceOp) List(cloudSpaceID int) (*MachineList, error) {
	cloudSpaceIDMap := make(map[string]interface{})
	cloudSpaceIDMap["cloudspaceId"] = cloudSpaceID
	cloudSpaceIDJSON, err := json.Marshal(cloudSpaceIDMap)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/machines/list", bytes.NewBuffer(cloudSpaceIDJSON))
	if err != nil {
		return nil, err
	}
	body, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	machines := new(MachineList)
	err = json.Unmarshal(body, &machines)
	if err != nil {
		return nil, err
	}

	return machines, nil
}

// Get individual machine
func (s *MachineServiceOp) Get(id string) (*MachineInfo, error) {
	machineIDMap := make(map[string]interface{})
	machineIDMap["machineId"], _ = strconv.Atoi(id)
	machineIDJson, err := json.Marshal(machineIDMap)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/machines/get", bytes.NewBuffer(machineIDJson))
	if err != nil {
		return nil, err
	}
	body, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	machineInfo := new(MachineInfo)
	err = json.Unmarshal(body, &machineInfo)
	if err != nil {
		return nil, err
	}

	return machineInfo, nil
}

// GetByName gets an individual machine from its name
func (s *MachineServiceOp) GetByName(name string, cloudspaceID string) (*MachineInfo, error) {
	cid, err := strconv.Atoi(cloudspaceID)
	if err != nil {
		return nil, err
	}
	machines, err := s.client.Machines.List(cid)
	if err != nil {
		return nil, err
	}
	for _, mc := range *machines {
		if mc.Name == name {
			return s.client.Machines.Get(strconv.Itoa(mc.ID))
		}
	}

	return nil, fmt.Errorf("Machine %s not found", name)
}

// GetByReferenceID gets an individual machine from its reference ID
func (s *MachineServiceOp) GetByReferenceID(referenceID string) (*MachineInfo, error) {
	referenceIDMap := make(map[string]interface{})
	referenceIDMap["referenceId"] = referenceID
	referenceIDJson, err := json.Marshal(referenceIDMap)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/machines/getByReferenceId", bytes.NewBuffer(referenceIDJson))
	if err != nil {
		return nil, err
	}
	body, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	return s.client.Machines.Get(string(body))
}

// Create a new machine
func (s *MachineServiceOp) Create(machineConfig *MachineConfig) (string, error) {
	machineJSON, err := json.Marshal(*machineConfig)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/machines/create", bytes.NewBuffer(machineJSON))
	if err != nil {
		return "", err
	}
	body, err := s.client.Do(req)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// Update an existing machine
func (s *MachineServiceOp) Update(machineConfig *MachineConfig) (string, error) {
	machineJSON, err := json.Marshal(*machineConfig)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/machines/update", bytes.NewBuffer(machineJSON))
	if err != nil {
		return "", err
	}
	body, err := s.client.Do(req)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// Resize an existing machine
func (s *MachineServiceOp) Resize(machineConfig *MachineConfig) (string, error) {
	machineJSON, err := json.Marshal(*machineConfig)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/machines/resize", bytes.NewBuffer(machineJSON))
	if err != nil {
		return "", err
	}
	body, err := s.client.Do(req)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// Delete an existing machine
func (s *MachineServiceOp) Delete(machineConfig *MachineConfig) error {
	machineJSON, err := json.Marshal(*machineConfig)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/machines/delete", bytes.NewBuffer(machineJSON))
	if err != nil {
		return err
	}
	_, err = s.client.Do(req)

	return err
}

// DeleteByID deletes an existing machine by ID
func (s *MachineServiceOp) DeleteByID(machineID int) error {
	machineMap := make(map[string]interface{})
	machineMap["machineId"] = machineID
	machineMap["permanently"] = true
	machineJSON, err := json.Marshal(machineMap)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/machines/delete", bytes.NewBuffer(machineJSON))
	if err != nil {
		return err
	}
	_, err = s.client.Do(req)

	return err
}

// Template creates an image of the existing machine by ID
func (s *MachineServiceOp) Template(machineID int, templateName string) error {
	machineMap := make(map[string]interface{})
	machineMap["machineId"] = machineID
	machineMap["templateName"] = templateName
	machineJSON, err := json.Marshal(machineMap)
	if err != nil {
		return nil
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/machines/createTemplate", bytes.NewBuffer(machineJSON))
	if err != nil {
		return err
	}
	_, err = s.client.Do(req)

	return err
}

// Shutdown shuts a machine down
func (s *MachineServiceOp) Shutdown(machineID int) error {
	machineMap := make(map[string]interface{})
	machineMap["machineId"] = machineID
	machineMap["force"] = false
	machineJSON, err := json.Marshal(machineMap)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/machines/stop", bytes.NewBuffer(machineJSON))
	if err != nil {
		return err
	}
	_, err = s.client.Do(req)

	return err
}

// AddExternalIP adds external IP
func (s *MachineServiceOp) AddExternalIP(machineID int, externalNetworkID int) error {
	machineMap := make(map[string]interface{})
	machineMap["machineId"] = machineID
	if externalNetworkID != 0 {
		machineMap["externalNetworkId"] = externalNetworkID
	}
	machineJSON, err := json.Marshal(machineMap)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/machines/attachExternalNetwork", bytes.NewBuffer(machineJSON))
	if err != nil {
		return err
	}
	_, err = s.client.Do(req)

	return err
}

// DeleteExternalIP removes external IP
func (s *MachineServiceOp) DeleteExternalIP(machineID int, externalNetworkID int, externalNetworkIP string) error {
	machineMap := make(map[string]interface{})
	machineMap["machineId"] = machineID
	if externalNetworkID > 0 {
		machineMap["externalNetworkId"] = externalNetworkID
		if len(externalNetworkIP) > 0 {
			machineMap["externalnetworkip"] = externalNetworkIP
		}
	}
	machineJSON, err := json.Marshal(machineMap)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/machines/detachExternalNetwork", bytes.NewBuffer(machineJSON))
	if err != nil {
		return err
	}
	_, err = s.client.Do(req)

	return err
}
