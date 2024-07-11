package controller

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"gorm.io/gorm"
	"io"
	"net/http"
	"nine-dubz/app/model"
	"nine-dubz/app/model/language"
	"nine-dubz/app/model/payload"
	"nine-dubz/app/usecase"
	"strconv"
)

type UserController struct {
	UserInteractor     usecase.UserInteractor
	TokenController    *TokenController
	HelperInteractor   *usecase.HelperInteractor
	LanguageController *LanguageController
	FileController     *FileController
}

func NewUserController(db *gorm.DB, tc *TokenController, lc *LanguageController, fc *FileController) *UserController {
	return &UserController{
		UserInteractor: usecase.UserInteractor{
			UserRepository: &UserRepository{
				DB: db,
			},
		},
		TokenController: tc,
		HelperInteractor: &usecase.HelperInteractor{
			HelperRepository: &HelperRepository{},
		},
		LanguageController: lc,
		FileController:     fc,
	}
}

func (uc *UserController) AddHandler(w http.ResponseWriter, r *http.Request) {
	user := &model.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	id, err := uc.UserInteractor.Add(user)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err, http.StatusNotFound, "Cannot add user"))
		return
	}

	render.JSON(w, r, id)
}

func (uc *UserController) GetHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseUint(chi.URLParam(r, "userId"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
	}

	user, err := uc.UserInteractor.Get(uint(userId))
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err, http.StatusBadRequest, "User not found"))
		return
	}

	render.JSON(w, r, user)
}

func (uc *UserController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	user := &model.User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		render.Render(w, r, ErrInvalidRequest(err, http.StatusBadRequest, "Invalid login fields"))
		return
	}

	loginPayload := payload.NewLoginPayload(user)
	tokenCookie, stringCode, err := uc.Login(loginPayload, true)
	if err != nil {
		text, _ := uc.LanguageController.GetStringByCode(r, stringCode.Code)
		render.Render(w, r, ErrInvalidRequest(err, http.StatusUnauthorized, text))
		return
	}

	http.SetCookie(w, tokenCookie)

	render.JSON(w, r, struct {
		IsSuccess bool `json:"isSuccess"`
	}{true})
}

func (uc *UserController) Login(loginPayload *payload.LoginPayload, usePassword bool) (*http.Cookie, *language.StringCode, error) {
	if usePassword {
		if ok := uc.UserInteractor.Login(loginPayload.User); !ok {
			return nil, language.NewStringCode("LOGIN_USER_NOT_FOUND"), errors.New("user not found")
		}
	} else {
		if ok := uc.UserInteractor.LoginWOPassword(loginPayload.User); !ok {
			return nil, language.NewStringCode("LOGIN_USER_NOT_FOUND"), errors.New("user not found")
		}
	}

	tokenCookie, err := uc.GetTokenCookie(loginPayload.User)
	if err != nil {
		return nil, language.NewStringCode("TOKEN_FAILED_TO_CREATE"), err
	}

	return tokenCookie, language.NewStringCode(""), err
}

func (uc *UserController) GetTokenCookie(user *model.User) (*http.Cookie, error) {
	token, claims, err := uc.TokenController.Create(user)
	if err != nil {
		return nil, err
	}

	tokenCookie := http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Expires:  claims.ExpiresAt.Time,
		MaxAge:   0,
		Secure:   false,
		HttpOnly: false,
		SameSite: http.SameSiteDefaultMode,
	}

	return &tokenCookie, nil
}

func (uc *UserController) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	user := &model.User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		render.Render(w, r, ErrInvalidRequest(err, http.StatusBadRequest, "Invalid registration fields"))
		return
	}

	registrationPayload := payload.NewRegistrationPayload(user)
	userId, stringCode, err := uc.Register(registrationPayload, false)
	if err != nil {
		text, _ := uc.LanguageController.GetStringByCode(r, stringCode.Code)
		render.Render(w, r, ErrInvalidRequest(err, http.StatusBadRequest, text))
		return
	}

	render.JSON(w, r, struct {
		UserId uint `json:"userId"`
	}{userId})
}

func (uc *UserController) Register(registrationPayload *payload.RegistrationPayload, skipFieldsValidation bool) (uint, *language.StringCode, error) {
	if !skipFieldsValidation {
		if err := uc.HelperInteractor.ValidateRegistrationFields(registrationPayload); err != nil {
			return 0, language.NewStringCode("REGISTRATION_INVALID_FIELDS"), err
		}
	}

	userId, err := uc.UserInteractor.Add(registrationPayload.User)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return 0, language.NewStringCode("EMAIL_ALREADY_EXISTS"), err
		}
		return 0, language.NewStringCode("EMAIL_ALREADY_EXISTS"), err
	}

	return userId, language.NewStringCode(""), nil
}

func (uc *UserController) CheckUserWithNameExistsHandler(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("userName")
	if userName == "" {
		render.Render(w, r, ErrInvalidRequest(errors.New("userName parameter is required"), http.StatusBadRequest, ""))
		return
	}

	if ok := uc.HelperInteractor.ValidateUserName(userName); !ok {
		render.Render(w, r, ErrInvalidRequest(errors.New("not valid username parameter"), http.StatusBadRequest, ""))
		return
	}

	isUserExists := uc.CheckUserWithNameExists(userName)

	render.JSON(w, r, struct {
		IsUserExists bool `json:"isUserExists"`
	}{isUserExists})
}

func (uc *UserController) CheckUserWithNameExists(userName string) bool {
	isUserExists := false
	_, err := uc.UserInteractor.GetByName(userName)
	if err == nil {
		isUserExists = true
	}

	return isUserExists
}

func (uc *UserController) GetUserShortHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(uint)

	user, err := uc.UserInteractor.Get(userId)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err, http.StatusNotFound, "User not found"))
		return
	}

	render.JSON(w, r, *payload.NewUserShortPayload(user))
}

func (uc *UserController) UpdateUserPictureHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		text, _ := uc.LanguageController.GetStringByCode(r, "REQUEST_MAX_SIZE_LIMIT")
		render.Render(w, r, ErrInvalidRequest(err, http.StatusBadRequest, text))
		return
	}

	file, fileHeader, err := r.FormFile("picture")
	if err != nil {
		text, _ := uc.LanguageController.GetStringByCode(r, "REQUEST_NO_PICTURE_FILE")
		render.Render(w, r, ErrInvalidRequest(err, http.StatusBadRequest, text))
		return
	}
	defer file.Close()

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err, http.StatusInternalServerError, "Failed to verify file type"))
		return
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err, http.StatusInternalServerError, ""))
		return
	}
	isCorrectType, _ := uc.FileController.VerifyFileType(buff, []string{"image/jpeg", "image/png", "image/gif", "image/webp"})
	if !isCorrectType {
		text, _ := uc.LanguageController.GetStringByCode(r, "REQUEST_WRONG_FILE_TYPE")
		render.Render(w, r, ErrInvalidRequest(err, http.StatusBadRequest, text))
		return
	}

	picture, err := uc.FileController.SaveFile("upload/profile_pictures", fileHeader.Filename, file)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err, http.StatusBadRequest, "Failed to save file"))
		return
	}

	userId := r.Context().Value("userId").(uint)
	err = uc.UserInteractor.Updates(&model.User{
		ID:      userId,
		Picture: picture,
	})
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err, http.StatusInternalServerError, "Failed to update user picture"))
		return
	}
}
