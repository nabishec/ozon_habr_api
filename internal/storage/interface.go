package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/nabishec/ozon_habr_api/internal/model"
)

type StorageImp interface {
	AddPost(ctx context.Context, newPost *model.NewPost) (*model.Post, error)
	AddComment(ctx context.Context, postID int64, newComment *model.NewComment) (*model.Comment, error)
	UpdateEnableCommentToPost(ctx context.Context, postID int64, authorID uuid.UUID, commentsEnabled bool) (*model.Post, error)
	GetAllPosts(ctx context.Context) ([]*model.Post, error)
	GetPost(ctx context.Context, postID int64) (*model.Post, error)
	GetCommentsBranch(postID int64, path string) ([]*model.Comment, error)
	GetCommentPath(parentID int64) (string, error)
}
