package exiftool

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-json"
)

type ExifTool struct {
	exiftoolPath string
}

func NewExifTool() *ExifTool {
	return &ExifTool{
		exiftoolPath: "exiftool",
	}
}

func (et *ExifTool) SaveMetadata(me *Metadata, path string) error {

	jsonData, err := json.MarshalIndent(me, "", "  ")
	if err != nil {
		return err
	}

	tempFile := path + ".tmp"
	if err := os.WriteFile(tempFile, jsonData, 0644); err != nil {
		return err
	}
	return os.Rename(tempFile, path)
}

func (et *ExifTool) GetMetadata(filename string) (*Metadata, error) {

	// First, let's see what keys are available by running without groups
	cmd := exec.Command(et.exiftoolPath, "-j", "-c", "%.6f", filename)

	output, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return nil, fmt.Errorf("exiftool failed: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("failed to execute exiftool: %v", err)
	}

	var rawMetadata []map[string]interface{}
	if err := json.Unmarshal(output, &rawMetadata); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	if len(rawMetadata) == 0 {
		return nil, fmt.Errorf("no metadata found")
	}

	// Print all available keys for debugging
	//fmt.Println("Available EXIF keys:")
	//for key := range rawMetadata[0] {
	//	fmt.Printf("  - %s\n", key)
	//}

	return et.parseMetadata(filename, rawMetadata[0]), nil
}

func (et *ExifTool) parseMetadata(filename string, rawData map[string]interface{}) *Metadata {
	metadata := &Metadata{
		RawData: rawData, // Store raw data for debugging
	}

	metadata.DateTimeOriginal = et.parseDate(rawData)

	// Parse FileInfo
	metadata.FileInfo = et.parseFileInfo(filename, rawData)

	if strings.Contains(metadata.FileInfo.MimeType, "video") {
		// Parse VideoInfo
		metadata.Video = et.parseVideoInfo(rawData)
	} else {
		// Parse ImageInfo
		metadata.Image = et.parseImageInfo(rawData)
	}

	// Parse CameraInfo
	metadata.Camera = et.parseCameraInfo(rawData)

	// Parse Location
	metadata.Location = et.parseLocationInfo(rawData)

	return metadata
}

func (et *ExifTool) parseDate(rawData map[string]interface{}) time.Time {

	// Try multiple date formats
	if dateTimeStr, ok := getString(rawData, "DateTimeOriginal", "CreateDate", "ModifyDate"); ok {

		// Try different date formats
		formats := []string{
			"2006:01:02 15:04:05",
			"2006-01-02T15:04:05Z",
			"2006-01-02 15:04:05",
			time.RFC3339,
		}

		for _, format := range formats {
			if dateTime, err := time.Parse(format, dateTimeStr); err == nil {
				return dateTime
			}
		}
	}

	return time.Time{}
}

func (et *ExifTool) parseFileInfo(filename string, rawData map[string]interface{}) FileInfo {
	fileInfo := FileInfo{
		BaseURL: filepath.Base(filename),
	}

	if size, ok := getString(rawData, "FileSize"); ok {
		fileInfo.FileSize = size
	}

	if fileType, ok := getString(rawData, "FileType"); ok {
		fileInfo.FileType = fileType
	}

	if mimeType, ok := getString(rawData, "MIMEType"); ok {
		fileInfo.MimeType = mimeType
	}

	return fileInfo
}

func (et *ExifTool) parseImageInfo(rawData map[string]interface{}) ImageInfo {

	imageInfo := ImageInfo{}

	// Try multiple possible keys for dimensions
	if width, ok := getInt(rawData, "ImageWidth", "ExifImageWidth", "Width"); ok {
		imageInfo.Width = width
	}

	if height, ok := getInt(rawData, "ImageHeight", "ExifImageHeight", "Height"); ok {
		imageInfo.Height = height
	}

	// Calculate megapixels
	if imageInfo.Width > 0 && imageInfo.Height > 0 {
		imageInfo.Megapixels = float64(imageInfo.Width*imageInfo.Height) / 1000000.0
	}

	if orientation, ok := getString(rawData, "Orientation"); ok {
		imageInfo.Orientation = orientation
	}

	if colorSpace, ok := getString(rawData, "ColorSpace"); ok {
		imageInfo.ColorSpace = colorSpace
	}

	if encodingProcess, ok := getString(rawData, "EncodingProcess"); ok {
		imageInfo.EncodingProcess = encodingProcess
	}

	return imageInfo
}

func (et *ExifTool) parseVideoInfo(rawData map[string]interface{}) VideoInfo {
	videoInfo := VideoInfo{}

	if duration, ok := getString(rawData, "Duration", "MediaDuration"); ok {
		videoInfo.MediaDuration = duration
	}

	if width, ok := getInt(rawData, "VideoWidth", "ImageWidth"); ok {
		videoInfo.Width = width
	}

	if height, ok := getInt(rawData, "VideoHeight", "ImageHeight"); ok {
		videoInfo.Height = height
	}

	if frameRate, ok := getFloat(rawData, "VideoFrameRate", "FrameRate"); ok {
		videoInfo.VideoFrameRate = frameRate
	}

	if bitrate, ok := getString(rawData, "AvgBitrate", "Bitrate"); ok {
		videoInfo.AvgBitrate = bitrate
	}

	if encoder, ok := getString(rawData, "Encoder"); ok {
		videoInfo.Encoder = encoder
	}

	if rotation, ok := getInt(rawData, "Rotation"); ok {
		videoInfo.Rotation = rotation
	}

	if audioFormat, ok := getString(rawData, "AudioFormat"); ok {
		videoInfo.AudioFormat = audioFormat
	}

	if channels, ok := getInt(rawData, "AudioChannels", "Channels"); ok {
		videoInfo.AudioChannels = channels
	}

	if sampleRate, ok := getInt(rawData, "AudioSampleRate", "SampleRate"); ok {
		videoInfo.AudioSampleRate = sampleRate
	}

	if bitsPerSample, ok := getInt(rawData, "AudioBitsPerSample", "BitsPerSample"); ok {
		videoInfo.AudioBitsPerSample = bitsPerSample
	}

	return videoInfo
}

func (et *ExifTool) parseCameraInfo(rawData map[string]interface{}) CameraInfo {
	cameraInfo := CameraInfo{}

	if cameraMake, ok := getString(rawData, "Make"); ok {
		cameraInfo.Make = cameraMake
	}

	if model, ok := getString(rawData, "Model"); ok {
		cameraInfo.Model = model
	}

	if software, ok := getString(rawData, "Software"); ok {
		cameraInfo.Software = software
	}

	if exposureTime, ok := getString(rawData, "ExposureTime"); ok {
		cameraInfo.ExposureTime = exposureTime
	}

	if fNumber, ok := getFloat(rawData, "FNumber"); ok {
		cameraInfo.FNumber = fNumber
	}

	if iso, ok := getInt(rawData, "ISO", "ISOSpeed"); ok {
		cameraInfo.ISO = iso
	}

	if focalLength, ok := getString(rawData, "FocalLength"); ok {
		cameraInfo.FocalLength = focalLength
	}

	if focalLength35mm, ok := getString(rawData, "FocalLengthIn35mmFormat"); ok {
		cameraInfo.FocalLength35mm = focalLength35mm
	}

	if flash, ok := getString(rawData, "Flash"); ok {
		cameraInfo.Flash = flash
	}

	if lightSource, ok := getString(rawData, "LightSource"); ok {
		cameraInfo.LightSource = lightSource
	}

	if exposureMode, ok := getString(rawData, "ExposureMode"); ok {
		cameraInfo.ExposureMode = exposureMode
	}

	if whiteBalance, ok := getString(rawData, "WhiteBalance"); ok {
		cameraInfo.WhiteBalance = whiteBalance
	}

	return cameraInfo
}

func (et *ExifTool) parseLocationInfo(rawData map[string]interface{}) Location {
	location := Location{}

	if latitude, ok := getFloat(rawData, "GPSLatitude"); ok {
		location.Latitude = latitude
	}

	if longitude, ok := getFloat(rawData, "GPSLongitude"); ok {
		location.Longitude = longitude
	}

	return location
}

// Helper functions
func getString(data map[string]interface{}, keys ...string) (string, bool) {
	for _, key := range keys {
		if value, ok := data[key]; ok {
			if str, ok := value.(string); ok {
				return str, true
			}
		}
	}
	return "", false
}

func getInt(data map[string]interface{}, keys ...string) (int, bool) {
	for _, key := range keys {
		if value, ok := data[key]; ok {
			switch v := value.(type) {
			case string:
				if intVal, err := strconv.Atoi(v); err == nil {
					return intVal, true
				}
				// Try parsing float first, then convert to int
				if floatVal, err := strconv.ParseFloat(v, 64); err == nil {
					return int(floatVal), true
				}
			case float64:
				return int(v), true
			case int:
				return v, true
			}
		}
	}
	return 0, false
}

func getFloat(data map[string]interface{}, keys ...string) (float64, bool) {
	for _, key := range keys {
		if value, ok := data[key]; ok {
			switch v := value.(type) {
			case string:
				if floatVal, err := strconv.ParseFloat(v, 64); err == nil {
					return floatVal, true
				}
			case float64:
				return v, true
			case int:
				return float64(v), true
			}
		}
	}
	return 0, false
}
