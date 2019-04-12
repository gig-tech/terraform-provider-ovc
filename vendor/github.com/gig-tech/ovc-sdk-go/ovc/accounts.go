package ovc

import (
	"encoding/json"
	"errors"
	"net/http"
)

// Account contains
type Account struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// AccountList is a list of accounts
// Returned when using the List method
type AccountList []struct {
	ID           int    `json:"id"`
	UpdateTime   int    `json:"updateTime"`
	CreationTime int    `json:"creationTime"`
	Name         string `json:"name"`
	ACL          []struct {
		Status      string `json:"status"`
		Right       string `json:"right"`
		Explicit    bool   `json:"explicit"`
		UserGroupID string `json:"userGroupId"`
		GUID        string `json:"guid"`
		Type        string `json:"type"`
	} `json:"acl"`
}

// AccountService is an interface for interfacing with the Account
// endpoints of the OVC API
// See: https://ch-lug-dc01-001.gig.tech/g8vdc/#/ApiDocs
type AccountService interface {
	GetIDByName(string) (int, error)
	List() (*AccountList, error)
}

var _ AccountService = &AccountServiceOp{}

// AccountServiceOp handles communication with the account related methods of the
// OVC API
type AccountServiceOp struct {
	client *Client
}

// GetIDByName returns the account ID based on the account name
func (s *AccountServiceOp) GetIDByName(account string) (int, error) {
	var accounts, err = s.List()
	if err != nil {
		return 0, err
	}
	for _, acc := range *accounts {
		if acc.Name == account {
			return acc.ID, nil
		}
	}
	return -1, errors.New("Account not found")
}

// List all accounts
func (s *AccountServiceOp) List() (*AccountList, error) {
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/accounts/list", nil)
	if err != nil {
		return nil, err
	}
	body, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	var accounts = new(AccountList)
	err = json.Unmarshal(body, &accounts)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}
