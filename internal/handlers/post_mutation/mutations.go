package postmutation

import (
	"fmt"

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
	op := "internal.storage.db.NewPost()"

	log.Debug().Msgf("%s start", op)

	post, err := h.postImp.AddPost(newPost)

	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	log.Debug().Msgf("%s end", op)
	return post, nil
}
