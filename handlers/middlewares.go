package handlers

import (
	"context"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

func (h *HandlerService) AuthGaurd(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("auth")
		secret := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil {
			WriteResponse(w, http.StatusForbidden, "forbidden resources")
			return
		}
		if !token.Valid {
			WriteResponse(w, http.StatusForbidden, "forbidden resources")
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		if claims["id"] == nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), "userId", claims["id"])
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
