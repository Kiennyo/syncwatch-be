package users

import "errors"

var errDuplicateEmail = errors.New("duplicate email")
var errUserNotFound = errors.New("user not found")
