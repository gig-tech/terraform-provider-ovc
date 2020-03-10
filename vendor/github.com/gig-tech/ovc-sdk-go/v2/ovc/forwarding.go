package ovc

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// PortForwardingConfig is used when creating a portforward
type PortForwardingConfig struct {
	CloudspaceID     int    `json:"cloudspaceId,omitempty"`
	SourcePublicIP   string `json:"sourcePublicIp,omitempty"`
	SourcePublicPort int    `json:"sourcePublicPort,omitempty"`
	SourceProtocol   string `json:"sourceProtocol,omitempty"`
	PublicIP         string `json:"publicIp,omitempty"`
	PublicPort       int    `json:"publicPort,omitempty"`
	MachineID        int    `json:"machineId,omitempty"`
	LocalPort        int    `json:"localPort,omitempty"`
	Protocol         string `json:"protocol,omitempty"`
	ID               int    `json:"id,omitempty"`
}

// PortForwardingList is a list of portforwards
// Returned when using the List method
type PortForwardingList []struct {
	Protocol    string `json:"protocol"`
	LocalPort   string `json:"localPort"`
	MachineName string `json:"machineName"`
	PublicIP    string `json:"publicIp"`
	LocalIP     string `json:"localIp"`
	MachineID   int    `json:"machineId"`
	PublicPort  string `json:"publicPort"`
	ID          int    `json:"id"`
}

// PortForwardingInfo is returned when using the get method
type PortForwardingInfo struct {
	Protocol    string `json:"protocol"`
	LocalPort   string `json:"localPort"`
	MachineName string `json:"machineName"`
	PublicIP    string `json:"publicIp"`
	LocalIP     string `json:"localIp"`
	MachineID   int    `json:"machineId"`
	PublicPort  string `json:"publicPort"`
	ID          int    `json:"id"`
}

// ForwardingService is an interface for interfacing with the portforwards
// endpoints of the OVC API
type ForwardingService interface {
	Create(*PortForwardingConfig) (int, error)
	List(*PortForwardingConfig) (*PortForwardingList, error)
	Delete(*PortForwardingConfig) error
	DeleteByPort(int, string, int) error
	Update(*PortForwardingConfig) error
	Get(*PortForwardingConfig) (*PortForwardingInfo, error)
}

// ForwardingServiceOp handles communication with the machine related methods of the
// OVC API
type ForwardingServiceOp struct {
	client *Client
}

// Get a portforward based on ID
func (s *ForwardingServiceOp) Get(portForwardingConfig *PortForwardingConfig) (*PortForwardingInfo, error) {
	portForwardingList, err := s.List(portForwardingConfig)
	if err != nil {
		return nil, err
	}
	for _, portforward := range *portForwardingList {
		if portforward.PublicPort == strconv.Itoa(portForwardingConfig.PublicPort) {
			return &PortForwardingInfo{
				Protocol:    portforward.Protocol,
				LocalPort:   portforward.LocalPort,
				MachineName: portforward.MachineName,
				PublicIP:    portforward.PublicIP,
				LocalIP:     portforward.LocalIP,
				MachineID:   portforward.MachineID,
				PublicPort:  portforward.PublicPort,
				ID:          portforward.ID,
			}, nil
		}
	}

	return nil, fmt.Errorf("Could not find a portforward with publicport %v", portForwardingConfig.PublicPort)
}

// Create a new portforward
func (s *ForwardingServiceOp) Create(portForwardingConfig *PortForwardingConfig) (int, error) {
	if portForwardingConfig.PublicPort == 0 {
		portForwardingConfig.PublicPort = s.getRandomPublicPort(portForwardingConfig)
	}

	_, err := s.client.Post("/cloudapi/portforwarding/create", *portForwardingConfig, OperationalActionTimeout)
	if err != nil {
		return 0, err
	}

	return portForwardingConfig.PublicPort, nil
}

// Update an existing portforward
func (s *ForwardingServiceOp) Update(portForwardingConfig *PortForwardingConfig) error {
	_, err := s.client.Post("/cloudapi/portforwarding/updateByPort", *portForwardingConfig, OperationalActionTimeout)
	return err
}

// Delete an existing portforward
func (s *ForwardingServiceOp) Delete(portForwardingConfig *PortForwardingConfig) error {
	_, err := s.client.Post("/cloudapi/portforwarding/deleteByPort", *portForwardingConfig, OperationalActionTimeout)
	return err
}

// List all portforwards
func (s *ForwardingServiceOp) List(portForwardingConfig *PortForwardingConfig) (*PortForwardingList, error) {
	body, err := s.client.Post("/cloudapi/portforwarding/list", *portForwardingConfig, ModelActionTimeout)
	if err != nil {
		return nil, err
	}

	portForwardingList := new(PortForwardingList)
	err = json.Unmarshal(body, &portForwardingList)
	if err != nil {
		return nil, err
	}

	return portForwardingList, nil
}

// DeleteByPort Deletes a portforward by publicIP, public port and cloudspace ID
func (s *ForwardingServiceOp) DeleteByPort(publicPort int, publicIP string, cloudSpaceID int) error {
	pfMap := make(map[string]interface{})
	pfMap["publicIp"] = publicIP
	pfMap["publicPort"] = publicPort
	pfMap["cloudspaceId"] = cloudSpaceID

	_, err := s.client.Post("/cloudapi/portforwarding/deleteByPort", pfMap, OperationalActionTimeout)
	return err
}

func (s *ForwardingServiceOp) getRandomPublicPort(portForwardingConfig *PortForwardingConfig) int {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	randInt := r.Intn(40000) + 2000
	for s.hasPublicPort(portForwardingConfig, randInt) {
		randInt = rand.Intn(40000) + 2000
	}

	return randInt
}

func (s *ForwardingServiceOp) hasPublicPort(portForwardingConfig *PortForwardingConfig, r int) bool {
	config := &PortForwardingConfig{
		CloudspaceID: portForwardingConfig.CloudspaceID,
	}
	list, err := s.List(config)
	if err != nil {
		return false
	}
	for _, port := range *list {
		if port.PublicPort == strconv.Itoa(r) {
			return true
		}
	}

	return false
}
