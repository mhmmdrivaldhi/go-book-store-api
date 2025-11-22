package repository

import (
	"errors"

	"github.com/mhmmmdrivaldhi/go-book-api/model"
	"gorm.io/gorm"
)

type BookRepository interface {
	CreateBook(book *model.Book) (*model.Book, error)
	FindAll() ([]model.Book, error)
	FindById(id int) (*model.Book, error)
	UpdateBook(id int, updateBook *model.Book) (*model.Book, error)
	DeleteBook(id int) error
}

type bookRepository struct {
	db *gorm.DB
}

func (bookRepo *bookRepository) CreateBook(book *model.Book) (*model.Book, error) {
	err := bookRepo.db.Create(&book).Error
	if err != nil {
		return nil, err
	}

	var res model.Book
	err = bookRepo.db.Preload("Category").First(&res, book.Id).Error
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (bookRepo *bookRepository) FindAll() ([]model.Book, error) {
	var books []model.Book

	err := bookRepo.db.Preload("Category").Find(&books).Error
	if err != nil {
		return nil, err
	}

	if len(books) == 0 {
		return nil, errors.New("no books found")
	}

	return books, nil
}

func (bookRepo *bookRepository) FindById(id int) (*model.Book, error) {
	var book model.Book

	err := bookRepo.db.Preload("Category").First(&book, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("book not found")
	} else if err != nil {
		return nil, err
	}

	return &book, nil
}

func (bookRepo *bookRepository) UpdateBook(id int, updateBook *model.Book) (*model.Book, error) {
	var book model.Book

	err := bookRepo.db.First(&book, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("book not found")
	} else if err != nil {
		return nil, err
	}

	err = bookRepo.db.Model(&book).Updates(updateBook).Error
	if err != nil {
		return nil, err
	}

	err = bookRepo.db.Preload("Category").First(&book, id).Error
	if err != nil {
		return nil, err
	}


	return &book, nil
}

func (bookRepo *bookRepository) DeleteBook(id int) error {
	err := bookRepo.db.Delete(&model.Book{}, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("book not found")
	} else if err != nil {
		return err
	}

	return nil
}

func NewBookRepository(db *gorm.DB) *bookRepository {
	return &bookRepository{db: db}
}