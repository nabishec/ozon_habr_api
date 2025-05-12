package postmutation

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nabishec/ozon_habr_api/internal/model"
	"github.com/nabishec/ozon_habr_api/internal/pkg/errs"
	"github.com/rs/zerolog/log"
)

type PostMutation struct {
	postMutImp PostMutImp
}

func NewPostMutation(postImp PostMutImp) *PostMutation {
	return &PostMutation{postMutImp: postImp}
}

func (h *PostMutation) AddPost(ctx context.Context, newPost *model.NewPost) (*model.Post, error) {
	op := "internal.handlers.postmutation.AddPost()"

	log.Debug().Msgf("%s start", op)

	post, err := h.postMutImp.AddPost(ctx, newPost)

	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	log.Debug().Msgf("%s end", op)
	return post, nil
}

func (h *PostMutation) UpdateEnableCommentToPost(ctx context.Context, postID int64, authorID uuid.UUID, commentsEnabled bool) (*model.Post, error) {
	op := "internal.handlers.postmutation.UpdateEnableCommentToPost()"

	log.Debug().Msgf("%s start", op)

	post, err := h.postMutImp.UpdateEnableCommentToPost(ctx, postID, authorID, commentsEnabled)

	if err != nil {
		if err == errs.ErrPostNotExist || err == errs.ErrUnauthorizedAccess {
			return nil, err
		}
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	log.Debug().Msgf("%s end", op)
	return post, nil
}
