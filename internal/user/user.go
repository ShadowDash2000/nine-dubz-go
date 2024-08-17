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
	"time"
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

func (uc *UseCase) Add(user *User) (uint, error) {
	user.Hash = helper.Hash([]byte(user.Name + user.Email + user.Password + time.Now().String()))

	return uc.UserInteractor.Add(user)
}

func (uc *UseCase) Login(user *User) uint {
	user.Password = helper.HashPassword(user.Password)
	return uc.UserInteractor.Login(user)
}

func (uc *UseCase) LoginWOPassword(user *User) uint {
	return uc.UserInteractor.LoginWOPassword(user)
}

func (uc *UseCase) Register(user *User) (uint, error) {
	if ok := helper.ValidateUserName(user.Name); !ok {
		return 0, errors.New("REGISTRATION_INVALID_USER_NAME")
	}
	if ok := helper.ValidateEmail(user.Email); !ok {
		return 0, errors.New("REGISTRATION_INVALID_EMAIL")
	}
	if ok := helper.ValidatePassword(user.Password); !ok {
		return 0, errors.New("REGISTRATION_INVALID_PASSWORD")
	}

	user.Hash = helper.Hash([]byte(user.Name + user.Email + user.Password + time.Now().String()))
	user.Password = helper.HashPassword(user.Password)
	userId, err := uc.UserInteractor.Add(user)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return 0, errors.New("REGISTRATION_ALREADY_EXIST")
		}
		return 0, errors.New("INTERNAL_ERROR")
	} else if userId > 0 {
		return userId, nil
	}

	return 0, errors.New("INTERNAL_ERROR")
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

func (uc *UseCase) ConfirmRegistration(email, hash string) (uint, bool) {
	user := &User{}
	err := uc.UserInteractor.GetWhere(user, map[string]interface{}{"active": false, "email": email, "hash": hash})
	if err != nil {
		return 0, false
	}

	user = &User{
		ID:     user.ID,
		Active: true,
		Hash:   helper.Hash([]byte(user.Name + user.Email + user.Password + time.Now().String())),
	}
	err = uc.UserInteractor.Updates(user)
	if err != nil {
		return 0, false
	}

	return user.ID, true
}

func (uc *UseCase) GetMultiple(where interface{}) ([]User, error) {
	return uc.UserInteractor.GetWhereMultiple(where)
}

func (uc *UseCase) GetDistinctMultiple(where, distinct interface{}) ([]User, error) {
	return uc.UserInteractor.GetMultiple(where, distinct)
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
		return errors.New("USER_UPDATE_PICTURE_INVALID_TYPE")
	}

	user, err := uc.UserInteractor.GetById(userId)
	if err != nil {
		return errors.New("USER_NOT_FOUND")
	}
	if user.Picture != nil {
		uc.FileUseCase.Delete(user.Picture.Name)
	}

	pictureSavePath := fmt.Sprintf("user/inner/%d", userId)
	picture, err := uc.FileUseCase.Create(file, header.Filename, pictureSavePath, header.Size, "public")
	if err != nil {
		return errors.New("INTERNAL_ERROR")
	}

	err = uc.UserInteractor.Updates(&User{
		ID:      userId,
		Picture: picture,
	})
	if err != nil {
		return errors.New("INTERNAL_ERROR")
	}

	return nil
}

func (uc *UseCase) Update(user *UpdateRequest) error {
	fieldsToUpdate := 0

	if user.Name != "" {
		if ok := helper.ValidateUserName(user.Name); !ok {
			return errors.New("REGISTRATION_INVALID_USER_NAME")
		} else {
			fieldsToUpdate = fieldsToUpdate + 1
		}
	}

	if fieldsToUpdate == 0 {
		return errors.New("USER_UPDATE_INVALID_FIELDS")
	}

	err := uc.UserInteractor.Updates(NewUpdateRequest(user))
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.New("USER_UPDATE_USERNAME_ALREADY_EXIST")
		}
		return errors.New("INTERNAL_ERROR")
	}

	return nil
}

func (uc *UseCase) CheckUserPermission(userId uint, routePattern string, method string) bool {
	roles, err := uc.UserInteractor.GetRolesByUserId(userId)
	if err != nil {
		return false
	}

	var rolesIds []uint
	for _, role := range roles {
		if role.Code == "admin" {
			return true
		}

		rolesIds = append(rolesIds, role.ID)
	}

	apiMethods, err := uc.RoleUseCase.GetApiMethodsByRolesIds(rolesIds)
	if err != nil {
		return false
	}

	isHavePermission := false
	for _, apiMethod := range apiMethods {
		if apiMethod.Path == routePattern && apiMethod.Method == method {
			isHavePermission = true
			break
		}
	}

	return isHavePermission
}
