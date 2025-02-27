package storage

import "github.com/nabishec/ozon_habr_api/internal/model"

type StorageImp interface {
	AddPost(newPost *model.NewPost) (*model.Post, error)
}
