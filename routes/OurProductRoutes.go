package routes

import (
	"fmt"
	"net/http"
	"practice/models"

	"github.com/gin-gonic/gin"
)

func OurProductInsertingRoutes(c *gin.Context) {
	var data models.Products
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message":  "Bad Request",
			"error":    err,
			"Message1": fmt.Sprintf("%v", err),
		})
		return
	}
	res, err := models.OurProductInserting(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Internal error",
			"error":   err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "Success",
		"data":    res,
	})
}

// sending product Category for selecting the user to move on
func GetProductDetails(c *gin.Context) {
	var data models.AllProduct
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Bad Request",
			"error":   err,
		})
		return
	}
	res, err := models.ProductDetails(data.Id)
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

// sending the information to front end  after fetching the id
func OurProductFetchSpecification(c *gin.Context) {
	var data models.Products
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Bad Request",
		})
		return
	}
	DetailedData, Similar, err := models.GetDetailedOurProducts(data.ProductId, data.ProductTypeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message":         "Success",
		"detailedData":    DetailedData,
		"SimilarProducts": Similar,
	})
}
