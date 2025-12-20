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
	authUsecase usecase.AuthUsecase
	userUsecase usecase.UserUsecase
	bookUsecase usecase.BookUsecase
	categoryUsecase usecase.CategoryUsecase
	cartUsecase usecase.CartUsecase
	orderUsecase usecase.OrderUsecase
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
	controller.NewCartController(s.cartUsecase, authGroup)
	controller.NewOrderController(s.orderUsecase, s.cartUsecase, authGroup)
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
		&model.Order{},
		&model.OrderItem{},
	)

	jwtService := service.NewJwtService(cfg.ApiConfig)

	bookRepository := repository.NewBookRepository(db)
	bookUsecase := usecase.NewBookUsecase(bookRepository)

	categoryRepository := repository.NewCategoryRepository(db)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepository)

	userRepository := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepository)

	authUsecase := usecase.NewAuthUsecase(userRepository, jwtService)

	// Initialize Redis Client
	redisClient := config.NewRedisClient(cfg)

	cartRepository := repository.NewCartRepository(redisClient)
	cartUsecase := usecase.NewCartUsecase(cartRepository, bookUsecase)

	orderRepository := repository.NewOrderRepository(db)
	orderUsecase := usecase.NewOrderUsecase(orderRepository)


	engine := gin.Default()
	host := fmt.Sprintf(":%s", cfg.AppPort)
	return &Server{
		authUsecase: authUsecase,
		jwtService: jwtService,
		userUsecase: userUsecase,
		bookUsecase: bookUsecase,
		categoryUsecase: categoryUsecase,
		cartUsecase: cartUsecase,
		orderUsecase: orderUsecase,
		engine: engine,
		host: host,
	}
}

