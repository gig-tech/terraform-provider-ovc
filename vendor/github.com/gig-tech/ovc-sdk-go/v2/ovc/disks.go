package ovc

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// DiskConfig is used when creating a disk
type DiskConfig struct {
	AccountID   int    `json:"accountId,omitempty"`
	GridID      int    `json:"gid,omitempty"`
	MachineID   int    `json:"machineId,omitempty"`
	DiskName    string `json:"diskName,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Size        int    `json:"size,omitempty"`
	Type        string `json:"type,omitempty"`
	SSDSize     int    `json:"ssdSize,omitempty"`
	IOPS        int    `json:"iops,omitempty"`
	DiskID      int    `json:"diskId,omitempty"`
	Detach      bool   `json:"detach,omitempty"`
	Permanently string `json:"permanently,omitempty"`
}

// DiskDeleteConfig is used when deleting a disk
type DiskDeleteConfig struct {
	DiskID      int  `json:"diskId"`
	Detach      bool `json:"detach"`
	Permanently bool `json:"permanently"`
}

// DiskAttachConfig is used when attatching a disk to a machine
type DiskAttachConfig struct {
	DiskID    int `json:"diskId"`
	MachineID int `json:"machineId"`
}

// DiskInfo contains all information related to a disk
type DiskInfo struct {
	ReferenceID         string        `json:"referenceId"`
	DiskPath            string        `json:"diskPath"`
	Images              []interface{} `json:"images"`
	GUID                int           `json:"guid"`
	ID                  int           `json:"id"`
	PCIBus              int           `json:"pci_bus"`
	PCISlot             int           `json:"pci_slot"`
	AccountID           int           `json:"accountId"`
	SizeUsed            int           `json:"sizeUsed"`
	Descr               string        `json:"descr"`
	GridID              int           `json:"gid"`
	Role                string        `json:"role"`
	Params              string        `json:"params"`
	Type                string        `json:"type"`
	Status              string        `json:"status"`
	RealityDeviceNumber int           `json:"realityDeviceNumber"`
	Passwd              string        `json:"passwd"`
	Iotune              struct {
		TotalIopsSec int `json:"total_iops_sec"`
	} `json:"iotune"`
	Name    string        `json:"name"`
	SizeMax int           `json:"sizeMax"`
	Meta    []interface{} `json:"_meta"`
	ACL     struct {
	} `json:"acl"`
	Iqn           string `json:"iqn"`
	BootPartition int    `json:"bootPartition"`
	Login         string `json:"login"`
	Order         int    `json:"order"`
	Ckey          string `json:"_ckey"`
}

// DiskList is a list of disks
// Returned when using the List method
type DiskList []struct {
	Username    interface{} `json:"username"`
	Status      string      `json:"status"`
	Description string      `json:"description"`
	Name        string      `json:"name"`
	Size        int         `json:"sizeMax"`
	Type        string      `json:"type"`
	ID          int         `json:"id"`
	AccountID   int         `json:"accountId"`
}

// DiskExposeProtocolNBD is a constant representing the storage protocol NBD
// (Network Block Device)
const DiskExposeProtocolNBD = "nbd"

// DiskExposeConfig is used to request that a given disk is exposed via a given
// cloudspace and protocol
type DiskExposeConfig struct {
	Protocol     string `json:"-"`
	DiskID       int    `json:"diskId"`
	CloudSpaceID int    `json:"cloudspaceId"`
	IOPS         int    `json:"iops"`
}

// DiskEndPointDescriptor is an interface that a type representing a storage
// protocol endpoint shall implement to facilitate accessing protocol specific
// information in DiskExposeInfo
type DiskEndPointDescriptor interface {
	Protocol() string
}

// NBDDiskEndPointDescriptor contains NBD specific information on how to
// access a disk exposed as a Network Block Device
type NBDDiskEndPointDescriptor struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
	Name    string `json:"name"`
	User    string `json:"user"`
	Psk     string `json:"psk"`
}

// Protocol implements DiskEndPointDescriptor.Protocol
func (*NBDDiskEndPointDescriptor) Protocol() string {
	return DiskExposeProtocolNBD
}

// DiskExposeInfo contains information on how to access a disk exposed using
// a given storage protocol.
type DiskExposeInfo struct {
	Protocol string                 `json:"protocol"`
	EndPoint DiskEndPointDescriptor `json:"endpoint"`
}

// UnmarshalJSON deserializes DiskInfo from a buffer containing a JSON
// representation of DiskInfo
func (i *DiskExposeInfo) UnmarshalJSON(b []byte) error {
	type ProtoProbe struct {
		Protocol string `json:"protocol"`
	}

	probe := ProtoProbe{}
	err := json.Unmarshal(b, &probe)
	if err != nil {
		return err
	}

	if probe.Protocol != DiskExposeProtocolNBD {
		return fmt.Errorf("Unknown disk expose protocol \"%s\"", probe.Protocol)
	}

	type NBDDesc struct {
		EndPoint NBDDiskEndPointDescriptor `json:"endpoint"`
	}

	i.Protocol = probe.Protocol

	desc := NBDDesc{}
	err = json.Unmarshal(b, &desc)
	if err != nil {
		return err
	}

	i.EndPoint = &desc.EndPoint
	return nil
}

// DiskUnexposeConfig is used to request that a given exposed disk is exposed
type DiskUnexposeConfig struct {
	DiskID int `json:"diskId"`
}

// DiskService is an interface for interfacing with the Disk
// endpoints of the OVC API
type DiskService interface {
	Resize(*DiskConfig) error
	List(int, string) (*DiskList, error)
	Get(string) (*DiskInfo, error)
	GetByName(string, int, string) (*DiskInfo, error)
	Create(*DiskConfig) (string, error)
	CreateAndAttach(*DiskConfig) (string, error)
	Attach(*DiskAttachConfig) error
	Detach(*DiskAttachConfig) error
	Update(*DiskConfig) error
	Delete(*DiskDeleteConfig) error
	Expose(*DiskExposeConfig) (*DiskExposeInfo, error)
	Unexpose(*DiskUnexposeConfig) error
}

// DiskServiceOp handles communication with the disk related methods of the
// OVC API
type DiskServiceOp struct {
	client *Client
}

// List all disks
func (s *DiskServiceOp) List(accountID int, diskType string) (*DiskList, error) {
	diskMap := make(map[string]interface{})
	diskMap["accountId"] = accountID
	if len(diskType) != 0 {
		diskMap["type"] = diskType
	}

	body, err := s.client.Post("/cloudapi/disks/list", diskMap, OperationalActionTimeout)
	if err != nil {
		return nil, err
	}

	disks := new(DiskList)
	err = json.Unmarshal(body, &disks)
	if err != nil {
		return nil, err
	}

	return disks, nil
}

// CreateAndAttach a new Disk and attaches it to a machine
func (s *DiskServiceOp) CreateAndAttach(diskConfig *DiskConfig) (string, error) {
	defer ReleaseLock(diskConfig.MachineID)
	GetLock(diskConfig.MachineID)
	body, err := s.client.Post("/cloudapi/machines/addDisk", *diskConfig, OperationalActionTimeout)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// Create a new Disk
func (s *DiskServiceOp) Create(diskConfig *DiskConfig) (string, error) {
	body, err := s.client.Post("/cloudapi/disks/create", *diskConfig, OperationalActionTimeout)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// Attach attaches an existing disk to a machine
func (s *DiskServiceOp) Attach(diskAttachConfig *DiskAttachConfig) error {
	defer ReleaseLock(diskAttachConfig.MachineID)
	GetLock(diskAttachConfig.MachineID)
	_, err := s.client.Post("/cloudapi/machines/attachDisk", *diskAttachConfig, OperationalActionTimeout)
	return err
}

// Detach detaches an existing disk from a machine
func (s *DiskServiceOp) Detach(diskAttachConfig *DiskAttachConfig) error {
	s.client.logger.Debugf("Detaching disk %d from machine %d.", diskAttachConfig.DiskID, diskAttachConfig.MachineID)
	defer ReleaseLock(diskAttachConfig.MachineID)
	GetLock(diskAttachConfig.MachineID)
	_, err := s.client.Post("/cloudapi/machines/detachDisk", *diskAttachConfig, OperationalActionTimeout)
	if err == nil {
		s.client.logger.Debugf("Detaching disk %d from machine %d completed.", diskAttachConfig.DiskID, diskAttachConfig.MachineID)
	} else {
		s.client.logger.Debugf("Detaching disk %d from machine %d failed.", diskAttachConfig.DiskID, diskAttachConfig.MachineID)
	}
	return err
}

// Update updates an existing disk
func (s *DiskServiceOp) Update(diskConfig *DiskConfig) error {
	switch {
	case diskConfig.Size != 0:
		_, err := s.client.Post("/cloudapi/disks/resize", *diskConfig, OperationalActionTimeout)
		if err != nil {
			return err
		}

		fallthrough

	case diskConfig.IOPS != 0:
		_, err := s.client.Post("/cloudapi/disks/limitIO", *diskConfig, OperationalActionTimeout)
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete an existing Disk
func (s *DiskServiceOp) Delete(diskConfig *DiskDeleteConfig) error {
	_, err := s.client.Post("/cloudapi/disks/delete", *diskConfig, OperationalActionTimeout)
	return err
}

// Get individual Disk
func (s *DiskServiceOp) Get(diskID string) (*DiskInfo, error) {
	diskIDMap := make(map[string]interface{})
	diskIDInt, err := strconv.Atoi(diskID)
	if err != nil {
		return nil, err
	}
	diskIDMap["diskId"] = diskIDInt

	body, err := s.client.Post("/cloudapi/disks/get", diskIDMap, ModelActionTimeout)
	if err != nil {
		return nil, err
	}
	diskInfo := new(DiskInfo)
	err = json.Unmarshal(body, &diskInfo)
	if err != nil {
		return nil, err
	}

	return diskInfo, nil
}

// GetByName gets a disk by its name
func (s *DiskServiceOp) GetByName(name string, accountID int, diskType string) (*DiskInfo, error) {
	disks, err := s.client.Disks.List(accountID, diskType)
	if err != nil {
		return nil, err
	}
	for _, dk := range *disks {
		if dk.Name == name {
			did := strconv.Itoa(dk.ID)
			return s.client.Disks.Get(did)
		}
	}

	return nil, errors.New("Could not find disk based on name")
}

// Resize resizes a disk. Can only increase the size of a disk
func (s *DiskServiceOp) Resize(diskConfig *DiskConfig) error {
	_, err := s.client.Post("/cloudapi/disks/resize", *diskConfig, OperationalActionTimeout)
	return err
}

// Expose a disk using the requested protocol (currently only NBD is supported)
// via the cloudspace specified in the DiskExposeConfig.
func (s *DiskServiceOp) Expose(diskExposeConfig *DiskExposeConfig) (*DiskExposeInfo, error) {
	jsonOut, err := s.client.Post("/cloudapi/disks/expose", *diskExposeConfig, OperationalActionTimeout)
	if err != nil {
		return nil, err
	}

	exposeInfo := &DiskExposeInfo{}
	err = json.Unmarshal(jsonOut, exposeInfo)
	if err != nil {
		return nil, err
	}

	return exposeInfo, nil
}

// Unexpose a previously exposed disk. Unexposing a non-exposed disk returns an
// error.
func (s *DiskServiceOp) Unexpose(diskUnexposeConfig *DiskUnexposeConfig) error {
	_, err := s.client.Post("/cloudapi/disks/unexpose", *diskUnexposeConfig, OperationalActionTimeout)
	return err
}
