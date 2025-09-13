package exiftool

import (
	"fmt"
	"log"
	"testing"

	"github.com/goccy/go-json"
)

func TestNewExif(t *testing.T) {

	exifTool := NewExifTool()

	// Example with image file
	imageMetadata, err := exifTool.GetMetadata("/app/tmp/test.jpg")
	if err != nil {
		log.Fatalf("Error reading image metadata: %v", err)
	}

	fmt.Printf("\nImage Metadata:\n")
	fmt.Printf("File: %s (%s)\n", imageMetadata.FileInfo.BaseURL, imageMetadata.FileInfo.FileSize)
	fmt.Printf("Type: %s (%s)\n", imageMetadata.FileInfo.FileType, imageMetadata.FileInfo.MimeType)
	fmt.Printf("Dimensions: %dx%d (%.1f MP)\n", imageMetadata.Image.Width, imageMetadata.Image.Height, imageMetadata.Image.Megapixels)
	fmt.Printf("Camera: %s %s\n", imageMetadata.Camera.Make, imageMetadata.Camera.Model)
	if !imageMetadata.DateTimeOriginal.IsZero() {
		fmt.Printf("Date: %s\n", imageMetadata.DateTimeOriginal.Format("2006-01-02 15:04:05"))
	}

	// Print JSON output
	jsonData, _ := json.MarshalIndent(imageMetadata, "", "  ")
	fmt.Printf("\nJSON Output:\n%s\n", string(jsonData))

	err = exifTool.SaveMetadata(imageMetadata, "/app/tmp/test.json")
	if err != nil {
		return
	}

	//// Print raw data for debugging (first 10 keys)
	//fmt.Println("\nFirst 10 raw keys and values:")
	//count := 0
	//for key, value := range imageMetadata.RawData {
	//	if count >= 10 {
	//		break
	//	}
	//	fmt.Printf("  %s: %v\n", key, value)
	//	count++
	//}

}
