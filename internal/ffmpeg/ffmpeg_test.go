package ffmpeg

import (
	"testing"
)

func TestExtractFrame(t *testing.T) {

	inputVideo := "/app/tmp/test.mp4"
	outputImage := "/app/tmp/video_cover5.jpg"

	// Call the function with the desired file paths.
	if err := ExtractFrame(inputVideo, outputImage); err != nil {
		t.Errorf("Failed to extract frame: %v", err)
	}
}
