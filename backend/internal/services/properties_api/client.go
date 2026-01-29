package properties_api

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Client handles communication with the Properties API
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new Properties API client
func NewClient() *Client {
	baseURL := os.Getenv("PROPERTIES_API_URL")
	if baseURL == "" {
		// Default for local development (Docker Desktop)
		baseURL = "http://host.docker.internal:8003"
	}

	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// QueryProperties queries the properties API with the given request
func (c *Client) QueryProperties(req *QueryRequest) (*PropertiesResponse, error) {
	url := fmt.Sprintf("%s/api/v1/properties", c.baseURL)

	// Convert request to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Log the request for debugging
	log.Printf("Properties API Request URL: %s", url)
	log.Printf("Properties API Request Body: %s", string(jsonData))

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// If Fields is nil (requesting all fields), unmarshal as map to preserve all fields
	if req.Fields == nil {
		var rawResp map[string]interface{}
		if err := json.Unmarshal(body, &rawResp); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}
		
		// Convert to PropertiesResponse format with map[string]interface{} for data
		propertiesResp := &PropertiesResponse{}
		if data, ok := rawResp["data"].([]interface{}); ok {
			dataSlice := make([]map[string]interface{}, 0, len(data))
			for _, item := range data {
				if itemMap, ok := item.(map[string]interface{}); ok {
					dataSlice = append(dataSlice, itemMap)
				}
			}
			propertiesResp.Data = dataSlice
		}
		if count, ok := rawResp["count"].(float64); ok {
			propertiesResp.Count = int(count)
		}
		if limit, ok := rawResp["limit"].(float64); ok {
			propertiesResp.Limit = int(limit)
		}
		if totalCount, ok := rawResp["total_count"].(float64); ok {
			propertiesResp.TotalCount = int(totalCount)
		}
		
		return propertiesResp, nil
	}

	// Parse response using Property struct for limited fields
	var propertiesResp PropertiesResponse
	if err := json.Unmarshal(body, &propertiesResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &propertiesResp, nil
}

// QueryPropertiesGET queries the properties API using GET method with query parameters
func (c *Client) QueryPropertiesGET(fields []string, includePhotos bool, limit int, orderBy string) (*PropertiesResponse, error) {
	url := fmt.Sprintf("%s/api/v1/properties", c.baseURL)

	// Build query string
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	q := req.URL.Query()
	if len(fields) > 0 {
		fieldsStr := ""
		for i, f := range fields {
			if i > 0 {
				fieldsStr += ","
			}
			fieldsStr += f
		}
		q.Add("fields", fieldsStr)
	}
	if includePhotos {
		q.Add("include_photos", "true")
		q.Add("photo_fields", "photo_id,uri_thumb,uri_800,is_primary")
	}
	if limit > 0 {
		q.Add("limit", fmt.Sprintf("%d", limit))
	}
	if orderBy != "" {
		q.Add("order_by", orderBy)
	}

	req.URL.RawQuery = q.Encode()

	// Log the request for debugging
	log.Printf("Properties API GET Request URL: %s", req.URL.String())

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var propertiesResp PropertiesResponse
	if err := json.Unmarshal(body, &propertiesResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &propertiesResp, nil
}

// GetInitialProperties fetches all non-closed properties from the cached initial endpoint
// Returns the raw JSON array (handles gzip-compressed response)
func (c *Client) GetInitialProperties() ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/properties/initial", c.baseURL)

	// Create HTTP request
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Accept gzip encoding
	httpReq.Header.Set("Accept-Encoding", "gzip")

	// Log the request for debugging
	log.Printf("Properties API Initial Request URL: %s", url)

	// Execute request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code first
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Determine if response is gzip-compressed
	var reader io.Reader = resp.Body
	contentEncoding := resp.Header.Get("Content-Encoding")
	if strings.Contains(strings.ToLower(contentEncoding), "gzip") {
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer gzipReader.Close()
		reader = gzipReader
	}

	// Read response body (decompressed if gzip)
	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response as JSON array
	var properties []map[string]interface{}
	if err := json.Unmarshal(body, &properties); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return properties, nil
}

// GetInitialPropertiesStream fetches the gzip-compressed response and streams it directly
// Returns the HTTP response so it can be streamed to the client without decompression
func (c *Client) GetInitialPropertiesStream() (*http.Response, error) {
	url := fmt.Sprintf("%s/api/v1/properties/initial", c.baseURL)

	// Create HTTP request
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Accept gzip encoding
	httpReq.Header.Set("Accept-Encoding", "gzip")

	// Log the request for debugging
	log.Printf("Properties API Initial Request URL: %s", url)

	// Execute request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	return resp, nil
}

// SearchRentals searches for rental properties by a list of unparsed addresses
// This uses the /api/v1/properties/rentals/search endpoint which groups by address
// and returns only the latest version per address
func (c *Client) SearchRentals(addresses []string, fields []string, includePhotos bool, photoFields []string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/properties/rentals/search", c.baseURL)

	// Build request body
	reqBody := map[string]interface{}{
		"addresses": addresses,
	}
	if fields != nil {
		reqBody["fields"] = fields
	}
	if includePhotos {
		reqBody["include_photos"] = true
		if photoFields != nil {
			reqBody["photo_fields"] = photoFields
		}
	}

	// Convert request to JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Log the request for debugging
	log.Printf("Properties API Rentals Search Request URL: %s", url)
	log.Printf("Properties API Rentals Search Request Body: %s", string(jsonData))

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var response struct {
		Data       []map[string]interface{} `json:"data"`
		Count      int                      `json:"count"`
		TotalCount int                      `json:"total_count"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.Data, nil
}

