package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/mahdi-cpp/upload-service/internal/helpers"
)

const baseURL = "http://localhost:50103/api/v1/upload/"

func TestChatCreate(t *testing.T) {

	var currentURL = baseURL + "create"

	respBody, err := helpers.MakeRequest(t, "POST", currentURL, nil, nil)
	if err != nil {
		t.Errorf("create request failed: %v", err)
	}

	var r DirectoryRequest
	if err := json.Unmarshal(respBody, &r); err != nil {
		t.Errorf("unmarshaling response: %v", err)
	}

	if r.ID == uuid.Nil {
		t.Errorf("id should not be nil")
	}

	fmt.Println(r.ID)

	currentURL = baseURL + "image"
	filePath := "/app/files/videos/02.jpg" // Replace with a valid path to an image file.

	// A context with a 30-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	body := &UploadRequest{
		Directory: r.ID,
	}

	if err := UploadImage(ctx, currentURL, filePath, body); err != nil {
		t.Errorf("Error uploading image: %v", err)
	}
}

func TestGenerateHash(t *testing.T) {

	startTime := time.Now()

	hash, err := helpers.CreateSHA256Hash("/app/files/videos/01.jpg")
	if err != nil {
		return
	}

	duration := time.Since(startTime)

	fmt.Printf("هش SHA-256 تولید شده: %s\n", hash)
	fmt.Printf("زمان لازم برای هش کردن: %v\n", duration)
}

// UploadImage uploads a file to the specified API endpoint.
// The function takes the file path, the API URL, and a context for cancellation.
// UploadImage uploads a file and a JSON payload to the specified API endpoint.
// The function takes the file path, the API URL, the directory UUID, and a context for cancellation.
func UploadImage(ctx context.Context, apiURL, filePath string, uploadRequest *UploadRequest) error {

	// Open the file to be uploaded.
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	// Create a new multipart writer.
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Marshal the struct into a JSON string.
	jsonData, err := json.Marshal(uploadRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON data: %w", err)
	}

	// Create a new form field for the JSON data.
	jsonPart, err := writer.CreateFormField("data")
	if err != nil {
		return fmt.Errorf("failed to create JSON form field: %w", err)
	}
	_, err = jsonPart.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed to write JSON data to form field: %w", err)
	}

	// Create a form file part for the image.
	filePart, err := writer.CreateFormFile("image", filepath.Base(filePath))
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}

	// Copy the file content into the form file part.
	if _, err := io.Copy(filePart, file); err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	// Close the multipart writer to finalize the body.
	writer.Close()

	// Create the HTTP request.
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, &requestBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set the Content-Type header with the boundary from the writer.
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request.
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body to ensure the connection is reused.
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Warning: failed to read response body: %v", err)
	}

	// Check for a successful response status code.
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server responded with status code %d", resp.StatusCode)
	}

	fmt.Printf("Successfully uploaded image from %s\n", filePath)
	return nil
}
