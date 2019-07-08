package ovc

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// ImageConfig is used when uploading an image
type ImageConfig struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	GridID    int    `json:"gid"`
	BootType  string `json:"boottype"`
	Type      string `json:"imagetype"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	AccountID int    `json:"accountId"`
}

// ImageList is a list of images
// Returned when using the List method
type ImageList []ImageInfo

// ImageInfo contains information about the image returned by API
type ImageInfo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Size        int    `json:"size"`
	Status      string `json:"status"`
	Type        string `json:"type"`
	AccountID   int    `json:"accountId"`
	Username    string `json:"username"`
}

// ImageService is an interface for interfacing with the images
// of the OVC API
type ImageService interface {
	Upload(*ImageConfig) error
	DeleteByID(int) error
	DeleteSystemImageByID(int, string) error
	List(int) (*ImageList, error)
}

// ImageServiceOp handles communication with the image related methods of the
// OVC API
type ImageServiceOp struct {
	client *Client
}

// Upload uploads an image to the system API
func (s *ImageServiceOp) Upload(imageConfig *ImageConfig) error {
	imageJSON, err := json.Marshal(*imageConfig)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudbroker/image/createImage", bytes.NewBuffer(imageJSON))
	if err != nil {
		return err
	}
	_, err = s.client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

// DeleteByID deletes an existing image by ID
func (s *ImageServiceOp) DeleteByID(imageID int) error {
	imageMap := make(map[string]interface{})
	imageMap["imageId"] = imageID
	imageMap["permanently"] = true
	imageJSON, err := json.Marshal(imageMap)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/images/delete", bytes.NewBuffer(imageJSON))
	if err != nil {
		return err
	}
	_, err = s.client.Do(req)

	return err
}

// DeleteSystemImageByID deletes an existing system image by ID
func (s *ImageServiceOp) DeleteSystemImageByID(imageID int, reason string) error {
	imageMap := make(map[string]interface{})
	imageMap["imageId"] = imageID
	imageMap["reason"] = reason
	imageMap["permanently"] = true
	imageJSON, err := json.Marshal(imageMap)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudbroker/image/delete", bytes.NewBuffer(imageJSON))
	if err != nil {
		return err
	}
	_, err = s.client.Do(req)

	return err
}

// List all system images
func (s *ImageServiceOp) List(accountID int) (*ImageList, error) {
	accountIDMap := make(map[string]interface{})
	accountIDMap["accountId"] = accountID
	accountIDJson, err := json.Marshal(accountIDMap)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/images/list", bytes.NewBuffer(accountIDJson))
	if err != nil {
		return nil, err
	}
	body, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	images := new(ImageList)
	err = json.Unmarshal(body, &images)
	if err != nil {
		return nil, err
	}

	return images, nil
}
