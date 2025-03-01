package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.66

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nabishec/ozon_habr_api/graph/model"
	internalmodel "github.com/nabishec/ozon_habr_api/internal/model"
	"github.com/nabishec/ozon_habr_api/internal/storage"
	"github.com/rs/zerolog/log"
)

// AddPost is the resolver for the addPost field.
func (r *mutationResolver) AddPost(ctx context.Context, postInput model.NewPost) (*model.Post, error) {
	const op = "graph.AddPost()"

	log.Debug().Msgf("%s start", op)

	newPost := postToInternalModel(&postInput)
	post, err := r.PostMutation.AddPost(newPost)

	if err != nil {
		log.Error().Err(err).Msgf("%s end with error", op)
		if err != storage.ErrPostNotExist || err != storage.ErrUnauthorizedAccess {
			err = errors.New("internal server error")
		}

		return nil, err
	}

	log.Debug().Msgf("%s end", op)
	return postFromInternalModel(post), err
}

func postToInternalModel(postInput *model.NewPost) *internalmodel.NewPost {
	return &internalmodel.NewPost{
		AuthorID:        postInput.AuthorID,
		Title:           postInput.Title,
		Text:            postInput.Text,
		CommentsEnabled: postInput.CommentsEnabled,
	}
}

func postFromInternalModel(postInput *internalmodel.Post) *model.Post {
	return &model.Post{
		ID:              postInput.ID,
		AuthorID:        postInput.AuthorID,
		Title:           postInput.Title,
		Text:            postInput.Text,
		CommentsEnabled: postInput.CommentsEnabled,
		CreateDate:      postInput.CreateDate,
	}
}

// AddComment is the resolver for the addComment field.
func (r *mutationResolver) AddComment(ctx context.Context, commentInput model.NewComment) (*model.Comment, error) {
	panic(fmt.Errorf("not implemented: AddComment - addComment"))
}

// UpdateEnableComment is the resolver for the updateEnableComment field.
func (r *mutationResolver) UpdateEnableComment(ctx context.Context, postID int64, authorID uuid.UUID, commentsEnabled bool) (*model.Post, error) {
	const op = "graph.UpdateEnableComment()"

	log.Debug().Msgf("%s start", op)

	post, err := r.PostMutation.UpdateEnableCommentToPost(postID, authorID, commentsEnabled)

	if err != nil {
		log.Error().Err(err).Msgf("%s end with error", op)
		if err != storage.ErrPostNotExist || err != storage.ErrUnauthorizedAccess {
			err = errors.New("internal server error")
		}
		return nil, err
	}

	log.Debug().Msgf("%s end", op)
	return postFromInternalModel(post), err
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
	panic(fmt.Errorf("not implemented: Posts - posts"))
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, postID int64) (*model.Post, error) {
	panic(fmt.Errorf("not implemented: Post - post"))
}

// CommentAdded is the resolver for the commentAdded field.
func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID int64) (<-chan *model.Comment, error) {
	panic(fmt.Errorf("not implemented: CommentAdded - commentAdded"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
