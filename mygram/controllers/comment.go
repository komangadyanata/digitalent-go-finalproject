package controllers

import (
	"encoding/json"
	"mygram/configs"
	"mygram/models"
	"mygram/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateComment(c *gin.Context) {
	db := configs.GetDB()
	contentType := utils.GetContentType(c)

	commentRequest := models.CreateCommentRequest{}
	userID := utils.GetUserIdFromToken(c)

	if contentType == appJSON {
		if err := c.ShouldBindJSON(&commentRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
	} else {
		if err := c.ShouldBind(&commentRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
	}

	comment := models.Comment{
		PhotoId: commentRequest.PhotoId,
		Message: commentRequest.Message,
		UserId:  userID,
	}

	err := db.Debug().Create(&comment).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	commentString, _ := json.Marshal(comment)
	commentResponse := models.CreateCommentResponse{}
	json.Unmarshal(commentString, &commentResponse)

	c.JSON(http.StatusCreated, commentResponse)
}

func GetComment(c *gin.Context) {
	db := configs.GetDB()
	comments := []models.Comment{}

	err := db.Debug().Preload("User").Preload("Photo").Order("id asc").Find(&comments).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	commentsString, _ := json.Marshal(comments)
	commentsResponse := []models.CommentResponse{}
	json.Unmarshal(commentsString, &commentsResponse)

	c.JSON(http.StatusOK, commentsResponse)
}

func UpdateComment(c *gin.Context) {
	db := configs.GetDB()
	contentType := utils.GetContentType(c)

	commentRequest := models.UpdateCommentRequest{}
	commentId, _ := strconv.Atoi(c.Param("commentId"))
	userID := utils.GetUserIdFromToken(c)

	if contentType == appJSON {
		if err := c.ShouldBindJSON(&commentRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
	} else {
		if err := c.ShouldBind(&commentRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
	}

	comment := models.Comment{}
	comment.ID = uint(commentId)
	comment.UserId = userID

	updateString, _ := json.Marshal(commentRequest)
	updateData := models.Comment{}
	json.Unmarshal(updateString, &updateData)

	err := db.Debug().Model(&comment).Updates(updateData).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}
	_ = db.First(&comment, comment.ID).Error

	commentString, _ := json.Marshal(comment)
	commentResponse := models.UpdateCommentResponse{}
	json.Unmarshal(commentString, &commentResponse)

	c.JSON(http.StatusOK, commentResponse)
}

func DeleteComment(c *gin.Context) {
	db := configs.GetDB()
	commentId, _ := strconv.Atoi(c.Param("commentId"))
	userID := utils.GetUserIdFromToken(c)

	comment := models.Comment{}
	comment.ID = uint(commentId)
	comment.UserId = userID

	err := db.Debug().Delete(&comment).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your comment has been successfully deleted",
	})
}
