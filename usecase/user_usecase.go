package usecase

import (
	"errors"

	"github.com/mhmmmdrivaldhi/go-book-api/model"
	"github.com/mhmmmdrivaldhi/go-book-api/model/dto"
	"github.com/mhmmmdrivaldhi/go-book-api/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	Create(req dto.CreateUserRequest) (*model.User, error)
	GetAllUsers() ([]dto.UserResponse, error)
	GetUserById(id int) (*dto.UserResponse, error)
	UpdateUser(id int, userUpdate dto.UpdateUserRequest) (*model.User, error)
	DeleteUser(id int) error
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func (uu *userUsecase) Create(req dto.CreateUserRequest) (*model.User, error) {
	_, err := uu.userRepo.FindByEmail(req.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	} 

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := model.User{
		Name: req.Name,
		Email: req.Email,
		Password: string(hashPassword),
		Role: req.Role,
	}

	create, err := uu.userRepo.Create(&user)
	if err != nil {
		return nil, err
	}

	return create, nil
}

func (uu *userUsecase) GetAllUsers() ([]dto.UserResponse, error) {
	var response []dto.UserResponse

	users, err := uu.userRepo.FindAll()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		userResponse := dto.UserResponse{
			ID: user.ID,
			Name: user.Name,
			Email: user.Email,
			Role: user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
		response = append(response, userResponse)
	}
	return response, nil
}

func (uu *userUsecase) GetUserById(id int) (*dto.UserResponse, error) {
	user, err := uu.userRepo.FindById(id)
	if err != nil {
		return nil, err
	}

	userResponse := &dto.UserResponse{
		ID: user.ID,
		Name: user.Name,
		Email: user.Email,
		Role: user.Role,
	}

	return userResponse, nil
}

func (uu *userUsecase) UpdateUser(id int, userUpdate dto.UpdateUserRequest) (*model.User, error) {
	user, err := uu.userRepo.FindById(id)
	if err != nil {
		return nil, err
	}

	_, err = uu.userRepo.FindByEmail(user.Email)
	if err != nil {
		return nil, err
	}

	if userUpdate.Name != nil {
		user.Name = *userUpdate.Name
	}

	if userUpdate.Email != nil {
		user.Email = *userUpdate.Email
	}

	if userUpdate.Password != nil {
		user.Password = *userUpdate.Password
	}

	if userUpdate.Role != nil {
		user.Role = *userUpdate.Role
	}
	
	updateUser, err := uu.userRepo.Update(id, user)
	if err != nil {
		return nil, err
	}

	return updateUser, nil
}

func (uu *userUsecase) DeleteUser(id int) error {
	_, err := uu.userRepo.FindById(id)
	if err != nil {
		return errors.New("user not found")
	}

	err = uu.userRepo.Delete(id)
	if err != nil {
		return errors.New("failed delete user")
	}

	return nil
}

func NewUserUsecase(userRepo repository.UserRepository) *userUsecase {
	return &userUsecase{userRepo: userRepo}
}
