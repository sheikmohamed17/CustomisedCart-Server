package routes

import (
	"net/http"
	"practice/models"

	"github.com/gin-gonic/gin"
)

// selecting Processor
func GettingProcessor(c *gin.Context) {
	var data models.Components
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": err,
		})
		return
	}
	res, err := models.GettingProcessor(data.ComponentType)
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

// data ComoponetList sending
func ComponentList(c *gin.Context) {
	var data models.ComponentList
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": err,
		})
		return
	}
	res, err := models.GettingComponentList(data.ComponentID, data.CompSpecID, data.SocketNumber)
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

// Customize Getting Data
// Testing Code for an Checking the data
func CustomizedDataRoutes(c *gin.Context) {
	var data models.Products
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Bad Request",
		})
		return
	}
	res, err := models.CustomizedData(data.SpecString)
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
func CustomizedDataRoutes1(c *gin.Context) {
	var data models.Products
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Bad Request",
		})
		return
	}
	res, err := models.CustomizedData1(data.SpecString, data.ProductId)
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

// Full Customization
func FullCustomized(c *gin.Context) {
	var data models.ComponentList1
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Bad Request",
		})
		return
	}
	res, err := models.FullCustomization(data.ComponentType, data.SocketNumber, data.SupportedRam, data.GraphicsCard)
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
