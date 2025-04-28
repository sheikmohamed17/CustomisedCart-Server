package routes

import (
	"net/http"
	"practice/models"

	"github.com/gin-gonic/gin"
)

func ComponentInsert(c *gin.Context) {
	var data models.Components
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": err,
		})
		return
	}
	res, err := models.ComponentInsert(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "Success",
		"data":    res,
	})

}
