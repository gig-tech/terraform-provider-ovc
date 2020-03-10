package ovc

import (
	"encoding/json"
	"errors"
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
type AccountService interface {
	GetIDByName(string) (int, error)
	List() (*AccountList, error)
}

// AccountServiceOp handles communication with the account related methods of the
// OVC API
type AccountServiceOp struct {
	client *Client
}

// GetIDByName returns the account ID based on the account name
func (s *AccountServiceOp) GetIDByName(account string) (int, error) {
	accounts, err := s.List()
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
	body, err := s.client.PostRaw("/cloudapi/accounts/list", nil, ModelActionTimeout)
	if err != nil {
		return nil, err
	}

	accounts := new(AccountList)
	err = json.Unmarshal(body, &accounts)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}
