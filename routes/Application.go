package routes

import (
	"net/http"
	"practice/models"

	"github.com/gin-gonic/gin"
)

// new added functionalities
func AppNamesList(c *gin.Context) {
	res, err := models.AppNamesList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "Success",
		"data":    res,
	})
}

func AppCategoryGetByID(c *gin.Context) {
	var data models.SpecificationWithApps
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Bad Request": err,
		})
		return
	}
	detailedData, SimilarProducts, err := models.AppCategoryGetByID(data.AppTypeID, data.SpecID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message":         "Success",
		"detailedData":    detailedData,
		"SimilarProducts": SimilarProducts,
	})
}

// Category Names Lists
func CategoryAppNamesList(c *gin.Context) {
	var data models.App
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Bad Request": err,
		})
		return
	}
	res, err := models.GetAllAppCategoryNames(data.AppTypeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "Success",
		"data":    res,
	})
}

// Multiple Apps Selection
func MultiAppsSelection(c *gin.Context) {
	var data models.SpecificationWithApps
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Bad Request": err,
		})
		return
	}
	detailedData, SimilarProducts, err := models.MultipleAppsSelection(data.AppsString, data.AppTypeName, data.AppTypeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message":         "Success",
		"detailedData":    detailedData,
		"SimilarProducts": SimilarProducts,
	})
}
