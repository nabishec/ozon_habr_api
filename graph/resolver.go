package graph

import postmutation "github.com/nabishec/ozon_habr_api/internal/handlers/post_mutation"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	PostMutation *postmutation.PostMutation
}

func NewResolver(postMutation *postmutation.PostMutation) *Resolver {
	return &Resolver{
		PostMutation: postMutation,
	}
}
