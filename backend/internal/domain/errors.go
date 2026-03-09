package domain

import "errors"

var (
	ErrNotFound       = errors.New("resource not found")
	ErrAlreadyExists  = errors.New("resource already exists")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrValidation     = errors.New("validation error")
	ErrInternal       = errors.New("internal error")
	ErrInvalidInput   = errors.New("invalid input")
	ErrExpiredToken   = errors.New("token expired")
	ErrInvalidToken   = errors.New("invalid token")
	ErrRateLimited    = errors.New("rate limited")
	ErrPaymentFailed  = errors.New("payment failed")
	ErrStorageFull    = errors.New("storage quota exceeded")
)
