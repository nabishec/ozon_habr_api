package inmemory

import (
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/nabishec/ozon_habr_api/internal/model"
	"github.com/nabishec/ozon_habr_api/internal/pkg/errs"
	"github.com/rs/zerolog/log"
)

type Storage struct {
	postsLastIndex   int64
	commentLastIndex int64
	comments         map[int64][]*model.Comment
	comment          map[int64]*model.Comment
	repliesByPath    map[string][]*model.Comment
	post             map[int64]*model.Post
	posts            []*model.Post
}

func NewStorage() *Storage {
	return &Storage{
		postsLastIndex:   0,
		commentLastIndex: 0,
		comments:         make(map[int64][]*model.Comment), // root comments
		comment:          make(map[int64]*model.Comment),
		repliesByPath:    make(map[string][]*model.Comment),
		post:             make(map[int64]*model.Post),
		posts:            make([]*model.Post, 0),
	}
}

func (r *Storage) AddPost(ctx context.Context, newPost *model.NewPost) (*model.Post, error) {
	op := "internal.storage.inmemory.AddPost()"

	log.Debug().Msgf("%s start", op)
	post := &model.Post{
		AuthorID:        newPost.AuthorID,
		Title:           newPost.Title,
		Text:            newPost.Text,
		CommentsEnabled: newPost.CommentsEnabled,
		CreateDate:      time.Now(),
	}

	select {
	case <-ctx.Done():
		log.Warn().Msgf("%s canceled", op)
		return nil, ctx.Err()
	default:
	}

	postID := r.postsLastIndex + 1
	r.postsLastIndex += 1
	post.ID = postID
	r.post[postID] = post
	r.posts = append(r.posts, post)

	log.Debug().Msgf("%s end", op)
	return post, nil
}

func (r *Storage) AddComment(ctx context.Context, postID int64, newComment *model.NewComment) (*model.Comment, error) {
	op := "internal.storage.inmemory.AddComment()"

	log.Debug().Msgf("%s start", op)

	comment := &model.Comment{
		AuthorID:   newComment.AuthorID,
		PostID:     postID,
		ParentID:   newComment.ParentID,
		Text:       newComment.Text,
		CreateDate: time.Now(),
	}

	select {
	case <-ctx.Done():
		log.Warn().Msgf("%s canceled", op)
		return nil, ctx.Err()
	default:
	}

	post, ok := r.post[postID]
	if !ok {
		return nil, errs.ErrPostNotExist
	}

	if !post.CommentsEnabled {
		return nil, errs.ErrCommentsNotEnabled
	}

	var parentPath string
	if comment.ParentID != nil {
		if parentComment, ok := r.comment[*comment.ParentID]; ok {
			parentPath = parentComment.Path
			if parentComment.PostID != postID {
				return nil, errs.ErrParentCommentNotExist
			}
			log.Info().Msgf("parentPath: %s", parentPath)
		} else {
			return nil, errs.ErrParentCommentNotExist
		}
		r.repliesByPath[parentPath] = append(r.repliesByPath[parentPath], comment)
		parentPath += "."

	} else {
		r.comments[postID] = append(r.comments[postID], comment)
	}

	commentID := r.commentLastIndex + 1
	r.commentLastIndex += 1
	comment.ID = commentID

	path := parentPath + strconv.FormatInt(comment.ID, 10)
	comment.Path = path

	r.comment[commentID] = comment

	log.Debug().Msgf("%s end", op)
	return comment, nil
}

func (r *Storage) UpdateEnableCommentToPost(ctx context.Context, postID int64, authorID uuid.UUID, commentsEnabled bool) (*model.Post, error) {
	op := "internal.storage.inmemory.UpdateEnableCommentToPost()"

	log.Debug().Msgf("%s start", op)

	select {
	case <-ctx.Done():
		log.Warn().Msgf("%s canceled", op)
		return nil, ctx.Err()
	default:
	}

	post, ok := r.post[postID]
	if !ok {
		return nil, errs.ErrPostNotExist
	}

	if post.AuthorID != authorID {
		return nil, errs.ErrUnauthorizedAccess
	}

	post.CommentsEnabled = commentsEnabled

	log.Debug().Msgf("%s end", op)
	return post, nil
}

func (r *Storage) GetAllPosts(ctx context.Context) ([]*model.Post, error) {
	op := "internal.storage.inmemory.GetAllPosts()"

	log.Debug().Msgf("%s start", op)

	select {
	case <-ctx.Done():
		log.Warn().Msgf("%s canceled", op)
		return nil, ctx.Err()
	default:
	}

	if len(r.posts) == 0 {
		return nil, errs.ErrPostsNotExist
	}

	log.Debug().Msgf("%s end", op)

	return r.posts, nil
}

func (r *Storage) GetPost(ctx context.Context, postID int64) (*model.Post, error) {
	op := "internal.storage.inmemory.GetPost()"

	log.Debug().Msgf("%s start", op)

	select {
	case <-ctx.Done():
		log.Warn().Msgf("%s canceled", op)
		return nil, ctx.Err()
	default:
	}

	post, ok := r.post[postID]

	if !ok {
		return nil, errs.ErrPostNotExist
	}

	log.Debug().Msgf("%s end", op)

	return post, nil
}

func (r *Storage) GetCommentsBranch(ctx context.Context, postID int64, path string) ([]*model.Comment, error) {
	op := "internal.storage.inmemory.GetCommentsBranch()"

	log.Debug().Msgf("%s start", op)

	select {
	case <-ctx.Done():
		log.Warn().Msgf("%s canceled", op)
		return nil, ctx.Err()
	default:
	}

	if path == "" {
		if v, ok := r.comments[postID]; ok {
			return v, nil
		}
		return nil, errs.ErrCommentsNotExist
	}

	comments, ok := r.repliesByPath[path]
	if !ok {
		return nil, errs.ErrPathNotExist
	}

	if len(comments) == 0 {
		return nil, errs.ErrCommentsNotExist
	}

	return comments, nil

}

func (r *Storage) GetCommentPath(ctx context.Context, commentID int64) (string, error) {
	op := "internal.storage.inmemory.GetCommentPath()"

	log.Debug().Msgf("%s start", op)

	select {
	case <-ctx.Done():
		log.Warn().Msgf("%s canceled", op)
		return "", ctx.Err()
	default:
	}

	comment, ok := r.comment[commentID]
	if !ok {
		return "", errs.ErrCommentsNotExist
	}

	path := comment.Path

	log.Debug().Msgf("%s end", op)
	return path, nil
}
