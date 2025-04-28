package models

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type ScoreID struct {
	ToolName  string `json:"toolName"`
	ToolScore string `json:"toolScore"`
	ID        int    `json:"id"`
}
type AppID struct {
	AppType  string `json:"appType"`
	AppName  string `json:"appName"`
	Category string `json:"appChoice"`
	ID       int    `json:"id"`
}

// Demo For spec Fetching
type SpecificationString struct {
	ComponentTypeID int    `json:"componentTypeID"`
	ComponentType   string `json:"componentType"`
	ComponentName   string `json:"componentName"`
	SpecString      string `json:"specString"`
	ID              int    `json:"id"`
	SpecName        string `json:"specName"`
	SocketNumber    string `json:"socketNumber"`
}

type Products struct {
	ProductId          int                   `json:"productId"`
	ProductTypeID      int                   `json:"productTypeId"`
	ProductType        string                `json:"productType"`
	ProductModel       string                `json:"productModel"`
	ProductImage       string                `json:"productImage"`
	ProductDescription string                `json:"productDescription"`
	Specification      []SpecificationString `json:"specification"`
	SpecString         string                `json:"specString"`
	ProductOverview    string                `json:"productOverview"`
	Scores             []ScoreID             `json:"scores"`
	Apps               []AppID               `json:"apps"`
	CreatedBy          string                `json:"createdBy"`
	CreateAt           string                `json:"createAt"`
	UpdatedBy          string                `json:"updatedBy"`
	UploadAt           string                `json:"uploadAt"`
	IsDeleted          bool                  `json:"isDeleted"`
}

// OurProductInserting function new added image also
func OurProductInserting(data Products) (Products, error) {
	var formattedScores []int
	var resultTool string
	var resultApps string
	var formattedApps []int

	// // Process the Base64 image
	// if data.ProductImage != "" {
	// 	log.Printf("Processing image for product type: %s", data.ProductType)

	// 	// Decode Base64
	// 	imageData, err := base64.StdEncoding.DecodeString(data.ProductImage)
	// 	if err != nil {
	// 		log.Printf("Failed to decode image: %v", err)
	// 		return data, errors.New("invalid image data")
	// 	}

	// 	// Define the directory where the image will be saved
	// 	uploadDir := "uploads"

	// 	// Ensure the directory exists, create it if not
	// 	err = os.MkdirAll(uploadDir, 0755) // 0755 is a standard permission for directories
	// 	if err != nil {
	// 		log.Printf("Failed to create directory: %v", err)
	// 		return data, errors.New("unable to create directory")
	// 	}

	// 	// Save the image to the directory
	// 	timestamp := time.Now().Format("20060102150405")
	// 	imageFileName := fmt.Sprintf("%s/%s_%s.jpg", uploadDir, data.ProductType, timestamp)
	// 	err = ioutil.WriteFile(imageFileName, imageData, 0644)
	// 	if err != nil {
	// 		log.Printf("Failed to save image: %v", err)
	// 		return data, errors.New("unable to save image")
	// 	}

	// 	// Save the image path back to the product struct
	// 	data.ProductImage = imageFileName
	// 	log.Printf("Image saved to: %s", imageFileName)
	// }

	// Process Scores
	for _, tools := range data.Scores {
		toolNameId := 0

		// Check if ToolName exists in the toolname table
		query := `SELECT Tid FROM toolname WHERE ToolName = ?`
		err := DB.QueryRow(query, tools.ToolName).Scan(&toolNameId)
		if err != nil {
			if err == sql.ErrNoRows {
				// Insert new ToolName
				insertQuery := `INSERT INTO toolname (ToolName) VALUES (?)`
				result, err := DB.Exec(insertQuery, tools.ToolName)
				if err != nil {
					log.Printf("Error inserting ToolName: %v", err)
					continue
				}
				LastID, err := result.LastInsertId()
				if err != nil {
					log.Printf("Error fetching LastInsertId for ToolName: %v", err)
					continue
				}
				toolNameId = int(LastID)
			} else {
				log.Printf("Error fetching ToolName: %v", err)
				continue
			}
		}
		// Check if ToolScore exists
		query1 := `SELECT Tid FROM toolscore WHERE Tid = ? AND toolscore = ?`
		err1 := DB.QueryRow(query1, toolNameId, tools.ToolScore).Scan(&toolNameId)
		if err1 != nil {
			if err1 == sql.ErrNoRows {
				// Insert new ToolScore
				insertQuery := `INSERT INTO toolscore (Tid, toolscore) VALUES (?, ?)`
				res, err := DB.Exec(insertQuery, toolNameId, tools.ToolScore)
				if err != nil {
					log.Printf("Error inserting ToolScore: %v", err)
					continue
				}
				LastID, err := res.LastInsertId()
				if err != nil {
					log.Printf("Error fetching LastInsertId for ToolScore: %v", err)
					continue
				}
				toolNameId = int(LastID)
				formattedScores = append(formattedScores, int(toolNameId))
			} else {
				log.Printf("Error checking ToolScore existence: %v", err1)
			}
		} else {
			formattedScores = append(formattedScores, int(toolNameId))
		}
	}

	// Format the Scores result
	// Change in Adding Commas
	for i, res := range formattedScores {
		if i > 0 {
			resultTool += ","
		}
		resultTool += fmt.Sprintf("%d", res)
	}
	fmt.Printf("Formatted Scores: %s\n", resultTool)

	// Process Apps
	for _, appsData := range data.Apps {
		newID := 0

		// Check if the AppType exists in AppCategory
		query := `SELECT Aid FROM appcategory WHERE AppType = ?`
		err := DB.QueryRow(query, appsData.AppType).Scan(&newID)
		if err != nil {
			if err == sql.ErrNoRows {
				// Insert new AppType
				queryInsert := `INSERT INTO appcategory (AppType, Category) VALUES(?, ?)`
				res, err := DB.Exec(queryInsert, appsData.AppType, appsData.Category)
				if err != nil {
					log.Printf("Error inserting AppType: %v", err)
					return data, nil
				}
				lastInsertID, err := res.LastInsertId()
				if err != nil {
					log.Printf("Error fetching last insert ID: %v", err)
					return data, nil
				}
				newID = int(lastInsertID)
			} else {
				log.Printf("Error querying AppType: %v", err)
				return data, nil
			}
		}
		// Check if AppName exists
		queryCheck := `SELECT id FROM appname WHERE Aid = ? AND AppName = ?`
		err = DB.QueryRow(queryCheck, newID, appsData.AppName).Scan(&newID)
		if err != nil {
			if err == sql.ErrNoRows {
				// Insert new AppName
				queryAppName := `INSERT INTO appname (Aid, AppName) VALUES(?, ?)`
				res, err := DB.Exec(queryAppName, newID, appsData.AppName)
				if err != nil {
					log.Printf("Error inserting AppName: %v", err)
					return data, nil
				}
				lastAppID, err := res.LastInsertId()
				if err != nil {
					log.Printf("Error fetching last insert ID: %v", err)
					return data, nil
				}
				newID = int(lastAppID)
				formattedApps = append(formattedApps, int(newID))
			}
		} else {
			formattedApps = append(formattedApps, newID)
		}
	}
	// Format Apps result
	for i, res := range formattedApps {
		if i > 0 { //it is an index value
			resultApps += ","
		}
		resultApps += fmt.Sprintf("%d", res)
	}
	// our product page
	var productId int
	query := `SELECT Pid FROM ourproduct1 WHERE ProductType =?`
	err1 := DB.QueryRow(query, data.ProductType).Scan(&productId)
	if err1 != nil {
		if err1 == sql.ErrNoRows {
			log.Printf("Inserting new Product type %v", data.ProductType)
			query := `INSERT INTO ourproduct1 (ProductType, ProductDescription,ProductImage) VALUES (?, ?, ?)`
			res, err := DB.Exec(query, data.ProductType, data.ProductDescription, data.ProductImage)
			if err != nil {
				log.Print("Error inserting product:", err)
			}
			LastId, err := res.LastInsertId()
			if err != nil {
				log.Print("Error fetching the last product ID:", err)
			}
			productId = int(LastId)
		} else {
			log.Print("Error executing query to find product:", err1)
		}
	} else {
		log.Printf("ProductType %v already exists", data.ProductType)
	}
	query1 := `SELECT Pid, ProductModel FROM ourproduct2 WHERE Pid =? AND ProductModel =?`
	err2 := DB.QueryRow(query1, productId, data.ProductModel).Scan(&productId, &data.ProductModel)
	if err2 != nil {
		if err2 == sql.ErrNoRows {
			log.Printf("Inserting new Product model %v", data.ProductModel)
			query := `INSERT INTO ourproduct2 (Pid, ProductModel,  specification, ProductOverview, score_id, app_id) VALUES(?, ?, ?, ?, ?, ?)`
			_, err := DB.Exec(query, productId, data.ProductModel, data.Specification, data.ProductOverview, resultTool, resultApps)
			if err != nil {
				log.Print("Error inserting product model:", err)
				return data, nil
			}
		} else {
			log.Print("Error executing query to find product model:", err2)
		}
	} else {
		log.Printf("Product Type ID %v already exists for Product model %v, skipping insertion", productId, data.ProductModel)
	}
	return data, nil
}

// Sending AllProduct Lists
type AllProduct struct {
	Id                 int    `json:"id"`
	ProductType        string `json:"productType"`
	ProductImage       string `json:"productImage"`
	ProductDescription string `json:"productDescription"`
}

func ProductDetails(ProductID int) ([]AllProduct, error) {
	var products []AllProduct // Define as a slice of Products
	if ProductID == 0 {
		Query := `SELECT Pid, ProductType, ProductDescription, ProductImage FROM ourproduct1`
		res, err := DB.Query(Query)
		if err != nil {
			log.Print("Error executing Query Selecting coloumn value", err)
			return nil, err
		}
		for res.Next() {
			var product AllProduct
			err := res.Scan(&product.Id, &product.ProductType, &product.ProductDescription, &product.ProductImage)
			if err != nil {
				log.Print("Error scanning result", err)
				return nil, err
			}
			products = append(products, product)
		}
		return products, nil
	} else {
		Query := `SELECT Pid, ProductType, ProductDescription, ProductImage FROM ourproduct1 WHERE Pid = ?`
		res, err := DB.Query(Query, ProductID)
		if err != nil {
			log.Print("Error executing Query Selecting coloumn value", err)
			return nil, err
		}
		for res.Next() {
			var product AllProduct
			err := res.Scan(&product.Id, &product.ProductType, &product.ProductDescription, &product.ProductImage)
			if err != nil {
				log.Print("Error scanning result", err)
				return nil, err
			}
			products = append(products, product)
		}
	}
	return products, nil
}

// Getting Specifications
func SpecConverting(SpecString string) ([]SpecificationString, error) {
	var result []SpecificationString
	// fmt.Printf("SpecString: %s\n", SpecString)
	if SpecString != "" {
		appIDs := strings.Split(SpecString, ",")
		// fmt.Printf("appIDs: %v\n", appIDs)
		for _, idStr := range appIDs {
			// fmt.Print("Entering the loop")
			idStr = strings.TrimSpace(idStr)
			if idStr == "" {
				continue
			}
			id, err := strconv.Atoi(idStr)
			if err != nil {
				log.Printf("Error converting String ID '%s' to integer: %v", idStr, err)
				continue
			}
			query := `SELECT c2.id, c1.id, c2.specName, c2.socketNumber, c1.ComponentName, c1.ComponentType FROM componentspec c2
			JOIN components c1
			ON c2.cid = c1.id
			WHERE c2.id = ?`
			var data SpecificationString
			err1 := DB.QueryRow(query, id).Scan(&data.ID, &data.ComponentTypeID, &data.SpecName, &data.SocketNumber, &data.ComponentName, &data.ComponentType)
			if err1 != nil {
				log.Println("Error occur in fetching the specification", err1)
				// return result, err1
				continue
			}
			// fmt.Printf("result: %v\n", result)
			result = append(result, data)
		}
		// fmt.Printf("result: %v\n", result)
	}
	// fmt.Print("Out of the loop")
	return result, nil
}

// Final Output Demo testing code
func GetDetailedOurProducts(productID, ProductTypeID int) (Products, []Products, error) {
	var product Products
	var products []Products
	var scoreIDsString string
	var appIDsString string
	var err error
	var rows *sql.Rows
	var specJSON string
	var TypeID int

	// Logging for debugging
	fmt.Printf("ProductID %v\n", productID)
	fmt.Printf("productTypeId %v\n", ProductTypeID)

	// Initialize a list to keep track of displayed product IDs (should be preserved across multiple searches)
	var displayedProductIDs []int

	// Add the current productID if it's not zero
	if productID != 0 {
		displayedProductIDs = append(displayedProductIDs, productID)
	}

	// Fetch product details based on ProductID
	if productID == 0 {
		// Query when productID is 0 (fetching product by ProductTypeID)
		DetailedProductQuery := `SELECT ourproduct2.id, ourproduct1.ProductType, ourproduct2.ProductModel, ourproduct2.ProductImage, 
				ourproduct1.ProductDescription, ourproduct2.Specification, ourproduct2.ProductOverview, 
				ourproduct2.score_id, ourproduct2.app_id, ourproduct2.Pid
				FROM ourproduct2
				JOIN ourproduct1 ON ourproduct1.Pid = ourproduct2.Pid
				WHERE ourproduct2.Pid =? 
				ORDER BY id DESC LIMIT 1`
		err = DB.QueryRow(DetailedProductQuery, ProductTypeID).Scan(&product.ProductId, &product.ProductType, &product.ProductModel,
			&product.ProductImage, &product.ProductDescription, &specJSON, &product.ProductOverview, &scoreIDsString, &appIDsString, &TypeID)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("No product found with ProductTypeID %d", ProductTypeID)
				return product, products, err
			}
			log.Printf("Error executing query for product with ProductTypeID %v: %v", ProductTypeID, err)
			return product, products, err
		}
	} else {
		// Query when productID is provided
		DetailedProduct := `SELECT ourproduct2.id, ourproduct1.ProductType, ourproduct2.ProductModel, ourproduct2.ProductImage, 
				ourproduct1.ProductDescription, ourproduct2.Specification, ourproduct2.ProductOverview, 
				ourproduct2.score_id, ourproduct2.app_id, ourproduct2.Pid
				FROM ourproduct2
				JOIN ourproduct1 ON ourproduct1.Pid = ourproduct2.Pid	
				WHERE ourproduct2.id =?`
		err = DB.QueryRow(DetailedProduct, productID).Scan(&product.ProductId, &product.ProductType, &product.ProductModel,
			&product.ProductImage, &product.ProductDescription, &specJSON, &product.ProductOverview, &scoreIDsString, &appIDsString, &TypeID)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("No product found with ID %d", productID)
				return product, products, err
			}
			log.Printf("Error executing query for product with ID %v: %v", productID, err)
			return product, products, err
		}
	}
	fmt.Printf("SpecJSON: %s\n", specJSON)
	res, err := SpecConverting(specJSON)
	if err != nil {
		fmt.Println("Error in converting the Specification:", err)
	}
	product.Specification = res

	// Getting the Score And App Details for the main product
	scores, apps, err := AppsStringsDetails(scoreIDsString, appIDsString)
	if err != nil {
		log.Print("Error Fetching an App or Score:", err)
		return product, products, err
	}
	product.Scores = scores
	product.Apps = apps

	// Similar Products
	var similarProductsQuery string
	if ProductTypeID != 0 {
		similarProductsQuery = `SELECT p2.id, p2.ProductModel, p2.specification, p2.ProductImage 
				FROM ourproduct1 p1 	
				JOIN ourproduct2 p2 ON p1.Pid = p2.Pid 
				WHERE p2.Pid = ? AND p2.id != ? ORDER BY id DESC LIMIT 3`
		rows, err = DB.Query(similarProductsQuery, ProductTypeID, product.ProductId)
	} else {
		similarProductsQuery = `SELECT p2.id, p2.ProductModel, p2.specification, p2.ProductImage 
				FROM ourproduct1 p1 
				JOIN ourproduct2 p2 ON p1.Pid = p2.Pid 
				WHERE p2.Pid = ? AND p2.id != ? ORDER BY id DESC LIMIT 3`
		rows, err = DB.Query(similarProductsQuery, TypeID, product.ProductId)
	}

	if err != nil {
		log.Printf("Error fetching similar products: %v", err)
		return product, products, err
	}
	defer rows.Close()

	// Process similar products
	for rows.Next() {
		var spec []byte
		var similarProduct Products
		err := rows.Scan(&similarProduct.ProductId, &similarProduct.ProductModel, &spec, &similarProduct.ProductImage)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return product, products, err
		}
		//Fetching The Specification Details
		res, err := SpecConverting(string(spec))
		if err != nil {
			fmt.Println("Error in converting the Specification:", err)
		}
		similarProduct.Specification = res
		similarProduct.ProductTypeID = ProductTypeID
		// Add the similar product to the products list
		products = append(products, similarProduct)
		displayedProductIDs = append(displayedProductIDs, similarProduct.ProductId) // Add to displayed IDs
	}
	if len(products) == 0 {
		// Prepare a list of product IDs to exclude from random products
		excludeQuery := ""
		fmt.Print("Entering the Random products section")
		for i, id := range displayedProductIDs {
			if i == 0 {
				excludeQuery = fmt.Sprintf("%d", id)
			} else {
				excludeQuery = fmt.Sprintf("%s, %d", excludeQuery, id)
			}
		}

		// Fetch random products excluding already displayed product
		if len(products) == 0 {
			RandomProductQuery := ""
			if len(displayedProductIDs) > 0 {
				// Create a comma-separated list of IDs
				excludeQuery := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(displayedProductIDs)), ","), "[]")
				RandomProductQuery = fmt.Sprintf(`
			SELECT 
				ourproduct2.id, ourproduct1.ProductType, ourproduct2.ProductModel, ourproduct2.ProductImage, 
				ourproduct1.ProductDescription, ourproduct2.Specification, ourproduct2.ProductOverview, 
				ourproduct2.score_id, ourproduct2.app_id, ourproduct2.Pid
			FROM 
				ourproduct2
			JOIN 
				ourproduct1 ON ourproduct1.Pid = ourproduct2.Pid
			WHERE 
				ourproduct2.id NOT IN (%s)	

			ORDER BY 
				RAND() LIMIT 3;`, excludeQuery)
				fmt.Print("Random Product If part")
			} else {
				// No IDs to exclude, fetch all random products
				RandomProductQuery = `
			SELECT 
				ourproduct2.id, ourproduct1.ProductType, ourproduct2.ProductModel, ourproduct2.ProductImage, 
				ourproduct1.ProductDescription, ourproduct2.Specification, ourproduct2.ProductOverview, 
				ourproduct2.score_id, ourproduct2.app_id, ourproduct2.Pid
			FROM 
				ourproduct2
			JOIN 
				ourproduct1 ON ourproduct1.Pid = ourproduct2.Pid
			ORDER BY 
				RAND() LIMIT 3;`
			}
			fmt.Print("Random Product Else part")
			rows, err = DB.Query(RandomProductQuery)
			if err != nil {
				log.Printf("Error fetching random products: %v", err)
				return product, products, err
			}
			defer rows.Close()

			// Process random products
			for rows.Next() {
				var Random Products
				var spec []byte
				err := rows.Scan(&Random.ProductId, &Random.ProductType, &Random.ProductModel, &Random.ProductImage, &Random.ProductDescription, &spec, &Random.ProductOverview, &scoreIDsString, &appIDsString, &TypeID)
				if err != nil {
					log.Printf("Error scanning row for random products: %v", err)
					return product, products, err
				}
				fmt.Print("Random Product Scanning")
				// Fetching The Specification Details
				res, err := SpecConverting(string(spec))
				if err != nil {
					fmt.Println("Error in converting the Specification:", err)
				}
				Random.Specification = res
				//Getting the Score And App Details
				scores, apps, err := AppsStringsDetails(scoreIDsString, appIDsString)	
				if err != nil {
					fmt.Println("Error in getting the Score and App Details: or No Apps and ", err)
					return product, products, err
				}
				Random.Scores = scores
				Random.Apps = apps
				products = append(products, Random)
			}
		}
	}
	return product, products, nil
}

// fetching the Specifications Details
func AppsStringsDetails(scoreIDsString, appIDsString string) ([]ScoreID, []AppID, error) {
	scores := []ScoreID{}
	apps := []AppID{}
	// Process Score IDs
	if scoreIDsString != "" {
		scoreIDs := strings.Split(scoreIDsString, ",")
		for _, idStr := range scoreIDs {
			idStr = strings.TrimSpace(idStr)
			if idStr == "" {
				continue
			}

			id, err := strconv.Atoi(idStr)
			if err != nil {
				log.Printf("Error converting score ID '%s' to integer: %v", idStr, err)
				continue
			}
			// Fetch toolscore details
			queryScore := `SELECT ToolName, ToolScore FROM toolscore
								INNER JOIN toolname ON toolscore.Tid = toolname.Tid WHERE toolscore.id = ?`
			row := DB.QueryRow(queryScore, id)
			var score ScoreID
			err = row.Scan(&score.ToolName, &score.ToolScore)
			if err == sql.ErrNoRows {
				log.Printf("No toolscore found for ID %d", id)
				continue
			} else if err != nil {
				log.Printf("Error fetching toolscore for ID %d: %v", id, err)
				continue
			}
			scores = append(scores, score)
		}
	} else {
		log.Println("No score IDs found for the product")
	}

	// Process App IDs
	if appIDsString != "" {
		appIDs := strings.Split(appIDsString, ",")
		for _, idStr := range appIDs {
			idStr = strings.TrimSpace(idStr)
			if idStr == "" {
				continue
			}
			id, err := strconv.Atoi(idStr)
			if err != nil {
				log.Printf("Error converting app ID '%s' to integer: %v", idStr, err)
				continue
			}
			// Fetch app details
			queryApp := `SELECT AppType, AppName, Category FROM appcategory
							INNER JOIN appname ON appcategory.Aid = appname.Aid WHERE appname.id = ?`
			row := DB.QueryRow(queryApp, id)
			var app AppID
			err = row.Scan(&app.AppType, &app.AppName, &app.Category)
			if err == sql.ErrNoRows {
				log.Printf("No app details found for ID %d", id)
				continue
			} else if err != nil {
				log.Printf("Error fetching app details for ID %d: %v", id, err)
				continue
			}
			apps = append(apps, app)
		}
	} else {
		log.Println("No app IDs found for the product")
	}
	return scores, apps, nil
}
