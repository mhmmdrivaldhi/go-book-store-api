package usecase

import (
	"errors"

	"github.com/mhmmmdrivaldhi/go-book-api/model/dto"
	"github.com/mhmmmdrivaldhi/go-book-api/repository"
	"github.com/mhmmmdrivaldhi/go-book-api/service"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
}

type authUsecase struct {
	userRepo repository.UserRepository
	jwtService service.JwtService
}

func (au *authUsecase) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := au.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	 err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}
	
	token, err := au.jwtService.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID: user.ID,
			Name: user.Name,
			Email: user.Email,
			Role: user.Role,
		},
	}, nil
}

func NewAuthUsecase(userRepo repository.UserRepository, jwtService service.JwtService) *authUsecase {
	return &authUsecase{
		userRepo: userRepo,
		jwtService: jwtService,
	}
}