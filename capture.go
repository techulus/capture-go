package capture

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

type RequestType string

const (
	RequestTypeImage    RequestType = "image"
	RequestTypePDF      RequestType = "pdf"
	RequestTypeContent  RequestType = "content"
	RequestTypeMetadata RequestType = "metadata"
	RequestTypeAnimated RequestType = "animated"
)

type RequestOptions map[string]interface{}

type CreateSessionOptions struct {
	MaxTtlSeconds      int  `json:"maxTtlSeconds,omitempty"`
	Proxy              bool `json:"proxy,omitempty"`
	BypassBotDetection bool `json:"bypassBotDetection,omitempty"`
}

type SessionActionPayload map[string]interface{}
type SessionActionResponse map[string]interface{}
type SessionResponse map[string]interface{}

type SessionsAPIError struct {
	StatusCode int
	Body       map[string]interface{}
}

func (e *SessionsAPIError) Error() string {
	if message, ok := e.Body["error"].(string); ok && message != "" {
		return message
	}
	return fmt.Sprintf("Capture Sessions API request failed with status %d", e.StatusCode)
}

type Capture struct {
	APIURL  string
	EdgeURL string
	Key     string
	Secret  string
	UseEdge bool
	Client  *http.Client
}

func New(key, secret string, options ...Option) *Capture {
	c := &Capture{
		APIURL:  "https://cdn.capture.page",
		EdgeURL: "https://edge.capture.page",
		Key:     key,
		Secret:  secret,
		UseEdge: false,
		Client:  &http.Client{},
	}

	for _, option := range options {
		option(c)
	}

	return c
}

type Option func(*Capture)

func WithEdge() Option {
	return func(c *Capture) {
		c.UseEdge = true
	}
}

func WithHTTPClient(client *http.Client) Option {
	return func(c *Capture) {
		c.Client = client
	}
}

func (c *Capture) generateToken(secret, query string) string {
	hash := md5.Sum([]byte(secret + query))
	return hex.EncodeToString(hash[:])
}

func (c *Capture) toQueryString(options RequestOptions) string {
	if options == nil {
		return ""
	}

	params := make(map[string]string)

	for key, value := range options {
		if value == nil || value == "" {
			continue
		}

		var strValue string
		switch v := value.(type) {
		case string:
			strValue = v
		case int:
			strValue = strconv.Itoa(v)
		case int64:
			strValue = strconv.FormatInt(v, 10)
		case float64:
			strValue = strconv.FormatFloat(v, 'f', -1, 64)
		case bool:
			strValue = strconv.FormatBool(v)
		default:
			strValue = fmt.Sprintf("%v", v)
		}

		params[key] = strValue
	}

	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var queryParts []string
	for _, key := range keys {
		value := params[key]
		encodedKey := url.QueryEscape(key)
		encodedValue := url.QueryEscape(value)
		queryParts = append(queryParts, encodedKey+"="+encodedValue)
	}

	return strings.Join(queryParts, "&")
}

func (c *Capture) buildURL(requestType RequestType, targetURL string, options RequestOptions) (string, error) {
	if c.Key == "" || c.Secret == "" {
		return "", fmt.Errorf("key and secret are required")
	}

	if targetURL == "" {
		return "", fmt.Errorf("url is required")
	}

	requestOptions := make(RequestOptions)
	if options != nil {
		for k, v := range options {
			requestOptions[k] = v
		}
	}
	requestOptions["url"] = targetURL

	query := c.toQueryString(requestOptions)

	token := c.generateToken(c.Secret, query)

	baseURL := c.APIURL
	if c.UseEdge {
		baseURL = c.EdgeURL
	}

	finalURL := fmt.Sprintf("%s/%s/%s/%s", baseURL, c.Key, token, requestType)
	if query != "" {
		finalURL += "?" + query
	}

	return finalURL, nil
}

func (c *Capture) BuildImageURL(targetURL string, options RequestOptions) (string, error) {
	return c.buildURL(RequestTypeImage, targetURL, options)
}

func (c *Capture) BuildPDFURL(targetURL string, options RequestOptions) (string, error) {
	return c.buildURL(RequestTypePDF, targetURL, options)
}

func (c *Capture) BuildContentURL(targetURL string, options RequestOptions) (string, error) {
	return c.buildURL(RequestTypeContent, targetURL, options)
}

func (c *Capture) BuildMetadataURL(targetURL string, options RequestOptions) (string, error) {
	return c.buildURL(RequestTypeMetadata, targetURL, options)
}

func (c *Capture) BuildAnimatedURL(targetURL string, options RequestOptions) (string, error) {
	return c.buildURL(RequestTypeAnimated, targetURL, options)
}

func (c *Capture) FetchImage(targetURL string, options RequestOptions) ([]byte, error) {
	url, err := c.BuildImageURL(targetURL, options)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return buf, nil
}

func (c *Capture) FetchPDF(targetURL string, options RequestOptions) ([]byte, error) {
	url, err := c.BuildPDFURL(targetURL, options)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch PDF: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return buf, nil
}

type ContentResponse struct {
	Success     bool   `json:"success"`
	HTML        string `json:"html"`
	TextContent string `json:"textContent"`
	Markdown    string `json:"markdown"`
}

func (c *Capture) FetchContent(targetURL string, options RequestOptions) (*ContentResponse, error) {
	url, err := c.BuildContentURL(targetURL, options)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch content: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	var contentResp ContentResponse
	if err := json.NewDecoder(resp.Body).Decode(&contentResp); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return &contentResp, nil
}

type MetadataResponse struct {
	Success  bool                   `json:"success"`
	Metadata map[string]interface{} `json:"metadata"`
}

func (c *Capture) FetchMetadata(targetURL string, options RequestOptions) (*MetadataResponse, error) {
	url, err := c.BuildMetadataURL(targetURL, options)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metadata: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	var metadataResp MetadataResponse
	if err := json.NewDecoder(resp.Body).Decode(&metadataResp); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return &metadataResp, nil
}

func (c *Capture) FetchAnimated(targetURL string, options RequestOptions) ([]byte, error) {
	url, err := c.BuildAnimatedURL(targetURL, options)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch animated: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return buf, nil
}

func (c *Capture) sessionsBearerToken() (string, error) {
	if c.Key == "" || c.Secret == "" {
		return "", fmt.Errorf("key and secret are required")
	}

	return base64.StdEncoding.EncodeToString([]byte(c.Key + ":" + c.Secret)), nil
}

func (c *Capture) sessionURL(path string) string {
	return strings.TrimRight(c.EdgeURL, "/") + "/v1/sessions" + path
}

func (c *Capture) sessionRequest(method, path string, body interface{}, out interface{}) error {
	token, err := c.sessionsBearerToken()
	if err != nil {
		return err
	}

	var requestBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to encode session request body: %w", err)
		}
		requestBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, c.sessionURL(path), requestBody)
	if err != nil {
		return fmt.Errorf("failed to build session request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute session request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read session response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		decoded := map[string]interface{}{}
		if len(respBody) > 0 {
			if err := json.Unmarshal(respBody, &decoded); err != nil {
				decoded["error"] = string(respBody)
			}
		}
		return &SessionsAPIError{StatusCode: resp.StatusCode, Body: decoded}
	}

	if len(respBody) == 0 {
		return nil
	}

	if err := json.Unmarshal(respBody, out); err != nil {
		return fmt.Errorf("failed to decode session response: %w", err)
	}

	return nil
}

func (c *Capture) CreateSession(options *CreateSessionOptions) (SessionResponse, error) {
	if options == nil {
		options = &CreateSessionOptions{}
	}

	var response SessionResponse
	if err := c.sessionRequest(http.MethodPost, "", options, &response); err != nil {
		return nil, err
	}
	return response, nil
}

func (c *Capture) GetSession(sessionID string) (SessionResponse, error) {
	var response SessionResponse
	if err := c.sessionRequest(http.MethodGet, "/"+url.PathEscape(sessionID), nil, &response); err != nil {
		return nil, err
	}
	return response, nil
}

func (c *Capture) CloseSession(sessionID string) (SessionResponse, error) {
	var response SessionResponse
	if err := c.sessionRequest(http.MethodDelete, "/"+url.PathEscape(sessionID), nil, &response); err != nil {
		return nil, err
	}
	return response, nil
}

func (c *Capture) ExecuteAction(sessionID, actionType string, payload SessionActionPayload) (SessionActionResponse, error) {
	if payload == nil {
		payload = SessionActionPayload{}
	}

	body := map[string]interface{}{
		"type":    actionType,
		"payload": payload,
	}
	var response SessionActionResponse
	if err := c.sessionRequest(http.MethodPost, "/"+url.PathEscape(sessionID)+"/actions", body, &response); err != nil {
		return nil, err
	}
	return response, nil
}
