package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

const baseUrl = "https://api.pinata.cloud/psa"

//IpfsClient is a client for the IPFS API
type IpfsClient struct {
	BaseUrl    string
	apiKey     string
	apiSecret  string
	HttpClient *http.Client
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
		BaseUrl:   baseUrl,
		apiKey:    apiKey,
		apiSecret: apiSecret,
		HttpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

//sendRequest sends a request to the IPFS API
func (c *IpfsClient) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("pinata_api_key", os.Getenv("PINATA_API_KEY"))
	req.Header.Set("pinata_secret_api_key", os.Getenv("PINATA_SECRET_API_KEY"))

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return err
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
	url := c.BaseUrl + "/pinning/pingJSONToIPFS"
	jsonData, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest("Post", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	res := Pin{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
