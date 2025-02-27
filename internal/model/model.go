package model

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID         int64              `json:"id" db:"comment_id"`
	AuthorID   uuid.UUID          `json:"authorID" db:"author_id"`
	PostID     int64              `json:"postID" db:"post_id"`
	ParentID   *int64             `json:"parentID,omitempty" db:"parent_id"`
	Text       string             `json:"text" db:"text"`
	CreateDate time.Time          `json:"createDate" db:"create_date"`
	Replies    *CommentConnection `json:"replies,omitempty"`
}

type CommentConnection struct {
	Edges    []*CommentEdge `json:"edges"`
	PageInfo *PageInfo      `json:"pageInfo"`
}

type CommentEdge struct {
	Node   *Comment `json:"node"`
	Cursor string   `json:"cursor"`
}

type Mutation struct {
}

type NewComment struct {
	AuthorID uuid.UUID `json:"authorID"`
	PostID   int64     `json:"postID"`
	ParentID *int64    `json:"parentID,omitempty"`
	Text     string    `json:"text"`
}

type NewPost struct {
	AuthorID        uuid.UUID `json:"authorID"`
	Title           string    `json:"title"`
	Text            string    `json:"text"`
	CommentsEnabled bool      `json:"commentsEnabled"`
}

type PageInfo struct {
	EndCursor   *string `json:"endCursor,omitempty"`
	HasNextPage bool    `json:"hasNextPage"`
}

type Post struct {
	ID              int64              `json:"id" db:"post_id"`
	AuthorID        uuid.UUID          `json:"authorID" db:"author_id"`
	Title           string             `json:"title" db:"title"`
	Text            string             `json:"text" db:"text"`
	CommentsEnabled bool               `json:"commentsEnabled" db:"comments_enabled"`
	Comments        *CommentConnection `json:"comments,omitempty"`
	CreateDate      time.Time          `json:"createDate" db:"create_date"`
}
