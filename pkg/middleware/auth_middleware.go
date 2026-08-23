// pkg/middleware/auth_middleware.go

package middleware

import (
    "context"
    "net/http"
    "strings"

    "github.com/everest/bheri/pkg/helpers"
    "github.com/everest/bheri/pkg/utils"
)

func Auth(
    next http.Handler,
) http.Handler {

    return http.HandlerFunc(
        func(w http.ResponseWriter, r *http.Request) {
            var tokenString string

            // 1. Try reading the token from the HTTP-only cookie first
            if cookie, err := r.Cookie("token"); err == nil && cookie.Value != "" {
                tokenString = cookie.Value
            } else {
                // 2. Fallback to Authorization: Bearer <token> header
                auth := r.Header.Get("Authorization")
                if auth != "" {
                    parts := strings.SplitN(auth, " ", 2)
                    if len(parts) == 2 && parts[0] == "Bearer" {
                        tokenString = parts[1]
                    }
                }
            }

            // If no token is found in either place, reject
            if tokenString == "" {
                utils.Error(
                    w,
                    http.StatusUnauthorized,
                    "authorization required",
                )
                return
            }

            // 3. Validate the JWT token and get the user ID
            userID, err := helpers.ValidateToken(tokenString)
            if err != nil {
                utils.Error(
                    w,
                    http.StatusUnauthorized,
                    "invalid or expired token",
                )
                return
            }

            // 4. Inject the validated user ID into the request context
            // Note: using string "user_id" to match your controller lookup
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
