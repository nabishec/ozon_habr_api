package postquery

import (
	"context"

	"github.com/nabishec/ozon_habr_api/internal/model"
)

//go:generate minimock -i PostQueryImp
type PostQueryImp interface {
	GetAllPosts(ctx context.Context) ([]*model.Post, error)
	GetPost(ctx context.Context, postID int64) (*model.Post, error)
}
