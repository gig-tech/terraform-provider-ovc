package ovc

import (
	"encoding/json"
	"net/http"
)

// LocationService represents Location service interface
type LocationService interface {
	List() (*LocationList, error)
}

// LocationServiceOp handles communication with the location related methods of the
// OVC API
type LocationServiceOp struct {
	client *Client
}

// LocationInfo represents the information of the location
type LocationInfo struct {
	Name   string `json:"name"`
	ID     int    `json:"id"`
	GUID   int    `json:"guid"`
	GridID int    `json:"gid"`
	Code   string `json:"locationCode"`
	Flag   string `json:"flag"`
}

// LocationList represents a list of location info
type LocationList []LocationInfo

// List lists all locations of the G8
func (s *LocationServiceOp) List() (*LocationList, error) {
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/locations/list", nil)
	if err != nil {
		return nil, err
	}
	body, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	locations := new(LocationList)
	err = json.Unmarshal(body, &locations)
	if err != nil {
		return nil, err
	}

	return locations, nil
}
