# Capture Go SDK

A Go SDK for the [Capture browser API](https://capture.page) that provides easy access to web page screenshots, PDF generation, content extraction, and metadata retrieval.

## Features

- **Image Capture**: Generate screenshots of web pages
- **PDF Generation**: Convert web pages to PDF documents
- **Content Extraction**: Extract HTML and text content from web pages
- **Metadata Retrieval**: Get metadata information from web pages
- **Edge Mode Support**: Use edge servers for faster processing
- **Custom HTTP Client**: Configure custom HTTP client settings
- **Type Safety**: Full Go type safety with proper error handling

## Installation

```bash
go get github.com/techulus/capture-go
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/techulus/capture-go"
)

func main() {
    // Initialize the Capture client
    key := os.Getenv("CAPTURE_KEY")
    secret := os.Getenv("CAPTURE_SECRET")
    
    c := capture.New(key, secret)
    
    // Capture a screenshot
    imageData, err := c.FetchImage("https://www.google.com", capture.RequestOptions{
        "vw":  1920,
        "vh": 1080,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Save the image
    os.WriteFile("screenshot.png", imageData, 0644)
    fmt.Println("Screenshot saved!")
}
```

## API Reference

> **Note**: This SDK implements the official Capture API. For the complete list of available options and their detailed descriptions, refer to the [official Capture API documentation](https://docs.capture.page/docs).

### Initialization

```go
// Basic initialization
c := capture.New(key, secret)

// With edge mode enabled
c := capture.New(key, secret, capture.WithEdge())

// With custom HTTP client
client := &http.Client{Timeout: 30 * time.Second}
c := capture.New(key, secret, capture.WithHTTPClient(client))
```

### URL Building

Build URLs for different capture types without making HTTP requests:

```go
// Build image URL
imageURL, err := c.BuildImageURL("https://example.com", capture.RequestOptions{
    "vw":  1920,
    "vh": 1080,
    "format": "png",
})

// Build PDF URL
pdfURL, err := c.BuildPDFURL("https://example.com", capture.RequestOptions{
    "format": "A4",
})

// Build content URL
contentURL, err := c.BuildContentURL("https://example.com", capture.RequestOptions{})

// Build metadata URL
metadataURL, err := c.BuildMetadataURL("https://example.com", capture.RequestOptions{})
```

### Fetching Data

Fetch actual data from the Capture API:

```go
// Fetch image as bytes
imageData, err := c.FetchImage("https://example.com", capture.RequestOptions{
    "vw":  800,
    "vh": 600,
})

// Fetch PDF as bytes
pdfData, err := c.FetchPDF("https://example.com", capture.RequestOptions{
    "format": "A4",
})

// Fetch content with structured response
contentResp, err := c.FetchContent("https://example.com", capture.RequestOptions{})
if err == nil {
    fmt.Printf("HTML: %s\n", contentResp.HTML)
    fmt.Printf("Text: %s\n", contentResp.TextContent)
}

// Fetch metadata with structured response
metadataResp, err := c.FetchMetadata("https://example.com", capture.RequestOptions{})
if err == nil {
    fmt.Printf("Metadata: %+v\n", metadataResp.Metadata)
}
```

### Request Options

The `RequestOptions` type is a map that accepts various configuration options. Different request types support different options:

#### Screenshot Options

```go
options := capture.RequestOptions{
    // Basic options
    "url":                "https://example.com",     // Target URL (required)
    "httpAuth":           "base64url_encoded_auth",  // HTTP Basic Authentication
    "userAgent":          "Custom User Agent",       // Custom user agent
    
    // Viewport and sizing
    "vw":                 1440,                       // Viewport Width (default: 1440)
    "vh":                 900,                        // Viewport Height (default: 900)
    "scaleFactor":        1,                          // Screen scale factor (default: 1)
    "width":              1920,                       // Clipping Width (default: Viewport Width)
    "height":             1080,                       // Clipping Height (default: Viewport Height)
    "top":                0,                          // Top offset for clipping (default: 0)
    "left":               0,                          // Left offset for clipping (default: 0)
    
    // Timing and waiting
    "delay":              2,                          // Delay in seconds (default: 0)
    "waitFor":            ".selector",                // Wait for CSS selector
    "waitForId":          "element-id",               // Wait for element ID
    
    // Capture behavior
    "full":               true,                       // Full page capture (default: false)
    "darkMode":           true,                       // Dark mode screenshot (default: false)
    "transparent":        true,                       // Transparent background (default: false)
    "selector":           ".specific-element",        // Screenshot specific element
    "selectorId":         "specific-id",              // Screenshot element by ID
    
    // Content blocking
    "blockCookieBanners": true,                       // Block cookie banners (default: false)
    "blockAds":           true,                       // Block ads (default: false)
    "bypassBotDetection": true,                       // Bypass bot detection (default: false)
    
    // Image processing
    "type":               "png",                      // Image type: png, jpeg, webp (default: png)
    "bestFormat":         true,                       // Best format (default: false)
    "resizeWidth":        800,                        // Resize width
    "resizeHeight":       600,                        // Resize height
    
    // Caching and reloading
    "timestamp":          "1234567890",               // Force reload
    "fresh":              true,                       // Fresh screenshot (default: false)
    
    // S3 integration
    "fileName":           "screenshot.png",           // S3 file name
    "s3Acl":              "public-read",              // S3 ACL
    "s3Redirect":         true,                       // Redirect to S3 URL (default: false)
    "skipUpload":         true,                       // Skip S3 upload (default: false)
}
```

#### PDF Options

```go
options := capture.RequestOptions{
    // Basic options
    "url":          "https://example.com",     // Target URL (required)
    "httpAuth":     "base64url_encoded_auth",  // HTTP Basic Authentication
    "userAgent":    "Custom User Agent",       // Custom user agent
    
    // Paper size and format
    "format":       "A4",                      // Paper format: Letter, Legal, Tabloid, Ledger, A0-A6 (default: A4)
    "width":        "8.5in",                   // Custom paper width with units
    "height":       "11in",                    // Custom paper height with units
    
    // Margins
    "marginTop":    "1in",                     // Top margin with units
    "marginRight":  "1in",                     // Right margin with units
    "marginBottom": "1in",                     // Bottom margin with units
    "marginLeft":   "1in",                     // Left margin with units
    
    // Rendering
    "scale":        1,                         // Scale of webpage rendering (default: 1)
    "landscape":    true,                      // Paper orientation (default: false)
    
    // Timing
    "delay":        2,                         // Delay in seconds (default: 0)
    "timestamp":    "1234567890",              // Force reload
    
    // S3 integration
    "fileName":     "document.pdf",            // S3 file name
    "s3Acl":        "public-read",             // S3 ACL
    "s3Redirect":   true,                      // Redirect to S3 URL (default: false)
}
```

#### Content Options

```go
options := capture.RequestOptions{
    "url":       "https://example.com",     // Target URL (required)
    "httpAuth":  "base64url_encoded_auth",  // HTTP Basic Authentication
    "userAgent": "Custom User Agent",       // Custom user agent
    "delay":     2,                         // Delay in seconds (default: 0)
    "waitFor":   ".selector",               // Wait for CSS selector
    "waitForId": "element-id",              // Wait for element ID
}
```

#### Metadata Options

```go
options := capture.RequestOptions{
    "url":       "https://example.com",     // Target URL (required)
    "httpAuth":  "base64url_encoded_auth",  // HTTP Basic Authentication
    "userAgent": "Custom User Agent",       // Custom user agent
    "delay":     2,                         // Delay in seconds (default: 0)
    "waitFor":   ".selector",               // Wait for CSS selector
    "waitForId": "element-id",              // Wait for element ID
}
```

### Request Type Differences

Each request type supports different options:

- **Screenshot**: Supports the most options including viewport settings, image processing, content blocking, and S3 integration
- **PDF**: Supports paper size, margins, orientation, and basic timing options
- **Content**: Supports basic options for authentication, timing, and element waiting
- **Metadata**: Supports the same basic options as Content

For detailed information about which options are supported by each request type, see the [official API documentation](https://docs.capture.page/docs).

### Response Types

#### ContentResponse
```go
type ContentResponse struct {
    Success     bool   `json:"success"`
    HTML        string `json:"html"`
    TextContent string `json:"textContent"`
}
```

#### MetadataResponse
```go
type MetadataResponse struct {
    Success  bool                   `json:"success"`
    Metadata map[string]interface{} `json:"metadata"`
}
```

## Configuration Options

### Edge Mode

Enable edge mode for faster processing:

```go
c := capture.New(key, secret, capture.WithEdge())
```

### Custom HTTP Client

Configure a custom HTTP client:

```go
client := &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        10,
        IdleConnTimeout:     30 * time.Second,
        DisableCompression:  true,
    },
}

c := capture.New(key, secret, capture.WithHTTPClient(client))
```

## Error Handling

The SDK provides comprehensive error handling:

```go
imageData, err := c.FetchImage("https://example.com", capture.RequestOptions{})
if err != nil {
    switch {
    case strings.Contains(err.Error(), "key and secret are required"):
        log.Fatal("Missing API credentials")
    case strings.Contains(err.Error(), "url is required"):
        log.Fatal("Missing target URL")
    case strings.Contains(err.Error(), "HTTP error"):
        log.Fatal("API request failed")
    default:
        log.Fatal("Unexpected error:", err)
    }
}
```

## Examples

See the `example/` directory for complete working examples:

- Basic usage examples
- Error handling
- File saving
- Different capture types
- Edge mode usage

Set these environment variables for authentication while running the examples:

```bash
export CAPTURE_KEY="your_api_key"
export CAPTURE_SECRET="your_api_secret"
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## Support

For support and questions:
- GitHub Issues: [https://github.com/techulus/capture-go/issues](https://github.com/techulus/capture-go/issues)
- Documentation: [https://capture.page/docs](https://capture.page/docs) 