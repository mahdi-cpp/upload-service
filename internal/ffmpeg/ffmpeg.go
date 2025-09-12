package ffmpeg

import (
	"log"
	"os"
	"os/exec"
)

// ExtractFrame extracts a single frame from an input video at a specified timestamp
// and saves it to the given output path.
//
// The command used is:
// ffmpeg -ss 00:01:30 -i <inputPath> -vframes 1 -q:v 2 -vf "scale=1280:-1" <outputPath>
func ExtractFrame(inputPath, outputPath string) error {
	// First, check if the ffmpeg executable is available in the system's PATH.
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		log.Printf("Error: ffmpeg not found. Please ensure it is installed and available in your system's PATH.")
		return err
	}

	// The command and its arguments are defined as a slice of strings.
	// This is the standard and safest way to pass arguments to an external command.
	args := []string{
		"-ss", "00:00:5",
		"-i", inputPath,
		"-vframes", "1",
		"-q:v", "2",
		"-vf", "scale=1280:-1",
		outputPath,
	}

	// Create a new command object.
	cmd := exec.Command(ffmpegPath, args...)

	// Optional: Set a specific working directory for the command.
	// cmd.Dir = "/path/to/your/directory"

	// Set the command's output to the standard logger for debugging purposes.
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("Executing command: %s %s", cmd.Path, args)

	// Run the command and return any error.
	if err := cmd.Run(); err != nil {
		log.Printf("Error executing ffmpeg command: %v", err)
		return err
	}

	log.Printf("Successfully extracted frame to %s", outputPath)
	return nil
}
