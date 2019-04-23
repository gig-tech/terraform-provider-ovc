package ovc

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

// SizesList is a list of sizes
// Returned when using the list method
type SizesList []struct {
	ID          int    `json:"id"`
	Vcpus       int    `json:"vcpus"`
	Disks       []int  `json:"disks"`
	Name        string `json:"name"`
	Memory      int    `json:"memory"`
	Description string `json:"description,omitempty"`
}

// Size contains all information related to a disk
type Size struct {
	ID          int    `json:"id"`
	Vcpus       int    `json:"vcpus"`
	Disks       []int  `json:"disks"`
	Name        string `json:"name"`
	Memory      int    `json:"memory"`
	Description string `json:"description,omitempty"`
}

// SizesService is an interface for interfacing with the Sizes
// endpoints of the OVC API
// See: https://ch-lug-dc01-001.gig.tech/g8vdc/#/ApiDocs
type SizesService interface {
	List(string) (*SizesList, error)
	GetByVcpusAndMemory(int, int, string) (*Size, error)
}

// SizesServiceOp handles communication with the size related methods of the
// OVC API
type SizesServiceOp struct {
	client *Client
}

var _ SizesService = &SizesServiceOp{}

// List all sizes
func (s *SizesServiceOp) List(cloudspaceID string) (*SizesList, error) {
	sizesMap := make(map[string]interface{})
	sizesMap["cloudspaceId"] = cloudspaceID
	sizesJSON, err := json.Marshal(sizesMap)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/sizes/list", bytes.NewBuffer(sizesJSON))
	if err != nil {
		return nil, err
	}
	body, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	var sizes = new(SizesList)
	err = json.Unmarshal(body, &sizes)
	if err != nil {
		return nil, err
	}
	return sizes, nil
}

// GetByVcpusAndMemory gets sizes by vcpus and memory
func (s *SizesServiceOp) GetByVcpusAndMemory(vcpus int, memory int, cloudspaceID string) (*Size, error) {
	sizes, err := s.client.Sizes.List(cloudspaceID)
	if err != nil {
		return nil, err
	}
	var size = new(Size)
	for _, sz := range *sizes {
		if sz.Vcpus == vcpus && sz.Memory == memory {
			*size = sz
			return size, nil
		}
	}
	return nil, errors.New("Could not find sizes")
}
