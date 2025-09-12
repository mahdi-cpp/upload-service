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
	"github.com/mahdi-cpp/upload-service/internal/helpers"
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
	//Hash      string    `json:"hash"`
}

// UploadImage api/v1/upload/image" [POST]
//func (h *UploadHandler) UploadImage(c *gin.Context) {
//
//	var request UploadRequest
//	if err := c.ShouldBindJSON(&request); err != nil {
//		helpers.AbortWithRequestInvalid(c)
//		return
//	}
//
//	// Single image upload
//	file, err := c.FormFile("file")
//	if err != nil {
//		fmt.Println(err)
//		c.JSON(http.StatusBadRequest, UploadResponse{
//			Message: "No file uploaded",
//			Error:   err.Error(),
//		})
//		return
//	}
//
//	//hash, err := helpers.CreateSHA256Hash("/app/files/videos/01.jpg")
//	//if err != nil {
//	//	return
//	//}
//	//
//	//if hash == "" {
//	//}
//
//	// Check if it's a JPEG
//	if !isJPEG(file) {
//		fmt.Println("not jpeg")
//		c.JSON(http.StatusBadRequest, UploadResponse{
//			Message: "Only jpg files are allowed",
//			Error:   "Invalid file type",
//		})
//		return
//	}
//
//	// Generate unique filename
//	imageID, err := generateId()
//	if err != nil {
//		fmt.Println(err)
//		c.JSON(http.StatusBadRequest, gin.H{
//			"message": "failed to upload image",
//			"error":   err.Error(),
//		})
//	}
//	fileDirectory := filepath.Join(h.UploadDir, imageID.String()+".jpg")
//
//	vips.Startup(nil)
//	defer vips.Shutdown()
//
//	// Save the file
//	if err := c.SaveUploadedFile(file, fileDirectory); err != nil {
//		fmt.Printf("failed to upload image: %v", err)
//		c.JSON(http.StatusInternalServerError, UploadResponse{Message: "Failed to save file", Error: err.Error()})
//		return
//	}
//
//	//if err := thumbnail.CreateSingleThumbnail(fileDirectory, imageID.String()+".jpg"); err != nil {
//	//	c.JSON(http.StatusInternalServerError, UploadResponse{Message: "Failed to create thumbnail file", Error: err.Error()})
//	//	log.Fatalf("An error occurred during thumbnail creation: %v", err)
//	//}
//
//	c.JSON(http.StatusOK, UploadResponse{
//		Message:  "File uploaded successfully",
//		Filename: imageID.String(),
//		Size:     file.Size,
//	})
//}

// generateID generates a new UUID v7.
func generateID() (uuid.UUID, error) {
	return uuid.NewV7()
}

// UploadImage handles the image upload via multipart/form-data.
func (h *UploadHandler) UploadImage(c *gin.Context) {

	// 1. Extract the JSON payload from the "data" form field.
	jsonData := c.PostForm("data")
	if jsonData == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'data' form field"})
		return
	}

	var request UploadRequest
	if err := json.Unmarshal([]byte(jsonData), &request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data in 'data' field"})
		return
	}

	// 2. Access the file from the "image" form field.
	file, err := c.FormFile("image")
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

	// Generate unique filename
	imageID, err := generateID()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to upload image",
			"error":   err.Error(),
		})
		return
	}

	fileDirectory := filepath.Join(config.UploadDir, request.Directory.String(), imageID.String()+".jpg")

	// Save the file
	if err := c.SaveUploadedFile(file, fileDirectory); err != nil {
		fmt.Printf("failed to save image: %v", err)
		c.JSON(http.StatusInternalServerError, UploadResponse{Message: "Failed to save file", Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, UploadResponse{
		Message:  "File uploaded successfully",
		Filename: imageID.String(),
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
