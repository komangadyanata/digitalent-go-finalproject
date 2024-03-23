package controllers

import (
	"encoding/json"
	"mygram/configs"
	"mygram/models"
	"mygram/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	appJSON = "application/json"
)

func UserRegister(c *gin.Context) {
	db := configs.GetDB()
	contentType := utils.GetContentType(c)

	userRequest := models.CreateUserRequest{}

	if contentType == appJSON {
		err := c.ShouldBindJSON(&userRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
	} else {
		err := c.ShouldBind(&userRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
	}

	strPass, err := utils.GenerateHashPassword(userRequest.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error Hash",
			"message": err.Error(),
		})
		return
	}

	user := models.User{
		Age:      userRequest.Age,
		Email:    userRequest.Email,
		Password: strPass,
		Username: userRequest.Username,
	}

	err = db.Debug().Create(&user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	userString, _ := json.Marshal(user)
	userResponse := models.CreateUserResponse{}
	json.Unmarshal(userString, &userResponse)

	c.JSON(http.StatusCreated, userResponse)
}

func UserLogin(c *gin.Context) {
	db := configs.GetDB()
	contentType := utils.GetContentType(c)

	userRequest := models.LoginUserRequest{}

	if contentType == appJSON {
		if err := c.ShouldBindJSON(&userRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
	} else {
		if err := c.ShouldBind(&userRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
	}

	password := userRequest.Password
	user := models.User{}

	err := db.Debug().Where("email = ?", userRequest.Email).Take(&user).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Invalid email/password",
		})
		return
	}

	comparePass := utils.CompareHashPassword(user.Password, password)
	if !comparePass {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Invalid email/password",
		})
		return
	}

	token := utils.GenerateToken(user.ID, user.Email)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func UpdateUser(c *gin.Context) {
	db := configs.GetDB()
	contentType := utils.GetContentType(c)

	userRequest := models.UpdateUserRequest{}
	userID := utils.GetUserIdFromToken(c)

	if contentType == appJSON {
		if err := c.ShouldBindJSON(&userRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
	} else {
		if err := c.ShouldBind(&userRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
	}

	user := models.User{}
	user.ID = userID

	updateString, _ := json.Marshal(userRequest)
	updateData := models.User{}
	json.Unmarshal(updateString, &updateData)

	err := db.Model(&user).Updates(updateData).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}
	_ = db.First(&user, user.ID).Error

	userString, _ := json.Marshal(user)
	userResponse := models.UpdateUserResponse{}
	json.Unmarshal(userString, &userResponse)

	c.JSON(http.StatusCreated, userResponse)
}

func DeleteUser(c *gin.Context) {
	db := configs.GetDB()

	userID := utils.GetUserIdFromToken(c)

	user := models.User{}
	user.ID = userID

	db.Where("user_id = ?", user.ID).Delete(&models.Comment{})
	db.Where("user_id = ?", user.ID).Delete(&models.Photo{})
	db.Where("user_id = ?", user.ID).Delete(&models.SocialMedia{})
	err := db.Delete(&user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your account has been successfully deleted",
	})
}
