package categoryservice

import (
	"errors"
	"fmt"
	"gamegolang/entity"
)

// CategoryRepository تعریف اینترفیس ریپازیتوری
type CategoryRepository interface {
	Create(req CreateRequest) (*entity.Category, error)
}

// Service ساختار سرویس دسته‌بندی
type Service struct {
	Repo CategoryRepository
}

// CreateRequest ساختار درخواست ایجاد دسته‌بندی
type CreateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Create ایجاد یک دسته‌بندی جدید
func (s Service) Create(req CreateRequest) (*entity.Category, error) {
	// اعتبارسنجی ورودی‌ها
	if req.Title == "" {
		return nil, errors.New("title is required")
	}

	// ایجاد دسته‌بندی در ریپازیتوری
	category, err := s.Repo.Create(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}
