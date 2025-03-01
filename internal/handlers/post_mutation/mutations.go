package postmutation

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/nabishec/ozon_habr_api/internal/model"
	"github.com/rs/zerolog/log"
)

type PostMutation struct {
	postImp PostImp
}

func NewPostMutation(postImp PostImp) *PostMutation {
	return &PostMutation{postImp: postImp}
}

func (h *PostMutation) AddPost(newPost *model.NewPost) (*model.Post, error) {
	op := "internal.storage.db.AddPost()"

	log.Debug().Msgf("%s start", op)

	post, err := h.postImp.AddPost(newPost)

	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	log.Debug().Msgf("%s end", op)
	return post, nil
}

func (h *PostMutation) UpdateEnableCommentToPost(postID int64, authorID uuid.UUID, commentsEnabled bool) (*model.Post, error) {
	op := "internal.storage.db.UpdateEnableCommentToPost()"

	log.Debug().Msgf("%s start", op)

	post, err := h.postImp.UpdateEnableCommentToPost(postID, authorID, commentsEnabled)

	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	log.Debug().Msgf("%s end", op)
	return post, nil
}
