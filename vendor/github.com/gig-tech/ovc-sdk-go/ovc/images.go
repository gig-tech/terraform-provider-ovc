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
	GID       int    `json:"gid"`
	BootType  string `json:"boottype"`
	Type      string `json:"imagetype"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	AccountID int    `json:"accountId"`
}

// ImageService is an interface for interfacing with the images
// of the OVC API
// See: https://ch-lug-dc01-001.gig.tech/system/
type ImageService interface {
	Upload(*ImageConfig) error
	DeleteByID(int) error
	DeleteSystemImageByID(int, string) error
}

// ImageServiceOp handles communication with the image related methods of the
// OVC API
type ImageServiceOp struct {
	client *Client
}

var _ ImageService = &ImageServiceOp{}

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
	if err != nil {
		return err
	}
	return nil
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
	if err != nil {
		return err
	}
	return nil
}
