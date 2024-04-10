package user

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/trsnaqe/gotask/services/auth"
	"github.com/trsnaqe/gotask/types"
	"github.com/trsnaqe/gotask/utils"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserPayload
	err := utils.ParseJSON(r, &payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", err.(validator.ValidationErrors)))
		return
	}
	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid email or password"))
		return
	}

	if !auth.CompareValue(u.Password, payload.Password) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid email or password"))
		return
	}

	tokens, err := auth.CreateTokens(u.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	hashedRefreshToken, err := auth.HashRefreshToken(tokens.RefreshToken)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.UpdateUser(u.ID, types.User{RefreshToken: &hashedRefreshToken})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, tokens)
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterUserPayload
	err := utils.ParseJSON(r, &payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", err.(validator.ValidationErrors)))
		return
	}
	_, err = h.store.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}

	hashedPassword, err := auth.HashValue(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.CreateUser(types.User{
		Email:    payload.Email,
		Password: hashedPassword,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	createdUser, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	tokens, err := auth.CreateTokens(createdUser.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	hashedRefreshToken, err := auth.HashRefreshToken(tokens.RefreshToken)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.UpdateUser(createdUser.ID, types.User{RefreshToken: &hashedRefreshToken})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, tokens)
}

// removes refresh token from user, effectively logging out
func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	log.Println("logout")
	log.Println(types.UserKey)

	id := auth.GetUserIDFromContext(r.Context())
	u, err := h.store.GetUserByID(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.UpdateUser(u.ID, types.User{RefreshToken: nil})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, nil)
}

// refreshes access and refresh token using refresh token
func (h *Handler) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	tokenString := utils.GetTokenFromRequest(r)

	userId := auth.GetUserIDFromContext(r.Context())
	u, err := h.store.GetUserByID(userId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if u.RefreshToken == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("no refresh token found"))
		return
	}

	if !auth.CompareRefreshToken(*u.RefreshToken, tokenString) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid refresh token"))
		return
	}

	_, err = auth.ValidateJWT(tokenString)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	tokens, err := auth.CreateTokens(userId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	hashedRefreshToken, err := auth.HashRefreshToken(tokens.RefreshToken)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	err = h.store.UpdateUser(userId, types.User{RefreshToken: &hashedRefreshToken})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, tokens)
}

// Changes user password using old password
func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var payload types.ChangePasswordPayload
	err := utils.ParseJSON(r, &payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", err.(validator.ValidationErrors)))
		return
	}

	userId := auth.GetUserIDFromContext(r.Context())
	err = h.store.ChangePassword(userId, payload.OldPassword, payload.NewPassword)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, nil)
}
