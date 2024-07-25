package token

type Interactor interface {
	Add(token *Token) error
	Delete(userId uint, tokenString string) (int64, error)
	GetByUserId(userId uint) (*Token, error)
	GetUserIdByToken(tokenString string) (*Token, error)
}
