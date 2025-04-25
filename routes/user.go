package routes

import (
	"example/models"
	"example/utils"
	"net/http"

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

	//structtaki password ile databasedeki password aynıysa token oluştur.
	token, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not authenticate user."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Login Successful!", "token": token})
}
