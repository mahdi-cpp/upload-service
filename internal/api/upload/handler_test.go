package upload

import (
	"bytes"
	"context"
	"fmt"
	"io"
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

// httpClient is a shared instance of the HTTP client for efficiency.
var httpClient = &http.Client{Timeout: 30 * time.Second}

func TestUploadMedias(t *testing.T) {

	var apiURL = baseURL + "create"

	respBody, err := helpers.MakeRequest(t, "POST", apiURL, nil, nil)
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
	apiURL = baseURL + "media"

	// A context with a 30-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//--- upload
	resp, err1 := uploadImage(ctx, httpClient, apiURL, r.ID, "/app/tmp/test.heic")
	if err1 != nil {
		t.Errorf("%v", err)
	}
	if resp == nil || resp.ID == uuid.Nil {
		t.Fatal("Expected a non-nil response, but got nil")
	}
	fmt.Println(resp.ID)

	//--- upload
	resp, err = uploadVideo(ctx, httpClient, apiURL, r.ID, "/app/tmp/test.mp4")
	if err != nil {
		t.Errorf("%v", err)
	}
	if resp == nil || resp.ID == uuid.Nil {
		t.Fatal("Expected a non-nil response, but got nil")
	}
	fmt.Println(resp.ID)
}

// uploadImage handles the preparation and upload of a single image file.
func uploadImage(ctx context.Context, client *http.Client, apiURL string, directoryID uuid.UUID, filePath string) (*Response, error) {

	// Create the Request struct with the necessary metadata.
	uploadReq := &Request{
		Directory: directoryID,
		IsVideo:   false,
	}

	// Call the MediaUpload function.
	resp, err := MediaUpload(ctx, client, apiURL, filePath, uploadReq)
	if err != nil {
		return nil, fmt.Errorf("media upload failed: %w", err)
	}

	return resp, nil
}

func uploadVideo(ctx context.Context, client *http.Client, apiURL string, directoryID uuid.UUID, filePath string) (*Response, error) {

	// Create the Request struct with the necessary metadata.
	uploadReq := &Request{
		Directory: directoryID,
		IsVideo:   true,
	}

	// Call the MediaUpload function.
	resp, err := MediaUpload(ctx, client, apiURL, filePath, uploadReq)
	if err != nil {
		return nil, fmt.Errorf("media upload failed: %w", err)
	}

	return resp, nil
}

// MediaUpload uploads a file and a JSON payload to the specified API endpoint.
// The function takes the HTTP client, file path, API URL, and a context for cancellation.
func MediaUpload(ctx context.Context, client *http.Client, apiURL, filePath string, uploadRequest *Request) (*Response, error) {

	// Open the file to be uploaded.
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	// Create a new multipart writer.
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Marshal the struct into a JSON string.
	jsonData, err := json.Marshal(uploadRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON data: %w", err)
	}

	// Create a new form field for the JSON data.
	if err := writer.WriteField("metadata", string(jsonData)); err != nil {
		return nil, fmt.Errorf("failed to write JSON data to form field: %w", err)
	}

	// Create a form file part for the media.
	filePart, err := writer.CreateFormFile("media", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	// Copy the file content into the form file part.
	if _, err := io.Copy(filePart, file); err != nil {
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}

	// Close the multipart writer to finalize the body.
	writer.Close()

	// Create the HTTP request.
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, &requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set the Content-Type header with the boundary from the writer.
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request.
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check for a successful response status code.
	if resp.StatusCode != http.StatusOK {
		// Read and log the server's error message.
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("server responded with status code %d, but could not read error body: %w", resp.StatusCode, err)
		}
		return nil, fmt.Errorf("server responded with status code %d and body: %s", resp.StatusCode, string(bodyBytes))
	}

	// Decode the JSON response body.
	var serverResponse Response
	if err := json.NewDecoder(resp.Body).Decode(&serverResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	fmt.Printf("Successfully uploaded image from %s\n", filePath)
	return &serverResponse, nil
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
