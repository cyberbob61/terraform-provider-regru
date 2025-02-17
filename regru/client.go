package regru

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
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

func (c Client) doRequest(request interface{}, path ...string) (*APIResponse, error) {
	endpoint := c.baseURL.JoinPath(path...)

	inputData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	fmt.Println("Request JSON:", string(inputData))

	var requestData map[string]any
	if err := json.Unmarshal(inputData, &requestData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request data: %w", err)
	}

	formData := url.Values{}
	for key, value := range requestData {
		switch v := value.(type) {
		case []any:
			for _, item := range v {
				if m, ok := item.(map[string]any); ok {
					if dname, exists := m["dname"]; exists {
						formData.Add("domain_name", fmt.Sprintf("%v", dname))
					}
				} else {
					formData.Add(key, fmt.Sprintf("%v", item))
				}
			}
		case map[string]any:
			if dname, exists := v["dname"]; exists {
				formData.Add("domain_name", fmt.Sprintf("%v", dname))
			} else {
				for subKey, subValue := range v {
					formData.Add(fmt.Sprintf("%s.%s", key, subKey), fmt.Sprintf("%v", subValue))
				}
			}
		default:
			formData.Add(key, fmt.Sprintf("%v", v))
		}
	}

	fmt.Println("Form Data:", formData.Encode())

	req, err := http.NewRequest(http.MethodPost, endpoint.String(), strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	fmt.Println("Request Headers:", req.Header)

	httpClient := c.HTTPClient
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(raw))

	if resp.StatusCode/100 != 2 {
		return nil, parseError(req, resp)
	}

	var apiResp APIResponse
	if err := json.Unmarshal(raw, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal API response: %w", err)
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
