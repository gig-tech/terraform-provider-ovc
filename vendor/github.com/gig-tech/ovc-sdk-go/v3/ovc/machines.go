package ovc

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// NIC basic information
type NIC struct {
	Status      string `json:"status"`
	MacAddress  string `json:"macAddress"`
	ReferenceID string `json:"referenceId"`
	DeviceName  string `json:"deviceName"`
	Type        string `json:"type"`
	Params      string `json:"params"`
	NetworkID   int    `json:"networkId"`
	GUID        string `json:"guid"`
	IPAddress   string `json:"ipAddress"`
}

// Machine basic information
type Machine struct {
	Status       string `json:"status"`
	StackID      int    `json:"stackId"`
	UpdateTime   int    `json:"updateTime"`
	ReferenceID  string `json:"referenceId"`
	Name         string `json:"name"`
	Nics         []NIC  `json:"nics"`
	SizeID       int    `json:"sizeId"`
	Disks        []int  `json:"disks"`
	CreationTime int    `json:"creationTime"`
	ImageID      int    `json:"imageId"`
	Storage      int    `json:"storage"`
	Vcpus        int    `json:"vcpus"`
	Memory       int    `json:"memory"`
	ID           int    `json:"id"`
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

// EmptyMachineConfig is used when creating a new "empty" machine.
// A machine is considered "empty" if it's not created from an
// existing image.
type EmptyMachineConfig struct {
	CloudspaceID int    `json:"cloudspaceId,omitempty"`
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	Memory       int    `json:"memory,omitempty"`
	Vcpus        int    `json:"vcpus,omitempty"`
	Disksize     int    `json:"disksize,omitempty"`
	DataDisks    []int  `json:"datadisks,omitempty"`
	Userdata     string `json:"userdata,omitempty"`
}

// ACL basic information
type ACL struct {
	Status       string `json:"status"`
	CanBeDeleted bool   `json:"canBeDeleted"`
	Right        string `json:"right"`
	Type         string `json:"type"`
	UserGroupID  string `json:"userGroupId"`
}

// MachineDisk basic information
type MachineDisk struct {
	Status  string `json:"status,omitempty"`
	SizeMax int    `json:"sizeMax,omitempty"`
	Name    string `json:"name,omitempty"`
	Descr   string `json:"descr,omitempty"`
	Type    string `json:"type"`
	ID      int    `json:"id"`
}

// UserAccount basic information
type UserAccount struct {
	GUID     string `json:"guid"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

// MachineInfo contains all information related to a cloudspace
type MachineInfo struct {
	CloudspaceID int           `json:"cloudspaceid"`
	Status       string        `json:"status"`
	UpdateTime   int           `json:"updateTime"`
	Hostname     string        `json:"hostname"`
	Locked       bool          `json:"locked"`
	Name         string        `json:"name"`
	CreationTime int           `json:"creationTime"`
	SizeID       int           `json:"sizeid"`
	Disks        []MachineDisk `json:"disks"`
	Storage      int           `json:"storage"`
	ACL          []ACL         `json:"acl"`
	OsImage      string        `json:"osImage"`
	Accounts     []UserAccount `json:"accounts"`
	Interfaces   []NIC         `json:"interfaces"`
	ImageID      int           `json:"imageid"`
	ID           int           `json:"id"`
	Memory       int           `json:"memory"`
	Vcpus        int           `json:"vcpus"`
	Description  *string       `json:"description"`
}

// MachineService is an interface for interfacing with the Machine
// endpoints of the OVC API
type MachineService interface {
	List(int) (*[]Machine, error)
	Get(int) (*MachineInfo, error)
	GetByName(string, int) (*MachineInfo, error)
	GetByReferenceID(string) (*MachineInfo, error)
	Create(*MachineConfig) (int, error)
	CreateEmpty(*EmptyMachineConfig) (int, error)
	Update(*MachineConfig) (string, error)
	Resize(*MachineConfig) (string, error)
	Delete(int, bool) error
	CreateImage(int, string) error
	Shutdown(int) error
	AddExternalIP(int, int) error
	DeleteExternalIP(int, int, string) error
	Stop(int, bool) error
	Start(int, int) error
}

// MachineServiceOp handles communication with the machine related methods of the
// OVC API
type MachineServiceOp struct {
	client *Client
}

// List all machines
func (s *MachineServiceOp) List(cloudSpaceID int) (*[]Machine, error) {
	cloudSpaceIDMap := make(map[string]interface{})
	cloudSpaceIDMap["cloudspaceId"] = cloudSpaceID

	body, err := s.client.Post("/cloudapi/machines/list", cloudSpaceIDMap, ModelActionTimeout)
	if err != nil {
		return nil, err
	}

	machines := new([]Machine)
	err = json.Unmarshal(body, &machines)
	if err != nil {
		return nil, err
	}

	return machines, nil
}

// Get individual machine
func (s *MachineServiceOp) Get(id int) (*MachineInfo, error) {
	defer ReleaseLock(id)
	GetLock(id)
	machineIDMap := make(map[string]interface{})
	var err error
	machineIDMap["machineId"] = id

	body, err := s.client.Post("/cloudapi/machines/get", machineIDMap, OperationalActionTimeout)
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
func (s *MachineServiceOp) GetByName(name string, cloudspaceID int) (*MachineInfo, error) {
	machines, err := s.client.Machines.List(cloudspaceID)
	if err != nil {
		return nil, err
	}
	for _, mc := range *machines {
		if mc.Name == name {
			return s.client.Machines.Get(mc.ID)
		}
	}

	return nil, fmt.Errorf("Machine %s not found", name)
}

// GetByReferenceID gets an individual machine from its reference ID
func (s *MachineServiceOp) GetByReferenceID(referenceID string) (*MachineInfo, error) {
	referenceIDMap := make(map[string]interface{})
	referenceIDMap["referenceId"] = referenceID

	body, err := s.client.Post("/cloudapi/machines/getByReferenceId", referenceIDMap, OperationalActionTimeout)
	if err != nil {
		return nil, err
	}

	machineID, err := strconv.Atoi(string(body))
	if err != nil {
		return nil, err
	}

	return s.client.Machines.Get(machineID)
}

// Create a new machine
func (s *MachineServiceOp) Create(machineConfig *MachineConfig) (int, error) {
	body, err := s.client.Post("/cloudapi/machines/create", *machineConfig, OperationalActionTimeout)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(body))
}

// CreateEmpty a new "empty" machine (= not based on an existing image)
func (s *MachineServiceOp) CreateEmpty(emptyMachineConfig *EmptyMachineConfig) (int, error) {
	body, err := s.client.Post("/cloudapi/machines/createEmptyMachine", *emptyMachineConfig, ModelActionTimeout)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(body))
}

// Update an existing machine
func (s *MachineServiceOp) Update(machineConfig *MachineConfig) (string, error) {
	body, err := s.client.Post("/cloudapi/machines/update", *machineConfig, ModelActionTimeout)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// Resize an existing machine
func (s *MachineServiceOp) Resize(machineConfig *MachineConfig) (string, error) {
	body, err := s.client.Post("/cloudapi/machines/resize", *machineConfig, OperationalActionTimeout)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// Delete deletes an existing machine
func (s *MachineServiceOp) Delete(id int, permanently bool) error {
	machineMap := make(map[string]interface{})
	machineMap["machineId"] = id
	machineMap["permanently"] = permanently

	_, err := s.client.Post("/cloudapi/machines/delete", machineMap, OperationalActionTimeout)
	return err
}

// Stop stops a machine
func (s *MachineServiceOp) Stop(id int, force bool) error {
	machineMap := make(map[string]interface{})
	machineMap["machineId"] = id
	machineMap["stop"] = force

	_, err := s.client.Post("/cloudapi/machines/stop", machineMap, OperationalActionTimeout)
	return err
}

// Start starts a machine, boots from ISO if diskID is given
func (s *MachineServiceOp) Start(id int, diskID int) error {
	machineMap := make(map[string]interface{})
	machineMap["machineId"] = id
	if diskID != 0 {
		machineMap["diskId"] = diskID
	}

	_, err := s.client.Post("/cloudapi/machines/start", machineMap, OperationalActionTimeout)
	return err
}

// CreateImage creates an image of the existing machine by ID
func (s *MachineServiceOp) CreateImage(id int, imageName string) error {
	machineMap := make(map[string]interface{})
	machineMap["machineId"] = id
	machineMap["templateName"] = imageName

	_, err := s.client.Post("/cloudapi/machines/createTemplate", machineMap, DataActionTimeout)
	return err
}

// Shutdown shuts a machine down
func (s *MachineServiceOp) Shutdown(id int) error {
	machineMap := make(map[string]interface{})
	machineMap["machineId"] = id
	machineMap["force"] = false

	_, err := s.client.Post("/cloudapi/machines/stop", machineMap, OperationalActionTimeout)
	return err
}

// AddExternalIP adds external IP
func (s *MachineServiceOp) AddExternalIP(id int, externalNetworkID int) error {
	machineMap := make(map[string]interface{})
	machineMap["machineId"] = id
	if externalNetworkID != 0 {
		machineMap["externalNetworkId"] = externalNetworkID
	}
	_, err := s.client.Post("/cloudapi/machines/attachExternalNetwork", machineMap, OperationalActionTimeout)
	return err
}

// DeleteExternalIP removes external IP
func (s *MachineServiceOp) DeleteExternalIP(id int, externalNetworkID int, externalNetworkIP string) error {
	machineMap := make(map[string]interface{})
	machineMap["machineId"] = id
	if externalNetworkID > 0 {
		machineMap["externalNetworkId"] = externalNetworkID
		if len(externalNetworkIP) > 0 {
			machineMap["externalnetworkip"] = externalNetworkIP
		}
	}

	_, err := s.client.Post("/cloudapi/machines/detachExternalNetwork", machineMap, OperationalActionTimeout)
	return err
}
