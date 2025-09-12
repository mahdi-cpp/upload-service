package handler

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/mahdi-cpp/upload-service/internal/config"
	"github.com/mahdi-cpp/upload-service/internal/exiftool"
	"github.com/mahdi-cpp/upload-service/internal/ffmpeg"
	"github.com/mahdi-cpp/upload-service/internal/helpers"
	"github.com/mahdi-cpp/upload-service/internal/thumbnail"
)

//https://chat.deepseek.com/a/chat/s/913cf162-1ad1-4857-8048-2990d3c959a4

type UploadHandler struct {
	UploadDir string
}

type UploadResponse struct {
	Message  string `json:"message"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	URL      string `json:"url"`
	Error    string `json:"error,omitempty"`
}

type DirectoryRequest struct {
	ID      uuid.UUID `json:"id"`
	Message string    `json:"message"`
	Errors  string    `json:"errors,omitempty"`
}

// CreateDirectory create directory for an entity
func (h *UploadHandler) CreateDirectory(c *gin.Context) {

	//var request DirectoryRequest
	//if err := c.ShouldBindJSON(&request); err != nil {
	//	helpers.AbortWithRequestInvalid(c)
	//	return
	//}

	directoryId, err := uuid.NewV7()
	if err != nil {
		c.JSON(http.StatusForbidden, UploadResponse{
			Message: "failed to create directory",
			Error:   err.Error(),
		})
		return
	}

	if err := helpers.CreateDirectory(filepath.Join(config.UploadDir, directoryId.String()), 0755); err != nil {
		c.JSON(http.StatusForbidden, UploadResponse{
			Message: "failed to create directory",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, DirectoryRequest{
		ID:      directoryId,
		Message: "successful create directory",
		Errors:  "",
	})
	return
}

type UploadRequest struct {
	Directory uuid.UUID `json:"directory"`
	IsVideo   bool      `json:"isVideo"`
	//Hash      string    `json:"hash"`
}

// UploadMedia handles the image upload via multipart/form-data.
func (h *UploadHandler) UploadMedia(c *gin.Context) {

	// 1. Extract the JSON payload from the "data" form field.
	jsonData := c.PostForm("metadata")
	if jsonData == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'data' form field"})
		return
	}

	var request UploadRequest
	if err := json.Unmarshal([]byte(jsonData), &request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data in 'data' field"})
		return
	}

	// 2. Access the file from the "media" form field.
	file, err := c.FormFile("media")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, UploadResponse{
			Message: "No file uploaded",
			Error:   err.Error(),
		})
		return
	}

	// Check if it's a JPEG
	//if !isJPEG(file) {
	//	fmt.Println("not jpeg")
	//	c.JSON(http.StatusBadRequest, UploadResponse{
	//		Message: "Only jpg files are allowed",
	//		Error:   "Invalid file type",
	//	})
	//	return
	//}

	//--- Generate unique filename
	mediaID, err := generateID()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to upload image",
			"error":   err.Error(),
		})
		return
	}

	//--- Save original file
	workDir := filepath.Join(config.UploadDir, request.Directory.String())

	if request.IsVideo {

		originalVideo := filepath.Join(workDir, mediaID.String()+".mp4")
		if err := c.SaveUploadedFile(file, originalVideo); err != nil {
			fmt.Printf("failed to save video: %v", err)
			c.JSON(http.StatusInternalServerError, UploadResponse{Message: "Failed to save image", Error: err.Error()})
			return
		}

		coverFile := filepath.Join(workDir, mediaID.String()+".jpg")
		if err := ffmpeg.ExtractFrame(originalVideo, coverFile); err != nil {
			c.JSON(http.StatusInternalServerError, UploadResponse{Message: "Failed to save video", Error: err.Error()})
		}

		//--- Thumbnail
		var sizes = []int{200, 400}
		for _, size := range sizes {
			thumbnailPath := filepath.Join(workDir, mediaID.String())
			if err := thumbnail.ProcessImage2(coverFile, thumbnailPath, size); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "failed to create thumbnail",
				})
			}
		}

		//--- Metadata
		exifTool := exiftool.NewExifTool()
		imageMetadata, err := exifTool.GetMetadata(originalVideo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to create metadata",
			})
		}
		err = exifTool.WriteItemToDisk(imageMetadata, filepath.Join(workDir, mediaID.String()+".json"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to save metadata to file",
			})
		}

	} else {

		original := filepath.Join(workDir, mediaID.String()+".jpg")

		if err := c.SaveUploadedFile(file, original); err != nil {
			fmt.Printf("failed to save image: %v", err)
			c.JSON(http.StatusInternalServerError, UploadResponse{Message: "Failed to save image", Error: err.Error()})
			return
		}

		_, err = os.ReadFile(original)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to create thumbnail",
			})
			return
		}

		//--- Thumbnail
		var sizes = []int{200}
		for _, size := range sizes {
			thumbnailPath := filepath.Join(workDir, mediaID.String())

			if err := thumbnail.ProcessImage2(original, thumbnailPath, size); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "failed to create thumbnail",
				})
			}
		}

		//--- Metadata
		exifTool := exiftool.NewExifTool()
		imageMetadata, err := exifTool.GetMetadata(original)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to create metadata",
			})
		}
		err = exifTool.WriteItemToDisk(imageMetadata, filepath.Join(workDir, mediaID.String()+".json"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to save metadata to file",
			})
		}
	}

	c.JSON(http.StatusOK, UploadResponse{
		Message:  "File uploaded successfully",
		Filename: mediaID.String(),
		Size:     file.Size,
	})
}

func (h *UploadHandler) UploadMultiple(c *gin.Context) {

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, UploadResponse{
			Message: "Failed to parse form",
			Error:   err.Error(),
		})
		return
	}

	files := form.File["files"]
	var responses []UploadResponse
	var errors []string

	for _, file := range files {
		// Check if it's a JPEG
		if !isJPEG(file) {
			errors = append(errors, fmt.Sprintf("%s: Not a JPEG file", file.Filename))
			continue
		}

		// Generate unique filename
		uniqueName, err := generateId()
		if err != nil {
			return
		}
		dst := filepath.Join(h.UploadDir, uniqueName.String())

		// Save the file
		if err := c.SaveUploadedFile(file, dst); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %s", file.Filename, err.Error()))
			continue
		}

		responses = append(responses, UploadResponse{
			Message:  "File uploaded successfully",
			Filename: uniqueName.String(),
			Size:     file.Size,
		})
	}

	if len(errors) > 0 {
		c.JSON(http.StatusPartialContent, gin.H{
			"message": "Some files failed to upload",
			"uploads": responses,
			"errors":  errors,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All files uploaded successfully",
		"uploads": responses,
	})
}

func (h *UploadHandler) ListFiles(c *gin.Context) {
	files, err := getJPEGFiles(h.UploadDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list files",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"files": files,
	})
}

// Helper functions
func isJPEG(file *multipart.FileHeader) bool {
	// Check content type
	contentType := file.Header.Get("Content-Type")
	if contentType != "image/jpeg" {
		return false
	}

	// Also check file extension for extra safety
	ext := strings.ToLower(filepath.Ext(file.Filename))
	return ext == ".jpg" || ext == ".jpeg"
}

func generateId() (uuid.UUID, error) {

	u7, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, fmt.Errorf("error generating UUIDv7: %w", err)
	}

	if u7 == uuid.Nil {
		return uuid.Nil, fmt.Errorf("error generating UUIDv7: ")
	}

	return u7, nil
}

func getJPEGFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			if ext == ".jpg" || ext == ".jpeg" {
				// Return relative path
				rel, err := filepath.Rel(dir, path)
				if err == nil {
					files = append(files, rel)
				}
			}
		}

		return nil
	})

	return files, err
}

// generateID generates a new UUID v7.
func generateID() (uuid.UUID, error) {
	return uuid.NewV7()
}
