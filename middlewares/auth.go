package middlewares

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/trsnaqe/gotask/services/auth"
	"github.com/trsnaqe/gotask/types"
	"github.com/trsnaqe/gotask/utils"
)

func AuthMiddleware(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := utils.GetTokenFromRequest(r)

		token, err := auth.ValidateJWT(tokenString)
		if err != nil {
			log.Printf("failed to validate token: %v", err)
			utils.PermissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("invalid token")
			utils.PermissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		str := claims["userID"].(string)

		userID, err := strconv.Atoi(str)
		if err != nil {
			log.Printf("failed to convert userID to int: %v", err)
			utils.PermissionDenied(w)
			return
		}

		u, err := store.GetUserByID(userID)
		if err != nil {
			log.Printf("failed to get user by id: %v", err)
			utils.PermissionDenied(w)
			return
		}

		// Add the user to the context
		ctx := r.Context()
		ctx = context.WithValue(ctx, types.UserKey, u.ID)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}
