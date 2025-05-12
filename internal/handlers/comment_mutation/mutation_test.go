package commentmutation

import (
	"context"
	"strings"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/nabishec/ozon_habr_api/internal/model"
	"github.com/nabishec/ozon_habr_api/internal/pkg/errs"
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
		if err != nil {
			t.Fatalf("got error - %v", err)
		}

		if comment.ID != expectedComment.ID || comment.Text != expectedComment.Text {
			t.Fatalf("expected %v, got %v", expectedComment, comment)
		}
	})

	t.Run("Error comment long", func(t *testing.T) {
		ctx := context.Background()
		postID := int64(1)
		newComment := &model.NewComment{
			Text: strings.Repeat("f", 2001),
		}

		comment, err := handler.AddComment(ctx, postID, newComment)
		if err != errs.ErrIncorrectCommentLength || comment != nil {
			t.Fatalf("unexpected error - %v", err)
		}

	})

	t.Run("Error comment empty", func(t *testing.T) {
		ctx := context.Background()
		postID := int64(1)
		newComment := &model.NewComment{
			Text: "",
		}

		comment, err := handler.AddComment(ctx, postID, newComment)
		if err != errs.ErrIncorrectCommentLength || comment != nil {
			t.Fatalf("unexpected error - %v", err)
		}

	})

	t.Run("Error post not exist", func(t *testing.T) {
		ctx := context.Background()
		postID := int64(1)
		newComment := &model.NewComment{
			Text: "Correct comment",
		}

		commentMutationImpMock.AddCommentMock.Expect(ctx, postID, newComment).Return(nil, errs.ErrPostNotExist)

		comment, err := handler.AddComment(ctx, postID, newComment)
		if err != errs.ErrPostNotExist || comment != nil {
			t.Fatalf("unexpected error - %v", err)
		}

	})

	t.Run("Error post not exist", func(t *testing.T) {
		ctx := context.Background()
		postID := int64(1)
		newComment := &model.NewComment{
			Text: "Correct comment",
		}

		commentMutationImpMock.AddCommentMock.Expect(ctx, postID, newComment).Return(nil, errs.ErrParentCommentNotExist)

		comment, err := handler.AddComment(ctx, postID, newComment)
		if err != errs.ErrParentCommentNotExist || comment != nil {
			t.Fatalf("unexpected error - %v", err)
		}

	})

	t.Run("Error post not exist", func(t *testing.T) {
		ctx := context.Background()
		postID := int64(1)
		newComment := &model.NewComment{
			Text: "Correct comment",
		}

		commentMutationImpMock.AddCommentMock.Expect(ctx, postID, newComment).Return(nil, errs.ErrCommentsNotEnabled)

		comment, err := handler.AddComment(ctx, postID, newComment)
		if err != errs.ErrCommentsNotEnabled || comment != nil {
			t.Fatalf("unexpected error - %v", err)
		}

	})

}
