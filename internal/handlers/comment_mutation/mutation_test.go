package commentmutation

import (
	"context"
	"strings"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/nabishec/ozon_habr_api/internal/model"
	"github.com/nabishec/ozon_habr_api/internal/pkg/errs"
	"github.com/stretchr/testify/assert"
)

func TestAddComment(t *testing.T) {

	mc := minimock.NewController(t)

	commentMutationImpMock := NewCommentMutationImpMock(mc)

	handler := CommentMutation{commentMutationImp: commentMutationImpMock}

	t.Run("Successfully add comment", func(t *testing.T) {
		ctx := context.Background()
		postID := int64(1)
		newComment := &model.NewComment{
			Text: "Correct comment",
		}

		expectedComment := &model.Comment{
			ID:   1,
			Text: newComment.Text,
		}

		commentMutationImpMock.AddCommentMock.Expect(ctx, postID, newComment).Return(expectedComment, nil)

		comment, err := handler.AddComment(ctx, postID, newComment)
		assert.NoError(t, err)

		assert.Equal(t, expectedComment, comment)
	})

	t.Run("Error comment long", func(t *testing.T) {
		ctx := context.Background()
		postID := int64(1)
		newComment := &model.NewComment{
			Text: strings.Repeat("f", 2001),
		}

		comment, err := handler.AddComment(ctx, postID, newComment)
		assert.Equal(t, errs.ErrIncorrectCommentLength, err)
		assert.Nil(t, comment)

	})

	t.Run("Error comment empty", func(t *testing.T) {
		ctx := context.Background()
		postID := int64(1)
		newComment := &model.NewComment{
			Text: "",
		}

		comment, err := handler.AddComment(ctx, postID, newComment)
		assert.Equal(t, errs.ErrIncorrectCommentLength, err)
		assert.Nil(t, comment)

	})

	t.Run("Error post not exist", func(t *testing.T) {
		ctx := context.Background()
		postID := int64(1)
		newComment := &model.NewComment{
			Text: "Correct comment",
		}

		commentMutationImpMock.AddCommentMock.Expect(ctx, postID, newComment).Return(nil, errs.ErrPostNotExist)

		comment, err := handler.AddComment(ctx, postID, newComment)
		assert.Equal(t, errs.ErrPostNotExist, err)
		assert.Nil(t, comment)

	})

	t.Run("Error parent comment not exist", func(t *testing.T) {
		ctx := context.Background()
		postID := int64(1)
		newComment := &model.NewComment{
			Text: "Correct comment",
		}

		commentMutationImpMock.AddCommentMock.Expect(ctx, postID, newComment).Return(nil, errs.ErrParentCommentNotExist)

		comment, err := handler.AddComment(ctx, postID, newComment)
		assert.Equal(t, errs.ErrParentCommentNotExist, err)
		assert.Nil(t, comment)

	})

	t.Run("Errorcomment not enabled", func(t *testing.T) {
		ctx := context.Background()
		postID := int64(1)
		newComment := &model.NewComment{
			Text: "Correct comment",
		}

		commentMutationImpMock.AddCommentMock.Expect(ctx, postID, newComment).Return(nil, errs.ErrCommentsNotEnabled)

		comment, err := handler.AddComment(ctx, postID, newComment)
		assert.Equal(t, errs.ErrCommentsNotEnabled, err)
		assert.Nil(t, comment)

	})

}
