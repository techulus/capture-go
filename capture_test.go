package capture

import (
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	key := "test_key"
	secret := "test_secret"

	c := New(key, secret)
	if c.Key != key {
		t.Errorf("Expected key %s, got %s", key, c.Key)
	}
	if c.Secret != secret {
		t.Errorf("Expected secret %s, got %s", secret, c.Secret)
	}
	if c.UseEdge {
		t.Error("Expected UseEdge to be false by default")
	}
	if c.Client == nil {
		t.Error("Expected HTTP client to be initialized")
	}
}

func TestNewWithEdge(t *testing.T) {
	c := New("key", "secret", WithEdge())
	if !c.UseEdge {
		t.Error("Expected UseEdge to be true when using WithEdge option")
	}
}

func TestNewWithHTTPClient(t *testing.T) {
	customClient := &http.Client{}
	c := New("key", "secret", WithHTTPClient(customClient))
	if c.Client != customClient {
		t.Error("Expected custom HTTP client to be set")
	}
}

func TestGenerateToken(t *testing.T) {
	c := New("key", "secret")
	secret := "test_secret"
	query := "url=example.com&width=1920"
	
	token := c.generateToken(secret, query)
	if token == "" {
		t.Error("Expected non-empty token")
	}
	
	// Token should be consistent for same inputs
	token2 := c.generateToken(secret, query)
	if token != token2 {
		t.Error("Expected consistent token generation")
	}
}

func TestToQueryString(t *testing.T) {
	c := New("key", "secret")
	
	tests := []struct {
		name     string
		options  RequestOptions
		expected string
	}{
		{
			name:     "empty options",
			options:  RequestOptions{},
			expected: "",
		},
		{
			name: "basic options",
			options: RequestOptions{
				"width":  1920,
				"height": 1080,
			},
			expected: "height=1080&width=1920",
		},
		{
			name: "with format",
			options: RequestOptions{
				"width":  1920,
				"format": "png",
			},
			expected: "format=png&width=1920",
		},
		{
			name: "with empty values (should be ignored)",
			options: RequestOptions{
				"width":  1920,
				"empty":  "",
				"zero":   0,
				"false":  false,
			},
			expected: "false=false&width=1920&zero=0",
		},
		{
			name: "with special characters",
			options: RequestOptions{
				"userAgent": "Custom Agent (v1.0)",
			},
			expected: "userAgent=Custom+Agent+%28v1.0%29",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := c.toQueryString(tt.options)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestBuildURL(t *testing.T) {
	targetURL := "https://example.com"
	
	tests := []struct {
		name        string
		requestType RequestType
		options     RequestOptions
		useEdge     bool
		expectError bool
	}{
		{
			name:        "image capture",
			requestType: RequestTypeImage,
			options:     RequestOptions{"width": 1920},
			useEdge:     false,
			expectError: false,
		},
		{
			name:        "pdf capture",
			requestType: RequestTypePDF,
			options:     RequestOptions{"format": "A4"},
			useEdge:     false,
			expectError: false,
		},
		{
			name:        "content capture",
			requestType: RequestTypeContent,
			options:     RequestOptions{},
			useEdge:     false,
			expectError: false,
		},
		{
			name:        "metadata capture",
			requestType: RequestTypeMetadata,
			options:     RequestOptions{},
			useEdge:     false,
			expectError: false,
		},
		{
			name:        "edge mode",
			requestType: RequestTypeImage,
			options:     RequestOptions{},
			useEdge:     true,
			expectError: false,
		},
		{
			name:        "missing key",
			requestType: RequestTypeImage,
			options:     RequestOptions{},
			useEdge:     false,
			expectError: true,
		},
		{
			name:        "missing secret",
			requestType: RequestTypeImage,
			options:     RequestOptions{},
			useEdge:     false,
			expectError: true,
		},
		{
			name:        "empty url",
			requestType: RequestTypeImage,
			options:     RequestOptions{},
			useEdge:     false,
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new client for each test to avoid state pollution
			testClient := New("test_key", "test_secret")
			if tt.useEdge {
				testClient.UseEdge = true
			}
			
			// Handle special cases for error testing
			if tt.name == "missing key" {
				testClient.Key = ""
			} else if tt.name == "missing secret" {
				testClient.Secret = ""
			} else if tt.name == "empty url" {
				targetURL = ""
			}
			
			url, err := testClient.buildURL(tt.requestType, targetURL, tt.options)
			
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if url == "" {
				t.Error("Expected non-empty URL")
			}
			
			// Verify URL structure
			expectedBase := testClient.APIURL
			if tt.useEdge {
				expectedBase = testClient.EdgeURL
			}
			
			if !contains(url, expectedBase) {
				t.Errorf("URL should contain base URL %s, got %s", expectedBase, url)
			}
			
			if !contains(url, string(tt.requestType)) {
				t.Errorf("URL should contain request type %s, got %s", tt.requestType, url)
			}
		})
	}
}

func TestBuildImageURL(t *testing.T) {
	c := New("test_key", "test_secret")
	url, err := c.BuildImageURL("https://example.com", RequestOptions{"width": 1920})
	
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	if url == "" {
		t.Error("Expected non-empty URL")
	}
	
	if !contains(url, string(RequestTypeImage)) {
		t.Errorf("URL should contain image type, got %s", url)
	}
}

func TestBuildPDFURL(t *testing.T) {
	c := New("test_key", "test_secret")
	url, err := c.BuildPDFURL("https://example.com", RequestOptions{"format": "A4"})
	
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	if url == "" {
		t.Error("Expected non-empty URL")
	}
	
	if !contains(url, string(RequestTypePDF)) {
		t.Errorf("URL should contain PDF type, got %s", url)
	}
}

func TestBuildContentURL(t *testing.T) {
	c := New("test_key", "test_secret")
	url, err := c.BuildContentURL("https://example.com", RequestOptions{})
	
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	if url == "" {
		t.Error("Expected non-empty URL")
	}
	
	if !contains(url, string(RequestTypeContent)) {
		t.Errorf("URL should contain content type, got %s", url)
	}
}

func TestBuildMetadataURL(t *testing.T) {
	c := New("test_key", "test_secret")
	url, err := c.BuildMetadataURL("https://example.com", RequestOptions{})
	
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	if url == "" {
		t.Error("Expected non-empty URL")
	}
	
	if !contains(url, string(RequestTypeMetadata)) {
		t.Errorf("URL should contain metadata type, got %s", url)
	}
}

func TestScreenshotOptions(t *testing.T) {
	c := New("test_key", "test_secret")
	
	// Test all valid screenshot options according to official docs
	options := RequestOptions{
		"url":                "https://example.com",
		"httpAuth":           "base64url_encoded_auth",
		"vw":                 1440,                    // Viewport Width (default: 1440)
		"vh":                 900,                     // Viewport Height (default: 900)
		"scaleFactor":        1,                       // Screen scale factor (default: 1)
		"top":                0,                       // Top offset for clipping (default: 0)
		"left":               0,                       // Left offset for clipping (default: 0)
		"width":              1920,                    // Clipping Width (default: Viewport Width)
		"height":             1080,                    // Clipping Height (default: Viewport Height)
		"waitFor":            ".selector",             // Wait for CSS selector
		"waitForId":          "element-id",            // Wait for element ID
		"delay":              2,                       // Delay in seconds (default: 0)
		"full":               true,                    // Full page capture (default: false)
		"darkMode":           true,                    // Dark mode screenshot (default: false)
		"blockCookieBanners": true,                    // Block cookie banners (default: false)
		"blockAds":           true,                    // Block ads (default: false)
		"bypassBotDetection": true,                    // Bypass bot detection (default: false)
		"selector":           ".specific-element",     // Screenshot specific element
		"selectorId":         "specific-id",           // Screenshot element by ID
		"transparent":        true,                    // Transparent background (default: false)
		"userAgent":          "Custom User Agent",     // Custom user agent
		"timestamp":          "1234567890",            // Force reload
		"fresh":              true,                    // Fresh screenshot (default: false)
		"resizeWidth":        800,                     // Resize width
		"resizeHeight":       600,                     // Resize height
		"fileName":           "screenshot.png",        // S3 file name
		"s3Acl":              "public-read",           // S3 ACL
		"s3Redirect":         true,                    // Redirect to S3 URL (default: false)
		"skipUpload":         true,                    // Skip S3 upload (default: false)
		"type":               "png",                   // Image type: png, jpeg, webp (default: png)
		"bestFormat":         true,                    // Best format (default: false)
	}
	
	url, err := c.BuildImageURL("https://example.com", options)
	if err != nil {
		t.Errorf("Unexpected error building screenshot URL: %v", err)
	}
	
	if url == "" {
		t.Error("Expected non-empty screenshot URL")
	}
	
	// Verify URL contains the request type
	if !contains(url, string(RequestTypeImage)) {
		t.Errorf("URL should contain image type, got %s", url)
	}
}

func TestPDFOptions(t *testing.T) {
	c := New("test_key", "test_secret")
	
	// Test all valid PDF options according to official docs
	options := RequestOptions{
		"url":          "https://example.com",
		"httpAuth":     "base64url_encoded_auth",
		"userAgent":    "Custom User Agent",
		"width":        "8.5in",              // Paper width with units
		"height":       "11in",               // Paper height with units
		"marginTop":    "1in",                // Top margin with units
		"marginRight":  "1in",                // Right margin with units
		"marginBottom": "1in",                // Bottom margin with units
		"marginLeft":   "1in",                // Left margin with units
		"scale":        1,                    // Scale of webpage rendering (default: 1)
		"landscape":    true,                 // Paper orientation (default: false)
		"delay":        2,                    // Delay in seconds (default: 0)
		"timestamp":    "1234567890",         // Force reload
		"format":       "A4",                 // Paper format (default: A4)
		"fileName":     "document.pdf",       // S3 file name
		"s3Acl":        "public-read",        // S3 ACL
		"s3Redirect":   true,                 // Redirect to S3 URL (default: false)
	}
	
	url, err := c.BuildPDFURL("https://example.com", options)
	if err != nil {
		t.Errorf("Unexpected error building PDF URL: %v", err)
	}
	
	if url == "" {
		t.Error("Expected non-empty PDF URL")
	}
	
	// Verify URL contains the request type
	if !contains(url, string(RequestTypePDF)) {
		t.Errorf("URL should contain PDF type, got %s", url)
	}
}

func TestContentOptions(t *testing.T) {
	c := New("test_key", "test_secret")
	
	// Test all valid content options according to official docs
	options := RequestOptions{
		"url":       "https://example.com",
		"httpAuth":  "base64url_encoded_auth",
		"userAgent": "Custom User Agent",
		"delay":     2,                    // Delay in seconds (default: 0)
		"waitFor":   ".selector",          // Wait for CSS selector
		"waitForId": "element-id",         // Wait for element ID
	}
	
	url, err := c.BuildContentURL("https://example.com", options)
	if err != nil {
		t.Errorf("Unexpected error building content URL: %v", err)
	}
	
	if url == "" {
		t.Error("Expected non-empty content URL")
	}
	
	// Verify URL contains the request type
	if !contains(url, string(RequestTypeContent)) {
		t.Errorf("URL should contain content type, got %s", url)
	}
}

func TestMetadataOptions(t *testing.T) {
	c := New("test_key", "test_secret")
	
	// Test all valid metadata options according to official docs
	// Metadata API only supports basic options
	options := RequestOptions{
		"url":       "https://example.com",
		"httpAuth":  "base64url_encoded_auth",
		"userAgent": "Custom User Agent",
		"delay":     2,                    // Delay in seconds (default: 0)
		"waitFor":   ".selector",          // Wait for CSS selector
		"waitForId": "element-id",         // Wait for element ID
	}
	
	url, err := c.BuildMetadataURL("https://example.com", options)
	if err != nil {
		t.Errorf("Unexpected error building metadata URL: %v", err)
	}
	
	if url == "" {
		t.Error("Expected non-empty metadata URL")
	}
	
	// Verify URL contains the request type
	if !contains(url, string(RequestTypeMetadata)) {
		t.Errorf("URL should contain metadata type, got %s", url)
	}
}

func TestInvalidOptions(t *testing.T) {
	c := New("test_key", "test_secret")
	
	// Test that invalid options are handled gracefully
	invalidOptions := RequestOptions{
		"invalidOption": "value",
		"anotherInvalid": 123,
	}
	
	url, err := c.BuildImageURL("https://example.com", invalidOptions)
	if err != nil {
		t.Errorf("Unexpected error with invalid options: %v", err)
	}
	
	if url == "" {
		t.Error("Expected non-empty URL even with invalid options")
	}
}

func TestOptionTypeHandling(t *testing.T) {
	c := New("test_key", "test_secret")
	
	// Test different data types for options
	options := RequestOptions{
		"width":      1920,                    // int
		"height":     1080,                    // int
		"delay":      2.5,                     // float64
		"full":       true,                    // bool
		"userAgent":  "Custom Agent",          // string
		"scale":      1.0,                     // float64
		"landscape":  false,                   // bool
	}
	
	url, err := c.BuildImageURL("https://example.com", options)
	if err != nil {
		t.Errorf("Unexpected error with mixed option types: %v", err)
	}
	
	if url == "" {
		t.Error("Expected non-empty URL with mixed option types")
	}
}

func TestEmptyAndZeroValues(t *testing.T) {
	c := New("test_key", "test_secret")
	
	// Test that empty, zero, and false values are properly handled
	options := RequestOptions{
		"width":      0,                       // Should be included
		"height":     0,                       // Should be included
		"delay":      0,                       // Should be included
		"full":       false,                   // Should be included
		"empty":      "",                      // Should be excluded
		"nil":        nil,                     // Should be excluded
	}
	
	url, err := c.BuildImageURL("https://example.com", options)
	if err != nil {
		t.Errorf("Unexpected error with empty/zero values: %v", err)
	}
	
	if url == "" {
		t.Error("Expected non-empty URL with empty/zero values")
	}
	
	// Verify that empty and nil values are excluded from query string
	if contains(url, "empty=") {
		t.Error("URL should not contain empty value")
	}
	if contains(url, "nil=") {
		t.Error("URL should not contain nil value")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		s[len(s)-len(substr):] == substr || 
		func() bool {
			for i := 1; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())))
} 