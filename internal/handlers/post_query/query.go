package postquery

import (
	"fmt"

	"github.com/nabishec/ozon_habr_api/internal/model"
	"github.com/nabishec/ozon_habr_api/internal/storage"
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
		if err == storage.ErrPostsNotExist {
			return nil, err
		}
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	log.Debug().Msgf("%s end", op)
	return posts, nil
}

func (h *PostQuery) GetPostWithComment(postID int64) (*model.Post, error) {
	op := "internal.storage.db.GetPostWithComment()"

	log.Debug().Msgf("%s start", op)

	post, err := h.postQueryImp.GetPost(postID)

	if err != nil {
		if err == storage.ErrPostNotExist {
			return nil, err
		}
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	log.Debug().Msgf("%s end", op)
	return post, nil
}
