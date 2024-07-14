package token

type Interactor interface {
	Add(token *Token) error
	GetByUserId(userId uint) (*Token, error)
	GetUserIdByToken(tokenString string) (*Token, error)
}
