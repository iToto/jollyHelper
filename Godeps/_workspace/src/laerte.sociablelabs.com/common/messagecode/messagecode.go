// Package messagecode provides a centralized error message and code
// assignment. Can be used to assign error codes/message paring.
package messagecode

//type DefaultErrorId string
type Merror error

// MessageCode used for the
var (
	DefMessage *MessageCode = NewMessageCode()
)

// MessageCode is the interface for the message object
type MessageCode struct {
	messagemap map[string]map[string]interface{}
}

// NewMessageCode creates and returns a new Message Code
//
func NewMessageCode() *MessageCode {
	Messages = make(map[string]map[string]interface{})
	defineErrorMessages()
	defineStatusMessages()

	r := &MessageCode{messagemap: Messages}
	return r
}

// Get gets the message for a given message id
func (e *MessageCode) Get(id string) map[string]interface{} {
	if m, ok := e.messagemap[id]; ok {
		return m
	}
	r := make(map[string]interface{})
	r["HTTP_CODE"] = HTTP_BAD_GATEWAY
	r["MSG"] = id
	return r
}
