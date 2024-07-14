package user

type RegistrationRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegistrationResponse struct {
	IsSuccess bool `json:"isSuccess"`
}

func NewRegistrationRequest(user *RegistrationRequest) *User {
	return &User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}
}

func NewRegistrationResponse(isSuccess bool) *RegistrationResponse {
	return &RegistrationResponse{
		IsSuccess: isSuccess,
	}
}
