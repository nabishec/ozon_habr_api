package postmutation

import (
	"context"
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/nabishec/ozon_habr_api/internal/model"
	"github.com/nabishec/ozon_habr_api/internal/pkg/errs"
	"github.com/stretchr/testify/assert"
)

func TestPostMutation(t *testing.T) {
	mc := minimock.NewController(t)

	postMutImpMock := NewPostMutImpMock(mc)
	handler := PostMutation{postMutImp: postMutImpMock}

	t.Run("Succesfully add post", func(t *testing.T) {
		ctx := context.Background()
		newPost := &model.NewPost{
			AuthorID:        uuid.New(),
			Title:           "Test Title",
			Text:            "Test Content",
			CommentsEnabled: true,
		}
		expectedPost := &model.Post{
			ID:              1,
			AuthorID:        newPost.AuthorID,
			Title:           newPost.Title,
			Text:            newPost.Text,
			CommentsEnabled: newPost.CommentsEnabled,
		}

		postMutImpMock.AddPostMock.Expect(ctx, newPost).Return(expectedPost, nil)
		post, err := handler.AddPost(ctx, newPost)
		assert.NoError(t, err)
		assert.Equal(t, post, expectedPost)
	})

	t.Run("Unexpexted error", func(t *testing.T) {
		ctx := context.Background()
		newPost := &model.NewPost{
			AuthorID:        uuid.New(),
			Title:           "Test Title",
			Text:            "Test Content",
			CommentsEnabled: true,
		}

		postMutImpMock.AddPostMock.Expect(ctx, newPost).Return(nil, errors.New("unexpected error"))
		post, err := handler.AddPost(ctx, newPost)
		assert.NotNil(t, err)
		assert.Nil(t, post)
	})

	t.Run("Succesfully update enable comment to post", func(t *testing.T) {
		ctx := context.Background()
		postID := int64(1)
		authorID := uuid.New()
		commentsEnabled := true
		expectedPost := &model.Post{
			ID:              postID,
			AuthorID:        authorID,
			Title:           "some title",
			Text:            "some text",
			CommentsEnabled: commentsEnabled,
		}

		postMutImpMock.UpdateEnableCommentToPostMock.Expect(ctx, postID, authorID, commentsEnabled).Return(expectedPost, nil)
		post, err := handler.UpdateEnableCommentToPost(ctx, postID, authorID, commentsEnabled)
		assert.NoError(t, err)
		assert.Equal(t, post, expectedPost)
	})

	t.Run("Error post not exist ", func(t *testing.T) {
		ctx := context.Background()
		postID := int64(1)
		authorID := uuid.New()
		commentsEnabled := true

		postMutImpMock.UpdateEnableCommentToPostMock.Expect(ctx, postID, authorID, commentsEnabled).Return(nil, errs.ErrPostNotExist)
		post, err := handler.UpdateEnableCommentToPost(ctx, postID, authorID, commentsEnabled)
		assert.Equal(t, err, errs.ErrPostNotExist)
		assert.Nil(t, post)
	})

	t.Run("Error author hass not access", func(t *testing.T) {
		ctx := context.Background()
		postID := int64(1)
		authorID := uuid.New()
		commentsEnabled := true

		postMutImpMock.UpdateEnableCommentToPostMock.Expect(ctx, postID, authorID, commentsEnabled).Return(nil, errs.ErrUnauthorizedAccess)
		post, err := handler.UpdateEnableCommentToPost(ctx, postID, authorID, commentsEnabled)
		assert.Equal(t, err, errs.ErrUnauthorizedAccess)
		assert.Nil(t, post)
	})

	t.Run("Unexpected error", func(t *testing.T) {
		ctx := context.Background()
		postID := int64(1)
		authorID := uuid.New()
		commentsEnabled := true

		postMutImpMock.UpdateEnableCommentToPostMock.Expect(ctx, postID, authorID, commentsEnabled).Return(nil, errors.New("unexpected error"))
		post, err := handler.UpdateEnableCommentToPost(ctx, postID, authorID, commentsEnabled)
		assert.NotNil(t, err)
		assert.Nil(t, post)
	})
}
