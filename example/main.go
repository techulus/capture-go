package main

import (
	"fmt"
	"log"
	"os"

	"github.com/techulus/capture-go"
)

func main() {
	// Initialize the Capture client
	// Replace with your actual API key and secret
	key := os.Getenv("CAPTURE_KEY")
	secret := os.Getenv("CAPTURE_SECRET")

	if key == "" || secret == "" {
		log.Fatal("Please set CAPTURE_KEY and CAPTURE_SECRET environment variables")
	}

	// Create a new Capture client
	c := capture.New(key, secret)

	// Example 1: Build URLs for different capture types
	targetURL := "https://techulus.xyz"
	
	// Build image URL
	imageURL, err := c.BuildImageURL(targetURL, capture.RequestOptions{
		"width":  1920,
		"height": 1080,
		"type":   "png",
	})
	if err != nil {
		log.Printf("Error building image URL: %v", err)
	} else {
		fmt.Printf("Image URL: %s\n", imageURL)
	}

	// Build PDF URL
	pdfURL, err := c.BuildPDFURL(targetURL, capture.RequestOptions{
		"format": "A4",
	})
	if err != nil {
		log.Printf("Error building PDF URL: %v", err)
	} else {
		fmt.Printf("PDF URL: %s\n", pdfURL)
	}

	// Example 2: Fetch content
	contentResp, err := c.FetchContent(targetURL, capture.RequestOptions{})
	if err != nil {
		log.Printf("Error fetching content: %v", err)
	} else {
		fmt.Printf("Content Success: %t\n", contentResp.Success)
		fmt.Printf("HTML Length: %d\n", len(contentResp.HTML))
		fmt.Printf("Text Content Length: %d\n", len(contentResp.TextContent))
	}

	// Example 3: Fetch metadata
	metadataResp, err := c.FetchMetadata(targetURL, capture.RequestOptions{})
	if err != nil {
		log.Printf("Error fetching metadata: %v", err)
	} else {
		fmt.Printf("Metadata Success: %t\n", metadataResp.Success)
		fmt.Printf("Metadata: %+v\n", metadataResp.Metadata)
	}

	// Example 4: Using edge mode
	edgeClient := capture.New(key, secret, capture.WithEdge())
	edgeImageURL, err := edgeClient.BuildImageURL(targetURL, capture.RequestOptions{})
	if err != nil {
		log.Printf("Error building edge image URL: %v", err)
	} else {
		fmt.Printf("Edge Image URL: %s\n", edgeImageURL)
	}

	// Example 5: Fetch image (saves to file)
	imageData, err := c.FetchImage(targetURL, capture.RequestOptions{
		"width":  800,
		"height": 600,
		"type":   "png",
	})
	if err != nil {
		log.Printf("Error fetching image: %v", err)
	} else {
		// Save image to file
		err = os.WriteFile("screenshot.png", imageData, 0644)
		if err != nil {
			log.Printf("Error saving image: %v", err)
		} else {
			fmt.Println("Screenshot saved as screenshot.png")
		}
	}

	// Example 6: Fetch PDF (saves to file)
	pdfData, err := c.FetchPDF(targetURL, capture.RequestOptions{
		"format": "A4",
	})
	if err != nil {
		log.Printf("Error fetching PDF: %v", err)
	} else {
		// Save PDF to file
		err = os.WriteFile("document.pdf", pdfData, 0644)
		if err != nil {
			log.Printf("Error saving PDF: %v", err)
		} else {
			fmt.Println("PDF saved as document.pdf")
		}
	}
} 