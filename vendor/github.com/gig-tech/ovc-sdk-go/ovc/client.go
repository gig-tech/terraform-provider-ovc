package ovc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	// ErrAuthentication represents an authentication error from the server 401
	ErrAuthentication = errors.New("OVC authentication error")
)

// Config used to connect to the API
type Config struct {
	URL          string
	ClientID     string
	ClientSecret string
	JWT          string
	Verbose      bool
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

	logger *logrus.Entry

	Machines         MachineService
	CloudSpaces      CloudSpaceService
	Accounts         AccountService
	Disks            DiskService
	Portforwards     ForwardingService
	Templates        TemplateService
	Sizes            SizesService
	Images           ImageService
	Ipsec            IpsecService
	ExternalNetworks ExternalNetworkService
	Locations        LocationService
}

// NewClient returns a OpenVCloud API Client
func NewClient(c *Config) (*Client, error) {
	if c.ClientID != "" && c.ClientSecret != "" && c.JWT != "" {
		return nil, fmt.Errorf("ClientID, ClientSecret and JWT are provided, please only set ClientID and ClientSecret or JWT")
	}

	var err error
	client := &Client{}
	tokenString := ""

	log := logrus.New()
	if c.Verbose {
		log.SetLevel(logrus.DebugLevel)
	} else {
		log.SetLevel(logrus.InfoLevel)
	}
	logEntry := log.WithField("source", "OpenvCloud client")

	if c.JWT == "" {
		if c.ClientID == "" && c.ClientSecret == "" {
			return nil, fmt.Errorf("no credentials were provided")
		}

		tokenString, err = jwtFromIYO(c)
		if err != nil {
			return nil, err
		}
	} else {
		tokenString = c.JWT
	}
	jwt, err := NewJWT(tokenString, "IYO", logEntry)
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

	client.ServerURL = c.URL + "/restmachine"
	client.JWT = jwt
	client.Access = username.(string) + "@itsyouonline"

	client.logger = logEntry

	client.Machines = &MachineServiceOp{client: client}
	client.CloudSpaces = &CloudSpaceServiceOp{client: client}
	client.Accounts = &AccountServiceOp{client: client}
	client.Disks = &DiskServiceOp{client: client}
	client.Portforwards = &ForwardingServiceOp{client: client}
	client.Templates = &TemplateServiceOp{client: client}
	client.Sizes = &SizesServiceOp{client: client}
	client.Images = &ImageServiceOp{client: client}
	client.Ipsec = &IpsecServiceOp{client: client}
	client.ExternalNetworks = &ExternalNetworkServiceOp{client: client}
	client.Locations = &LocationServiceOp{client: client}

	return client, nil
}

// async adds "async=true" flag to all API calls
func async(req *http.Request) (*http.Request, error) {
	// fetch request body to the string
	jsonMap := make(map[string]interface{})

	if req.Body != nil {
		reqBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(reqBody, &jsonMap)
		if err != nil {
			return nil, err
		}
	}

	jsonMap["_async"] = true
	configJSON, err := json.Marshal(jsonMap)
	if err != nil {
		return nil, err
	}
	newReq, err := http.NewRequest(req.Method, req.URL.String(), bytes.NewBuffer(configJSON))
	if err != nil {
		return nil, err
	}
	return newReq, nil
}

// Do sends and API Request and returns the body as an array of bytes
func (c *Client) Do(req *http.Request) ([]byte, error) {
	req, err := async(req) // make request asynchronous
	if err != nil {
		return nil, err
	}
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	taskID := string(body)

	c.logger.Debug("OVC call: " + req.URL.Path)
	c.logger.Debug("OVC response status code: " + resp.Status)
	c.logger.Debug("OVC response body: " + string(taskID))

	switch {
	case resp.StatusCode == 401:
		return nil, ErrAuthentication
	case resp.StatusCode > 202:
		return body, errors.New(taskID)
	}

	// remove quotes from taskID if contains any
	taskID = strings.Replace(taskID, "\"", "", -1)

	// create request to get result of an async API call by job id
	taskJSON, err := json.Marshal(
		struct {
			TaskID string `json:"taskguid"`
		}{
			TaskID: taskID,
		},
	)
	if err != nil {
		return nil, err
	}

	var taskResp *http.Response
	result := make([]interface{}, 0)
	start, timeout := time.Now(), 10*time.Minute

	// wait for result of the async task
	for {
		taskReq, err := http.NewRequest("POST", c.ServerURL+"/system/task/get", bytes.NewBuffer(taskJSON))
		if err != nil {
			return nil, err
		}
		taskReq.Header.Set("Authorization", fmt.Sprintf("bearer %s", tokenString))
		taskReq.Header.Set("Content-Type", "application/json")
		taskResp, err = client.Do(taskReq)
		if taskResp != nil {
			defer taskResp.Body.Close()
		}
		if err != nil {
			return nil, err
		}
		switch {
		case taskResp.StatusCode == 401:
			return nil, ErrAuthentication
		case taskResp.StatusCode == 404:
			// task may have not been registered yet
			continue
		case taskResp.StatusCode > 202:
			return nil, errors.New(taskID)
		}
		resultBody, err := ioutil.ReadAll(taskResp.Body)
		if err != nil {
			return nil, err
		}
		if len(resultBody) != 0 {
			// if body is not empty, parse result
			err = json.Unmarshal(resultBody, &result)
			if err != nil {
				return resultBody, err
			}
			if len(result) != 0 {
				// result is not empty if can be parsed to a []interface{}
				break
			}
		}
		if now := time.Now(); now.Sub(start) > timeout {
			return nil, fmt.Errorf("job timeout %s", taskID)
		}
		time.Sleep(2 * time.Second)
	}

	success, ok := result[0].(bool)
	if !ok {
		return nil, fmt.Errorf("Task response is incorrect taskId %v \n expected response in form [True/False, taskResult], received: \n %v", string(taskID), result)
	}
	if !success {
		return nil, fmt.Errorf("Task was not successfull taskID: %v:\n %v", string(taskID), result[1])
	}
	finalBody, err := json.Marshal(result[1])
	if err != nil {
		return finalBody, err
	}
	return finalBody, nil
}

// GetLocation parses the URL to return the location of the API
func (c *Client) GetLocation() string {
	u, _ := url.Parse(c.ServerURL)
	hostName := u.Hostname()
	return hostName[:strings.IndexByte(hostName, '.')]
}

// jwtFromIYO fetches a JWT into the itsyouonline platform using the config struct
func jwtFromIYO(c *Config) (string, error) {
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
