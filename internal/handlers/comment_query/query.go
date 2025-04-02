package commentquery

import (
	"fmt"

	"github.com/nabishec/ozon_habr_api/internal/model"
	"github.com/nabishec/ozon_habr_api/internal/storage"
	"github.com/rs/zerolog/log"
)

type CommentQuery struct {
	commentQueryImp CommentQueryImp
}

func NewCommentQuery(commentQueryImp CommentQueryImp) *CommentQuery {
	return &CommentQuery{commentQueryImp: commentQueryImp}
}

func (h *CommentQuery) GetCommentsBranchToPost(postID int64, path string) ([]*model.Comment, error) {
	op := "internal.storage.db.GetPostWithComment()"

	log.Debug().Msgf("%s start", op)

	comments, err := h.commentQueryImp.GetCommentsBranch(postID, path)
	if err != nil {
		if err == storage.ErrCommentsNotExist {
			return nil, err
		} else if err == storage.ErrPathNotExist {
			return nil, err
		}
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	log.Debug().Msgf("%s end", op)
	return comments, nil
}

func (h *CommentQuery) GetPathToComments(parentID int64) (string, error) {
	op := "internal.storage.db.GetPostWithComment()"

	log.Debug().Msgf("%s start", op)

	path, err := h.commentQueryImp.GetCommentPath(parentID)
	if err != nil {
		if err == storage.ErrCommentsNotExist {
			return "", storage.ErrCommentsNotExist
		}
		return "", fmt.Errorf("%s:%w", op, err)
	}

	log.Debug().Msgf("%s start", op)

	return path, nil
}
