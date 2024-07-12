package payload

type GoogleToken struct {
	IssuedTo      string `json:"issued_to"`
	Audience      string `json:"audience"`
	UserId        string `json:"user_id"`
	Scope         string `json:"scope"`
	ExpiresIn     int    `json:"expires_in"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	AccessType    string `json:"access_type"`
}
