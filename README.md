# Capture Go SDK & CLI

Go SDK and CLI for the [Capture API](https://capture.page) - screenshots, PDFs, content extraction, and more.

## Installation

**SDK:**
```bash
go get github.com/techulus/capture-go
```

**CLI:**
```bash
go install github.com/techulus/capture-go/cmd/capture@latest
```

Or via Homebrew:
```bash
brew tap techulus/tap
brew install capture
```

## CLI Usage

Set your credentials:
```bash
export CAPTURE_KEY="your_api_key"
export CAPTURE_SECRET="your_api_secret"
```

Commands:
```bash
capture screenshot https://example.com -o screenshot.png
capture screenshot https://example.com -X vw=1920 -X vh=1080 -X fullPage=true -o full.png

capture pdf https://example.com -o document.pdf
capture pdf https://example.com -X format=A4 -X landscape=true -o landscape.pdf

capture content https://example.com --format markdown
capture content https://example.com --format html -o page.html

capture metadata https://example.com --pretty

capture animated https://example.com -X duration=5 -o recording.gif

capture sessions create --max-ttl-seconds 300 --pretty
capture sessions get sess_123 --pretty
capture sessions action sess_123 goto --payload-json '{"url":"https://example.com"}'
capture sessions action sess_123 screenshot -X fullPage=true --pretty
capture sessions close sess_123 --pretty
```

Use `--edge` for faster response, `--dry-run` to preview the request URL.

See [docs.capture.page](https://docs.capture.page/) for all available options.

## SDK Usage

```go
package main

import (
    "os"
    "github.com/techulus/capture-go"
)

func main() {
    c := capture.New(os.Getenv("CAPTURE_KEY"), os.Getenv("CAPTURE_SECRET"))

    // Screenshot
    img, _ := c.FetchImage("https://example.com", capture.RequestOptions{
        "vw": 1920,
        "vh": 1080,
    })
    os.WriteFile("screenshot.png", img, 0644)

    // PDF
    pdf, _ := c.FetchPDF("https://example.com", capture.RequestOptions{
        "format": "A4",
    })
    os.WriteFile("document.pdf", pdf, 0644)

    // Content
    content, _ := c.FetchContent("https://example.com", capture.RequestOptions{})
    println(content.Markdown)

    // Metadata
    meta, _ := c.FetchMetadata("https://example.com", capture.RequestOptions{})
    println(meta.Metadata["title"])

    // Animated
    gif, _ := c.FetchAnimated("https://example.com", capture.RequestOptions{
        "duration": 5,
    })
    os.WriteFile("recording.gif", gif, 0644)

    // Browser sessions
    created, _ := c.CreateSession(&capture.CreateSessionOptions{
        MaxTtlSeconds: 300,
    })
    session := created["session"].(map[string]interface{})
    sessionID := session["id"].(string)

    c.ExecuteAction(sessionID, "goto", capture.SessionActionPayload{
        "url": "https://example.com",
    })
    c.ExecuteAction(sessionID, "screenshot", capture.SessionActionPayload{
        "fullPage": true,
    })
    c.CloseSession(sessionID)
}
```

### Options

```go
// Edge mode (faster)
c := capture.New(key, secret, capture.WithEdge())

// Custom HTTP client
c := capture.New(key, secret, capture.WithHTTPClient(&http.Client{
    Timeout: 60 * time.Second,
}))

// Build URL without fetching
url, _ := c.BuildImageURL("https://example.com", capture.RequestOptions{})
```

See [docs.capture.page](https://docs.capture.page/) for all available request options.

## Links

- [Documentation](https://docs.capture.page/)
- [Issues](https://github.com/techulus/capture-go/issues)
