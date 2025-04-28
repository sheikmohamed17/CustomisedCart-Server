package routes

import (
	"net/http"
	"practice/models"

	"github.com/gin-gonic/gin"
)

func ToolNameList(c *gin.Context) {
	res, err := models.ToolNameList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "Sucess",
		"data":    res,
	})
}

// specification get by toolName
func SpecGetByToolNameRoutes(c *gin.Context) {
	var data models.SpecificationWithTool
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Bad Request",
		})
		return
	}
	res, err := models.ScoreGetByToolName(data.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Internal Error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "Success",
		"data":    res,
	})

}

// Getting Score from our product Table
func SpecWithTool(c *gin.Context) {
	var data models.SpecificationWithTool
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Bad Request",
		})
		return
	}
	res, err := models.ScoreCategoryGetByID(data.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Internal Error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "Success",
		"data":    res,
	})

}

// new Code
func SpecificationWithTools(c *gin.Context) {
	var data models.SpecificationWithTools
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Bad Request",
		})
		return
	}
	detailedData, SimilarProducts, err := models.DetailedBenchmark(data.ToolNameID, data.SpecID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Internal Error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message":         "Success",
		"detailedData":    detailedData,
		"SimilarProducts": SimilarProducts,
	})
}

// Demo Testing Code
func SpecGetByScoreRange(c *gin.Context) {
	var data models.SpecificationWithTools1
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Bad Request",
		})
		return
	}
	detailedData, SimilarProducts, err := models.DetailedBenchmark1(data.ToolNameID, data.SpecID, data.StartingRange, data.EndingRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Internal Error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message":         "Success",
		"detailedData":    detailedData,
		"SimilarProducts": SimilarProducts,
	})
}
