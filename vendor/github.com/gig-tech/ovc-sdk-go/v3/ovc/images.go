package ovc

import (
	"encoding/json"
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
	Delete(int) error
	DeleteSystemImage(int, string) error
	List(int) (*[]ImageInfo, error)
}

// ImageServiceOp handles communication with the image related methods of the
// OVC API
type ImageServiceOp struct {
	client *Client
}

// Upload uploads an image to the system API
func (s *ImageServiceOp) Upload(imageConfig *ImageConfig) error {
	_, err := s.client.Post("/cloudbroker/image/createImage", *imageConfig, DataActionTimeout)
	return err
}

// Delete deletes an existing image by ID
func (s *ImageServiceOp) Delete(id int) error {
	imageMap := make(map[string]interface{})
	imageMap["imageId"] = id
	imageMap["permanently"] = true

	_, err := s.client.Post("/cloudapi/images/delete", imageMap, OperationalActionTimeout)
	return err
}

// DeleteSystemImage deletes an existing system image by ID
func (s *ImageServiceOp) DeleteSystemImage(id int, reason string) error {
	imageMap := make(map[string]interface{})
	imageMap["imageId"] = id
	imageMap["reason"] = reason
	imageMap["permanently"] = true

	_, err := s.client.Post("/cloudbroker/image/delete", imageMap, OperationalActionTimeout)
	return err
}

// List all system images
func (s *ImageServiceOp) List(accountID int) (*[]ImageInfo, error) {
	accountIDMap := make(map[string]interface{})
	accountIDMap["accountId"] = accountID

	body, err := s.client.Post("/cloudapi/images/list", accountIDMap, ModelActionTimeout)
	if err != nil {
		return nil, err
	}
	images := new([]ImageInfo)
	err = json.Unmarshal(body, &images)
	if err != nil {
		return nil, err
	}

	return images, nil
}
