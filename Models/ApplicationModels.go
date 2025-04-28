package models

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Application struct {
	ID       int    `json:"appId"`
	AppType  string `json:"appType"`
	AppImage string `json:"appImage"`
}

// new functionalities for sending front end data
func AppNamesList() ([]Application, error) {
	var data []Application
	Query := `SELECT Aid, AppType, AppImage FROM appcategory`
	res, err := DB.Query(Query)
	if err != nil {
		log.Print("error Executing Query", err)
		return nil, err
	}
	for res.Next() {
		var data1 Application
		err := res.Scan(&data1.ID, &data1.AppType, &data1.AppImage)
		if err != nil {
			log.Print("error occur executing query", err)
			return nil, err
		}
		fmt.Printf("AppImage%v", data1.AppImage)
		data = append(data, Application{
			ID:       data1.ID,
			AppType:  data1.AppType,
			AppImage: data1.AppImage,
		})
	}
	return data, nil
}

// Working Code for Application
type App struct {
	AppTypeID       int    `json:"appTypeId"`
	AppTypeIDstring string `json:"appTypeIdString"`
	AppID           int    `json:"appId"`
	AppName         string `json:"appName"`
}
type ScoreId struct {
	ToolName  string `json:"toolName"`
	ToolScore int    `json:"toolScore"`
}
type SpecificationStringFetch struct {
	ComponentTypeID int    `json:"componentTypeID"`
	ComponentType   string `json:"componentType"`
	ComponentName   string `json:"componentName"`
	SpecString      string `json:"specString"`
	ID              int    `json:"id"`
	SpecName        string `json:"specName"`
	SocketNumber    string `json:"socketNumber"`
}
type SpecificationWithApps struct {
	AppTypeID     int `json:"appTypeID"`
	ProductTypeId int `json:"productTypeId"`
	// ProductID          int                        `json:"productID"`
	SpecID             int                        `json:"specID"`
	ProductModel       string                     `json:"productModel"`
	Specification      []SpecificationStringFetch `json:"specification"`
	ProductDescription string                     `json:"productDescription"`
	ProductOverview    string                     `json:"productOverview"`
	ProductImage       string                     `json:"productImage"`
	Apps               []App                      `json:"apps"`
	FullApps           []App                      `json:"fullApps"`
	Score              []ScoreId                  `json:"scores"`
	AppTypeName        string                     `json:"appTypeName"`
	AppsString         string                     `json:"appsString"`
}

// Working code Currently
func AppCategoryGetByID(AppTypeID, SpecID int) (detailedData SpecificationWithApps, SimilarProducts []SpecificationWithApps, err error) {
	var result []SpecificationWithApps
	// Debugging logs
	fmt.Printf("Received AppTypeID: %d, SpecID: %d\n", AppTypeID, SpecID)
	var displayedProductIDs []int
	var rows *sql.Rows
	if SpecID != 0 {
		displayedProductIDs = append(displayedProductIDs, SpecID)
	}
	var productRows *sql.Rows
	var (
		productTypeId int
		// productId          int
		appIDsString       string
		scoreIDsString     string
		productModel       string
		productOverview    string
		productImage       string
		productDescription string
		specification      string
	)
	// Query for all products if no SpecID is provided
	var productQuery string
	if SpecID == 0 {
		fmt.Println("Fetching all products as SpecID is 0...")
		productQuery = `SELECT
			id,
			Pid,
			app_id,
			score_id,
			ProductModel,
			ProductOverview,
			ProductImage,
			productDescription,
			specification
		FROM
			ourproduct2`
		productRows, err = DB.Query(productQuery)
	} else {
		// Query for all products to identify detailedData and similar products
		fmt.Println("Fetching specific product and similar products...")

		productQuery = `SELECT id, Pid, app_id, score_id, ProductModel, ProductOverview, ProductImage, productDescription, specification FROM ourproduct2`

		productRows, err = DB.Query(productQuery)
		if err != nil {
			log.Printf("Error executing product query: %v", err)
			return SpecificationWithApps{}, nil, err
		}
	}
	if err != nil {
		log.Printf("Error executing product query: %v \n", err)
		return SpecificationWithApps{}, nil, err
	}
	defer productRows.Close()

	// Process query results
	for productRows.Next() {
		var data SpecificationWithApps
		err := productRows.Scan(&data.SpecID, &productTypeId, &appIDsString, &scoreIDsString, &productModel, &productOverview, &productImage, &productDescription, &specification)
		if err != nil {
			log.Printf("Error scanning product data: %v \n", err)
			continue
		}
		res, err := SpecConverting(specification)
		if err != nil {
			log.Printf("Error parsing specification JSON for product ID %d: %v \n", data.SpecID, err)
			continue
		}
		data.Specification = convertToSpecificationStringFetch(res)

		//Fetching App Details
		apps := fetchAppsWithCondition(appIDsString, AppTypeID)
		scores := fetchScores(scoreIDsString)

		// Add product to result
		if len(apps) > 0 {
			result = append(result, SpecificationWithApps{
				AppTypeID:          AppTypeID,
				ProductTypeId:      productTypeId,
				SpecID:             data.SpecID,
				ProductModel:       productModel,
				ProductDescription: productDescription,
				ProductOverview:    productOverview,
				ProductImage:       productImage,
				Specification:      data.Specification,
				Apps:               apps,
				Score:              scores,
			})
		}
	}
	// Check if no products were found
	if len(result) == 0 {
		return SpecificationWithApps{}, nil, fmt.Errorf("no specifications available")
	}
	fmt.Printf("SpecID Checking: %v\n", SpecID)
	// Find detailedData and similar products
	if SpecID != 0 {
		for _, spec := range result {
			fmt.Printf("SpecID: %d\n", spec.SpecID)
			if spec.SpecID == SpecID {
				// fmt.Printf("DetailedData: %v\n", spec)
				detailedData = spec
				//Getting All the Supported Apps
				res, err := GetAllApps(SpecID)
				if err != nil {
					log.Print("error in Fetch DetailSpec", err)
				}
				detailedData.FullApps = res

			} else {
				fmt.Println("Executing a random product DetailedData query")
				// fmt.Printf("SpecID: %d\n", spec.SpecID)
				var Apps string
				var Scores string
				// Query to fetch detailed data for a random product
				detailedDataRandomProductQuery := `SELECT id, Pid, ProductModel, ProductOverview, app_id, score_id, ProductImage, ProductDescription, specification  FROM ourproduct2 WHERE id = ?`
				err = DB.QueryRow(detailedDataRandomProductQuery, SpecID).Scan(
					&detailedData.SpecID,
					&detailedData.ProductTypeId,
					&detailedData.ProductModel,
					&detailedData.ProductOverview,
					&Apps,
					&Scores,
					&detailedData.ProductImage,
					&detailedData.ProductDescription,
					&specification,
				)
				if err != nil {
					log.Printf("Error fetching detailed data for ID %d: %v \n", SpecID, err)
					continue
				}
				res, err := SpecConverting(specification)
				if err != nil {
					log.Printf("Invalid specification JSON for product ID %d: %v \n", SpecID, err)
					continue
				}
				detailedData.Specification = convertToSpecificationStringFetch(res)
				//Fetching the Apps Details
				apps1 := fetchAppsWithCondition(Apps, AppTypeID)
				//Fetching score Details
				scores1 := fetchScores(Scores)
				detailedData.Apps = apps1
				detailedData.Score = scores1
				SimilarProducts = append(SimilarProducts, spec)
				//fetching All Supported Apps
				fmt.Printf("DetailedData Spec id: %v\n", detailedData.SpecID)
				AllAppsList, err := GetAllApps(detailedData.SpecID)
				if err != nil {
					log.Print("error in Fetch DetailSpec", err)
				}
				detailedData.FullApps = AllAppsList
			}
		}
	} else if SpecID == 0 {
		fmt.Println("SpecID is 0")
		detailedData = result[0]
		// Fetching the Apps Details
		fmt.Printf("DetailedData Spec id 0 problem: %v\n", detailedData.SpecID)
		fullApps, err := GetAllApps(detailedData.SpecID)
		if err != nil {
			log.Print("error in Fetch DetailSpec", err)
		}
		detailedData.FullApps = fullApps
		SimilarProducts = append(SimilarProducts, result[1:]...)
	}
	//Fetching Random Products
	if len(SimilarProducts) == 0 {
		fmt.Printf("Random Products Entry Message: %v\n", SimilarProducts)
		// Prepare a list of product IDs to exclude from random products
		excludeQuery := ""
		for i, id := range displayedProductIDs {
			if i == 0 {
				excludeQuery = fmt.Sprintf("%d \n", id)
			} else {
				excludeQuery = fmt.Sprintf("%s, %d \n", excludeQuery, id)
			}
		}
		// Fetch random products excluding already displayed product
		if len(SimilarProducts) == 0 {
			RandomProductQuery := ""
			if len(displayedProductIDs) > 0 {
				// Create a comma-separated list of IDs
				excludeQuery := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(displayedProductIDs)), ","), "[]")
				fmt.Println("Executing an random products query")
				RandomProductQuery = `
			SELECT
			id,
			Pid,
			app_id,
			score_id,
			ProductModel,
			ProductOverview,
			ProductImage,
			specification,
			ProductDescription
		FROM
			ourproduct2
		WHERE
			id NOT IN (` + excludeQuery + `)
			ORDER BY RAND() 
			LIMIT 3;`
				fmt.Print("Random Product first part")
			} else {
				// No IDs to exclude, fetch all random products
				RandomProductQuery = `
			SELECT
			id,
			Pid,
			app_id,
			score_id,
			ProductModel,
			ProductOverview,
			ProductImage,
			specification,
			ProductDescription
		FROM
			ourproduct2 
		ORDER BY RAND() 
		LIMIT 3;`
				fmt.Print("Random Product else Part")
			}
			rows, err = DB.Query(RandomProductQuery)
			if err != nil {
				log.Printf("Error fetching random products: %v \n", err)
				return detailedData, SimilarProducts, err
			}
			defer rows.Close()

			// Process random products
			for rows.Next() {
				var Random SpecificationWithApps
				err := rows.Scan(&Random.SpecID, &productTypeId, &appIDsString, &scoreIDsString, &productModel, &productOverview, &productImage, &specification, &productDescription)
				if err != nil {
					log.Printf("Error scanning row for random products: %v \n", err)
					return detailedData, SimilarProducts, err
				}
				// Unmarshal specification for the random product
				res, err := SpecConverting(specification)
				if err != nil {
					log.Printf("Invalid specification JSON for product ID %d: %v \n", Random.SpecID, err)
					continue
				}
				fmt.Print("Random Product Entry \n")
				Random.Specification = convertToSpecificationStringFetch(res)
				Random.ProductTypeId = productTypeId
				// Random.SpecID = SpecID
				Random.ProductModel = productModel
				Random.ProductImage = productImage
				SimilarProducts = append(SimilarProducts, Random)
			}
		}
	}
	return detailedData, SimilarProducts, nil
}
func convertToSpecificationStringFetch(input []SpecificationString) []SpecificationStringFetch {
	var output []SpecificationStringFetch
	for _, spec := range input {
		output = append(output, SpecificationStringFetch(spec))
	}
	return output
}

// Function for Creating Fetch Similar Products
func fetchAppsWithCondition(appIDsString string, AppTypeID int) []App {
	var apps []App
	fmt.Printf("appidstring: %v\n AppTypeID: %v\n", appIDsString, AppTypeID)
	if appIDsString == "" {
		return apps
	}

	for _, idStr := range strings.Split(appIDsString, ",") {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}

		appID, err := strconv.Atoi(idStr)
		if err != nil {
			continue
		}

		var appName string
		query := `SELECT AppName FROM appname WHERE id=?`
		if AppTypeID != 0 {
			query += ` AND Aid=?`
			err = DB.QueryRow(query, appID, AppTypeID).Scan(&appName)
		} else {
			err = DB.QueryRow(query, appID).Scan(&appName)
		}
		if err == nil {
			apps = append(apps, App{AppID: appID,
				AppName: appName})
		}
	}
	return apps
}

func fetchApps(appIDsString string) ([]App, error) {
	var apps []App
	if appIDsString == "" {
		return apps, nil
	}

	for _, idStr := range strings.Split(appIDsString, ",") {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}

		appID, err := strconv.Atoi(idStr)
		if err != nil {
			continue
		}
		var appName string
		query := `SELECT AppName FROM appname WHERE id=?`
		err = DB.QueryRow(query, appID).Scan(&appName)
		if err == nil {
			apps = append(apps, App{AppID: appID, AppName: appName})
		}
	}
	return apps, nil
}

func fetchScores(scoreIDsString string) []ScoreId {
	var scores []ScoreId
	if scoreIDsString == "" {
		return scores
	}

	for _, idStr := range strings.Split(scoreIDsString, ",") {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}

		scoreID, err := strconv.Atoi(idStr)
		if err != nil {
			continue
		}

		query := `SELECT ToolName, ToolScore FROM toolscore INNER JOIN toolname ON toolscore.Tid = toolname.Tid WHERE toolscore.id = ?`
		var score ScoreId
		if err := DB.QueryRow(query, scoreID).Scan(&score.ToolName, &score.ToolScore); err == nil {
			scores = append(scores, score)
		}
	}

	return scores
}

// Fetching All Apps
func GetAllApps(SpecID int) ([]App, error) {
	var data []App
	var appIDsString string
	query := `SELECT app_id FROM ourproduct2 WHERE id=?`
	err := DB.QueryRow(query, SpecID).Scan(&appIDsString)
	if err != nil {
		log.Print("Error Executing All Apps Query", err)
		return data, err
	}
	fmt.Printf("AppID string: %s\n", appIDsString)
	res, err := fetchApps(appIDsString)
	if err != nil {
		log.Print("error fetch Apps", err)

	}
	data = append(data, res...)
	return data, err
}

type Apps struct {
	AppID   int    `json:"appId"`
	AppName string `json:"appName"`
}

// Getting All AppCategory Names
func GetAllAppCategoryNames(AppTypeID int) ([]Apps, error) {
	var data []Apps
	query := `SELECT id, AppName FROM appname WHERE Aid=?`
	res, err := DB.Query(query, AppTypeID)
	if err != nil {
		log.Print("Error Executing All Apps Query", err)
		return data, err
	}
	for res.Next() {
		var data1 Apps
		err := res.Scan(&data1.AppID, &data1.AppName)
		if err != nil {
			log.Print("Error Scanning App Name", err)
			return data, err
		}
		data = append(data, data1)
	}
	return data, nil
}

// Test code
func MultipleAppsSelection(AppsString, AppTypeName1 string, AppTypeID int) (detailData SpecificationWithApps, SimilarData []SpecificationWithApps, err error) {
	var matchingSpecs []SpecificationWithApps
	var ArrayInt []int
	var specificationString string
	var ScoreIDsString string
	fmt.Printf("AppTypeID %v\n", AppTypeID)
	fmt.Printf("AppTypeName: %s\n", AppTypeName1)
	fmt.Printf("AppsString: %s\n", AppsString)
	if AppsString == "" {
		return detailData, SimilarData, nil
	}
	// Convert AppTypeIDstring to an array of integers
	for _, idStr := range strings.Split(AppsString, ",") {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}
		appID, err := strconv.Atoi(idStr)
		if err != nil {
			continue
		}
		ArrayInt = append(ArrayInt, appID)
	}

	// fmt.Printf("Input ArrayValues: %v\n", ArrayInt)

	// Map to count SpecID matches
	specCount := make(map[int]int)

	// Fetching OurProduct details
	query := `SELECT id, ProductModel, specification, ProductOverview, ProductImage, ProductDescription, score_id, app_id FROM ourproduct2`
	res, err := DB.Query(query)
	if err != nil {
		log.Print("Error Executing All Apps Query:", err)
		return detailData, SimilarData, err
	}
	defer res.Close()

	for res.Next() {
		var data1 SpecificationWithApps
		err := res.Scan(&data1.SpecID, &data1.ProductModel, &specificationString, &data1.ProductOverview, &data1.ProductImage, &data1.ProductDescription, &ScoreIDsString, &data1.AppsString)
		if err != nil {
			log.Print("Error Scanning App Name1:", err)
			return detailData, SimilarData, err
		}
		res, err := SpecConverting(specificationString)
		if err != nil {
			log.Print("Error Scanning App Name2:", err)
			return detailData, SimilarData, err
		}
		// binding Specifications
		data1.Specification = convertToSpecificationStringFetch(res)
		//Getting FullAll Apps
		// fmt.Printf("Specification Ids: %v\n", data1.SpecID)
		fullApps, errs := GetAllApps(data1.SpecID)
		if errs != nil {
			log.Print("Error Scanning App Name3:", err)
			return detailData, SimilarData, err
		}
		data1.FullApps = fullApps
		//Getting All Apps
		data1.Apps = fetchAppsWithCondition(AppsString, AppTypeID)
		//Getting All Sccores
		data1.Score = fetchScores(ScoreIDsString)
		//Binding AppID
		data1.AppTypeID = AppTypeID
		data1.AppTypeName = AppTypeName1
		// Convert AppsString to an array of integers
		var appIDs []int
		for _, appStr := range strings.Split(data1.AppsString, ",") {
			appStr = strings.TrimSpace(appStr)
			if appStr == "" {
				continue
			}
			appID, err := strconv.Atoi(appStr)
			if err != nil {
				log.Printf("Error converting appString to int: %v", err)
				continue
			}
			appIDs = append(appIDs, appID)
		}

		// Count matches between input ArrayInt and appIDs
		matchCount := 0
		for _, inputID := range ArrayInt {
			for _, dbID := range appIDs {
				if inputID == dbID {
					matchCount++
				}
			}
		}
		// Only include the spec if it has matches
		if matchCount > 0 {
			specCount[data1.SpecID] = matchCount
			matchingSpecs = append(matchingSpecs, data1)
		}
	}

	// Find the SpecID with the highest count
	highestSpecID := 0
	maxCount := 0
	for specID, count := range specCount {
		if count > maxCount {
			highestSpecID = specID
			maxCount = count
		}
	}
	fmt.Printf("Highest SpecID: %d with Count: %d\n", highestSpecID, maxCount)

	// Filter the results to include only the highest matching SpecID
	var highestMatchingSpec SpecificationWithApps
	var similarSpecs []SpecificationWithApps
	for _, spec := range matchingSpecs {
		if spec.SpecID == highestSpecID {
			highestMatchingSpec = spec
		} else {
			similarSpecs = append(similarSpecs, spec)
		}
	}
	return highestMatchingSpec, similarSpecs, nil
}
