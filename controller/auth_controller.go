package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mhmmmdrivaldhi/go-book-api/model/dto"
	"github.com/mhmmmdrivaldhi/go-book-api/usecase"
)

type AuthController struct {
	authUsecase usecase.AuthUsecase
	validate *validator.Validate
}

func (ac *AuthController) Login(ctx *gin.Context) {
	var reqLogin dto.LoginRequest

	err := ctx.ShouldBindJSON(&reqLogin)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	err = ac.validate.Struct(reqLogin)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	resp, err := ac.authUsecase.Login(reqLogin)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.GeneralResponse{
		Message: "successfully login",
		Data: resp,
	})
}

func NewAuthController(authUsecase usecase.AuthUsecase, rg *gin.RouterGroup) {
	authController :=  &AuthController{
		authUsecase: authUsecase,
		validate: validator.New(),
	}

	rg.POST("/login", authController.Login)
}