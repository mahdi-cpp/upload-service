package main

import (
	"fmt"
	"time"

	"github.com/mahdi-cpp/photos-api/tools/thumbnail_v2"
)

func main() {

	start := time.Now()

	err := thumbnail_v2.Create("018f3a8b-1b32-729a-f7e5-5467c1b2d3e4")
	if err != nil {
		return
	}

	elapsed := time.Since(start)
	fmt.Printf("Thumbnail creation took %s to run.\n", elapsed)

	//metadata, err := exiftool_v1.Start("/app/iris/com.iris.asset/users/assets/0198c111-0fe1-7e2d-b38b-62b1a1d89907.jpg")
	//if err != nil {
	//	return
	//}
	//
	//fmt.Println("\nExtracted Metadata:")
	//fmt.Printf("  File Size: %s\n", metadata.FileSize)
	//fmt.Printf("  File Type: %s\n", metadata.FileType)
	//fmt.Printf("  Make: %s\n", metadata.Make)
	//fmt.Printf("  Model: %s\n", metadata.Model)
	//fmt.Printf("  Orientation: %s\n", metadata.Orientation)
	//fmt.Printf("  Create Date: %s\n", metadata.CreateDate)
}
