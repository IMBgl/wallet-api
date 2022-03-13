package service

import "errors"

var ErrEmailAlreadyInUse = errors.New("Email already in use")
var ErrInvalidCredentials = errors.New("Invalid credentials")
var ErrUserNotFound = errors.New("User not found")
var ErrNotFound = errors.New("Not found")
var ErrCategoryNotFound = errors.New("Category not found")
var ErrWalletNotFound = errors.New("Wallet not found")
