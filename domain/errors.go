package domain

import "errors"

var (
	ErrNotFound            = errors.New("resource not found")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrInvalidPaymentMode  = errors.New("invalid payment mode")
	ErrInvalidCategory     = errors.New("invalid category")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrInvalidInput        = errors.New("invalid input")
)
