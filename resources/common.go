// Package service provides the business logic for all routes/endpoints
package resources

import (
	"github.com/gin-gonic/gin"
	"laerte.sociablelabs.com/common"
	"laerte.sociablelabs.com/common/messagecode"
)

// sendError a common function for logging & sending errors
func sendError(err *error, code string, c *gin.Context) {
	r := messagecode.DefMessage.Get(code)
	c.JSON(r["HTTP_CODE"].(int), &common.HttpResp{"status": r["MSG"].(string)})
	return
}

// sendResponse a common function for sending a response
func sendResponse(data interface{}, code string, c *gin.Context) {
	r := messagecode.DefMessage.Get(code)

	var body interface{}
	if data != nil {
		body = data
	} else {
		body = &common.HttpResp{"status": r["MSG"]}
	}
	c.JSON(r["HTTP_CODE"].(int), body)
	return
}
