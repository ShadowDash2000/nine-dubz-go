package user

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"io"
	"mime/multipart"
	"nine-dubz/internal/apimethod"
	"nine-dubz/internal/file"
	"nine-dubz/internal/helper"
	"nine-dubz/internal/mail"
	"nine-dubz/internal/role"
	"nine-dubz/internal/token"
)

type UseCase struct {
	UserInteractor   Interactor
	ApiMethodUseCase *apimethod.UseCase
	TokenUseCase     *token.UseCase
	RoleUseCase      *role.UseCase
	FileUseCase      *file.UseCase
	MailUseCase      *mail.UseCase
}

func New(db *gorm.DB, tuc *token.UseCase, ruc *role.UseCase, fuc *file.UseCase, muc *mail.UseCase) *UseCase {
	return &UseCase{
		UserInteractor: &Repository{
			DB: db,
		},
		ApiMethodUseCase: apimethod.New(db),
		TokenUseCase:     tuc,
		RoleUseCase:      ruc,
		FileUseCase:      fuc,
		MailUseCase:      muc,
	}
}

func (uc *UseCase) Add(user *User) uint {
	return uc.UserInteractor.Add(user)
}

func (uc *UseCase) Login(user *User) uint {
	user.Password = helper.HashPassword(user.Password)
	return uc.UserInteractor.Login(user)
}

func (uc *UseCase) LoginWOPassword(user *User) uint {
	return uc.UserInteractor.LoginWOPassword(user)
}

func (uc *UseCase) Register(user *User) uint {
	if err := helper.ValidateRegistrationFields(user.Name, user.Email, user.Password); err != nil {
		return 0
	}

	user.Password = helper.HashPassword(user.Password)
	return uc.UserInteractor.Add(user)
}

func (uc *UseCase) CheckUserWithNameExists(userName string) bool {
	isUserExists := false
	_, err := uc.UserInteractor.GetByName(userName)
	if err == nil {
		isUserExists = true
	}

	return isUserExists
}

func (uc *UseCase) SendRegistrationEmail(to, subject, content string) error {
	return uc.MailUseCase.SendMail(to, subject, content)
}

func (uc *UseCase) ConfirmRegistration(email, hash string) bool {
	user := &User{}
	err := uc.UserInteractor.GetWhere(user, map[string]interface{}{"active": false, "email": email, "hash": hash})
	if err != nil {
		return false
	}

	user = &User{
		ID:     user.ID,
		Active: true,
	}
	err = uc.UserInteractor.Updates(user)
	if err != nil {
		return false
	}

	return true
}

func (uc *UseCase) GetById(id uint) (*User, error) {
	return uc.UserInteractor.GetById(id)
}

func (uc *UseCase) UpdatePicture(userId uint, file multipart.File, header *multipart.FileHeader) error {
	buff := make([]byte, 512)
	_, err := file.Read(buff)
	if err != nil {
		return err
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	isCorrectType, _ := uc.FileUseCase.VerifyFileType(buff, []string{"image/jpeg", "image/png", "image/gif", "image/webp"})
	if !isCorrectType {
		return errors.New("invalid file type")
	}

	picture, err := uc.FileUseCase.SaveFile("upload/profile_pictures", header.Filename, file)
	if err != nil {
		return err
	}

	err = uc.UserInteractor.Updates(&User{
		ID:      userId,
		Picture: picture,
	})
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) CheckUserPermission(tokenString string, routePattern string, method string) (uint, bool) {
	userId, err := uc.TokenUseCase.GetUserIdByToken(tokenString)
	if err != nil {
		return 0, false
	}
	fmt.Println(userId)
	roles, err := uc.UserInteractor.GetRolesByUserId(userId)
	if err != nil {
		return 0, false
	}
	fmt.Println(roles)
	var rolesIds []uint
	for _, role := range roles {
		if role.Code == "admin" {
			return userId, true
		}

		rolesIds = append(rolesIds, role.ID)
	}
	fmt.Println(rolesIds)
	apiMethods, err := uc.RoleUseCase.GetApiMethodsByRolesIds(rolesIds)
	if err != nil {
		return 0, false
	}
	fmt.Println(apiMethods)
	isHavePermission := false
	for _, apiMethod := range apiMethods {
		if apiMethod.Path == routePattern && apiMethod.Method == method {
			isHavePermission = true
			break
		}
	}

	return userId, isHavePermission
}
