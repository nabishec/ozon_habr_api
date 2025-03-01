package postquery

import (
	"fmt"

	"github.com/nabishec/ozon_habr_api/internal/model"
	"github.com/rs/zerolog/log"
)

type PostQuery struct {
	postQueryImp PostQueryImp
}

func NewPostQuery(postQueryImp PostQueryImp) *PostQuery {
	return &PostQuery{postQueryImp: postQueryImp}
}

func (h *PostQuery) GetAllPosts() ([]*model.Post, error) {
	op := "internal.storage.db.GetAllPosts()"

	log.Debug().Msgf("%s start", op)

	posts, err := h.postQueryImp.GetAllPosts()

	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	log.Debug().Msgf("%s end", op)
	return posts, nil
}
