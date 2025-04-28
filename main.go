package main

import (
	"fmt"
	"log"

	// "net/http"        // Import the models package
	"practice/models" // Import the routes package
	"practice/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database connection
	if err := models.InitDB(); err != nil {
		log.Fatal("Database initialization failed:", err)
	}
	defer models.CloseDB()
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length"},
	}))
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		c.Writer.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		c.Next()
	})
	//Admin Login Credentials
	r.POST("/adminLogin", routes.RegistrationCheck)
	//Tool Score Table
	r.GET("/toolName", routes.ToolNameList)
	r.POST("/SpecToolScore", routes.SpecWithTool)
	r.POST("/ToolNameSpec", routes.SpecGetByToolNameRoutes)
	//new Updated Code
	r.POST("/ToolsSpecification", routes.SpecificationWithTools)

	//Component Table inserting
	r.POST("/componentInsert", routes.ComponentInsert)

	//Our product Inserting
	r.POST("/ourProductInsert", routes.OurProductInsertingRoutes)

	r.POST("/OurProducts", routes.GetProductDetails)
	//sending data frontend
	r.POST("/productSpec", routes.OurProductFetchSpecification)
	// r.POST("/specFetching", routes.GetOurProductDataFetchById)

	//Application table
	r.POST("/usageBasedSpecificationDetails", routes.AppCategoryGetByID)
	r.GET("/appCategory", routes.AppNamesList)

	// Customize Pc Section (AND) ComponentInserting
	//component Selecting
	r.POST("/ComponentType", routes.GettingProcessor)
	//Sending platform list like "i3 12th GEN"/"Customize Functions"
	r.POST("/GettingComponentList", routes.ComponentList)

	// Getting AppCategory Lists like "Gaming ="Gta,Gof Lists""
	r.POST("/GetAppCategoryList", routes.CategoryAppNamesList)

	//Multiple App Selecting
	r.POST("/MultiAppsSelection", routes.MultiAppsSelection)
	//Selecting MotherBoard Component

	// r.POST("/GettingMotherBoard", routes.GettingMotherBoard)
	//EmailFunctions Sending Users
	r.POST("/getQuote", routes.EmailSending)
	r.GET("/GetUserDetails", routes.GetUserDetails)
	//Demo
	r.POST("/demo", routes.CustomizedDataRoutes)
	r.POST("/demo1", routes.CustomizedDataRoutes1)

	//Full Customization functions
	r.POST("/fullCustomization", routes.FullCustomized)
	//getting Score By setting Range
	r.POST("/scoreByRange", routes.SpecGetByScoreRange)
	fmt.Println("Server is running on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
