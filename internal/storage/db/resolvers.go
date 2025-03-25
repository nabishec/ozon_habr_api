package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/nabishec/ozon_habr_api/internal/model"
	"github.com/nabishec/ozon_habr_api/internal/storage"
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

func (r *Storage) UpdateEnableCommentToPost(postID int64, authorID uuid.UUID, commentsEnabled bool) (*model.Post, error) {
	op := "internal.storage.db.UpdateEnableCommentToPost()"

	tx, err := r.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	post := &model.Post{
		ID: postID,
	}

	queryGetPost := `SELECT author_id, title, text, comments_enabled, create_date
						FROM Posts
						WHERE post_id = $1 `

	queryUpdatePost := `UPDATE Posts
							SET comments_enabled = $1
							WHERE post_id = $2`

	err = tx.Get(post, queryGetPost, postID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrPostNotExist
		}
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	if post.AuthorID != authorID {
		return nil, storage.ErrUnauthorizedAccess
	}

	_, err = tx.Exec(queryUpdatePost, commentsEnabled, postID)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	post.CommentsEnabled = commentsEnabled

	return post, nil
}

func (r *Storage) GetAllPosts() ([]*model.Post, error) {
	op := "internal.storage.db.GetAllPosts()"

	queryGetAllPosts := `SELECT post_id, author_id, title, text, comments_enabled, create_date
						FROM Posts
						ORDER BY create_date`

	var posts []*model.Post
	err := r.db.Select(&posts, queryGetAllPosts)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrPostsNotExist
		}
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return posts, nil
}

func (r *Storage) GetPost(postID int64) (*model.Post, error) {
	op := "internal.storage.db.GetPost()"

	queryGetPost := `SELECT post_id, author_id, title, text, comments_enabled, create_date
						FROM Posts
						WHERE post_id = $1`

	var post = new(model.Post)
	err := r.db.Get(post, queryGetPost, postID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrPostNotExist
		}
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return post, nil
}
