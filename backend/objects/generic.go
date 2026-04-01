package objects

type GenericResponse struct {
	Success bool   `json:"success"`
	ErrMsg  string `json:"errMsg"`
}

var SuccessGeneric = GenericResponse{
	Success: true,
	ErrMsg:  "",
}

type GenericError struct {
	Success    bool   `json:"success"`
	ErrMsg     string `json:"errMsg"`
	StatusCode int    `json:"statusCode"`
}

var InternalError = GenericError{
	Success:    false,
	ErrMsg:     "internal server error",
	StatusCode: 500,
}

var MalformedError = GenericError{
	Success:    false,
	ErrMsg:     "malformed request",
	StatusCode: 500,
}

func (g GenericError) Error() string {
	return g.ErrMsg
}
