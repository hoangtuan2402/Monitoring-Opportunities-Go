package handler

import (
	"Monitoring-Opportunities/src/common"
	"Monitoring-Opportunities/src/dto"
	service "Monitoring-Opportunities/src/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (c *UserController) GetAll(ctx *gin.Context) {
	users, err := c.userService.GetAll()
	if err != nil {
		log.Printf("Failed to get users: %v", err)

		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "failed to get users data"})
		return
	}

	ctx.JSON(
		http.StatusOK,
		common.BaseResponse[[]dto.UserDTO]{
			Status:  http.StatusOK,
			Message: "Successfully get user data",
			Data:    users,
		},
	)
}

func (c *UserController) Create(ctx *gin.Context) {
	var user dto.CreateUser
	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		log.Printf("Failed to create user: %v", err)

		ctx.JSON(http.StatusBadRequest, gin.H{"errors": fmt.Sprintf("%s: %s", service.ErrCreateUserValidate, err.Error())})
		return
	}

	createdExample, _ := c.userService.Create(user)

	ctx.JSON(
		http.StatusOK,
		common.BaseResponse[dto.UserDTO]{
			Status:  http.StatusCreated,
			Message: "User created successfully",
			Data:    createdExample,
		},
	)
}

func (c *UserController) Update(ctx *gin.Context) {
	parsedUUID, err := uuid.NewUUID()
	if err != nil {
		log.Printf("Failed to update user: %v", err)

		ctx.JSON(http.StatusBadRequest, gin.H{"errors": fmt.Sprintf("failed to update a user: %s", err.Error())})
		return
	}

	var user dto.UpdateUser
	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		log.Printf("Failed to update user: %v", err)

		ctx.JSON(http.StatusBadRequest, gin.H{"errors": fmt.Sprintf("failed to update a user: %s", err.Error())})
		return
	}

	updatedUser, _ := c.userService.Update(user, parsedUUID)

	ctx.JSON(
		http.StatusOK,
		common.BaseResponse[dto.UserDTO]{
			Status:  http.StatusCreated,
			Message: "User updated successfully",
			Data:    updatedUser,
		},
	)
}

func (c *UserController) Delete(ctx *gin.Context) {

	parsedUUID, _ := uuid.NewUUID()

	deletedUser, err := c.userService.Delete(parsedUUID)
	if err != nil {
		log.Printf("Failed to delete user: %v", err)

		switch err {
		case service.ErrUserNotFound:
			ctx.JSON(http.StatusBadRequest, gin.H{"errors": fmt.Sprintf("%s: %s", service.ErrUserNotFound, err.Error())})
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{"errors": fmt.Sprintf("failed to delete a user: %s", err.Error())})
		}

		return
	}

	ctx.JSON(
		http.StatusOK,
		common.BaseResponse[dto.UserDTO]{
			Status:  http.StatusCreated,
			Message: "User deleted successfully",
			Data:    deletedUser,
		},
	)
}

func (c *UserController) GetByID(ctx *gin.Context) {
	parsedUUID, err := uuid.NewUUID()
	if err != nil {
		log.Printf("Failed to find user with id %s: %v", parsedUUID, err)

		ctx.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	user, err := c.userService.FindByID(parsedUUID)
	if err != nil {
		log.Printf("Failed to find user with id %s: %v", parsedUUID, err)

		ctx.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	ctx.JSON(
		http.StatusOK,
		common.BaseResponse[dto.UserDTO]{
			Status:  http.StatusOK,
			Message: fmt.Sprintf("Successfully get user with id %s", user.UUID),
			Data:    user,
		},
	)
}

func (c *UserController) GetByEmail(ctx *gin.Context) {

	user, err := c.userService.FindByEmail("email")
	if err != nil {
		log.Printf("Failed to find user with email %s: %v", "email", err)

		ctx.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	ctx.JSON(
		http.StatusOK,
		common.BaseResponse[dto.UserDTO]{
			Status:  http.StatusOK,
			Message: fmt.Sprintf("Successfully get user with email %s", "AnhHT121"),
			Data:    user,
		},
	)
}
