package commentquery

import (
	"fmt"

	"context"

	"github.com/nabishec/ozon_habr_api/internal/model"
	"github.com/nabishec/ozon_habr_api/internal/pkg/errs"
	"github.com/rs/zerolog/log"
)

type CommentQuery struct {
	commentQueryImp CommentQueryImp
}

func NewCommentQuery(commentQueryImp CommentQueryImp) *CommentQuery {
	return &CommentQuery{commentQueryImp: commentQueryImp}
}

func (h *CommentQuery) GetCommentsBranchToPost(ctx context.Context, postID int64, path string) ([]*model.Comment, error) {
	op := "internal.handlers.commentquery.GetCommentsBranchToPost()"

	log.Debug().Msgf("%s start", op)

	comments, err := h.commentQueryImp.GetCommentsBranch(ctx, postID, path)
	if err != nil {
		if err == errs.ErrCommentsNotExist || err == errs.ErrPathNotExist {
			return nil, err
		}
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	log.Debug().Msgf("%s end", op)
	return comments, nil
}

func (h *CommentQuery) GetPathToComments(ctx context.Context, parentID int64) (string, error) {
	op := "internal.handlers.commentquery.GetPathToComments()"

	log.Debug().Msgf("%s start", op)

	path, err := h.commentQueryImp.GetCommentPath(ctx, parentID)
	if err != nil {
		if err == errs.ErrCommentsNotExist {
			return "", err
		}
		return "", fmt.Errorf("%s:%w", op, err)
	}

	log.Debug().Msgf("%s start", op)

	return path, nil
}
