package forms

import (
	"github.com/google/uuid"
	"time"
)

type SignupRequest struct {
	Name     string `form:"name" json:"name" binding:"required,max=120"`
	Email    string `form:"email" json:"email" binding:"required,email,max=120"`
	Password string `form:"password" json:"password" binding:"required,min=6"`
}

type SigninRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type SigninResponse struct {
	SessionID             uuid.UUID `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}
