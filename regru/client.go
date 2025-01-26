package regru

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	username string
	password string

	baseURL    *url.URL
	HTTPClient *http.Client
}

func loadTLSConfig(certFile, keyFile string) (*tls.Config, error) {
	// Load client cert
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load client certificate: %w", err)
	}

	// Setup TLS config
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true, // Disable server certificate verification
	}

	return tlsConfig, nil
}

func NewClient(username, password, apiEndpoint, certFile, keyFile string) (*Client, error) {
	baseURL, err := url.Parse(apiEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse API endpoint: %w", err)
	}

	tlsConfig, err := loadTLSConfig(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS config: %w", err)
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := &Client{
		username:   username,
		password:   password,
		baseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: 5 * time.Second, Transport: transport},
	}

	return client, nil
}

func (c Client) doRequest(request any, path ...string) (*APIResponse, error) {
	endpoint := c.baseURL.JoinPath(path...)

	inputData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to create input data: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint.String(), bytes.NewReader(inputData))
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %w", err)
	}

	//req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	log.Printf("Request: ----start-----")
	log.Printf("Method: %s", http.MethodPost)
	log.Printf("Endpoint: %s", endpoint.String())
	log.Printf("Header: %s", req.Header)
	log.Printf("Input: %s", string(inputData))
	log.Printf("Request: ---end------")

	httpClient := c.HTTPClient
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("Response Body: %s", string(raw))

	if resp.StatusCode/100 != 2 {
		return nil, parseError(req, resp)
	}

	var apiResp APIResponse
	err = json.Unmarshal(raw, &apiResp)
	if err != nil {
		return nil, err
	}

	return &apiResp, nil
}

func parseError(_ *http.Request, resp *http.Response) error {
	raw, _ := io.ReadAll(resp.Body)

	var errAPI APIResponse
	err := json.Unmarshal(raw, &errAPI)
	if err != nil {
		return err
	}

	return fmt.Errorf("status code: %d, %w", resp.StatusCode, errAPI)
}
