package usecase

import (
	"errors"

	"github.com/mhmmmdrivaldhi/go-book-api/model"
	"github.com/mhmmmdrivaldhi/go-book-api/model/dto"
	"github.com/mhmmmdrivaldhi/go-book-api/repository"
)

type BookUsecase interface {
	Create(req dto.CreateBookRequest) (*model.Book, error)
	GetAll() ([]dto.BookResponse, error)
	GetById(id int) (*dto.BookResponse, error)
	Update(id int, reqUpdate dto.UpdateBookRequest) (*model.Book, error)
	Delete(id int,) error
}

type bookUsecase struct {
	bookRepo repository.BookRepository
}

func (bu *bookUsecase) Create(req dto.CreateBookRequest,) (*model.Book, error) {
	books := &model.Book{
		Title:       req.Title,
		Description: req.Description,
		Author:      req.Author,
		Price:       req.Price,
		Rating:      req.Rating,
		CategoryID:  req.CategoryID,
	}

	create, err := bu.bookRepo.CreateBook(books)
	if err != nil {
		return nil, err
	}

	return create, nil
}

func (bu *bookUsecase) GetAll() ([]dto.BookResponse, error) {
	var response []dto.BookResponse

	books, err := bu.bookRepo.FindAll()
	if err != nil {
		return nil, err
	}

	for _, book := range books {
		var categoryResponse *dto.CategoryResponse
        if book.Category.ID != 0 { 
            categoryResponse = &dto.CategoryResponse{
                ID:   book.Category.ID,
                Name: book.Category.Name,
            }
        }

		bookResponse := dto.BookResponse{
			Category: categoryResponse,
			Id: book.Id,
			Title: book.Title,
			Description: book.Description,
			Author: book.Author,
			Price: book.Price,
			Rating: book.Rating,
			CreatedAt: book.CreatedAt,
			UpdatedAt: book.UpdatedAt,
		}
		response = append(response, bookResponse)
	}
	return response, nil
}

func (bu *bookUsecase) GetById(id int) (*dto.BookResponse, error) {
	book, err := bu.bookRepo.FindById(id)
	if err != nil {
		return nil, err
	}

	var categoryResponse *dto.CategoryResponse
	if book.Category.ID != 0 {
		categoryResponse = &dto.CategoryResponse{
			ID: book.Category.ID,
			Name: book.Category.Name,
		}
	}

	bookResponse := &dto.BookResponse{
		Category: categoryResponse,
		Id: book.Id,
		Title: book.Title,
		Description: book.Description,
		Author: book.Author,
		Price: book.Price,
		Rating: book.Rating,
		CreatedAt: book.CreatedAt,
		UpdatedAt: book.UpdatedAt,
	}
	return bookResponse, nil
}

func (bu *bookUsecase) Update(id int, reqUpdate dto.UpdateBookRequest) (*model.Book, error) {
	exists, err := bu.bookRepo.FindById(id)
	if err != nil {
		return nil, err
	}
	
	if reqUpdate.Title != nil {
		exists.Title = *reqUpdate.Title
	}

	if reqUpdate.Description != nil {
		exists.Description = *reqUpdate.Description
	}

	if reqUpdate.Author != nil {
		exists.Author = *reqUpdate.Author
	
	}

	if reqUpdate.Price != nil {
		exists.Price = *reqUpdate.Price
	}

	if reqUpdate.Rating != nil {
		exists.Rating = *reqUpdate.Rating
	}

	updateBook, err := bu.bookRepo.UpdateBook(id, exists)
	if err != nil {
		return nil, err
	}

	return updateBook, nil
}

func (bu *bookUsecase) Delete(id int,) error {
	_, err := bu.bookRepo.FindById(id)
	if err != nil {
		return errors.New("book not found")
	}

	err = bu.bookRepo.DeleteBook(id)
	if err != nil {
		return errors.New("failed to delete book")
	}

	return  nil
}

func NewBookUsecase(bookRepo repository.BookRepository) *bookUsecase {
	return &bookUsecase{bookRepo: bookRepo}
}
