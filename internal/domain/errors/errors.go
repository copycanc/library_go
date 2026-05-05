package errors

import "errors"

var (
	ErrBooksListIsEmpty = errors.New("books list is empty")
	ErrUserExist        = errors.New("user already exists")
	ErrInvalidCreds     = errors.New("invalid email or password")
	ErrBookNotFound     = errors.New("book not found")
)
