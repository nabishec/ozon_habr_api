package storage

import (
	"github.com/google/uuid"
	"github.com/nabishec/ozon_habr_api/internal/model"
)

type StorageImp interface {
	AddPost(newPost *model.NewPost) (*model.Post, error)
	UpdateEnableCommentToPost(postID int64, authorID uuid.UUID, commentsEnabled bool) (*model.Post, error)
	GetAllPosts() ([]*model.Post, error)
	GetPost(postID int64) (*model.Post, error)
}
