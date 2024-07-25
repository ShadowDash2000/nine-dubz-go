package token

import "gorm.io/gorm"

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) Add(token *Token) error {
	return r.DB.Create(token).Error
}

func (r *Repository) Delete(userId uint, tokenString string) (int64, error) {
	result := r.DB.Where("user_id = ? AND token = ?", userId, tokenString).Delete(&Token{})

	return result.RowsAffected, result.Error
}

func (r *Repository) GetByUserId(userId uint) (*Token, error) {
	token := &Token{}
	result := r.DB.Where("user_id = ?", userId).First(token)

	return token, result.Error
}

func (r *Repository) GetUserIdByToken(tokenString string) (*Token, error) {
	token := &Token{}
	result := r.DB.Where("token = ?", tokenString).First(token)

	return token, result.Error
}
