package server

import (
	"net/http"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/nabishec/ozon_habr_api/graph"
	commentmutation "github.com/nabishec/ozon_habr_api/internal/handlers/comment_mutation"
	commentquery "github.com/nabishec/ozon_habr_api/internal/handlers/comment_query"
	postmutation "github.com/nabishec/ozon_habr_api/internal/handlers/post_mutation"
	postquery "github.com/nabishec/ozon_habr_api/internal/handlers/post_query"
	"github.com/nabishec/ozon_habr_api/internal/storage"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func RunServer(storage storage.StorageImp) {
	op := "cmd.server.RunServer()"
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = defaultPort
	}

	postMutation := postmutation.NewPostMutation(storage)
	postQuery := postquery.NewPostQuery(storage)
	commentMutation := commentmutation.NewCommentMutation(storage)
	commentQuery := commentquery.NewCommentQuery(storage)

	resolver := graph.NewResolver(postMutation, postQuery, commentMutation, commentQuery)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.GRAPHQL{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Info().Msgf("Connect to http://localhost:%s/ for GraphQL playground", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Error().AnErr(op, err).Msg("Failed to start server")
		os.Exit(1)
	}

	log.Error().Msg("Unknown error")
}
