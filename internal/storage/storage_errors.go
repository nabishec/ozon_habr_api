package storage

import "errors"

var (
	ErrPostNotExist       = errors.New("post not exist")
	ErrUnauthorizedAccess = errors.New("user doesn't have access rights")
	ErrPostsNotExist      = errors.New("no posts have been created yet")
	ErrCommentsNotExist   = errors.New("no comments have been created yet")
	ErrPostNotCached      = errors.New("post not cached yet")
	ErrPathNotExist       = errors.New("path not exist")
)
