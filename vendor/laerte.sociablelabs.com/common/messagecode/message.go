// Package messagecode provides a centralized error message and code
// assignment. Can be used to assign error codes/message paring.
package messagecode

import (
	"net/http"
)

// Messages is the interface used for converting a key to customized message
// eg: E_INVALID_REQUEST -> The request is malformed
var Messages map[string]map[string]interface{}

// Binding of HTTP response codes and messages
const (
	HTTP_BAD_REQUEST             int = http.StatusBadRequest
	HTTP_UNAUTHORIZED                = http.StatusUnauthorized
	HTTP_BAD_GATEWAY                 = http.StatusBadGateway
	HTTP_INTERNAL_ERROR              = http.StatusInternalServerError
	HTTP_OK                          = http.StatusOK
	HTTP_CREATED                     = http.StatusCreated
	HTTP_UPDATED                     = http.StatusNoContent
	HTTP_DELETED                     = http.StatusNoContent
	HTTP_DISABLED                    = http.StatusNoContent
	HTTP_ACCESS_DENIED               = http.StatusForbidden
	HTTP_TEMPORARILY_UNAVAILABLE     = http.StatusServiceUnavailable
	HTTP_NOT_FOUND                   = http.StatusNotFound

	E_INVALID_REQUEST           string = "invalid_request"
	E_UNAUTHORIZED_CLIENT              = "unauthorized_client"
	E_ACCESS_DENIED                    = "access_denied"
	E_UNSUPPORTED_RESPONSE_TYPE        = "unsupported_response_type"
	E_SERVER_ERROR                     = "server_error"
	E_TEMPORARILY_UNAVAILABLE          = "temporarily_unavailable"

	S_RESOURCE_CREATED  string = "resource_created"
	S_RESOURCE_OK              = "resource_ok"
	S_RESOURCE_UPDATED         = "resource_updated"
	S_RESOURCE_DISABLED        = "resource_disabled"
	S_RESOURCE_NOTFOUND        = "resource_not_found"
)

// defineErrorMessages defines all error codes and messages
// parings.
func defineErrorMessages() {
	Messages[E_INVALID_REQUEST] = make(map[string]interface{})
	Messages[E_INVALID_REQUEST]["HTTP_CODE"] = HTTP_BAD_REQUEST
	Messages[E_INVALID_REQUEST]["MSG"] = "The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed."

	Messages[E_UNAUTHORIZED_CLIENT] = make(map[string]interface{})
	Messages[E_UNAUTHORIZED_CLIENT]["MSG"] = "The client is not authorized to request using this method."
	Messages[E_UNAUTHORIZED_CLIENT]["HTTP_CODE"] = HTTP_UNAUTHORIZED

	Messages[E_ACCESS_DENIED] = make(map[string]interface{})
	Messages[E_ACCESS_DENIED]["MSG"] = ""
	Messages[E_ACCESS_DENIED]["HTTP_CODE"] = HTTP_ACCESS_DENIED

	Messages[E_UNSUPPORTED_RESPONSE_TYPE] = make(map[string]interface{})
	Messages[E_UNSUPPORTED_RESPONSE_TYPE]["MSG"] = ""
	Messages[E_UNSUPPORTED_RESPONSE_TYPE]["HTTP_CODE"] = HTTP_BAD_REQUEST

	Messages[E_SERVER_ERROR] = make(map[string]interface{})
	Messages[E_SERVER_ERROR]["MSG"] = ""
	Messages[E_SERVER_ERROR]["HTTP_CODE"] = HTTP_INTERNAL_ERROR

	Messages[E_TEMPORARILY_UNAVAILABLE] = make(map[string]interface{})
	Messages[E_TEMPORARILY_UNAVAILABLE]["MSG"] = ""
	Messages[E_TEMPORARILY_UNAVAILABLE]["HTTP_CODE"] = HTTP_TEMPORARILY_UNAVAILABLE
}

// defineStatusMessages defines all status codes and messages
// parings.
func defineStatusMessages() {
	Messages[S_RESOURCE_OK] = make(map[string]interface{})
	Messages[S_RESOURCE_OK]["HTTP_CODE"] = HTTP_OK
	Messages[S_RESOURCE_OK]["MSG"] = ""

	Messages[S_RESOURCE_CREATED] = make(map[string]interface{})
	Messages[S_RESOURCE_CREATED]["HTTP_CODE"] = HTTP_CREATED
	Messages[S_RESOURCE_CREATED]["MSG"] = ""

	Messages[S_RESOURCE_UPDATED] = make(map[string]interface{})
	Messages[S_RESOURCE_UPDATED]["MSG"] = ""
	Messages[S_RESOURCE_UPDATED]["HTTP_CODE"] = HTTP_UPDATED

	Messages[S_RESOURCE_DISABLED] = make(map[string]interface{})
	Messages[S_RESOURCE_DISABLED]["MSG"] = ""
	Messages[S_RESOURCE_DISABLED]["HTTP_CODE"] = HTTP_DISABLED

	Messages[S_RESOURCE_NOTFOUND] = make(map[string]interface{})
	Messages[S_RESOURCE_NOTFOUND]["MSG"] = ""
	Messages[S_RESOURCE_NOTFOUND]["HTTP_CODE"] = HTTP_NOT_FOUND

}

/*
func Init() {
	Messages = make(map[string]map[string]interface{})
	defineErrorMessages()
	defineStatusMessages()
}
*/
