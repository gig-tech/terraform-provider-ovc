package ovc

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/dgrijalva/jwt-go"
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
	JWT       string
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
	req.Header.Set("Authorization", "bearer "+c.JWT)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Println("Status code: " + resp.Status)
	log.Println("Body: " + string(body))
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

	client := &Client{}
	tokenString := ""
	claims := jwt.MapClaims{}

	if c.JWT == "" {
		tokenString = NewLogin(c)
	} else {
		tokenString = c.JWT
	}
	jwt.ParseWithClaims(tokenString, claims, nil)
	username, ok := claims["username"]
	if !ok {
		return nil, fmt.Errorf("JWT does not contain a username claim")
	}

	client.ServerURL = url + "/restmachine"
	client.JWT = tokenString
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
func NewLogin(c *Config) string {

	authForm := url.Values{}
	authForm.Add("grant_type", "client_credentials")
	authForm.Add("client_id", c.ClientID)
	authForm.Add("client_secret", c.ClientSecret)
	authForm.Add("response_type", "id_token")
	req, _ := http.NewRequest("POST", "https://itsyou.online/v1/oauth/access_token", strings.NewReader(authForm.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error performing login request")
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body")
	}
	jwt := string(bodyBytes)
	defer resp.Body.Close()
	return jwt
}

// GetLocation parses the URL to return the location of the API
func (c *Client) GetLocation() string {
	u, _ := url.Parse(c.ServerURL)
	hostName := u.Hostname()
	return hostName[:strings.IndexByte(hostName, '.')]
}
