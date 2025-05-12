package db

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/nabishec/ozon_habr_api/internal/model"
	"github.com/nabishec/ozon_habr_api/internal/pkg/errs"
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

func (r *Storage) AddPost(ctx context.Context, newPost *model.NewPost) (*model.Post, error) {
	op := "internal.storage.db.AddPost()"

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

	err := r.db.QueryRowContext(ctx, queryNewPost, post.AuthorID, post.Title, post.Text, post.CommentsEnabled, post.CreateDate).Scan(&post.ID)

	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	log.Debug().Msgf("%s end", op)
	return post, err
}

func (r *Storage) AddComment(ctx context.Context, postID int64, newComment *model.NewComment) (*model.Comment, error) {
	op := "internal.storage.db.AddComment()"

	log.Debug().Msgf("%s start", op)
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	defer func() {
		if err != nil && err != errs.ErrCommentsNotEnabled && err != errs.ErrPostNotExist && err != errs.ErrParentCommentNotExist {

			errRB := tx.Rollback()
			if errRB != nil {
				log.Error().Err(errRB).Msg(" roll back transaction failed")
			}

		}
	}()

	comment := &model.Comment{
		AuthorID:   newComment.AuthorID,
		PostID:     postID,
		ParentID:   newComment.ParentID,
		Text:       newComment.Text,
		CreateDate: time.Now(),
	}

	queryGetCommentEnabledForPost := `SELECT comments_enabled 
							FROM Posts 
							WHERE post_id = $1`

	queryGetParentPath := `SELECT path
							FROM Comments
							WHERE comment_id = $1 AND post_id = $2`

	queryNewComment := `INSERT INTO Comments (author_id, post_id, parent_id, path, replies_level, text, create_date)
						VALUES ($1, $2, $3, $4, $5, $6, $7)
						RETURNING comment_id `

	queryUpdateCommentPath := `UPDATE Comments
							SET path = $1, replies_level = $2
							WHERE comment_id = $3`

	var commentEnabled bool
	err = tx.GetContext(ctx, &commentEnabled, queryGetCommentEnabledForPost, postID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrPostNotExist
		}
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	if commentEnabled == false {
		return nil, errs.ErrCommentsNotEnabled
	}

	var parentPath string
	if comment.ParentID != nil {
		err = tx.GetContext(ctx, &parentPath, queryGetParentPath, comment.ParentID, comment.PostID)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, errs.ErrParentCommentNotExist
			}
			return nil, fmt.Errorf("%s:%w", op, err)
		}
		parentPath += "."
	}
	// we specify 0 as the path and 1 as rep. level, because we know that the comment ID cannot be zero, and the query did not return errors due to not null
	err = tx.QueryRowContext(ctx, queryNewComment, comment.AuthorID, comment.PostID, comment.ParentID, "0", 1, comment.Text, comment.CreateDate).Scan(&comment.ID)

	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	path := parentPath + strconv.FormatInt(comment.ID, 10)
	repliceLevel := strings.Count(path, ".") + 1

	_, err = tx.ExecContext(ctx, queryUpdateCommentPath, path, repliceLevel, comment.ID)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	comment.Path = path

	commentsBranch, err := r.getCommentsToPostFromCashe(ctx, comment.PostID, parentPath)
	if err == nil {
		commentsBranch = append(commentsBranch, comment)
		err = r.setCommentsBranchToPostInCache(ctx, commentsBranch, comment.PostID, parentPath[:len(parentPath)-1])
		if err != nil {
			log.Warn().Err(err).Msg("cache returned error")
			err = nil
		}
	} else {
		if err != cache.ErrCacheMiss {
			log.Warn().Err(err).Msg("cache returned error") // logging cache error and return nil error because it is not critical
		}
		err = nil
	}
	log.Debug().Msgf("%s end", op)
	return comment, nil
}

func (r *Storage) setCommentsBranchToPostInCache(ctx context.Context, commentsBranch []*model.Comment, postID int64, path string) error {
	op := "internal.storage.db.SetCommentsBranchToPostInCache()"
	log.Debug().Msgf("%s start", op)

	var err error
	if path == "" {
		err = r.cache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   "post:" + strconv.FormatInt(postID, 10),
			Value: commentsBranch,
			TTL:   30 * time.Minute,
		})
	} else {
		err = r.cache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   "comments:" + path,
			Value: commentsBranch,
			TTL:   30 * time.Minute,
		})
	}

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug().Msgf("%s end", op)
	return nil
}
func (r *Storage) UpdateEnableCommentToPost(ctx context.Context, postID int64, authorID uuid.UUID, commentsEnabled bool) (*model.Post, error) {
	op := "internal.storage.db.UpdateEnableCommentToPost()"

	log.Debug().Msgf("%s start", op)

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	defer func() {
		if err != nil && err != errs.ErrPostNotExist && err != errs.ErrUnauthorizedAccess {
			errRB := tx.Rollback()
			if errRB != nil {
				log.Error().Err(errRB).Msg(" roll back transaction failed")
			}
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

	err = tx.GetContext(ctx, post, queryGetPost, postID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrPostNotExist
		}
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	if post.AuthorID != authorID {
		return nil, errs.ErrUnauthorizedAccess
	}

	_, err = tx.ExecContext(ctx, queryUpdatePost, commentsEnabled, postID)
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

func (r *Storage) GetAllPosts(ctx context.Context) ([]*model.Post, error) {
	op := "internal.storage.db.GetAllPosts()"

	log.Debug().Msgf("%s start", op)

	queryGetAllPosts := `SELECT post_id, author_id, title, text, comments_enabled, create_date
						FROM Posts
						ORDER BY create_date`

	var posts []*model.Post
	err := r.db.SelectContext(ctx, &posts, queryGetAllPosts)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrPostsNotExist
		}
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	log.Debug().Msgf("%s end", op)

	return posts, nil
}

func (r *Storage) GetPost(ctx context.Context, postID int64) (*model.Post, error) {
	op := "internal.storage.db.GetPost()"

	log.Debug().Msgf("%s start", op)

	queryGetPost := `SELECT post_id, author_id, title, text, comments_enabled, create_date
						FROM Posts
						WHERE post_id = $1`

	var post = new(model.Post)
	err := r.db.GetContext(ctx, post, queryGetPost, postID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrPostNotExist
		}
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	log.Debug().Msgf("%s end", op)

	return post, nil
}

func (r *Storage) GetCommentsBranch(ctx context.Context, postID int64, path string) ([]*model.Comment, error) {
	op := "internal.storage.db.GetCommentsBranch()"

	log.Debug().Msgf("%s start", op)

	allComments, err := r.getCommentsToPostFromCashe(ctx, postID, path)
	if err != nil {
		if err == errs.ErrPathNotExist {
			return nil, err
		}
		log.Warn().Err(err).Msg("Cache returned error")
	} else {
		return allComments, nil
	}

	allComments = make([]*model.Comment, 0)

	queryGetCommentsToPost := `SELECT comment_id, author_id, post_id, parent_id, path, text, create_date
								FROM Comments
								WHERE post_id = $1
								ORDER BY string_to_array(path::text, '.')::int[],
								create_date DESC`

	err = r.db.SelectContext(ctx, &allComments, queryGetCommentsToPost, postID)

	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	if len(allComments) == 0 {
		return nil, errs.ErrCommentsNotExist
	}

	commentsMap, rootComments := createCommentMap(allComments)

	err = r.setCommentsToPostInCache(ctx, commentsMap, rootComments, postID)
	if err != nil {
		log.Warn().Err(err).Msg("Cache returned error")
	}

	if path == "" {
		log.Debug().Msgf("%s end", op)
		return rootComments, nil
	} else {
		if v, ok := commentsMap[path]; ok != true {
			return nil, errs.ErrPathNotExist
		} else {
			log.Debug().Msgf("%s end", op)
			return v, nil
		}
	}

}

func createCommentMap(allComments []*model.Comment) (map[string][]*model.Comment, []*model.Comment) {

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

func (r *Storage) getCommentsToPostFromCashe(ctx context.Context, postID int64, path string) ([]*model.Comment, error) {
	op := "internal.storage.db.GetCommentsToPostFromCashe()"
	log.Debug().Msgf("%s start", op)

	var rootComments []*model.Comment
	err := r.cache.Get(ctx, "post:"+strconv.FormatInt(postID, 10), &rootComments)
	if err != nil {
		if err == cache.ErrCacheMiss {
			return nil, err
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(rootComments) == 0 {
		return nil, errs.ErrPostNotCached
	}

	if path == "" {

		log.Debug().Msgf("%s end", op)
		return rootComments, nil
	}

	var commentsBranch []*model.Comment

	err = r.cache.Get(ctx, "comments:"+path, &commentsBranch)
	if err != nil {
		if err == cache.ErrCacheMiss && len(rootComments) > 0 {
			return nil, nil // it work when try get replies of comment without replies
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(commentsBranch) == 0 {
		return nil, errs.ErrPathNotExist
	}

	log.Debug().Msgf("%s end", op)
	return commentsBranch, nil

}

func (r *Storage) setCommentsToPostInCache(ctx context.Context, commentsMap map[string][]*model.Comment, rootComments []*model.Comment, postID int64) error {
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
				Ctx:   ctx,
				Key:   "comments:" + path,
				Value: comments,
				TTL:   30 * time.Minute,
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
		Ctx:   ctx,
		Key:   "post:" + strconv.FormatInt(postID, 10),
		Value: rootComments,
		TTL:   30 * time.Minute,
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug().Msgf("%s end", op)
	return nil
}

func (r *Storage) GetCommentPath(ctx context.Context, commentID int64) (string, error) {
	op := "internal.storage.db.GetCommentPath()"

	log.Debug().Msgf("%s start", op)

	var path string

	queryGetCommentsToPost := `SELECT path
								FROM Comments
								WHERE  comment_id = $1`

	err := r.db.GetContext(ctx, &path, queryGetCommentsToPost, commentID)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", errs.ErrCommentsNotExist
		}
		return "", fmt.Errorf("%s:%w", op, err)
	}

	log.Debug().Msgf("%s end", op)
	return path, nil
}
