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

type SendResponse struct {
	Response string `json:"response"`
}

type UserInfoResponse struct {
	AccountSecret string `json:"account_secret"`
}
