package entities

type GenericResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Code  int64  `json:"code"`
	Error string `json:"error"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
