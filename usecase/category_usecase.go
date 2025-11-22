package usecase

import (
	"errors"

	"github.com/mhmmmdrivaldhi/go-book-api/model"
	"github.com/mhmmmdrivaldhi/go-book-api/model/dto"
	"github.com/mhmmmdrivaldhi/go-book-api/repository"
)

type CategoryUsecase interface {
	Create(req dto.CreateCategoryRequest) (*model.Category, error)
	GetAll() ([]dto.CategoryResponse, error)
	GetById(id int) (*dto.CategoryResponse, error)
	Update(id int, updateCategory *dto.UpdateCategoryRequest) (*model.Category, error)
	Delete(id int) error 
}

type categoryUsecase struct {
	categoryRepo repository.CategoryRepository
}

func (cu *categoryUsecase) Create(req dto.CreateCategoryRequest) (*model.Category, error) {
	category := &model.Category{
		Name: req.Name,
	}

	book, err := cu.categoryRepo.CreateCategory(category)
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (cu *categoryUsecase) GetAll() ([]dto.CategoryResponse, error) {
	var response [] dto.CategoryResponse

	categories, err := cu.categoryRepo.FindAll()
	if err != nil {
		return nil, err
	}

	for _, category := range categories {
		categoryResponse := dto.CategoryResponse{
			ID: category.ID,
			Name: category.Name,
		}
		response = append(response, categoryResponse)
	}

	return response, nil
}

func (cu *categoryUsecase) GetById(id int) (*dto.CategoryResponse, error) {
	category, err := cu.categoryRepo.FindById(id)
	if err != nil {
		return nil, err
	}

	response := &dto.CategoryResponse{
		ID: category.ID,
		Name: category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}

	return response, nil
}

func (cu *categoryUsecase) Update(id int, updateCategory *dto.UpdateCategoryRequest) (*model.Category, error) {
	category, err := cu.categoryRepo.FindById(id)
	if err != nil {
		return nil, err
	}

	if updateCategory.Name != nil {
		category.Name = *updateCategory.Name
	}

	update, err := cu.categoryRepo.UpdateCategory(id, category)
	if err != nil {
		return nil, err
	
	}

	return update, nil
}

func (cu *categoryUsecase) Delete(id int) error {
	_, err := cu.categoryRepo.FindById(id)
	if err != nil {
		return errors.New("category not found")
	}

	err = cu.categoryRepo.DeleteCategory(id)
	if err != nil {
		return errors.New("failed to delete category")
	} 

	return err
}

func NewCategoryUsecase(categoryRepo repository.CategoryRepository) *categoryUsecase {
	return &categoryUsecase{categoryRepo: categoryRepo}
}