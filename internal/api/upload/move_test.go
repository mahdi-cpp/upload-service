package upload

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestMoveFile(t *testing.T) {

	file := "019943ad-681f-72b5-b371-29e7cdd7b63e.jpg"

	start := time.Now()

	sourceFile := "/app/iris/services/uploads/019943ad-6617-7911-afff-c95153bdad6c/" + file
	destinationPath := "/app/iris/com.iris.messages/chats/018f3a8b-1b32-7295-a2c7-87654b4d4567/assets/" + file

	err := os.Rename(sourceFile, destinationPath)
	if err != nil {
		fmt.Println("Error moving file:", err)
		return
	}

	// Get the duration since start
	duration := time.Since(start)
	// Convert the duration to milliseconds
	durationInMs := duration.Milliseconds()
	fmt.Printf("The operation took %d milliseconds\n", durationInMs)

}
