package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/LuisDavid01/Websockets-go/internal/store"
	"github.com/LuisDavid01/Websockets-go/internal/utils"
)

type registerUserReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserHandler struct {
	userStore store.IUser
	logger    *log.Logger
}

func NewUserHandler(userStore store.IUser, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

func (h *UserHandler) validateRegisterUser(req *registerUserReq) error {
	if req.Username == "" {
		return errors.New("Err Invalid Username")
	}
	if req.Email == "" {
		return errors.New("Err Invalid Email")
	}
	if req.Password == "" {
		return errors.New("Err Invalid Password")
	}

	if len(req.Password) < 8 {
		return errors.New("Err Password must be at least 8 characters long")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("Err Invalid Email format")
	}
	return nil

}

// Handles user register
func (h *UserHandler) HandlerRegisterUser(w http.ResponseWriter, r *http.Request) {
	var req registerUserReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("Error decoding the payload: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"ERROR": err.Error()})
		return
	}

	err = h.validateRegisterUser(&req)
	if err != nil {
		h.logger.Printf("Invalid request: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"ERROR": err.Error()})
	}

	user := &store.User{
		Username: req.Username,
		Email:    req.Email,
	}

	err = user.Password.Set(req.Password)
	if err != nil {
		h.logger.Printf("Error setting user password: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"ERROR": err.Error()})
		return
	}
	utils.WriteJson(w, http.StatusCreated, utils.Envelope{"user": user})
}
