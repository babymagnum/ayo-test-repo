package response

type LoginResponse struct {
	BaseResponse
	Data LoginData `json:"data"`
}

type LoginData struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
	Email string `json:"email"`
}
