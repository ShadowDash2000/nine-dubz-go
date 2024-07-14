package user

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	IsSuccess bool `json:"isSuccess"`
}

func NewLoginRequest(user *LoginRequest) *User {
	return &User{
		Active:   true,
		Email:    user.Email,
		Password: user.Password,
	}
}

func NewLoginResponse(isSuccess bool) *LoginResponse {
	return &LoginResponse{
		IsSuccess: isSuccess,
	}
}
