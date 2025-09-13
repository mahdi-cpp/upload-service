package upload

import "github.com/google/uuid"

type Handler struct {
	UploadDir string
}

type Response struct {
	Message  string    `json:"message,omitempty"`
	Filename string    `json:"filename,omitempty"`
	ID       uuid.UUID `json:"id,omitempty"`
	Error    string    `json:"error,omitempty"`
}

type DirectoryRequest struct {
	ID      uuid.UUID `json:"id"`
	Message string    `json:"message"`
	Errors  string    `json:"errors,omitempty"`
}

type Request struct {
	Directory uuid.UUID `json:"directory"`
	IsVideo   bool      `json:"isVideo"`
	//Hash      string    `json:"hash"`
}
