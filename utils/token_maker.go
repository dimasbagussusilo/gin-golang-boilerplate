package utils

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Different types of error returned by the VerifyToken function
var (
	ErrInvalidToken = errors.New("tokens is invalid")
	ErrExpiredToken = errors.New("tokens has expired")
)

// TokenPayload contains the payload data of the tokens
type TokenPayload struct {
	Id        uuid.UUID `json:"id"`
	UserId    int       `json:"user_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload creates a new tokens payload with a specific username and duration
func NewPayload(userId int, duration time.Duration) (*TokenPayload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &TokenPayload{
		Id:        tokenId,
		UserId:    userId,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

// Valid checks if the tokens payload is valid or not
func (payload *TokenPayload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

// TokenMaker is an interface for managing tokens
type TokenMaker interface {
	// CreateToken creates a new tokens for a specific username and duration
	CreateToken(userId int, duration time.Duration) (string, *TokenPayload, error)

	// VerifyToken checks if the tokens is valid or not
	VerifyToken(token string) (*TokenPayload, error)
}
