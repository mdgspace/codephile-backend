package errors

type ErrorResponse struct {
	ErrorType string `json:"error_type"`
	Err       string `json:"error_desc"`
}

func InternalServerError(error string) ErrorResponse {
	return ErrorResponse{
		ErrorType: "server_error",
		Err:       error,
	}
}

func BadInputError(error string) ErrorResponse {
	return ErrorResponse{
		ErrorType: "bad_input",
		Err:       error,
	}
}

func AlreayExistsError(error string) ErrorResponse {
	return ErrorResponse{
		ErrorType: "conflict",
		Err:       error,
	}
}
func NotFoundError(error string) ErrorResponse {
	return ErrorResponse{
		ErrorType: "not_found",
		Err:       error,
	}
}
