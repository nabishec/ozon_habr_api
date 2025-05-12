package commentmutation

import (
	"context"

	"github.com/nabishec/ozon_habr_api/internal/model"
)

type CommentMutationImp interface {
	AddComment(ctx context.Context, postID int64, newComment *model.NewComment) (*model.Comment, error)
}
