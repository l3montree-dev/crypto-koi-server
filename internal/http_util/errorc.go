package http_util

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.com/l3montree/microservices/libs/orchardclient"
)

// ErrorC - Custom error.
type ErrorC struct {
	Status int
	Err    string `json:"error"`
}

func (errorc *ErrorC) Error() string {
	return fmt.Sprintf("Status: %d, Message: %s", errorc.Status, errorc.Err)
}

func (errorc ErrorC) ToJSON() []byte {
	b, err := json.Marshal(errorc)
	orchardclient.FailOnError(err, "could not marshal ErrorC to JSON")
	return b
}

// returns a function to simply log the error
func WriteHttpError(writer http.ResponseWriter, status int, message string, args ...interface{}) {
	(writer).WriteHeader(status)
	(writer).Write(NewError(status, message, args...).ToJSON())
}

func NewError(status int, message string, args ...interface{}) ErrorC {
	if len(args) > 0 {
		return ErrorC{
			Status: status,
			Err:    fmt.Sprintf(message, args...),
		}
	}
	return ErrorC{
		Status: status,
		Err:    message,
	}
}
