package ovc

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var (
	// ErrAuthentication represents an authentication error from the server 401
	ErrAuthentication = errors.New("OVC authentication error")
)

// Config used to connect to the API
type Config struct {
	Hostname     string
	ClientID     string
	ClientSecret string
	JWT          string
}

// Credentials used to authenticate
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Client struct
type Client struct {
	JWT       *JWT
	ServerURL string
	Access    string

	Machines     MachineService
	CloudSpaces  CloudSpaceService
	Accounts     AccountService
	Disks        DiskService
	Portforwards ForwardingService
	Templates    TemplateService
	Sizes        SizesService
	Images       ImageService
}

// Do sends and API Request and returns the body as an array of bytes
func (c *Client) Do(req *http.Request) ([]byte, error) {
	var body []byte
	client := &http.Client{}
	tokenString, err := c.JWT.Get()
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", tokenString))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Println("[DEBUG] OVC response status code: " + resp.Status)
	log.Println("[DEBUG] OVC response body: " + string(body))
	switch {
	case resp.StatusCode == 401:
		return nil, ErrAuthentication
	case resp.StatusCode > 202:
		return body, errors.New(string(body))
	}

	return body, nil
}

// NewClient returns a OpenVCloud API Client
func NewClient(c *Config, url string) (*Client, error) {
	if c.ClientID != "" && c.ClientSecret != "" && c.JWT != "" {
		return nil, fmt.Errorf("ClientID, ClientSecret and JWT are provided, please only set ClientID and ClientSecret or JWT")
	}

	var err error
	client := &Client{}
	tokenString := ""

	if c.JWT == "" {
		tokenString, err = NewLogin(c)
		if err != nil {
			return nil, err
		}
	} else {
		tokenString = c.JWT
	}
	jwt, err := NewJWT(tokenString, "IYO")
	if err != nil {
		return nil, err
	}

	username, err := jwt.Claim("username")
	if err != nil {
		if err == ErrClaimNotPresent {
			return nil, fmt.Errorf("Username not in JWT claims")
		}
		return nil, err
	}

	client.ServerURL = url + "/restmachine"
	client.JWT = jwt
	client.Access = username.(string) + "@itsyouonline"

	client.Machines = &MachineServiceOp{client: client}
	client.CloudSpaces = &CloudSpaceServiceOp{client: client}
	client.Accounts = &AccountServiceOp{client: client}
	client.Disks = &DiskServiceOp{client: client}
	client.Portforwards = &ForwardingServiceOp{client: client}
	client.Templates = &TemplateServiceOp{client: client}
	client.Sizes = &SizesServiceOp{client: client}
	client.Images = &ImageServiceOp{client: client}

	return client, nil
}

// NewLogin logs into the itsyouonline platform using the comfig struct
func NewLogin(c *Config) (string, error) {
	authForm := url.Values{}
	authForm.Add("grant_type", "client_credentials")
	authForm.Add("client_id", c.ClientID)
	authForm.Add("client_secret", c.ClientSecret)
	authForm.Add("response_type", "id_token")
	req, _ := http.NewRequest("POST", "https://itsyou.online/v1/oauth/access_token", strings.NewReader(authForm.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", fmt.Errorf("Error fetching JWT: %s", err)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading JWT request body: %s", err)
	}
	bodyStr := string(bodyBytes)

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Failed to fetch JWT: %s", bodyStr)
	}

	return bodyStr, nil
}

// GetLocation parses the URL to return the location of the API
func (c *Client) GetLocation() string {
	u, _ := url.Parse(c.ServerURL)
	hostName := u.Hostname()
	return hostName[:strings.IndexByte(hostName, '.')]
}
