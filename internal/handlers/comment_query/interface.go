package commentquery

import (
	"context"

	"github.com/nabishec/ozon_habr_api/internal/model"
)

//go:generate minimock -i CommentQueryImp
type CommentQueryImp interface {
	GetCommentsBranch(ctx context.Context, postID int64, path string) ([]*model.Comment, error)
	GetCommentPath(ctx context.Context, parentID int64) (string, error)
}
