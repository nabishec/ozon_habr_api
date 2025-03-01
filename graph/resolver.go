package graph

import (
	postmutation "github.com/nabishec/ozon_habr_api/internal/handlers/post_mutation"
	postquery "github.com/nabishec/ozon_habr_api/internal/handlers/post_query"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	PostMutation *postmutation.PostMutation
	PostQuery    *postquery.PostQuery
}

func NewResolver(postMutation *postmutation.PostMutation, postQuery *postquery.PostQuery) *Resolver {
	return &Resolver{
		PostMutation: postMutation,
		PostQuery:    postQuery,
	}
}
