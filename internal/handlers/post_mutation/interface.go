package postmutation

import "github.com/nabishec/ozon_habr_api/internal/model"

type PostImp interface {
	AddPost(newPost *model.NewPost) (*model.Post, error)
}
