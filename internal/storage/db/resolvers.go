package db

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"sync"
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

	log.Debug().Msgf("%s start", op)

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

	log.Debug().Msgf("%s end", op)
	return post, nil
}

func (r *Storage) GetAllPosts() ([]*model.Post, error) {
	op := "internal.storage.db.GetAllPosts()"

	log.Debug().Msgf("%s start", op)

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

	log.Debug().Msgf("%s end", op)

	return posts, nil
}

func (r *Storage) GetPost(postID int64) (*model.Post, error) {
	op := "internal.storage.db.GetPost()"

	log.Debug().Msgf("%s start", op)

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

	log.Debug().Msgf("%s end", op)

	return post, nil
}

func (r *Storage) GetCommentsBranch(postID int64, path string) ([]*model.Comment, error) {
	op := "internal.storage.db.GetCommentsBranch()"

	log.Debug().Msgf("%s start", op)

	allComments, err := r.GetCommentsToPostFromCashe(postID, path)
	if err != nil {
		if err == storage.ErrPathNotExist {
			return nil, err
		}
		log.Warn().Err(err).Msg("Cache returned error")
	}
	// TODO: check if this ip is not necessary
	if len(allComments) != 0 {
		return allComments, nil
	}

	allComments = make([]*model.Comment, 0)

	queryGetCommentsToPost := `SELECT comment_id, author_id, post_id, parent_id, path, text, create_date
								FROM Comments
								WHERE post_id = $1
								ORDER BY path,
								create_date DESC`

	err = r.db.Select(&allComments, queryGetCommentsToPost, postID)

	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	if len(allComments) == 0 {
		return nil, storage.ErrCommentsNotExist
	}

	commentsMap, rootComments := CreateCommentMap(allComments)

	err = r.SetCommentToPostInCache(commentsMap, rootComments, postID)
	if err != nil {
		log.Warn().Err(err).Msg("Cache returned error")
	}

	if path == "" {
		log.Debug().Msgf("%s end", op)
		return rootComments, nil
	} else {
		if v, ok := commentsMap[path]; ok != true {
			return nil, storage.ErrPathNotExist
		} else {
			log.Debug().Msgf("%s end", op)
			return v, nil
		}
	}

}

func CreateCommentMap(allComments []*model.Comment) (map[string][]*model.Comment, []*model.Comment) {

	var rootComments []*model.Comment
	var commentsMap = make(map[string][]*model.Comment)
	var pathMap = make(map[int64]string, len(allComments))
	for _, v := range allComments {
		pathMap[v.ID] = v.Path
		if v.ParentID == nil {
			rootComments = append(rootComments, v)
		} else {

			commentsMap[pathMap[*v.ParentID]] = append(commentsMap[pathMap[*v.ParentID]], v) //we know that parent comment exists because of the query is ordered by path
		}
	}

	log.Debug().Msg("Comment map created successfully")
	return commentsMap, rootComments
}

func (r *Storage) GetCommentsToPostFromCashe(postID int64, path string) ([]*model.Comment, error) {
	op := "internal.storage.db.GetCommentsToPostFromCashe()"
	log.Debug().Msgf("%s start", op)

	var rootComments []*model.Comment
	err := r.cache.Get(context.Background(), "post:"+strconv.FormatInt(postID, 10), &rootComments)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(rootComments) == 0 {
		return nil, storage.ErrPostNotCached
	}

	if path == "" {

		log.Debug().Msgf("%s end", op)
		return rootComments, nil
	}

	var commentsBranch []*model.Comment

	err = r.cache.Get(context.Background(), "comments:"+path, &commentsBranch)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(commentsBranch) == 0 {
		return nil, storage.ErrPathNotExist
	}

	log.Debug().Msgf("%s end", op)
	return commentsBranch, nil

}

func (r *Storage) SetCommentToPostInCache(commentsMap map[string][]*model.Comment, rootComments []*model.Comment, postID int64) error {
	op := "internal.storage.db.SetCommentToPostInCache()"
	log.Debug().Msgf("%s start", op)

	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error

	for path, comments := range commentsMap {
		wg.Add(1)
		go func(path string, comments []*model.Comment) {
			defer wg.Done()

			err := r.cache.Set(&cache.Item{
				Key:   "comments:" + path,
				Value: comments,
				TTL:   time.Hour,
			})

			mu.Lock()
			defer mu.Unlock()
			if err != nil && firstErr == nil {
				firstErr = err
				return
			}
		}(path, comments)
	}
	wg.Wait()

	if firstErr != nil {
		return fmt.Errorf("%s: %w", op, firstErr)
	}

	err := r.cache.Set(&cache.Item{
		Key:   "post:" + strconv.FormatInt(postID, 10),
		Value: rootComments,
		TTL:   time.Hour,
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug().Msgf("%s end", op)
	return nil
}

func (r *Storage) GetCommentPath(commentID int64) (string, error) {
	op := "internal.storage.db.GetCommentPath()"

	log.Debug().Msgf("%s start", op)

	var path string

	queryGetCommentsToPost := `SELECT path
								FROM Comments
								WHERE  comment_id = $1`

	err := r.db.Get(&path, queryGetCommentsToPost, commentID)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", storage.ErrCommentsNotExist
		}
		return "", fmt.Errorf("%s:%w", op, err)
	}

	log.Debug().Msgf("%s start", op)
	return path, nil
}
