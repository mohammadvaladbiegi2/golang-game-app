package category_service

import (
	"errors"
	"fmt"
	"gamegolang/entity"
)

type CategoryRepository interface {
	CreateCategory(req CreateCategoryRequestStruct) (entity.Category, error)
}

type CategoryService struct {
	Repo CategoryRepository
}

type CreateCategoryRequestStruct struct {
	Title       string
	Description string
}

type CreateCategoryResponseStruct struct {
	Category entity.Category
}

func (receiver CategoryService) CreateCategory(req CreateCategoryRequestStruct) (CreateCategoryResponseStruct, error) {

	if req.Title == "" {
		return CreateCategoryResponseStruct{}, errors.New("title is required")
	}

	newCategory, Cerror := receiver.Repo.CreateCategory(req)
	if Cerror != nil {
		return CreateCategoryResponseStruct{}, fmt.Errorf("create category error: %w", Cerror)
	}

	return CreateCategoryResponseStruct{
		Category: newCategory,
	}, nil

}
