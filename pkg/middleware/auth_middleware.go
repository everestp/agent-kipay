// pkg/middleware/auth_middleware.go

package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/everest/bheri/pkg/utils"
)

type contextKey string

const userIDKey contextKey = "user_id"

func Auth(
	next http.Handler,
) http.Handler {

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			auth := r.Header.Get("Authorization")

			if auth == "" {
				utils.Error(
					w,
					http.StatusUnauthorized,
					"authorization required",
				)
				return
			}

			parts := strings.SplitN(auth, " ", 2)

			if len(parts) != 2 ||
				parts[0] != "Bearer" {

				utils.Error(
					w,
					http.StatusUnauthorized,
					"invalid authorization",
				)
				return
			}

			userID := parts[1]

			ctx := context.WithValue(
				r.Context(),
				"user_id",
				userID,
			)

			next.ServeHTTP(
				w,
				r.WithContext(ctx),
			)
		},
	)
}

func Logger(
	next http.Handler,
) http.Handler {

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			next.ServeHTTP(w, r)
		},
	)
}
