package middleware

import (
	"RMS/handler"
	"RMS/models"
	"RMS/utilities"
	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")

		claims := models.Claims{}

		tkn, err1 := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
			return handler.JwtKey, nil
		})
		if err1 != nil {
			if err1 == jwt.ErrSignatureInvalid {
				logrus.Printf("Signature invalid:%v", err1)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			logrus.Printf("ParseErr:%v", err1)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// FetchRole, err := helper.FetchRole(userID)
		// if err !=nil{
		//	w.WriteHeader(http.StatusBadRequest)
		//	log.Printf("Middleware: fetch role error:%v",err)
		//	return
		// }
		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			logrus.Printf("token is invalid")
			return
		}
		userID := claims.ID
		role := claims.Role

		value := models.ContextValues{ID: userID, Role: role}
		ctx := context.WithValue(r.Context(), utilities.UserContextKey, value)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contextValues, ok := r.Context().Value(utilities.UserContextKey).(models.ContextValues)

		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			logrus.Printf("AdminMiddleware:Context for ID:%v", ok)
			return
		}

		if contextValues.Role != "admin" {
			w.WriteHeader(http.StatusUnauthorized)
			logrus.Printf("Role invalid")
			_, err := w.Write([]byte("ERROR: Role mismatch"))

			if err != nil {
				return
			}

			return
		}

		next.ServeHTTP(w, r)
	})
}

func SubAdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contextValues, ok := r.Context().Value(utilities.UserContextKey).(models.ContextValues)

		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			logrus.Printf("SubAdminMiddleware:Context for ID:%v", ok)
			return
		}

		if contextValues.Role != "sub_admin" {
			w.WriteHeader(http.StatusUnauthorized)
			logrus.Printf("Role invalid")
			_, err := w.Write([]byte("ERROR: Role mismatch"))

			if err != nil {
				return
			}

			return
		}
		next.ServeHTTP(w, r)
	})
}
