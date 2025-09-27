package exiftool

import (
	"time"
)

type Metadata struct {
	FileInfo FileInfo   `json:"fileInfo"`
	Image    ImageInfo  `json:"image"`
	Video    VideoInfo  `json:"video"`
	Camera   CameraInfo `json:"camera"`
	Location Location   `json:"location"`
}

type FileInfo struct {
	BaseURL  string `json:"baseURL"`
	FileSize string `json:"fileSize"`
	FileType string `json:"fileType"`
	MimeType string `json:"mimeType"`
}

type ImageInfo struct {
	Width           int    `json:"width,omitempty"`
	Height          int    `json:"height,omitempty"`
	Megapixels      int    `json:"megapixels,omitempty"`
	Orientation     string `json:"orientation,omitempty"`
	ColorSpace      string `json:"colorSpace,omitempty"`
	EncodingProcess string `json:"encodingProcess,omitempty"`
}
type CameraInfo struct {
	Make             string    `json:"make,omitempty"`
	Model            string    `json:"model,omitempty"`
	Software         string    `json:"software,omitempty"`
	DateTimeOriginal time.Time `json:"dateTimeOriginal,omitempty"`
	ExposureTime     string    `json:"exposureTime,omitempty"`
	FNumber          float64   `json:"fNumber,omitempty"` // دیافراگم معمولاً float است
	ISO              int       `json:"iso,omitempty"`     // ISO معمولاً عدد صحیح است
	FocalLength      string    `json:"focalLength,omitempty"`
	FocalLength35mm  string    `json:"focalLength35mm,omitempty"`
	Flash            string    `json:"flash,omitempty"`
	LightSource      string    `json:"lightSource,omitempty"`
	ExposureMode     string    `json:"exposureMode,omitempty"`
	WhiteBalance     string    `json:"whiteBalance,omitempty"`
}

type VideoInfo struct {
	MediaDuration      string  `json:"mediaDuration,omitempty"`  // Video duration
	Width              int     `json:"width,omitempty"`          // Video frame width in pixels
	Height             int     `json:"height,omitempty"`         // Video frame height in pixels
	VideoFrameRate     float64 `json:"videoFrameRate,omitempty"` // Video frame rate
	AvgBitrate         string  `json:"avgBitrate,omitempty"`     // Average bitrate (quality and data volume)
	Encoder            string  `json:"encoder,omitempty"`        // Video encoding software
	Rotation           int     `json:"rotation,omitempty"`
	AudioFormat        string  `json:"audioFormat,omitempty"`
	AudioChannels      int     `json:"audioChannels,omitempty"`   // Number of audio channels (e.g., 2 for stereo)
	AudioSampleRate    int     `json:"audioSampleRate,omitempty"` // Audio sample rate
	AudioBitsPerSample int     `json:"audioBitsPerSample,omitempty"`
}

type Location struct {
	Latitude   float64 `json:"location,omitempty"`
	Longitude  float64 `json:"longitude,omitempty"`
	Country    string  `json:"country,omitempty"`
	Province   string  `json:"province,omitempty"`
	County     string  `json:"county,omitempty"`
	City       string  `json:"city,omitempty"`
	Village    string  `json:"village,omitempty"`
	Electronic int     `json:"electronic,omitempty"`
}
