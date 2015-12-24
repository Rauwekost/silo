package http

//Error ...
type Error struct {
	Status  int    `json:"-"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

func ErrInvalidRequest(detail string) *Error {
	return &Error{Status: 400, Message: "invalid request", Detail: detail}
}

func ErrNotAllowed(detail string) *Error {
	return &Error{Status: 401, Message: "not allowed", Detail: detail}
}

func ErrForbidden(detail string) *Error {
	return &Error{Status: 403, Message: "forbiden", Detail: detail}
}

func ErrNotFound(detail string) *Error {
	return &Error{Status: 404, Message: "not found", Detail: detail}
}

func ErrInternalServer(detail string) *Error {
	return &Error{Status: 500, Message: "internal server error", Detail: detail}
}
