package ovc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tidwall/limiter"
)

var (
	// ErrAuthentication represents an authentication error from the server 401
	ErrAuthentication = errors.New("OVC authentication error")
	// ErrNotFound represents a resource not found error from the server 404
	ErrNotFound = errors.New("Resource not found")
)

// ResponseTimeout for specific api requests
type ResponseTimeout time.Duration

const (
	// ModelActionTimeout is used for actions that only interfere with the model in the G8
	ModelActionTimeout ResponseTimeout = ResponseTimeout(time.Minute)
	// OperationalActionTimeout is used for actions that tamper with deploy resources
	OperationalActionTimeout ResponseTimeout = ResponseTimeout(time.Minute * 10)
	// DataActionTimeout is used for actions that involve moving data
	DataActionTimeout ResponseTimeout = ResponseTimeout(time.Hour * 24)
)

// Config used to connect to the API
type Config struct {
	URL          string
	ClientID     string
	ClientSecret string
	JWT          string
	// Deprecated: only used if no Logger is passed in.
	// Use an appropriately configured logger instead.
	Verbose bool
	Logger  Logger
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

	logger         Logger
	requestLimiter *limiter.Limiter

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

func setupLogger(c *Config) Logger {
	if c.Logger == nil {
		logger := logrus.New()
		if c.Verbose {
			logger.SetLevel(logrus.DebugLevel)
		} else {
			logger.SetLevel(logrus.InfoLevel)
		}
		return LogrusAdapter{FieldLogger: logger.WithField("source", "OpenvCloud client")}
	}
	return c.Logger
}

// NewClient returns an OpenVCloud API Client
func NewClient(c *Config) (*Client, error) {
	logger := setupLogger(c)

	if c.ClientID != "" && c.ClientSecret != "" && c.JWT != "" {
		return nil, fmt.Errorf("ClientID, ClientSecret and JWT are provided, please only set ClientID and ClientSecret or JWT")
	}

	var err error
	client := &Client{}
	tokenString := ""

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
	jwt, err := NewJWTFromIYO(tokenString, logger)
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

	client.logger = logger

	requestLimitConfiguration, found := os.LookupEnv("G8_API_CONCURRENT_REQUESTS")
	limit := 5
	if found {
		limit, err = strconv.Atoi(requestLimitConfiguration)
		if err != nil {
			return nil, err
		}
	}
	client.requestLimiter = limiter.New(limit)

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
func (c *Client) async(req *http.Request) ([]byte, error) {
	// fetch request body to the string
	jsonMap := make(map[string]interface{})

	if req.Body != nil {
		reqBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			c.logger.Errorf("Failed to read body from http request: %s", err)
			return nil, err
		}
		err = json.Unmarshal(reqBody, &jsonMap)
		if err != nil {
			c.logger.Errorf("req.Method: %s", req.Method)
			c.logger.Errorf("req.URL: %s", req.URL.String())
			c.logger.Errorf("Failed to marshal json body into an object: %s\njson body:\n%s", err, string(reqBody))
			return nil, err
		}
	}

	jsonMap["_async"] = true
	return json.Marshal(jsonMap)
}

func (c *Client) doHTTPRequest(client *http.Client, method string, url string, body io.Reader) (*http.Response, error) {
	defer c.requestLimiter.End()
	c.requestLimiter.Begin()
	asyncReq, err := http.NewRequest(method, url, body)
	if err != nil {
		c.logger.Errorf("Failed to create async request: %s", err)
		return nil, err
	}
	tokenString, err := c.JWT.Get()
	if err != nil {
		c.logger.Errorf("Could not make JWT: %s", err)
		return nil, err
	}
	asyncReq.Header.Set("Authorization", fmt.Sprintf("bearer %s", tokenString))
	asyncReq.Header.Set("Content-Type", "application/json")
	return client.Do(asyncReq)
}

// Do sends and API Request and returns the body as an array of bytes
func (c *Client) do(req *http.Request, timeout ResponseTimeout) ([]byte, error) {
	var requestTimeoutMultiplier int = 0
	var requestErrorCount int = 0
	var taskID string
	client := &http.Client{}
	asyncBody, err := c.async(req)
	if err != nil {
		c.logger.Errorf("Failed to make request body async")
		return nil, err
	}
	// Try to issue request, but retry if it would fail due 2 2 many concurrent requests
	for {
		resp, err := c.doHTTPRequest(client, req.Method, req.URL.String(), bytes.NewBuffer(asyncBody))
		if resp != nil {
			defer resp.Body.Close()
		}
		if err != nil {
			c.logger.Errorf("Error doing G8 Api request: %s", err)
			if requestErrorCount < 20 {
				time.Sleep(time.Duration(requestTimeoutMultiplier) * time.Second)
				requestErrorCount++
				continue
			} else {
				c.logger.Errorf("Could not do G8 Api request: %s", err)
				return nil, err
			}
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.logger.Errorf("Could not read response body: %s", err)
			return nil, err
		}
		taskID = string(body)

		c.logger.Debugf("OVC call: %s", req.URL.Path)
		c.logger.Debugf("OVC response status code: %s", resp.StatusCode)
		c.logger.Debugf("OVC response status: %s", resp.Status)
		c.logger.Debugf("OVC response body: %s", string(taskID))

		switch {
		case resp.StatusCode == http.StatusBadRequest && taskID == "<html>\r\n<head><title>400 Bad Request</title></head>\r\n<body>\r\n<center><h1>400 Bad Request</h1></center>\r\n<hr><center>nginx/1.17.6</center>\r\n</body>\r\n</html>\r\n":
			// Sometimes nginx returns 400 for no reason
			requestTimeoutMultiplier++
			time.Sleep(time.Duration(requestTimeoutMultiplier) * time.Second)
			continue
		case resp.StatusCode == http.StatusUnauthorized:
			c.logger.Errorf("Unauthorized: %s", ErrAuthentication)
			return nil, ErrAuthentication
		case resp.StatusCode == http.StatusTooManyRequests:
			requestTimeoutMultiplier++
			time.Sleep(time.Duration(requestTimeoutMultiplier) * time.Second)
			continue
		case resp.StatusCode > http.StatusAccepted:
			c.logger.Errorf("Request failed with error: %s", err)
			return body, errors.New(taskID)
		}
		break
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
		c.logger.Errorf("Could not marshal json body into object: %s", err)
		return nil, err
	}

	var taskResp *http.Response
	result := make([]interface{}, 0)
	start := time.Now()

	// wait for result of the async task
	requestErrorCount = 0
	fourOFourCount := 0
	for {
		if now := time.Now(); now.Sub(start) > time.Duration(timeout) {
			err = fmt.Errorf("job timeout %s", taskID)
			c.logger.Errorf("Task failed to complete within the timeout: %s", err)
			return nil, err
		}

		taskResp, err = c.doHTTPRequest(client, http.MethodPost, c.ServerURL+"/system/task/get", bytes.NewBuffer(taskJSON))
		if taskResp != nil {
			defer taskResp.Body.Close()
		}
		if err != nil {
			c.logger.Errorf("Error getting task result: %s", err)
			if requestErrorCount < 20 {
				time.Sleep(2 * time.Second)
				requestErrorCount++
				continue
			} else {
				c.logger.Errorf("Could not get task result: %s", err)
				return nil, err
			}
		}

		c.logger.Debugf("OVC call: %s", req.URL.Path)
		c.logger.Debugf("OVC task call: %s", taskID)
		c.logger.Debugf("OVC response status code: %s", taskResp.StatusCode)
		c.logger.Debugf("OVC response status: %s", taskResp.Status)
		resultBody, err := ioutil.ReadAll(taskResp.Body)
		if err != nil {
			c.logger.Errorf("Could not read response body: %s", err)
			return nil, err
		}
		c.logger.Debugf("OVC response: %s", string(resultBody))

		switch {
		case taskResp.StatusCode == http.StatusUnauthorized:
			c.logger.Errorf("Unauthorized: %s", ErrAuthentication)
			return nil, ErrAuthentication
		case taskResp.StatusCode == http.StatusNotFound:
			if fourOFourCount == 0 {
				fourOFourCount++
				c.logger.Error("Oops we hit a race condition bug in the API server prior 2.5.6")
				time.Sleep(2 * time.Second)
				continue
			} else {
				c.logger.Errorf("Task not found: %s", ErrNotFound)
				return nil, ErrNotFound
			}
		case taskResp.StatusCode == http.StatusBadRequest:
			// Sometimes nginx returns 400 for no reason
			c.logger.Error("Received 400, probably nginx issue.")
			time.Sleep(2 * time.Second)
			continue
		case taskResp.StatusCode == http.StatusTooManyRequests:
			c.logger.Error("Oops spamming the API to hard. G8 return 429")
			time.Sleep(2 * time.Second)
			continue
		case taskResp.StatusCode > http.StatusTooManyRequests:
			err = errors.New(taskID)
			c.logger.Errorf("Task failed: %s", err)
			return nil, err
		}
		if len(resultBody) != 0 {
			// if body is not empty, parse result
			err = json.Unmarshal(resultBody, &result)
			if err != nil {
				c.logger.Errorf("Could not marshal json body into object: %s", err)
				return resultBody, err
			}
			if len(result) != 0 {
				// result is not empty if can be parsed to a []interface{}
				break
			}
		}
		time.Sleep(2 * time.Second)
	}

	success, ok := result[0].(bool)
	if !ok {
		err = fmt.Errorf("Task response is incorrect taskId %v \n expected response in form [True/False, taskResult], received: \n %v", string(taskID), result)
		c.logger.Errorf("%s", err)
		return nil, err
	}
	if !success {
		err = fmt.Errorf("Task was not successfull taskID: %v:\n %v", string(taskID), result[1])
		c.logger.Errorf("%s", err)
		return nil, err
	}
	finalBody, err := json.Marshal(result[1])
	if err != nil {
		c.logger.Errorf("Could not marshal result object into json: %s", err)
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

// PostRaw POSTs a request with `raw` as data (nil is permitted) to `c.ServerUrl + endpoint`
func (c *Client) PostRaw(endpoint string, raw io.Reader, timeout ResponseTimeout) ([]byte, error) {
	req, err := http.NewRequest("POST", c.ServerURL+endpoint, raw)
	if err != nil {
		return nil, err
	}
	return c.do(req, timeout)
}

// Post marshals `in` to JSON and POSTs a request to `c.ServerUrl + endpoint`
func (c *Client) Post(endpoint string, in interface{}, timeout ResponseTimeout) ([]byte, error) {
	jsonIn, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	return c.PostRaw(endpoint, bytes.NewBuffer(jsonIn), timeout)
}
