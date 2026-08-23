// modules/user/controller/user_controller.go

package controller

import (
    "encoding/json"
    "net/http"
    "time"

    "github.com/everest/bheri/modules/user/dto"
    "github.com/everest/bheri/modules/user/service"
    "github.com/everest/bheri/pkg/utils"
)

type UserController struct {
    service service.UserService
}

func NewUserController(
    service service.UserService,
) *UserController {
    return &UserController{
        service: service,
    }
}

func (c *UserController) Register(
    w http.ResponseWriter,
    r *http.Request,
) {
    var request dto.RegisterRequest

    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        utils.Error(
            w,
            http.StatusBadRequest,
            "invalid request body",
        )
        return
    }

    user, token, err := c.service.Register(
        r.Context(),
        request,
    )

    if err != nil {
        utils.Error(
            w,
            http.StatusBadRequest,
            err.Error(),
        )
        return
    }

    // Set token in HTTP-only cookie (secure for browsers)
    setAuthCookie(w, token)

    // Return token in JSON response as well (if your frontend needs it)
    utils.Success(
        w,
        http.StatusCreated,
        dto.LoginResponse{ // Or dto.RegisterResponse if you have a separate one
            Token: token,
            User: dto.UserResponse{
                ID:        user.ID,
                Name:      user.Name,
                Email:     user.Email,
                Status:    user.Status,
                CreatedAt: user.CreatedAt.String(),
            },
        },
    )
}

func (c *UserController) Login(
    w http.ResponseWriter,
    r *http.Request,
) {
    var request dto.LoginRequest

    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        utils.Error(
            w,
            http.StatusBadRequest,
            "invalid request body",
        )
        return
    }

    user, token, err := c.service.Login(
        r.Context(),
        request,
    )

    if err != nil {
        utils.Error(
            w,
            http.StatusUnauthorized,
            "invalid credentials",
        )
        return
    }

    // Set token in HTTP-only cookie
    setAuthCookie(w, token)

    // Populate the Token field in the response DTO
    utils.Success(
        w,
        http.StatusOK,
        dto.LoginResponse{
            Token: token, // <-- Make sure this line is here!
            User: dto.UserResponse{
                ID:        user.ID,
                Name:      user.Name,
                Email:     user.Email,
                Status:    user.Status,
                CreatedAt: user.CreatedAt.String(),
            },
        },
    )
}
func (c *UserController) Logout(
    w http.ResponseWriter,
    r *http.Request,
) {
    // Clear the auth cookie by setting MaxAge to -1
    http.SetCookie(w, &http.Cookie{
        Name:     "token",
        Value:    "",
        Path:     "/",
        HttpOnly: true,
        Secure:   false, // Set to true if using HTTPS in production
        SameSite: http.SameSiteLaxMode,
        Expires:  time.Unix(0, 0),
        MaxAge:   -1,
    })

    utils.Success(
        w,
        http.StatusOK,
        map[string]string{"message": "logged out successfully"},
    )
}

func (c *UserController) Me(
    w http.ResponseWriter,
    r *http.Request,
) {
    userID := r.Context().Value("user_id")

    if userID == nil {
        utils.Error(
            w,
            http.StatusUnauthorized,
            "unauthorized",
        )
        return
    }

    user, err := c.service.GetByID(
        r.Context(),
        userID.(string),
    )

    if err != nil {
        utils.Error(
            w,
            http.StatusNotFound,
            "user not found",
        )
        return
    }

    utils.Success(
        w,
        http.StatusOK,
        user,
    )
}

// Helper function to manage cookie options cleanly
func setAuthCookie(w http.ResponseWriter, token string) {
    http.SetCookie(w, &http.Cookie{
        Name:     "token",
        Value:    token,
        Path:     "/",
        HttpOnly: true,
        Secure:   false, // Change to true in production (requires HTTPS)
        SameSite: http.SameSiteLaxMode,
        MaxAge:   86400 * 7, // 7 days expiration
    })
}
