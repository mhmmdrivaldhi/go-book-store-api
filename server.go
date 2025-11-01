package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mhmmmdrivaldhi/go-book-api/config"
	"github.com/mhmmmdrivaldhi/go-book-api/controller"
	"github.com/mhmmmdrivaldhi/go-book-api/model"
	"github.com/mhmmmdrivaldhi/go-book-api/repository"
	"github.com/mhmmmdrivaldhi/go-book-api/usecase"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	bookUsecase usecase.BookUsecase
	engine *gin.Engine
	host string
}

func (s *Server) InitRoute() {
	v1 := s.engine.Group("/v1")

	controller.NewBookController(s.bookUsecase, v1).Route()
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

	db.AutoMigrate(&model.Book{})

	bookRepository := repository.NewBookRepository(db)
	bookUsecase := usecase.NewBookUsecase(bookRepository)


	engine := gin.Default()
	host := fmt.Sprintf(":%s", cfg.AppPort)
	return &Server{
		bookUsecase: bookUsecase,
		engine: engine,
		host: host,
	}
}

