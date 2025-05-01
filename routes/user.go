package routes

import (
	"example/models"
	"example/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func signup(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request"})
		return
	}

	err = user.Save() // store that user in the databese
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save users. Try again later."})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "User created successfuly!"})
}

func login(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user) // bind the incoming request body to that struct
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request"})
		return
	}

	err = user.ValidateCredentials() //kimlik bilgilerini kontrole diyor
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Could not authenticate user."})
		return
	}

	myUser, err := models.FindUserByEmail(user.Email)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch user."})
		return
	}

	//structtaki password ile databasedeki password aynıysa token oluştur.
	token, err := utils.GenerateToken(user.Email, myUser.ID, myUser.Role)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not authenticate user."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Login Successful!", "token": token})
}

func getUsers(context *gin.Context) {
	role := context.GetString("role")
	if role != "admin" {
		context.JSON(http.StatusForbidden, gin.H{"message": "No permission to get users"})
		return
	}

	users, err := models.GetAllUsers()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": " Could not get users"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"users": users})
}

func getUserById(context *gin.Context) {
	strid := context.Param("id")
	id, err := strconv.ParseInt(strid, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse user id."})
		return
	}

	role := context.GetString("role")
	if role != "admin" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "no permission to ger user"})
		return
	}

	user, err := models.GetUserById(id)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Could not fetch user."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"user": user})
}

func updateUser(context *gin.Context) {
	idStr := context.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse user id"})
		return
	}

	var user models.User
	err = context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON body"})
		return
	}

	if context.GetString("role") != "admin" && context.GetInt64("userid") != id {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user"})
		return
	}

	err = user.UserUpdate(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not updated users. Try again later."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "User updated sucessfuly"})
}
