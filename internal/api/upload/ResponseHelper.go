package upload

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mahdi-cpp/upload-service/internal/exiftool"
)

// ResponseHelper handles standardized API responses
type ResponseHelper struct{}

// NewResponseHelper creates a new response helper instance
func NewResponseHelper() *ResponseHelper {
	return &ResponseHelper{}
}

// Send sends a standardized JSON response
func (rh *ResponseHelper) Send(c *gin.Context, statusCode int, message string, err error, id uuid.UUID) {
	response := Response{
		Message: message,
		ID:      id,
	}

	if err != nil {
		response.Error = err.Error()
	}

	c.JSON(statusCode, response)
}

// SendError sends an error response with the appropriate status code
func (rh *ResponseHelper) SendError(c *gin.Context, statusCode int, message string, err error) {
	rh.Send(c, statusCode, message, err, uuid.Nil)
}

// SendSuccess sends a successful response
func (rh *ResponseHelper) SendSuccess(c *gin.Context, message string, id uuid.UUID) {
	rh.Send(c, http.StatusOK, message, nil, id)
}
func (rh *ResponseHelper) SendSuccessMetadata(c *gin.Context, metadata *exiftool.Metadata) {
	c.JSON(http.StatusOK, metadata)
}

// Initialize the response helper (you can also inject this via dependency injection)
var responseHelper = NewResponseHelper()
