package ovc

import (
	"encoding/json"
	"errors"
)

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
type SizesService interface {
	List(int) (*[]Size, error)
	GetByVcpusAndMemory(int, int, int) (*Size, error)
}

// SizesServiceOp handles communication with the size related methods of the
// OVC API
type SizesServiceOp struct {
	client *Client
}

// List all sizes
func (s *SizesServiceOp) List(cloudspaceID int) (*[]Size, error) {
	sizesMap := make(map[string]interface{})
	sizesMap["cloudspaceId"] = cloudspaceID

	body, err := s.client.Post("/cloudapi/sizes/list", sizesMap, ModelActionTimeout)
	if err != nil {
		return nil, err
	}
	sizes := new([]Size)
	err = json.Unmarshal(body, &sizes)
	if err != nil {
		return nil, err
	}

	return sizes, nil
}

// GetByVcpusAndMemory gets sizes by vcpus and memory
func (s *SizesServiceOp) GetByVcpusAndMemory(vcpus int, memory int, cloudspaceID int) (*Size, error) {
	sizes, err := s.client.Sizes.List(cloudspaceID)
	if err != nil {
		return nil, err
	}
	size := new(Size)
	for _, sz := range *sizes {
		if sz.Vcpus == vcpus && sz.Memory == memory {
			*size = sz
			return size, nil
		}
	}

	return nil, errors.New("Could not find sizes")
}
