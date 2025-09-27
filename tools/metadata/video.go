package metadata

import (
	"encoding/json"
	"fmt"
)

// VideoMetadata struct for holding important video information
type VideoMetadata struct {
	FileName        string  `json:"fileName"`        // File name
	FileSize        string  `json:"fileSize"`        // File size
	MIMEType        string  `json:"mimeType"`        // File MIME type
	Duration        string  `json:"duration"`        // Video duration
	Width           int     `json:"width"`           // Video frame width in pixels
	Height          int     `json:"height"`          // Video frame height in pixels
	VideoFrameRate  float64 `json:"videoFrameRate"`  // Video frame rate
	AvgBitrate      string  `json:"avgBitrate"`      // Average bitrate (quality and data volume)
	AudioChannels   int     `json:"audioChannels"`   // Number of audio channels (e.g., 2 for stereo)
	AudioSampleRate int     `json:"audioSampleRate"` // Audio sample rate
	Encoder         string  `json:"encoder"`         // Video encoding software
}

// FileDetails struct for basic file system information
type FileDetails struct {
	SourceFile          string  `json:"sourceFile"`
	ExifToolVersion     float64 `json:"exifToolVersion"`
	FileModifyDate      string  `json:"fileModifyDate"`
	FileAccessDate      string  `json:"fileAccessDate"`
	FileInodeChangeDate string  `json:"fileInodeChangeDate"`
	FilePermissions     string  `json:"filePermissions"`
	FileTypeExtension   string  `json:"fileTypeExtension"`
}

// GeneralVideoDetails struct for overall video container properties
type GeneralVideoDetails struct {
	MajorBrand         string   `json:"majorBrand"`
	MinorVersion       string   `json:"minorVersion"`
	CompatibleBrands   []string `json:"compatibleBrands"`
	MediaDataSize      int      `json:"mediaDataSize"`
	MediaDataOffset    int      `json:"mediaDataOffset"`
	MovieHeaderVersion int      `json:"movieHeaderVersion"`
	CreateDate         string   `json:"createDate"`
	ModifyDate         string   `json:"modifyDate"`
	TimeScale          int      `json:"timeScale"`
	PreferredRate      int      `json:"preferredRate"`
	PreferredVolume    string   `json:"preferredVolume"`
	PreviewTime        string   `json:"previewTime"`
	PreviewDuration    string   `json:"previewDuration"`
	PosterTime         string   `json:"posterTime"`
	SelectionTime      string   `json:"selectionTime"`
	SelectionDuration  string   `json:"selectionDuration"`
	CurrentTime        string   `json:"currentTime"`
	NextTrackID        int      `json:"nextTrackID"`
	Rotation           int      `json:"rotation"`
}

// VideoTrackDetails struct for specific video track properties
type VideoTrackDetails struct {
	TrackHeaderVersion int    `json:"trackHeaderVersion"`
	TrackCreateDate    string `json:"trackCreateDate"`
	TrackModifyDate    string `json:"trackModifyDate"`
	TrackID            int    `json:"trackID"`
	TrackDuration      string `json:"trackDuration"`
	TrackLayer         int    `json:"trackLayer"`
	TrackVolume        string `json:"trackVolume"`
	GraphicsMode       string `json:"graphicsMode"`
	OpColor            string `json:"opColor"`
	CompressorID       string `json:"compressorID"`
	SourceImageWidth   int    `json:"sourceImageWidth"`
	SourceImageHeight  int    `json:"sourceImageHeight"`
	XResolution        int    `json:"xResolution"`
	YResolution        int    `json:"yResolution"`
	BitDepth           int    `json:"bitDepth"`
	MatrixStructure    string `json:"matrixStructure"`
}

// AudioTrackDetails struct for specific audio track properties
type AudioTrackDetails struct {
	MediaHeaderVersion int    `json:"mediaHeaderVersion"`
	MediaCreateDate    string `json:"mediaCreateDate"`
	MediaModifyDate    string `json:"mediaModifyDate"`
	MediaTimeScale     int    `json:"mediaTimeScale"`
	MediaDuration      string `json:"mediaDuration"`
	MediaLanguageCode  string `json:"mediaLanguageCode"`
	HandlerDescription string `json:"handlerDescription"`
	Balance            int    `json:"balance"`
	AudioFormat        string `json:"audioFormat"`
	AudioBitsPerSample int    `json:"audioBitsPerSample"`
}

// HandlerDetails struct for handler information
type HandlerDetails struct {
	HandlerType     string `json:"handlerType"`
	HandlerVendorID string `json:"handlerVendorID"`
}

// FullVideoMetadata struct that combines all video metadata categories
type FullVideoMetadata struct {
	Core           VideoMetadata       `json:"core"`
	FileDetails    FileDetails         `json:"fileDetails"`
	GeneralDetails GeneralVideoDetails `json:"generalDetails"`
	VideoTrack     VideoTrackDetails   `json:"videoTrack"`
	AudioTrack     AudioTrackDetails   `json:"audioTrack"`
	HandlerDetails HandlerDetails      `json:"handlerDetails"`
}

func main() {
	// An example of video information
	fullMetadata := FullVideoMetadata{
		Core: VideoMetadata{
			FileName:        "output8.mp4",
			FileSize:        "16 MB",
			MIMEType:        "video/mp4",
			Duration:        "11.01 s",
			Width:           1280,
			Height:          720,
			VideoFrameRate:  30.001,
			AvgBitrate:      "12.3 Mbps",
			AudioChannels:   2,
			AudioSampleRate: 48000,
			Encoder:         "Lavf58.29.100",
		},
		FileDetails: FileDetails{
			SourceFile:          "output8.mp4",
			ExifToolVersion:     11.88,
			FileModifyDate:      "2024:06:14 03:04:49+03:30",
			FileAccessDate:      "2025:08:24 01:00:56+03:30",
			FileInodeChangeDate: "2025:08:24 01:00:41+03:30",
			FilePermissions:     "rwxrwxrwx",
			FileTypeExtension:   "mp4",
		},
		GeneralDetails: GeneralVideoDetails{
			MajorBrand:         "MP4 Base Media v1 [IS0 14496-12:2003]",
			MinorVersion:       "0.2.0",
			CompatibleBrands:   []string{"isom", "iso2", "avc1", "mp41"},
			MediaDataSize:      16857304,
			MediaDataOffset:    48,
			MovieHeaderVersion: 0,
			CreateDate:         "0000:00:00 00:00:00",
			ModifyDate:         "0000:00:00 00:00:00",
			TimeScale:          1000,
			PreferredRate:      1,
			PreferredVolume:    "100.00%",
			PreviewTime:        "0 s",
			PreviewDuration:    "0 s",
			PosterTime:         "0 s",
			SelectionTime:      "0 s",
			SelectionDuration:  "0 s",
			CurrentTime:        "0 s",
			NextTrackID:        3,
			Rotation:           0,
		},
		VideoTrack: VideoTrackDetails{
			TrackHeaderVersion: 0,
			TrackCreateDate:    "0000:00:00 00:00:00",
			TrackModifyDate:    "0000:00:00 00:00:00",
			TrackID:            1,
			TrackDuration:      "11.00 s",
			TrackLayer:         0,
			TrackVolume:        "0.00%",
			GraphicsMode:       "srcCopy",
			OpColor:            "0 0 0",
			CompressorID:       "avc1",
			SourceImageWidth:   1280,
			SourceImageHeight:  720,
			XResolution:        72,
			YResolution:        72,
			BitDepth:           24,
			MatrixStructure:    "1 0 0 0 1 0 0 0 1",
		},
		AudioTrack: AudioTrackDetails{
			MediaHeaderVersion: 0,
			MediaCreateDate:    "0000:00:00 00:00:00",
			MediaModifyDate:    "0000:00:00 00:00:00",
			MediaTimeScale:     48000,
			MediaDuration:      "11.01 s",
			MediaLanguageCode:  "eng",
			HandlerDescription: "SoundHandle",
			Balance:            0,
			AudioFormat:        "mp4a",
			AudioBitsPerSample: 16,
		},
		HandlerDetails: HandlerDetails{
			HandlerType:     "Metadata",
			HandlerVendorID: "Apple",
		},
	}

	// Convert the struct to JSON
	jsonData, err := json.MarshalIndent(fullMetadata, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}

	// Print the generated JSON
	fmt.Println(string(jsonData))
}
