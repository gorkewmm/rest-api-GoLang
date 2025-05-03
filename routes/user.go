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

	context.JSON(http.StatusOK, gin.H{
		"message": "Login Successful!",
		"token":   token,
		"role":    myUser.Role, // 👈 role bilgisini ekle!
	})
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

func deleteUser(context *gin.Context) {
	idStr := context.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse user id"})
		return
	}

	if context.GetInt64("userid") != id && context.GetString("role") != "admin" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "No permission to delete user"})
		return
	}

	err = models.DeleteUser(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not delete user"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "User deleted sucessfuly"})
}

func changePassword(context *gin.Context) {
	idStr := context.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse user id"})
		return
	}

	userid := context.GetInt64("userid")
	if userid != id {
		context.JSON(http.StatusForbidden, gin.H{"message": "You can only change your own password"})
		return
	}

	user, err := models.GetUserById(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not get user"})
		return
	}

	var requestBody struct {
		Password    string `json:"password" binding:"required"`    // mevcut şifre
		NewPassword string `json:"newPassword" binding:"required"` // yeni şifre
	}
	err = context.ShouldBindJSON(&requestBody)

	if err != nil || requestBody.NewPassword == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Missing or invalid new password"})
		return
	}

	bool := utils.CheckPasswordHash(requestBody.Password, user.Password)
	if !bool {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Incorrect current password"})
		return
	}

	var myuser models.User
	myuser.Password = requestBody.NewPassword

	err = myuser.ChangePassword(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Password update failed"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Password changed successfully!"})
}
