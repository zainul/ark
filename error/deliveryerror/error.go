package deliveryerror

import (
	"github.com/zainul/ark/error/usererror"
)

// Error ...
type Error struct {
	ErrorCode string `json:"error_code"`
	ErrorMsg  string `json:"error_message"`
	ErrorCase string `json:"error_case"`
}

// GetError ...
func GetError(errCode string, caseError error) *Error {
	if caseError != nil {
		err := Error{}
		errMsg := usererror.GetCode(errCode)
		err.ErrorCode = errCode
		err.ErrorMsg = errMsg
		err.ErrorCase = caseError.Error()
		return &err
	}
	return nil
}
