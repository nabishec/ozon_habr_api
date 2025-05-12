package postmutation

import (
	"context"

	"github.com/google/uuid"
	"github.com/nabishec/ozon_habr_api/internal/model"
)

type PostMutImp interface {
	AddPost(ctx context.Context, newPost *model.NewPost) (*model.Post, error)
	UpdateEnableCommentToPost(ctx context.Context, postID int64, authorID uuid.UUID, commentsEnabled bool) (*model.Post, error)
}
