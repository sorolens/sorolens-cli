package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// Client is a typed HTTP client for the Sorolens REST API.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// New creates a Client. httpClient may be nil to use http.DefaultClient.
func New(baseURL string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{baseURL: baseURL, httpClient: httpClient}
}

func (c *Client) get(ctx context.Context, path string, out interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	return c.do(req, out)
}

func (c *Client) post(ctx context.Context, path string, body interface{}, out interface{}) error {
	buf, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(buf))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return c.do(req, out)
}

func (c *Client) do(req *http.Request, out interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		apiErr := &APIError{StatusCode: resp.StatusCode}
		_ = json.Unmarshal(raw, apiErr)
		if apiErr.Message == "" {
			apiErr.Message = fmt.Sprintf("HTTP %d", resp.StatusCode)
		}
		return apiErr
	}

	return json.Unmarshal(raw, out)
}

// GetContract fetches summary info for a contract.
func (c *Client) GetContract(ctx context.Context, id string) (*Contract, error) {
	var v Contract
	if err := c.get(ctx, "/api/v1/contracts/"+url.PathEscape(id), &v); err != nil {
		return nil, err
	}
	return &v, nil
}

// GetEvents fetches events for a contract.
func (c *Client) GetEvents(ctx context.Context, id, eventType, cursor string, limit int) (*EventsResponse, error) {
	q := url.Values{}
	if eventType != "" {
		q.Set("type", eventType)
	}
	if cursor != "" {
		q.Set("cursor", cursor)
	}
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	path := "/api/v1/contracts/" + url.PathEscape(id) + "/events"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	var v EventsResponse
	if err := c.get(ctx, path, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

// GetStorage fetches storage entries for a contract.
func (c *Client) GetStorage(ctx context.Context, id string) (*StorageResponse, error) {
	var v StorageResponse
	if err := c.get(ctx, "/api/v1/contracts/"+url.PathEscape(id)+"/storage", &v); err != nil {
		return nil, err
	}
	return &v, nil
}

// Track registers a contract for tracking.
func (c *Client) Track(ctx context.Context, contractID, alias string) (*TrackResponse, error) {
	req := TrackRequest{ContractID: contractID, Alias: alias}
	var v TrackResponse
	if err := c.post(ctx, "/api/v1/contracts", req, &v); err != nil {
		return nil, err
	}
	return &v, nil
}
