package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mhmmmdrivaldhi/go-book-api/middleware"
	"github.com/mhmmmdrivaldhi/go-book-api/model/dto"
	"github.com/mhmmmdrivaldhi/go-book-api/usecase"
)

type bookController struct {
	bookUsecase usecase.BookUsecase
}

func (bc *bookController) CreateBook(ctx *gin.Context) {
	var req dto.CreateBookRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	createBook, err := bc.bookUsecase.Create(req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, dto.GeneralResponse{
		Message: "successfully created book",
		Data: createBook,
	})
}

func (bc *bookController) GetAllBook(ctx *gin.Context) {
	books, err := bc.bookUsecase.GetAll()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.GeneralResponse{
		Message: "successfully get all books",
		Data: books,
	})
}

func (bc *bookController) GetBookById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	book, err := bc.bookUsecase.GetById(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.GeneralResponse{
		Message: "successfully get book by id",
		Data: book,
	})
}

func (bc *bookController) UpdateBook(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	var reqUpdate dto.UpdateBookRequest
	err = ctx.ShouldBindJSON(&reqUpdate)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	
	updateBook, err := bc.bookUsecase.Update(id, reqUpdate)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.GeneralResponse{
		Message: "successfully updated book",
		Data: updateBook,
	})	
}

func (bc *bookController) DeleteBook(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	err = bc.bookUsecase.Delete(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, dto.GeneralResponse{
		Message: "succesfully deleted book",
	})
}

func NewBookController(bu usecase.BookUsecase, rg *gin.RouterGroup) *bookController{
	controller := &bookController{bookUsecase: bu}

	// public routes
	rg.GET("/book", controller.GetAllBook)
	rg.GET("/book/:id", controller.GetBookById)

	// allowed roles routes
	protected := rg.Group("")
	protected.Use(middleware.RoleMiddleware("admin", "seller"))

	protected.POST("/book", controller.CreateBook)
	protected.PUT("/book/:id", controller.UpdateBook)
	protected.DELETE("/book/:id", controller.DeleteBook)

	return controller
}