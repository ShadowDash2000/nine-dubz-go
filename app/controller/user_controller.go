package controller

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"gorm.io/gorm"
	"net/http"
	"nine-dubz/app/model"
	"nine-dubz/app/model/payload"
	"nine-dubz/app/usecase"
	"strconv"
)

type UserController struct {
	UserInteractor   usecase.UserInteractor
	TokenController  *TokenController
	HelperInteractor *usecase.HelperInteractor
}

func NewUserController(db *gorm.DB, tokenController *TokenController) *UserController {
	return &UserController{
		UserInteractor: usecase.UserInteractor{
			UserRepository: &UserRepository{
				DB: db,
			},
		},
		TokenController: tokenController,
		HelperInteractor: &usecase.HelperInteractor{
			HelperRepository: &HelperRepository{},
		},
	}
}

func (uc *UserController) Add(w http.ResponseWriter, r *http.Request) {
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

func (uc *UserController) Get(w http.ResponseWriter, r *http.Request) {
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

func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
	user := &model.User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		render.Render(w, r, ErrInvalidRequest(err, http.StatusBadRequest, "Invalid login fields"))
		return
	}

	loginPayload := payload.NewLoginPayload(user)
	if ok := uc.UserInteractor.Login(loginPayload.User); !ok {
		render.Render(w, r, ErrInvalidRequest(nil, http.StatusUnauthorized, "No user found"))
		return
	}

	token, claims, err := uc.TokenController.Create(loginPayload.User)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err, http.StatusUnauthorized, "Failed to create token"))
		return
	}

	tokenCookie := http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Expires:  claims.ExpiresAt.Time,
		MaxAge:   0,
		Secure:   false,
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &tokenCookie)

	render.JSON(w, r, struct {
		IsSuccess bool `json:"isSuccess"`
	}{true})
}

func (uc *UserController) Register(w http.ResponseWriter, r *http.Request) {
	user := &model.User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		render.Render(w, r, ErrInvalidRequest(err, http.StatusBadRequest, "Invalid registration fields"))
		return
	}

	registrationPayload := payload.NewRegistrationPayload(user)
	if err := uc.HelperInteractor.ValidateRegistrationFields(registrationPayload); err != nil {
		render.Render(w, r, ErrInvalidRequest(err, http.StatusBadRequest, "Invalid registration fields"))
		return
	}

	userId, err := uc.UserInteractor.Add(registrationPayload.User)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			render.Render(w, r, ErrInvalidRequest(err, http.StatusBadRequest, "User with this email already exists"))
			return
		}
		render.Render(w, r, ErrInvalidRequest(err, http.StatusBadRequest, "Cannot register user"))
		return
	}

	render.JSON(w, r, struct {
		UserId uint `json:"userId"`
	}{userId})
}

func (uc *UserController) CheckUserWithNameExists(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("userName")
	if userName == "" {
		render.Render(w, r, ErrInvalidRequest(errors.New("userName parameter is required"), http.StatusBadRequest, ""))
		return
	}

	if ok := uc.HelperInteractor.ValidateUserName(userName); !ok {
		render.Render(w, r, ErrInvalidRequest(errors.New("not valid username parameter"), http.StatusBadRequest, ""))
		return
	}

	isUserExists := false
	_, err := uc.UserInteractor.GetByName(userName)
	if err == nil {
		isUserExists = true
	}

	render.JSON(w, r, struct {
		IsUserExists bool `json:"isUserExists"`
	}{isUserExists})
}

func (uc *UserController) GetUserShort(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(uint)

	user, err := uc.UserInteractor.Get(userId)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err, http.StatusNotFound, "User not found"))
	}

	render.JSON(w, r, *payload.NewUserShortPayload(user))
}
