package commentquery

import (
	"context"
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/nabishec/ozon_habr_api/internal/model"
	"github.com/nabishec/ozon_habr_api/internal/pkg/errs"
	"github.com/stretchr/testify/assert"
)

func TestCommentQuery(t *testing.T) {
	mc := minimock.NewController(t)

	commentQueryImpMock := NewCommentQueryImpMock(mc)
	handler := CommentQuery{commentQueryImp: commentQueryImpMock}

	t.Run("Successfully get comments branch", func(t *testing.T) {
		ctx := context.Background()
		postID := int64(1)
		path := "1.2.3"

		commentQueryImpMock.GetCommentsBranchMock.Expect(ctx, postID, path).Return([]*model.Comment{}, nil)
		comments, err := handler.GetCommentsBranchToPost(ctx, postID, path)
		assert.NoError(t, err)
		assert.NotNil(t, comments)
	})

	t.Run("Error comments to post not exist", func(t *testing.T) {
		ctx := context.Background()
		postID := int64(123) // number of post that doesn't exist
		path := "1.2.3"

		commentQueryImpMock.GetCommentsBranchMock.Expect(ctx, postID, path).Return(nil, errs.ErrCommentsNotExist)
		comments, err := handler.GetCommentsBranchToPost(ctx, postID, path)
		assert.Nil(t, comments)
		assert.Equal(t, err, errs.ErrCommentsNotExist)
	})

	t.Run("Error path not exist", func(t *testing.T) {
		ctx := context.Background()
		postID := int64(1)
		path := "1.2;l.3" // invalid path

		commentQueryImpMock.GetCommentsBranchMock.Expect(ctx, postID, path).Return(nil, errs.ErrPathNotExist)
		comments, err := handler.GetCommentsBranchToPost(ctx, postID, path)
		assert.Nil(t, comments)
		assert.Equal(t, err, errs.ErrPathNotExist)
	})

	t.Run("Unexpected error from get  comment branch", func(t *testing.T) {
		ctx := context.Background()
		postID := int64(1)
		path := "1.2.3"

		commentQueryImpMock.GetCommentsBranchMock.Expect(ctx, postID, path).Return(nil, errors.New("unexpected error"))
		comments, err := handler.GetCommentsBranchToPost(ctx, postID, path)
		assert.NotNil(t, err)
		assert.Nil(t, comments)
	})

	t.Run("Successfully get path to  comments", func(t *testing.T) {
		ctx := context.Background()
		postID := int64(1)

		commentQueryImpMock.GetCommentPathMock.Expect(ctx, postID).Return("correct path", nil)
		path, err := handler.GetPathToComments(ctx, postID)
		assert.NoError(t, err)
		assert.NotNil(t, path)
	})

	t.Run("Error comments not exist", func(t *testing.T) {
		ctx := context.Background()
		postID := int64(1)

		commentQueryImpMock.GetCommentPathMock.Expect(ctx, postID).Return("", errs.ErrCommentsNotExist)
		path, err := handler.GetPathToComments(ctx, postID)
		assert.Equal(t, err, errs.ErrCommentsNotExist)
		assert.Equal(t, path, "")
	})

	t.Run("Unexpected error from get path to comments", func(t *testing.T) {
		ctx := context.Background()
		postID := int64(1)

		commentQueryImpMock.GetCommentPathMock.Expect(ctx, postID).Return("", errors.New("unexpected error"))
		path, err := handler.GetPathToComments(ctx, postID)
		assert.NotNil(t, err)
		assert.Equal(t, path, "")
	})

}
