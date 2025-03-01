package storage

import "errors"

var (
	ErrPostNotExist       = errors.New("post not exist")
	ErrUnauthorizedAccess = errors.New("user doesn't have access rights")
	ErrPostsNotExist      = errors.New("no posts have been created yet")
)
