package routes

import (
	emails "CustomizedCart/Emails"
	"net/http"

	"github.com/gin-gonic/gin"
)

func EmailSending(c *gin.Context) {
	var data emails.EmailDetails
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Bad Request",
		})
		return
	}
	res, err := emails.EmailSending(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message":  "internal error",
			"Message1": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "Success",
		"data":    res,
	})

}

func GetUserDetails(c *gin.Context) {
	// var data emails.EmailDetails
	res, err := emails.GetUserDetails()
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

// API Emails
// func DemoEmailCheck(c *gin.Context) {
// 	var data emails.EmailValidationResponse
// 	err := c.ShouldBindJSON(&data)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"Message": "Bad Request",
// 		})
// 		return
// 	}
// 	res, err := emails.validateEmail(data)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"Message": err,
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"Message": "Success",
// 		"data":    res,
// 	})
// }
