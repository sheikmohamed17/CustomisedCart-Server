package routes

import (
	"fmt"
	"net/http"
	"practice/models"

	"github.com/gin-gonic/gin"
)

func RegistrationCheck(c *gin.Context) {
	var data models.Users
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Request Error",
			"message": err.Error(),
		})
		return
	}
	fmt.Printf("Users data admin routes user data  %s Password %s", data.Email, data.Password)
	// db := c.MustGet("db").(*sql.DB)
	res, err := models.RegistrationCheck(data.Email, data.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "internal error",
			"error":   err,
		})
		return
	}
	// fmt.Printf("user Email %s", data.Useremail)
	c.JSON(http.StatusOK, gin.H{
		"message": "inserted",
		"data":    res,
	})
}
