package utils

import (
	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success      bool        `json:"success"`
	Message      string      `json:"message"`
	Data         interface{} `json:"data"`
	Error        bool        `json:"error"`
	ErrorMessage *string     `json:"error_message"`
	ErrorCode    *string     `json:"error_code"`
	Errors       interface{} `json:"errors,omitempty"`
}

func Success(c *gin.Context, data interface{}, message string, status int) {
	status = normalizeStatus(status)
	c.JSON(status, APIResponse{
		Success:      true,
		Message:      message,
		Data:         data,
		Error:        false,
		ErrorMessage: nil,
		ErrorCode:    nil,
	})
}

func Error(c *gin.Context, errorMessage string, errorCode string, status int, data interface{}) {
	status = normalizeStatus(status)
	errCode := &errorCode
	errMessage := &errorMessage
	c.JSON(status, APIResponse{
		Success:      false,
		Message:      errorMessage,
		Data:         data,
		Error:        true,
		ErrorMessage: errMessage,
		ErrorCode:    errCode,
		Errors: []gin.H{
			{
				"code":   errorCode,
				"detail": errorMessage,
				"status": status,
			},
		},
	})
}

func normalizeStatus(status int) int {
	if status == 502 {
		return 503
	}
	return status
}
