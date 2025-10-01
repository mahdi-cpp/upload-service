package upload

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/mahdi-cpp/upload-service/internal/config"
	"github.com/mahdi-cpp/upload-service/internal/exiftool"
	"github.com/mahdi-cpp/upload-service/internal/ffmpeg"
	"github.com/mahdi-cpp/upload-service/internal/helpers"
	"github.com/mahdi-cpp/upload-service/internal/thumbnail"
)

func (h *Handler) CreateDirectory(c *gin.Context) {

	// Generate unique workDir
	directoryId, err := helpers.GenerateUUID()
	if err != nil {
		responseHelper.SendError(c, http.StatusInternalServerError, "Failed to generate directory ID", err)
		return
	}

	workDir := filepath.Join(config.UploadDir, directoryId.String())
	if err := os.MkdirAll(workDir, 0755); err != nil {
		responseHelper.SendError(c, http.StatusForbidden, "Failed to create directory", err)
		return
	}

	responseHelper.SendSuccess(c, "successful create directory", directoryId)
	return
}

func (h *Handler) UploadMedia(c *gin.Context) {

	// 1. Extract the JSON payload from the "data" form field.
	jsonData := c.PostForm("metadata")
	if jsonData == "" {
		responseHelper.SendError(c, http.StatusBadRequest, "Missing 'metadata' form field", nil)
		return
	}

	var request Request
	if err := json.Unmarshal([]byte(jsonData), &request); err != nil {
		responseHelper.SendError(c, http.StatusBadRequest, "Invalid JSON data in 'metadata'", err)
		return
	}

	// 2. Access the file from the "media" form field.
	file, err := c.FormFile("media")
	if err != nil {
		responseHelper.SendError(c, http.StatusBadRequest, "No file uploaded", err)
		return
	}

	// Generate unique filename
	mediaID, err := helpers.GenerateUUID()
	if err != nil {
		responseHelper.SendError(c, http.StatusInternalServerError, "Failed to generate media ID", err)
		return
	}

	workDir := filepath.Join(config.UploadDir, request.Directory.String())
	if err := os.MkdirAll(workDir, 0755); err != nil {
		responseHelper.SendError(c, http.StatusInternalServerError, "Failed to create directory", err)
		return
	}

	var metadata *exiftool.Metadata

	// Process media based on type
	if request.IsVideo {
		metadata, err = h.processVideo(c, file, mediaID, workDir)
		if err != nil {
			responseHelper.SendError(c, http.StatusInternalServerError, "Failed to process video", err)
			return
		}
	} else {
		metadata, err = h.processImage(c, file, mediaID, workDir)
		if err != nil {
			responseHelper.SendError(c, http.StatusInternalServerError, "Failed to process image", err)
			return
		}
	}

	responseHelper.SendSuccessMetadata(c, metadata)
}

func (h *Handler) processVideo(c *gin.Context, file *multipart.FileHeader, mediaID uuid.UUID, workDir string) (*exiftool.Metadata, error) {

	originalVideo := filepath.Join(workDir, mediaID.String()+".mp4")
	if err := c.SaveUploadedFile(file, originalVideo); err != nil {
		return nil, fmt.Errorf("save video: %w", err)
	}

	coverFile := filepath.Join(workDir, mediaID.String()+".jpg")
	if err := ffmpeg.ExtractFrame(originalVideo, coverFile); err != nil {
		return nil, fmt.Errorf("extract frame: %w", err)
	}

	sizes := []int{270, 400}
	for _, size := range sizes {
		thumbnailPath := filepath.Join(workDir, mediaID.String())
		if err := thumbnail.ProcessImage2(coverFile, thumbnailPath, size); err != nil {
			return nil, fmt.Errorf("generate thumbnail %d: %w", size, err)
		}
	}

	return h.saveMetadata(originalVideo, mediaID, workDir)
}

func (h *Handler) processImage(c *gin.Context, file *multipart.FileHeader, mediaID uuid.UUID, workDir string) (*exiftool.Metadata, error) {
	original := filepath.Join(workDir, mediaID.String()+".jpg")
	if err := c.SaveUploadedFile(file, original); err != nil {
		return nil, fmt.Errorf("save image: %w", err)
	}

	sizes := []int{270}
	for _, size := range sizes {
		thumbnailPath := filepath.Join(workDir, mediaID.String())
		if err := thumbnail.ProcessImage2(original, thumbnailPath, size); err != nil {
			return nil, fmt.Errorf("generate thumbnail %d: %w", size, err)
		}
	}

	return h.saveMetadata(original, mediaID, workDir)
}

func (h *Handler) saveMetadata(mediaPath string, mediaID uuid.UUID, workDir string) (*exiftool.Metadata, error) {
	exifTool := exiftool.NewExifTool()
	//defer exifTool.Close() // Assuming ExifTool has a Close method for cleanup

	metadata, err := exifTool.GetMetadata(mediaPath)
	if err != nil {
		return nil, fmt.Errorf("get metadata: %w", err)
	}

	//metadataPath := filepath.Join(workDir, mediaID.String()+".json")
	//if err := exifTool.SaveMetadata(metadata, metadataPath); err != nil {
	//	return nil, fmt.Errorf("save metadata: %w", err)
	//}

	//metadata.ID = mediaID

	return metadata, nil
}

//
//func (h *Handler) UploadMedia(c *gin.Context) {
//
//	// 1. Extract the JSON payload from the "data" form field.
//	jsonData := c.PostForm("metadata")
//	if jsonData == "" {
//		fmt.Println("1")
//		c.JSON(http.StatusBadRequest, Response{
//			Message: "Missing 'metadata' form field",
//			Error:   error(nil).Error(),
//		})
//		return
//	}
//
//	var request Request
//	if err := json.Unmarshal([]byte(jsonData), &request); err != nil {
//		fmt.Println("2")
//		c.JSON(http.StatusBadRequest, Response{
//			Message: "Invalid JSON data in 'data",
//			Error:   err.Error(),
//		})
//		return
//	}
//
//	// 2. Access the file from the "media" form field.
//	file, err := c.FormFile("media")
//	if err != nil {
//		fmt.Println("3")
//		c.JSON(http.StatusBadRequest, Response{
//			Message: "No file uploaded",
//			Error:   err.Error(),
//		})
//		return
//	}
//
//	// Check if it's a JPEG
//	//if !isJPEG(file) {
//	//	fmt.Println("not jpeg")
//	//	c.JSON(http.StatusBadRequest, UploadResponse{
//	//		Message: "Only jpg files are allowed",
//	//		Error:   "Invalid file type",
//	//	})
//	//	return
//	//}
//
//	//--- Generate unique filename
//	mediaID, err := helpers.GenerateUUID()
//	if err != nil {
//		fmt.Println("4")
//		c.JSON(http.StatusInternalServerError, Response{
//			Message: "failed to upload image",
//			Error:   err.Error(),
//		})
//		return
//	}
//
//	//--- Save original file
//	workDir := filepath.Join(config.UploadDir, request.Directory.String())
//
//	if request.IsVideo {
//
//		originalVideo := filepath.Join(workDir, mediaID.String()+".mp4")
//		if err := c.SaveUploadedFile(file, originalVideo); err != nil {
//			fmt.Println("5")
//			c.JSON(http.StatusInternalServerError, Response{
//				Message: "Failed to save image",
//				Error:   err.Error(),
//			})
//			return
//		}
//
//		coverFile := filepath.Join(workDir, mediaID.String()+".jpg")
//		if err := ffmpeg.ExtractFrame(originalVideo, coverFile); err != nil {
//			fmt.Println("6")
//			c.JSON(http.StatusInternalServerError, Response{
//				Message: "Failed to save jpeg",
//				Error:   err.Error(),
//			})
//			return
//		}
//
//		//--- Thumbnail
//		var sizes = []int{200, 400}
//		for _, size := range sizes {
//			thumbnailPath := filepath.Join(workDir, mediaID.String())
//			if err := thumbnail.ProcessImage2(coverFile, thumbnailPath, size); err != nil {
//				fmt.Println("7")
//				c.JSON(http.StatusInternalServerError, Response{
//					Message: "Failed to save video",
//					Error:   err.Error(),
//				})
//				return
//			}
//		}
//
//		//--- Metadata
//		exifTool := exiftool.NewExifTool()
//		imageMetadata, err := exifTool.GetMetadata(originalVideo)
//		if err != nil {
//			fmt.Println("8")
//			c.JSON(http.StatusInternalServerError, Response{
//				Message: "failed to create metadata",
//				Error:   err.Error(),
//			})
//			return
//		}
//		err = exifTool.SaveMetadata(imageMetadata, filepath.Join(workDir, mediaID.String()+".json"))
//		if err != nil {
//			fmt.Println("9")
//			c.JSON(http.StatusInternalServerError, Response{
//				Message: "failed to save metadata to file",
//				Error:   err.Error(),
//			})
//			return
//		}
//
//	} else {
//
//		original := filepath.Join(workDir, mediaID.String()+".jpg")
//
//		if err := c.SaveUploadedFile(file, original); err != nil {
//			fmt.Println("10")
//			c.JSON(http.StatusInternalServerError, Response{
//				Message: "failed to save metadata to file",
//				Error:   err.Error(),
//			})
//			return
//		}
//
//		_, err = os.ReadFile(original)
//		if err != nil {
//			fmt.Println("11")
//			c.JSON(http.StatusInternalServerError, Response{
//				Message: "failed to create thumbnail",
//				Error:   err.Error(),
//			})
//			return
//		}
//
//		//--- Thumbnail
//		var sizes = []int{200}
//		for _, size := range sizes {
//			thumbnailPath := filepath.Join(workDir, mediaID.String())
//			if err := thumbnail.ProcessImage2(original, thumbnailPath, size); err != nil {
//				fmt.Println("12")
//				c.JSON(http.StatusInternalServerError, Response{
//					Message: "failed to create thumbnail",
//					Error:   err.Error(),
//				})
//				return
//			}
//		}
//
//		//--- Metadata
//		exifTool := exiftool.NewExifTool()
//		imageMetadata, err := exifTool.GetMetadata(original)
//		if err != nil {
//			fmt.Println("13")
//			c.JSON(http.StatusInternalServerError, Response{
//				Message: "failed to create metadata",
//				Error:   err.Error(),
//			})
//			return
//		}
//		err = exifTool.SaveMetadata(imageMetadata, filepath.Join(workDir, mediaID.String()+".json"))
//		if err != nil {
//			fmt.Println("14")
//			c.JSON(http.StatusInternalServerError, Response{
//				Message: "failed to save metadata to file",
//				Error:   err.Error(),
//			})
//			return
//		}
//	}
//
//	c.JSON(http.StatusOK, Response{
//		Message: "file uploaded successfully",
//		ID:      mediaID,
//	})
//}
//
//func (h *Handler) ListFiles(c *gin.Context) {
//
//	files, err := getJPEGFiles(h.UploadDir)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"error": "Failed to list files",
//		})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"files": files,
//	})
//}
