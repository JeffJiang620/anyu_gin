package apierr

type APIError struct {
	code string
	desc string
}

func NewAPIError(code, desc string) *APIError {
	return &APIError{code: code, desc: desc}
}

func (e *APIError) Error() string {
	return e.desc
}

func (e *APIError) Code() string {
	return e.code
}
