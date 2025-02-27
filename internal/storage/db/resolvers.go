package db

import (
	"fmt"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/jmoiron/sqlx"
	"github.com/nabishec/ozon_habr_api/internal/model"
	"github.com/rs/zerolog/log"
)

type Storage struct {
	db *sqlx.DB

	cache *cache.Cache
}

func NewStorage(db *sqlx.DB, cache *cache.Cache) *Storage {
	return &Storage{
		db:    db,
		cache: cache,
	}
}

func (r *Storage) AddPost(newPost *model.NewPost) (*model.Post, error) {
	op := "internal.storage.db.NewPost()"

	log.Debug().Msgf("%s start", op)

	post := &model.Post{
		AuthorID:        newPost.AuthorID,
		Title:           newPost.Title,
		Text:            newPost.Text,
		CommentsEnabled: newPost.CommentsEnabled,
		CreateDate:      time.Now(),
	}

	queryNewPost := `INSERT INTO Posts (author_id, title, text, comments_enabled, create_date)
						VALUES ($1, $2, $3, $4, $5)
						RETURNING post_id `

	err := r.db.QueryRow(queryNewPost, post.AuthorID, post.Title, post.Text, post.CommentsEnabled, post.CreateDate).Scan(&post.ID)

	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	log.Debug().Msgf("%s end", op)
	return post, err
}
