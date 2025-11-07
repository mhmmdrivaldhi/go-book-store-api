package repository

import (
	"errors"

	"github.com/mhmmmdrivaldhi/go-book-api/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) (*model.User, error)
	FindAll() ([]model.User, error)
	FindById(id int) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	Update(id int, updateUser *model.User) (*model.User, error)
	Delete(id int) error
}

type userRepository struct {
	db *gorm.DB
}

func (userRepo *userRepository) Create(user *model.User) (*model.User, error) {
	err := userRepo.db.Create(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (userRepo *userRepository) FindAll() ([]model.User, error) {
	var users []model.User

	err := userRepo.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	
	if len(users) == 0 {
		return nil, errors.New("no users found")
	}

	return users, nil
}

func (userRepo *userRepository) FindById(id int) (*model.User, error) {
	var user model.User

	err := userRepo.db.First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (userRepo *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User

	err := userRepo.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("email address not found")
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (userRepo *userRepository) Update(id int, updateUser *model.User) (*model.User, error) {
	var user model.User

	err := userRepo.db.First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}

	err = userRepo.db.Model(&user).Updates(updateUser).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (userRepo *userRepository) Delete(id int) error {
	err := userRepo.db.Delete(model.User{}, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("user not found")
	} else if err != nil {
		return err
	}

	return nil
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}