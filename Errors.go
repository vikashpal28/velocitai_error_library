package myerrors

import "errors"

// dto of an API error, including a code, message, operation, and underlying error.
type ApiError struct{
	Code Code
	Message string
	Op string
	Err error
}

func (e *ApiError) Error() string {
	if e.Op != "" {
		if e.Err != nil{
			return e.Op + ": " + e.Message + ": " + e.Err.Error()
		}
		return e.Op + ": " + e.Message
	}
	return e.Message
}

// Unwrap returns the underlying error, if any.
func (e *ApiError) Unwrap() error {
	return e.Err
}

func (e *ApiError) Is(target error) bool {
	t , ok := target.(*ApiError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// New creates a new ApiError with the given code and message.
func New(code Code, message string) *ApiError {
	return &ApiError{
		Code: code,
		Message: message,
	}
}

func Wrap(err error, code Code, message string) *ApiError {
	return &ApiError{
		Code: code,
		Message: message,
		Err: err,
	}
}

// WithOp sets the operation name for the ApiError and returns a new instance with the updated operation.
func (e *ApiError) WithOp(op string) *ApiError {
	cp := *e
	cp.Op = op
	return &cp
}

func (e *ApiError) HTTPStatus() int {
	return e.Code.HTTPStatus()
}

func CodeOf(err error) Code {
	var ae *ApiError
	if errors.As(err, &ae) {
		return ae.Code
	}
	return CodeUnknown
}

func NotFound(message string) *ApiError {
	return &ApiError{
		Code: CodeNotFound,
		Message: message,
	}
}

func ValidationFailed(message string) *ApiError {
	return &ApiError{
		Code: CodeValidation,
		Message: message,
	}
}

func Unauthorized(message string) *ApiError {
	return &ApiError{
		Code: CodeUnauthorized,
		Message: message,
	}
}


func Forbidden(message string) *ApiError {
	return New(CodeForbidden, message)
}
 
func Internal(message string) *ApiError {
	return New(CodeInternal, message)
}