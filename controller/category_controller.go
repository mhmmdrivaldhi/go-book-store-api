package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mhmmmdrivaldhi/go-book-api/middleware"
	"github.com/mhmmmdrivaldhi/go-book-api/model/dto"
	"github.com/mhmmmdrivaldhi/go-book-api/usecase"
)

type categoryController struct {
	categoryUsecase usecase.CategoryUsecase
}

func (cc *categoryController) CreateCategory(ctx *gin.Context) {
	var req dto.CreateCategoryRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	createCategory, err := cc.categoryUsecase.Create(req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, dto.GeneralResponse{
		Message: "successfully created category",
		Data: createCategory,
	})
}

func (cc *categoryController) GetAllCategories(ctx *gin.Context) {
	categories, err := cc.categoryUsecase.GetAll()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.GeneralResponse{
		Message: "successfully get all categories",
		Data: categories,
	})
}

func (cc *categoryController) GetCategoryById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	category, err := cc.categoryUsecase.GetById(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, dto.ErrorResponse{Error: "category not found"})
		return
	}

	ctx.JSON(http.StatusOK, dto.GeneralResponse{
		Message: "successfully get category by id",
		Data: category,
	})
}

func (cc *categoryController) UpdateCategory(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	var reqUpdate dto.UpdateCategoryRequest 

	err = ctx.ShouldBindJSON(&reqUpdate)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	updateCategory, err := cc.categoryUsecase.Update(id, &reqUpdate)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.GeneralResponse{
		Message: "successfully updated category",
		Data: updateCategory,
	})
}

func (cc *categoryController) DeleteCategory(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	err = cc.categoryUsecase.Delete(id) 
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return 
	}

	ctx.JSON(http.StatusNoContent, dto.GeneralResponse{
		Message: "successfully deleted category",
	})
}

func NewCategoryController(cu usecase.CategoryUsecase, rg *gin.RouterGroup) *categoryController {
	controller := &categoryController{categoryUsecase: cu}

	// public routes
	rg.GET("/category", controller.GetAllCategories)
	rg.GET("/category/:id", controller.GetCategoryById)

	// routes with middleware protected
	protected := rg.Group("")
	protected.Use(middleware.RoleMiddleware("admin"))

	protected.POST("/category", controller.CreateCategory)
	protected.PUT("/category/:id", controller.UpdateCategory)
	protected.DELETE("/category/:id", controller.DeleteCategory)

	return controller
}