package postquery

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

func TestPostQuery(t *testing.T) {
	mc := minimock.NewController(t)

	postQueryImpMock := NewPostQueryImpMock(mc)
	handler := PostQuery{postQueryImp: postQueryImpMock}

	t.Run("Succesfully get all posts", func(t *testing.T) {
		ctx := context.Background()
		expectedPosts := []*model.Post{
			{
				ID:              1,
				AuthorID:        uuid.New(),
				Title:           "Test Title",
				Text:            "Test Content",
				CommentsEnabled: true,
			},
		}

		postQueryImpMock.GetAllPostsMock.Expect(ctx).Return(expectedPosts, nil)
		posts, err := handler.GetAllPosts(ctx)
		assert.NoError(t, err)
		assert.Equal(t, posts, expectedPosts)
	})

	t.Run("Error posts not exist", func(t *testing.T) {
		ctx := context.Background()

		postQueryImpMock.GetAllPostsMock.Expect(ctx).Return(nil, errs.ErrPostsNotExist)
		posts, err := handler.GetAllPosts(ctx)
		assert.Nil(t, posts)
		assert.Equal(t, err, errs.ErrPostsNotExist)
	})

	t.Run("Unexpected error", func(t *testing.T) {
		ctx := context.Background()

		postQueryImpMock.GetAllPostsMock.Expect(ctx).Return(nil, errors.New("unexpected error"))
		posts, err := handler.GetAllPosts(ctx)
		assert.Nil(t, posts)
		assert.NotNil(t, err)
	})

	t.Run("Succesfully get  post", func(t *testing.T) {
		ctx := context.Background()
		postID := int64(1)
		expectedPost := &model.Post{

			ID:              1,
			AuthorID:        uuid.New(),
			Title:           "Test Title",
			Text:            "Test Content",
			CommentsEnabled: true,
		}

		postQueryImpMock.GetPostMock.Expect(ctx, postID).Return(expectedPost, nil)
		post, err := handler.GetPost(ctx, postID)
		assert.NoError(t, err)
		assert.Equal(t, post, expectedPost)
	})

	t.Run("Error post not exist", func(t *testing.T) {
		ctx := context.Background()
		postID := int64(123) // not existing post ID

		postQueryImpMock.GetPostMock.Expect(ctx, postID).Return(nil, errs.ErrPostNotExist)
		post, err := handler.GetPost(ctx, postID)
		assert.Nil(t, post)
		assert.Equal(t, err, errs.ErrPostNotExist)
	})

	t.Run("Unexpected error", func(t *testing.T) {
		ctx := context.Background()
		postID := int64(1)

		postQueryImpMock.GetPostMock.Expect(ctx, postID).Return(nil, errors.New("unexpected error"))
		post, err := handler.GetPost(ctx, postID)
		assert.Nil(t, post)
		assert.NotNil(t, err)
	})
}
