package shared

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

const baseURL = "https://api.pinata.cloud"

//IpfsClient is a client for the IPFS API
type IpfsClient struct {
	BaseURL    string
	apiKey     string
	apiSecret  string
	HTTPClient *http.Client
}

type ipfsErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ipfsSuccessResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

//Pin struct represents a pin
type Pin struct {
	IpfsHash    string    `json:"ipfsHash"`
	PinSize     int       `json:"pinSize"`
	Timestamp   time.Time `json:"timestamp"`
	IsDuplicate bool      `json:"isDuplicate"`
}

//NewIpfsClient creates a new IPFS client
func NewIpfsClient(apiKey string, apiSecret string) *IpfsClient {
	return &IpfsClient{
		BaseURL:   baseURL,
		apiKey:    apiKey,
		apiSecret: apiSecret,
		HTTPClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

//SendRequest sends a request to the IPFS API
func (c *IpfsClient) SendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("pinata_api_key", os.Getenv("PINATA_API_KEY"))
	req.Header.Set("pinata_secret_api_key", os.Getenv("PINATA_SECRET_API_KEY"))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error sending round trip http request: %v", err)
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes ipfsErrorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}
		return fmt.Errorf("unkown error, status code %d", res.StatusCode)
	}

	fullResponse := ipfsSuccessResponse{
		Data: v,
	}

	if err = json.NewDecoder(res.Body).Decode(&fullResponse); err != nil {
		return err
	}

	return nil
}

//PinJSON pin a json object to IPFS
func (c *IpfsClient) PinJSON(data interface{}) (*Pin, error) {
	url := c.BaseURL + "/pinning/pinJSONToIPFS"
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling data: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}

	res := Pin{}

	err = c.SendRequest(req, &res)
	if err != nil {
		return nil, fmt.Errorf("Error sending request: %v", err)
	}
	return &res, nil
}
