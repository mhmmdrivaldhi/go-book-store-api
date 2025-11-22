package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mhmmmdrivaldhi/go-book-api/config"
	"github.com/mhmmmdrivaldhi/go-book-api/controller"
	"github.com/mhmmmdrivaldhi/go-book-api/middleware"
	"github.com/mhmmmdrivaldhi/go-book-api/model"
	"github.com/mhmmmdrivaldhi/go-book-api/repository"
	"github.com/mhmmmdrivaldhi/go-book-api/service"
	"github.com/mhmmmdrivaldhi/go-book-api/usecase"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	bookUsecase usecase.BookUsecase
	userUsecase usecase.UserUsecase
	categoryUsecase usecase.CategoryUsecase
	authUsecase usecase.AuthUsecase
	jwtService  service.JwtService
	engine *gin.Engine
	host string
}

func (s *Server) InitRoute() {
	auth := s.engine.Group("/api/auth")
	controller.NewAuthController(s.authUsecase, auth)

	v1 := s.engine.Group("/api/v1")

	authMiddleware := middleware.NewAuthMiddleware(s.jwtService)

	// public routes
	controller.NewUserController(s.userUsecase, v1).Route()

	// routes with authentication & authorization
	authGroup := v1.Group("")
	authGroup.Use(authMiddleware.RequireToken())

	controller.NewBookController(s.bookUsecase, authGroup)
	controller.NewCategoryController(s.categoryUsecase, authGroup)
}

func (s *Server) Run() {
	s.InitRoute()

	err :=  s.engine.Run(s.host)
	if err != nil {
		panic(fmt.Errorf("server not running on host %s, because error %v", s.host, err.Error()))
	}
}

func NewServer() *Server {
	cfg, _ := config.NewConfig()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Error Connect to Database")
	} else {
		fmt.Printf("successfully connect to database %s\n", cfg.Database)
	}

	db.AutoMigrate(
		&model.User{},
		&model.Book{},
		&model.Category{},
	)

	jwtService := service.NewJwtService(cfg.ApiConfig)

	bookRepository := repository.NewBookRepository(db)
	bookUsecase := usecase.NewBookUsecase(bookRepository)

	categoryRepository := repository.NewCategoryRepository(db)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepository)

	userRepository := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepository)

	authUsecase := usecase.NewAuthUsecase(userRepository, jwtService)


	engine := gin.Default()
	host := fmt.Sprintf(":%s", cfg.AppPort)
	return &Server{
		bookUsecase: bookUsecase,
		categoryUsecase: categoryUsecase,
		userUsecase: userUsecase,
		authUsecase: authUsecase,
		jwtService: jwtService,
		engine: engine,
		host: host,
	}
}

