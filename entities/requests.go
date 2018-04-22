package entities

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CallRequest struct {
	Token                 string `json:"token"`
	DeviceName            string `json:"device_name"`
	FunctionName          string `json:"function_name"`
	WaitForDeviceResponse bool   `json:"wait_for_device_response"`
}
