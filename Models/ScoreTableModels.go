package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type User struct {
	ID        int    `json:"id"`
	ToolName  string `json:"toolName"`
	ToolImage string `json:"toolImage"`
	MinScore  int    `json:"minScore"`
	MaxScore  int    `json:"maxScore"`
}

func ToolNameList() ([]User, error) {
	var data []User
	Query := `SELECT t1.Tid, t1.ToolName, t1.ToolImage, MIN(t2.ToolScore), MAX(t2.ToolScore) FROM toolname t1 JOIN ToolScore t2 ON t2.Tid = t1.Tid GROUP BY t1.Tid`
	res, err := DB.Query(Query)
	if err != nil {
		log.Print("error in executingQuery")
		return nil, err
	}
	for res.Next() {
		var data1 User
		err := res.Scan(&data1.ID, &data1.ToolName, &data1.ToolImage, &data1.MinScore, &data1.MaxScore)
		if err != nil {
			log.Print("executing Query", err)
			return nil, err
		}
		data = append(data, User{
			ID:        data1.ID,
			ToolName:  data1.ToolName,
			ToolImage: data1.ToolImage,
			MinScore:  data1.MinScore,
			MaxScore:  data1.MaxScore,
		})
	}
	return data, nil
}

// new code Benchmark Last Page Functionalities
type SpecificationTool struct {
	Processor   string `json:"processor"`
	Motherboard string `json:"motherBoard"`
	Ram         string `json:"ram"`
	GPU         string `json:"gpu"`
	SMPS        string `json:"smps"`
	Case        string `json:"case"`
	Cooler      string `json:"cooler"`
	Storage     string `json:"storage"`
}
type Score struct {
	ToolId    int    `json:"toolid"`
	ToolScore int    `json:"toolscore"`
	ToolName  string `json:"toolName"`
}

type SpecificationWithTool struct {
	ID            int               `json:"id"`
	ProductModel  string            `json:"productModel"`
	Specification SpecificationTool `json:"specification"`
	Scores        []Score           `json:"scores"`
}

// new code Benchmark Getting Spec Using Score ID
func ScoreGetByToolName(ScoreID int) ([]SpecificationWithTool, error) {
	var result []SpecificationWithTool

	// Query product data
	productQuery := `SELECT id, score_id, ProductModel, specification FROM ourproduct2`
	productRes, err := DB.Query(productQuery)
	if err != nil {
		log.Printf("Error executing product query: %v", err)
		return nil, err
	}
	defer productRes.Close()

	for productRes.Next() {
		var (
			productId      int
			scoreIDsString string
			productModel   string
			specification  string
		)
		// Scan product data
		err := productRes.Scan(&productId, &scoreIDsString, &productModel, &specification)
		if err != nil {
			log.Printf("Error scanning product data: %v", err)
			continue
		}

		// Parse specification JSON
		var spec SpecificationTool
		if err := json.Unmarshal([]byte(specification), &spec); err != nil {
			log.Printf("Error parsing specification JSON: %v", err)
			continue
		}
		// Prepare a list of scores for the current specification
		scores := []Score{}
		if scoreIDsString != "" {
			scoreIDs := strings.Split(scoreIDsString, ",")
			for _, idStr := range scoreIDs {
				idStr = strings.TrimSpace(idStr)
				if idStr == "" {
					continue
				}

				scoreID, err := strconv.Atoi(idStr)
				if err != nil {
					log.Printf("Error converting score ID '%s' to integer: %v", idStr, err)
					continue
				}
				// Fetch score details
				// scoreQuery := `SELECT id, ToolScore FROM toolscore WHERE id=? AND Tid=?`
				scoreQuery := `SELECT ts.id, tn.ToolName, ts.ToolScore FROM toolscore ts INNER JOIN toolname tn ON ts.Tid = tn.Tid WHERE ts.id=? AND tn.Tid=?`
				var toolScore int
				var toolName string
				err = DB.QueryRow(scoreQuery, scoreID, ScoreID).Scan(&scoreID, &toolName, &toolScore)
				if err != nil {
					if err != sql.ErrNoRows {
						log.Printf("Error fetching score details: %v", err)
					}
					continue
				}

				scores = append(scores, Score{
					ToolId:    scoreID,
					ToolName:  toolName,
					ToolScore: toolScore,
				})
			}
		}
		// Append the specification and associated scores to the result
		if len(scores) > 0 {
			result = append(result, SpecificationWithTool{
				ID:            productId,
				ProductModel:  productModel,
				Specification: spec,
				Scores:        scores,
			})
		}
	}
	return result, nil
}

// Using Spec ID Fetching the data From our Product Table
func ScoreCategoryGetByID(SpecId int) ([]SpecificationWithTool, error) {
	var result []SpecificationWithTool

	// Query product data
	productQuery := `SELECT id, score_id, ProductModel, specification FROM ourproduct2 WHERE id=?`
	productRes, err := DB.Query(productQuery, SpecId)
	if err != nil {
		log.Printf("Error executing product query: %v", err)
		return nil, err
	}
	defer productRes.Close()

	for productRes.Next() {
		var (
			productId      int
			scoreIDsString string
			productModel   string
			specification  string
		)

		// Scan product data
		err := productRes.Scan(&productId, &scoreIDsString, &productModel, &specification)
		if err != nil {
			log.Printf("Error scanning product data: %v", err)
			continue
		}

		// Parse specification JSON
		var spec SpecificationTool
		if err := json.Unmarshal([]byte(specification), &spec); err != nil {
			log.Printf("Error parsing specification JSON: %v", err)
			continue
		}

		// Prepare a list of scores for the current specification
		scores := []Score{}

		if scoreIDsString != "" {
			scoreIDs := strings.Split(scoreIDsString, ",")
			for _, idStr := range scoreIDs {
				idStr = strings.TrimSpace(idStr)
				if idStr == "" {
					continue
				}

				scoreID, err := strconv.Atoi(idStr)
				if err != nil {
					log.Printf("Error converting score ID '%s' to integer: %v", idStr, err)
					continue
				}

				// Fetch score details
				scoreQuery := `SELECT ts.id, tn.ToolName, ts.ToolScore FROM toolscore ts INNER JOIN toolname tn ON ts.Tid = tn.Tid WHERE ts.id=?`
				var toolScore int
				var toolName string
				err = DB.QueryRow(scoreQuery, scoreID).Scan(&scoreID, &toolName, &toolScore)
				if err != nil {
					if err != sql.ErrNoRows {
						log.Printf("Error fetching score details: %v", err)
					}
					continue
				}

				scores = append(scores, Score{
					ToolId:    scoreID,
					ToolName:  toolName,
					ToolScore: toolScore,
				})
			}
		}

		// Append the specification and associated scores to the result
		if len(scores) > 0 {
			result = append(result, SpecificationWithTool{
				ID:           productId,
				ProductModel: productModel,
				// Specification: spec,
				Scores: scores,
			})
		}
	}
	return result, nil
}

// type App struct {
// 	AppID   int    `json:"appId"`
// 	AppName string `json:"appName"`
// }

//	type Scores struct {
//		ToolName  string `json:"toolName"`
//		ToolScore int    `json:"toolScore"`
//	}
type SpecificationStringTool struct {
	ComponentTypeID int    `json:"componentTypeID"`
	ComponentType   string `json:"componentType"`
	ComponentName   string `json:"componentName"`
	SpecString      string `json:"specString"`
	ID              int    `json:"id"`
	SpecName        string `json:"specName"`
	SocketNumber    string `json:"socketNumber"`
}
type SpecificationWithTools struct {
	ToolNameID         int                       `json:"toolNameID"`
	SpecID             int                       `json:"SpecID"`
	ProductModel       string                    `json:"productModel"`
	Specification      []SpecificationStringTool `json:"specification"`
	ProductDescription string                    `json:"productDescription"`
	ProductOverview    string                    `json:"productOverview"`
	AllScores          []ScoreId                 `json:"AllScores"`
	ProductImage       string                    `json:"productImage"`
	Apps               []App                     `json:"apps"`
	Score              []ScoreId                 `json:"scores"`
}

// processing the Specification Details
func convertToSpecificationToolFetch(input []SpecificationString) []SpecificationStringTool {
	var output []SpecificationStringTool
	for _, spec := range input {
		output = append(output, SpecificationStringTool(spec))
	}
	return output
}

// Code for Benchmarking Details
func DetailedBenchmark(ToolNameID, SpecID int) (detailedData SpecificationWithTools, SimilarProducts []SpecificationWithTools, err error) {
	var result []SpecificationWithTools
	var productRows *sql.Rows
	var rows *sql.Rows
	fmt.Printf("ToolNameID %v\n SpecID %v \n", ToolNameID, SpecID)
	var displayedProductIDs []int
	if SpecID != 0 {
		displayedProductIDs = append(displayedProductIDs, SpecID)
	}
	var (
		// productId          int
		appIDsString       string
		scoreIDsString     string
		productModel       string
		productOverview    string
		productImage       string
		productDescription string
		specification      string
	)
	var DetailedQuery string
	fmt.Printf("ToolID: %d, SpecID: %d\n", ToolNameID, SpecID)
	//Query for fetching the product details
	DetailedQuery = `SELECT
			id,
			app_id,
			score_id,
			ProductModel,
			ProductOverview,
			ProductImage,
			ProductDescription,
			specification	
		FROM
			ourproduct2 `
	productRows, err = DB.Query(DetailedQuery)
	if err != nil {
		log.Printf("Error executing product query: %v", err)
	}
	for productRows.Next() {
		var data SpecificationWithTools
		err := productRows.Scan(&data.SpecID, &appIDsString, &scoreIDsString, &productModel, &productOverview, &productImage, &productDescription, &specification)
		if err != nil {
			log.Printf("Error scanning product data: %v", err)
			continue
		}
		// Parse specification JSON
		res, err := SpecConverting(specification)
		if err != nil {
			log.Printf("Error converting specification JSON: %v", err)
			continue
		}
		data.Specification = convertToSpecificationToolFetch(res)
		//Getting Apps and Score
		scores := ScoreFetchingCondition(scoreIDsString, ToolNameID)
		//Fetching Apps
		apps, err := fetchApps(appIDsString)
		if err != nil {
			log.Print("error Executing in App fetching functionalities")
		}
		//Binding the data
		if len(scores) > 0 && len(apps) > 0 {
			result = append(result, SpecificationWithTools{
				ToolNameID:         ToolNameID,
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
	if len(result) == 0 {
		fmt.Printf("No Data Found")
		return SpecificationWithTools{}, nil, nil
	}
	if SpecID != 0 {
		for _, spec := range result {
			if spec.SpecID == SpecID {
				detailedData = spec
				fmt.Printf("Spec id DetailData %d\n", spec.SpecID)
				res, err := GetAllScores(SpecID)
				if err != nil {
					log.Printf("Error Executing in Score fetching functionalities")
				}
				detailedData.AllScores = res
			} else {
				fmt.Printf("Spec id !=0 %d\n", spec.SpecID)
				var score string
				var apps string
				detailedDataRandomProductQuery := `SELECT id, ProductModel, ProductOverview, ProductImage, ProductDescription, specification, score_id, app_id  FROM ourproduct2 WHERE id = ?`
				err = DB.QueryRow(detailedDataRandomProductQuery, SpecID).Scan(
					&detailedData.SpecID,
					&detailedData.ProductModel,
					&detailedData.ProductOverview,
					&detailedData.ProductImage,
					&detailedData.ProductDescription,
					&specification,
					&score,
					&apps,
				)
				fmt.Printf("score , apps%v %v\n", score, apps)
				if err != nil {
					log.Printf("Error fetching detailed data for ID %d: %v \n", SpecID, err)
					continue
				}
				res, err := SpecConverting(specification)
				if err != nil {
					log.Printf("Error converting specification JSON: %v", err)
					continue
				}
				detailedData.Specification = convertToSpecificationToolFetch(res)
				//Fetching App And String Details
				scores := ScoreFetchingCondition(score, ToolNameID)
				//fetching Apps
				App, err := fetchApps(apps)
				if err != nil {
					log.Print("error Fetching the App details")
				}
				detailedData.Apps = App
				detailedData.Score = scores
				// Getting All scores
				fmt.Printf("Detailed SpecID %v\n", detailedData.SpecID)
				result, err := GetAllScores(detailedData.SpecID)
				if err != nil {
					log.Printf("Error Fetching the Scores")
				}
				detailedData.AllScores = result
				SimilarProducts = append(SimilarProducts, spec)
			}
		}
	} else if SpecID == 0 {
		detailedData = result[0]
		res, err := GetAllScores(detailedData.SpecID)
		if err != nil {
			log.Printf("Error Executing in Score fetching functionalities")
		}
		detailedData.AllScores = res
		SimilarProducts = result[1:]
	}
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
				var Random SpecificationWithTools
				err := rows.Scan(&Random.SpecID, &appIDsString, &scoreIDsString, &productModel, &productOverview, &productImage, &specification, &productDescription)
				if err != nil {
					log.Printf("Error scanning row for random products: %v \n", err)
					return detailedData, SimilarProducts, err
				}
				// Convert specification to SpecificationToolFetch format
				specApp, err := SpecConverting(specification)
				if err != nil {
					log.Printf("Invalid specification JSON for product ID %d: %v \n", Random.SpecID, err)
					continue
				}
				fmt.Print("Random Product Entry \n")
				Random.Specification = convertToSpecificationToolFetch(specApp)
				// Random.SpecID = productId
				Random.ProductModel = productModel
				Random.ProductImage = productImage
				SimilarProducts = append(SimilarProducts, Random)
			}
		}
	}
	return detailedData, SimilarProducts, nil
}

// fetching apps and scores function
func ScoreFetchingCondition(scoreIDsString string, ToolNameID int) []ScoreId {
	// Score Details
	// fmt.Printf("scoreIDsString Score fetching Condition %v\n", scoreIDsString)
	scores := []ScoreId{}
	if scoreIDsString != "" {
		scoreIDs := strings.Split(scoreIDsString, ",")
		for _, idStr := range scoreIDs {
			idStr = strings.TrimSpace(idStr)
			if idStr == "" {
				continue
			}
			scoreID, err := strconv.Atoi(idStr)
			if err != nil {
				log.Printf("Error converting score ID '%s': %v", idStr, err)
				continue
			}
			//Appending score details
			var score ScoreId
			queryScore := `SELECT t1.ToolName, t2.ToolScore FROM toolscore t2
						INNER JOIN toolname t1 ON t2.Tid = t1.Tid WHERE t2.id = ? AND t2.Tid = ?`
			err = DB.QueryRow(queryScore, scoreID, ToolNameID).Scan(&score.ToolName, &score.ToolScore)
			if err == sql.ErrNoRows {
				log.Printf("No toolscore found for ID %d", scoreID)
				continue
			}
			scores = append(scores, score)
		}
	}
	return scores
}

// Converting string to scores
func ConvertStringToScores(ScoreIDsString string) ([]ScoreId, error) {
	scores := []ScoreId{}
	if ScoreIDsString != "" {
		scoreIDs := strings.Split(ScoreIDsString, ",")
		for _, idStr := range scoreIDs {
			idStr = strings.TrimSpace(idStr)
			if idStr == "" {
				continue
			}
			scoreID, err := strconv.Atoi(idStr)
			if err != nil {
				log.Printf("Error converting score ID '%s': %v", idStr, err)
				continue
			}
			//Appending score details
			var score ScoreId
			queryScore := `SELECT ToolName, ToolScore FROM toolscore
			INNER JOIN toolname ON toolscore.Tid = toolname.Tid WHERE toolscore.id = ?`
			err = DB.QueryRow(queryScore, scoreID).Scan(&score.ToolName, &score.ToolScore)
			if err == sql.ErrNoRows {
				log.Printf("No toolscore found for ID %d", scoreID)
				continue
			}
			scores = append(scores, score)
		}
	}
	return scores, nil
}

// Fetching All Score
func GetAllScores(SpecID int) ([]ScoreId, error) {
	var data []ScoreId
	var ScoreIDsString string
	fmt.Printf("Received specID%v\n", SpecID)
	query := `SELECT score_id FROM ourproduct2 WHERE id=?`
	err := DB.QueryRow(query, SpecID).Scan(&ScoreIDsString)
	if err != nil {
		log.Print("Error Executing All Apps Query", err)
		return data, err
	}
	fmt.Printf("ScoreID string: %s\n", ScoreIDsString)
	res, err := ConvertStringToScores(ScoreIDsString)
	if err != nil {
		log.Print("error fetch Apps", err)

	}
	data = append(data, res...)
	return data, err
}

// Updated function
// Setting the Range to fetching the Specifications
// For Demo testing
func ScoreFetchingCondition1(scoreIDsString string, ToolNameID, StartingRange, EndingRange int) []ScoreId {
	// Score Details
	// fmt.Printf("scoreIDsString Score fetching Condition %v\n", scoreIDsString)
	fmt.Printf("ToolNameID %v\n", ToolNameID)
	fmt.Printf("StartingRange %v\nEndingRange %v\n", StartingRange, EndingRange)
	scores := []ScoreId{}
	if scoreIDsString != "" {
		scoreIDs := strings.Split(scoreIDsString, ",")
		for _, idStr := range scoreIDs {
			idStr = strings.TrimSpace(idStr)
			if idStr == "" {
				continue
			}
			scoreID, err := strconv.Atoi(idStr)
			if err != nil {
				log.Printf("Error converting score ID '%s': %v", idStr, err)
				continue
			}
			//Appending score details
			var score ScoreId
			queryScore := `SELECT t1.ToolName, t2.ToolScore FROM toolscore t2
						INNER JOIN toolname t1 ON t2.Tid = t1.Tid WHERE t2.id = ? AND t2.Tid = ?
						AND t2.ToolScore BETWEEN ? AND ?`
			err = DB.QueryRow(queryScore, scoreID, ToolNameID, StartingRange, EndingRange).Scan(&score.ToolName, &score.ToolScore)
			if err == sql.ErrNoRows {
				log.Printf("No toolscore found for ID %d", scoreID)
				continue
			}
			scores = append(scores, score)
		}
	}
	return scores
}

type SpecificationWithTools1 struct {
	ToolNameID         int                       `json:"toolNameID"`
	StartingRange      int                       `json:"startingRange"`
	EndingRange        int                       `json:"endingRange"`
	SpecID             int                       `json:"SpecID"`
	ProductModel       string                    `json:"productModel"`
	Specification      []SpecificationStringTool `json:"specification"`
	ProductDescription string                    `json:"productDescription"`
	ProductOverview    string                    `json:"productOverview"`
	AllScores          []ScoreId                 `json:"AllScores"`
	ProductImage       string                    `json:"productImage"`
	Apps               []App                     `json:"apps"`
	Score              []ScoreId                 `json:"scores"`
}

// Range Finding Functionality
// Demo code for range finding
func DetailedBenchmark1(ToolNameID, SpecID, StartingRange, EndingRange int) (detailedData SpecificationWithTools, SimilarProducts []SpecificationWithTools, err error) {
	var result []SpecificationWithTools
	var productRows *sql.Rows
	fmt.Printf("ToolNameID %v\n SpecID %v \n", ToolNameID, SpecID)
	var (
		appIDsString, scoreIDsString, productModel, productOverview, productImage, productDescription, specification string
	)

	// Construct SQL query to fetch products
	DetailedQuery := `SELECT id, app_id, score_id, ProductModel, ProductOverview, ProductImage, ProductDescription, specification FROM ourproduct2`
	productRows, err = DB.Query(DetailedQuery)
	if err != nil {
		log.Printf("Error executing product query: %v", err)
		return SpecificationWithTools{}, nil, err
	}
	defer productRows.Close()

	// Process product data
	for productRows.Next() {
		var data SpecificationWithTools
		err := productRows.Scan(&data.SpecID, &appIDsString, &scoreIDsString, &productModel, &productOverview, &productImage, &productDescription, &specification)
		if err != nil {
			log.Printf("Error scanning product data: %v", err)
			continue
		}

		// Convert specifications
		res, err := SpecConverting(specification)
		if err != nil {
			log.Printf("Error converting specification JSON: %v", err)
			continue
		}
		data.Specification = convertToSpecificationToolFetch(res)

		// Fetch Scores and Apps
		scores := ScoreFetchingCondition1(scoreIDsString, ToolNameID, StartingRange, EndingRange)
		apps, err := fetchApps(appIDsString)
		if err != nil {
			log.Printf("Error fetching app details: %v", err)
		}

		// Only append if data is valid
		if len(scores) > 0 && len(apps) > 0 {
			data.ToolNameID = ToolNameID
			data.ProductModel = productModel
			data.ProductDescription = productDescription
			data.ProductOverview = productOverview
			data.ProductImage = productImage
			data.Apps = apps
			data.Score = scores
			result = append(result, data)
		}
	}
	if len(result) == 0 {
		fmt.Println("No Data Found")
		return SpecificationWithTools{}, nil, nil
	}

	// Find the highest scoring product
	var highestScoreProduct SpecificationWithTools
	highestScore := 0

	if SpecID == 0 {
		for _, spec := range result {
			totalScore := 0
			for _, score := range spec.Score {
				totalScore += score.ToolScore // Assuming ToolScore is an integer
			}

			if totalScore > highestScore {
				highestScore = totalScore
				highestScoreProduct = spec
			}
		}
		detailedData = highestScoreProduct
	} else if SpecID != 0 {
		detailedData.SpecID = SpecID

	}

	// Fetch all scores for the highest score product
	res, err := GetAllScores(detailedData.SpecID)
	if err != nil {
		log.Printf("Error fetching all scores for SpecID %d: %v", detailedData.SpecID, err)
	}
	detailedData.AllScores = res

	// Store other products as similar products
	for _, spec := range result {
		if spec.SpecID != detailedData.SpecID {
			SimilarProducts = append(SimilarProducts, spec)
		}
	}
	return detailedData, SimilarProducts, nil
}
