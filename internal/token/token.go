package token

import "gorm.io/gorm"

type UseCase struct {
	TokenInteractor Interactor
}

func New(db *gorm.DB) *UseCase {
	return &UseCase{
		TokenInteractor: &Repository{
			DB: db,
		},
	}
}

func (uc *UseCase) Add(userId uint, tokenString string) error {
	token := &Token{
		UserId: userId,
		Token:  tokenString,
	}

	return uc.TokenInteractor.Add(token)
}

func (uc *UseCase) GetByUserId(userId uint) (*Token, error) {
	return uc.TokenInteractor.GetByUserId(userId)
}

func (uc *UseCase) GetUserIdByToken(tokenString string) (uint, error) {
	token, err := uc.TokenInteractor.GetUserIdByToken(tokenString)
	if err != nil {
		return 0, err
	}

	return token.UserId, nil
}
