package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mhmmmdrivaldhi/go-book-api/model/dto"
	"github.com/mhmmmdrivaldhi/go-book-api/usecase"
)

type userController struct {
	rg *gin.RouterGroup
	userUsecase usecase.UserUsecase
}

func (uc *userController) Route() {
	uc.rg.POST("/user", uc.createUser)
	uc.rg.GET("/user", uc.getAllUser)
	uc.rg.GET("/user/:id", uc.getUserById)
	uc.rg.PUT("/user/:id", uc.updateUser)
	uc.rg.DELETE("/user/:id", uc.deleteUser)
}

func (uc *userController) createUser(ctx *gin.Context) {
	var req dto.CreateUserRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	created, err := uc.userUsecase.Create(req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, dto.GeneralResponse{
		Message: "Successfully Created User",
		Data: created,
	})
}

func (uc *userController) getAllUser(ctx *gin.Context) {
	users, err := uc.userUsecase.GetAllUsers()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.GeneralResponse{
		Message: "Successfully Get Data All Users",
		Data: users,
	})
}

func (uc *userController) getUserById(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	
	user, err := uc.userUsecase.GetUserById(userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.GeneralResponse{
		Message: "Successfully Get User By ID",
		Data: user,
	})
}

func (uc *userController) updateUser(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}

	var reqUpdate dto.UpdateUserRequest
	err = ctx.ShouldBindJSON(&reqUpdate)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	userUpdate, err := uc.userUsecase.UpdateUser(userId, reqUpdate)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.GeneralResponse{
		Message: "Successfully Update Data User",
		Data: userUpdate,
	})
}

func (uc *userController) deleteUser(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	
	}

	err = uc.userUsecase.DeleteUser(userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.GeneralResponse{
		Message: "Successfully Delete Data User",
	})
}

func NewUserController(userUsecase usecase.UserUsecase, rg *gin.RouterGroup) *userController {
	return &userController{
		userUsecase: userUsecase,
		rg: rg,
	}
}