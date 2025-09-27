package exiftool_v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

// ImageMetadata struct to hold the extracted EXIF data.
type ImageMetadata struct {
	FileSize    string
	FileType    string
	Make        string
	Model       string
	Orientation string
	CreateDate  string
}

func Start(filePath string) (*ImageMetadata, error) {

	// Construct the exiftool_v1 command to get JSON output with numerical values.
	// "-n" flag displays numerical values for tags where applicable.
	cmd := exec.Command("exiftool_v1", "-json", "-n", filePath)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Execute the exiftool_v1 command.
	err := cmd.Run()
	if err != nil {
		// Return the error instead of exiting the program
		return nil, fmt.Errorf("error running exiftool_v1 for '%s': %w, stderr: %s", filePath, err, stderr.String())
	}

	// ExifTool returns an array of JSON objects. We expect at least one for a single file.
	var exifResults []map[string]interface{}
	err = json.Unmarshal(stdout.Bytes(), &exifResults)
	if err != nil {
		log.Printf("Error unmarshalling JSON output from exiftool_v1: %v", err)
		return nil, err
	}

	// Check if any metadata was extracted.
	if len(exifResults) == 0 {
		fmt.Printf("No metadata found for file '%s'.\n", filePath)
		return nil, err
	}

	// The first element in the array contains the metadata for our file.
	rawExifData := exifResults[0]

	// Create an instance of our custom ImageMetadata struct.
	var metadata ImageMetadata

	// Populate the struct fields from the raw ExifTool output.
	// The fields might be of different types (string, float, int) so we need type assertions.

	// FileSize
	if fs, ok := rawExifData["FileSize"].(string); ok {
		metadata.FileSize = fs
	} else if fsNum, ok := rawExifData["FileSize"].(float64); ok {
		metadata.FileSize = formatBytes(int64(fsNum)) // Convert bytes to human-readable format
	} else {
		metadata.FileSize = "Unknown"
	}

	// FileType
	if ft, ok := rawExifData["FileType"].(string); ok {
		metadata.FileType = ft
	} else {
		metadata.FileType = "Unknown"
	}

	// Make
	if makeVal, ok := rawExifData["Make"].(string); ok {
		metadata.Make = makeVal
	} else {
		metadata.Make = "Unknown"
	}

	// Model
	if modelVal, ok := rawExifData["Model"].(string); ok {
		metadata.Model = modelVal
	} else {
		metadata.Model = "Unknown"
	}

	// Orientation - Note: exiftool_v1 -n provides numerical orientation, we map it to string here.
	if orientationNum, ok := rawExifData["Orientation"].(float64); ok {
		metadata.Orientation = getOrientationString(int(orientationNum))
	} else if orientationStr, ok := rawExifData["Orientation"].(string); ok {
		metadata.Orientation = orientationStr
	} else {
		metadata.Orientation = "Unknown"
	}

	// CreateDate - ExifTool provides various date fields, DateTimeOriginal or CreateDate are common.
	if createDate, ok := rawExifData["CreateDate"].(string); ok {
		metadata.CreateDate = createDate
	} else if dateTimeOriginal, ok := rawExifData["DateTimeOriginal"].(string); ok {
		metadata.CreateDate = dateTimeOriginal
	} else {
		metadata.CreateDate = "Unknown"
	}

	return &metadata, nil
}

// Helper function to convert numerical orientation to a human-readable string.
func getOrientationString(orientation int) string {
	switch orientation {
	case 1:
		return "Normal"
	case 2:
		return "Horizontal Flip"
	case 3:
		return "Rotate 180"
	case 4:
		return "Vertical Flip"
	case 5:
		return "Rotate 90 CW & Horizontal Flip"
	case 6:
		return "Rotate 90 CW"
	case 7:
		return "Rotate 90 CCW & Horizontal Flip"
	case 8:
		return "Rotate 90 CCW"
	default:
		return fmt.Sprintf("Unknown (%d)", orientation)
	}
}

// Helper function to format bytes into human-readable sizes (KB, MB, GB).
func formatBytes(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
