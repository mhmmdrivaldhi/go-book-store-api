package repository

import (
	"errors"

	"github.com/mhmmmdrivaldhi/go-book-api/model"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	CreateCategory(category *model.Category) (*model.Category, error)
	FindAll() ([]model.Category, error)
	FindById(id int) (*model.Category, error)
	UpdateCategory(id int, updateCategory *model.Category) (*model.Category, error)
	DeleteCategory(id int) error
}

type categoryRepository struct {
	db *gorm.DB
}

func (cr *categoryRepository) CreateCategory(category *model.Category) (*model.Category, error) {
	err := cr.db.Create(category).Error
	if err != nil {
		return nil, errors.New("failed to create category")
	}

	return category, nil
}

func (cr *categoryRepository) FindAll() ([]model.Category, error) {
	var categories []model.Category

	err := cr.db.Find(&categories).Error
	if err != nil {
		return nil, errors.New("failed to find all categories")
	}

	return categories, nil
}

func (cr *categoryRepository) FindById(id int) (*model.Category, error) {
	var category model.Category

	err := cr.db.First(&category, id).Error
	if err != nil {
		return nil, errors.New("failed to find category by id")
	}

	return &category, nil
}

func (cr *categoryRepository) UpdateCategory(id int, updateCategory *model.Category) (*model.Category, error) {
	var category model.Category

	err := cr.db.First(&category, id).Error
	if err != nil {
		return nil, errors.New("failed to find category by id")
	}

	err = cr.db.Model(&category).Updates(updateCategory).Error
	if err != nil {
		return nil, errors.New("failed to update category")
	}

	return &category, nil 
}

func (cr *categoryRepository) DeleteCategory(id int) error {
	err := cr.db.Delete(&model.Category{}, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("category not found")
	} else if err != nil {
		return errors.New("failed to delete category")
	}

	return nil
}

func NewCategoryRepository(db *gorm.DB) *categoryRepository {
	return &categoryRepository{db: db}
} 