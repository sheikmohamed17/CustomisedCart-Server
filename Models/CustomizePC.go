package models

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Component struct {
	ID            int    `json:"id"`
	SpecName      string `json:"specName"`
	ComponentName string `json:"componentName"`
}

// sending the list of Processor
func GettingProcessor(ComponentType string) ([]Component, error) {
	var data []Component
	fmt.Println("Component Type", ComponentType)
	Query := `SELECT 
    c2.id, 
    c2.SpecName,
	c1.ComponentName
FROM 
    components c1 
JOIN 
    componentspec c2 
ON 
    c1.id = c2.cid
WHERE 	
    c1.ComponentType = ?;`
	res, err := DB.Query(Query, ComponentType)
	if err != nil {
		log.Print("error Executing on Getting Query", err)
		return data, nil
	}
	for res.Next() {
		var data1 Component
		err := res.Scan(&data1.ID, &data1.SpecName, &data1.ComponentName)
		if err != nil {
			log.Print("Error an Executing Scanning the data", err)
			return data, nil
		}
		data = append(data, data1)

	}
	return data, nil
}

// ComponentList Fetching
type ComponentList struct {
	ComponentID   int    `json:"id"`
	CompSpecID    int    `json:"CompSpecID"`
	SpecName      string `json:"specName"`
	SocketNumber  string `json:"socketNumber"`
	ComponentType string `json:"componentType"`
	// SupportedRam  string `json:"supportedRam"`
	// GraphicsCard  string `json:"graphicsCard"` //PCI Pins
}

// For Using an Customize PC Component Selecting list
// getting Data from Front end it search already Present data
func GettingComponentList(ComponentID, SpecId int, SocketNumber string) ([]ComponentList, error) {
	var Result []ComponentList
	var rows *sql.Rows
	var err error

	fmt.Printf("Id: %v, SocketNumber: %v\n", ComponentID, SocketNumber)
	fmt.Printf("ComponentSpecId:%v\n", SpecId)
	// Query based on whether SocketNumber is provided or not
	if SocketNumber != "" {
		query := `SELECT c2.id, c2.specName, c2.socketNumber, c1.ComponentType FROM componentspec c2 JOIN components c1 
		ON c2.cid = c1.id
		WHERE cid = ? AND socketNumber = ?`
		rows, err = DB.Query(query, ComponentID, SocketNumber)
		if err != nil {
			log.Printf("Error in First Query: %v\n", err)
			return nil, err
		}
	} else {
		query := `SELECT c2.id, c2.specName, c2.socketNumber, c1.ComponentType FROM componentspec c2 JOIN components c1 
		ON c2.cid = c1.id	
		WHERE cid = ?`
		rows, err = DB.Query(query, ComponentID)
		if err != nil {
			log.Printf("Error in Second Query: %v\n", err)
			return nil, err
		}
	}
	// Ensure rows are closed after processing
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	// Iterate through the rows and populate Result
	for rows.Next() {
		var data1 ComponentList
		if err := rows.Scan(&data1.CompSpecID, &data1.SpecName, &data1.SocketNumber, &data1.ComponentType); err != nil {
			log.Printf("Error Scanning Row: %v\n", err)
			return nil, err
		}
		if data1.CompSpecID != SpecId {
			Result = append(Result, data1)
		}
	}
	// Check for any errors encountered during iteration
	if err := rows.Err(); err != nil {
		log.Printf("Error Iterating Rows: %v\n", err)
		return nil, err
	}
	return Result, nil
}

// Customize pc Section sending data to front end
func CustomizedData(CustomizedProduct string) (Product Products, err error) {
	var data Products
	fmt.Printf("CustomizedProduct: %v\n", CustomizedProduct)
	if CustomizedProduct != "" {
		query := `SELECT id, specification FROM ourproduct2`
		rows, err := DB.Query(query)
		if err != nil {
			log.Printf("Error fetching product details first Query: %v", err)
			return data, err
		}
		defer rows.Close()
		for rows.Next() {
			var id int
			var spec string
			err := rows.Scan(&id, &spec)
			if err != nil {
				log.Printf("Error scanning row: %v", err)
				continue
			}
			fmt.Printf("Spec: %v\n, Ids%v\n", spec, id)
			if spec == CustomizedProduct {
				fmt.Printf("Spec Entering %v\n", spec)
				var AppList string
				var ScoreList string
				query := `SELECT id, Pid, ProductModel, ProductImage, ProductDescription, specification, app_id, score_id, ProductOverview FROM ourproduct2 WHERE id = ?`
				var SpecString string
				err = DB.QueryRow(query, id).Scan(&data.ProductId, &data.ProductTypeID, &data.ProductModel, &data.ProductImage, &data.ProductDescription, &SpecString, &AppList, &ScoreList, &data.ProductOverview)
				if err != nil {
					if err == sql.ErrNoRows {
						log.Printf("No rows found")
						return data, nil
					}
					log.Printf("Error fetching product details Second Query: %v", err)
					return data, err
				}
				// getting Specificaton details

				fmt.Printf("Spec: %v\n", SpecString)
				res, err := SpecConverting(SpecString)
				if err != nil {
					fmt.Println("Error in converting the Specification:", err)
				}
				data.Specification = res
				// Calling Apps & Scores
				score, app, err1 := AppsStringsFetching(AppList, ScoreList)
				if err1 != nil {
					log.Print("Error Executing Query in Fetch App score Details")
				}
				data.Scores = score
				data.Apps = app
				return data, nil
			} else {
				log.Printf("Product not found")
				continue
			}
		}
		// if rows.Err() != nil {
		// 	log.Printf("No Specification Available: %v", rows.Err())
		// 	return data, rows.Err()
		// }
	}
	return data, nil
}

func AppsStringsFetching(scoreIDsString, appIDsString string) ([]ScoreID, []AppID, error) {
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

// testing code
func CustomizedData1(CustomizedProduct string, ProductId int) (Product Products, err error) {
	var data Products
	fmt.Printf("CustomizedProduct: %v\n", CustomizedProduct)
	fmt.Printf("ProductId %v\n", ProductId)
	if CustomizedProduct != "" {
		query := `SELECT SpecID, SpecString FROM processing_table WHERE SpecID = ?`
		rows, err := DB.Query(query, ProductId)
		if err != nil {
			log.Printf("Error fetching product details first Query: %v", err)
			return data, err
		}
		defer rows.Close()
		for rows.Next() {
			var id int
			var spec string
			err := rows.Scan(&id, &spec)
			if err != nil {
				log.Printf("Error scanning row: %v", err)
				continue
			}
			if spec == CustomizedProduct {
				var AppList string
				var ScoreList string
				query :=
					`SELECT 
					p1.id, 
					p1.Pid, 
					p1.ProductModel, 
					p1.ProductImage, 
					p1.ProductDescription, 
					p1.specification, 
					p2.apps_string, 
					p2.score_string, 
					p1.ProductOverview 
				FROM 
					ourproduct2 p1 
				JOIN 
					processing_table p2 
				ON 
					p1.id = p2.SpecID 
				WHERE 
					p1.id = ?;`
				var SpecString string
				err = DB.QueryRow(query, id).Scan(&data.ProductId, &data.ProductTypeID, &data.ProductModel, &data.ProductImage, &data.ProductDescription, &SpecString, &AppList, &ScoreList, &data.ProductOverview)
				if err != nil {
					if err == sql.ErrNoRows {
						log.Printf("No rows found")
						return data, nil
					}
					log.Printf("Error fetching product details Second Query: %v", err)
					return data, err
				}
				// getting Specificaton details

				fmt.Printf("Spec: %v\n", SpecString)
				res, err := SpecConverting(SpecString)
				if err != nil {
					fmt.Println("Error in converting the Specification:", err)
				}
				data.Specification = res
				// Calling Apps & Scores
				score, app, err1 := AppsStringsFetching(AppList, ScoreList)
				if err1 != nil {
					log.Print("Error Executing Query in Fetch App score Details")
				}
				data.Scores = score
				data.Apps = app
				return data, nil
			} else {
				log.Printf("Product not found")
				continue
			}
		}

	} else {
		log.Printf("Enter Valid Specification")
		return data, fmt.Errorf("enter Valid Specification")
	}
	// returning the Specification fetching Front end
	query := `SELECT id, ProductModel, ProductImage, ProductDescription, ProductOverview FROM ourproduct2 WHERE id = ?`
	err = DB.QueryRow(query, ProductId).Scan(&data.ProductId, &data.ProductModel, &data.ProductImage, &data.ProductDescription, &data.ProductOverview)
	if err != nil {
		log.Printf("Error fetching product details in Scanning: %v", err)
		// return data, err
	}
	res, err := SpecConverting(CustomizedProduct)
	if err != nil {
		fmt.Println("Error in converting the Specification:", err)
	}
	data.Specification = res
	return data, err
}

type ComponentList1 struct {
	ID            int    `json:"id"`
	SpecName      string `json:"specName"`
	SocketNumber  string `json:"socketNumber"`
	ComponentType string `json:"componentType"`
	SupportedRam  string `json:"supportedRam"`
	GraphicsCard  string `json:"graphicsCard"` //PCI Pins
}

// Full Customization Functionality
func FullCustomization(ComponentType, SocketNumber, SupportedRam, GraphicsCard string) ([]ComponentList1, error) {
	var data []ComponentList1
	var rows *sql.Rows
	var err error
	//Using Socket Number fetching the Component Type
	if SocketNumber != "" {
		Query := `SELECT c2.id, c2.specName, c2.socketNumber, c2.SupportedRam, c2.PCIeSlot FROM componentspec c2 JOIN components c1 ON c2.cid=c1.id WHERE ComponentType = ? AND socketNumber=?`
		rows, err = DB.Query(Query, ComponentType, SocketNumber)
		if err != nil {
			log.Print("Error Executing Query", err)
			// return data, err
		}
	} else {
		Query := `SELECT c2.id, c2.specName, c2.socketNumber, c2.SupportedRam, c2.PCIeSlot FROM componentspec c2 JOIN components c1 ON c2.cid=c1.id WHERE ComponentType = ?`
		rows, err = DB.Query(Query, ComponentType)
		if err != nil {
			log.Print("Error Executing Query", err)
		}
	}
	if ComponentType != "" && SupportedRam != "" {
		// Query := `SELECT c2.id, c2.specName, c2.socketNumber, c2.SupportedRam, c2.PCIeSlot FROM componentspec c2 JOIN components c1 ON c2.cid=c1.id WHERE ComponentType =?`
		// rows, err = DB.Query(Query, ComponentType)
		// if err != nil {
		// 	log.Print("Error Executing Query", err)
		// 	// return data, err
		// }
		Query := `SELECT c2.id, c2.specName, c2.socketNumber, c2.SupportedRam, c2.PCIeSlot FROM componentspec c2 JOIN components c1 ON c2.cid=c1.id WHERE c1.ComponentType =? AND c2.SupportedRam=?`
		rows, err = DB.Query(Query, ComponentType, SupportedRam)
		if err != nil {
			log.Print("Error Executing Query", err)
			// return data, err
		}
	} else if ComponentType != "" && GraphicsCard != "" {
		Query := `SELECT c2.id, c2.specName, c2.socketNumber, c2.SupportedRam, c2.PCIeSlot FROM componentspec c2 JOIN components c1 ON c2.cid=c1.id WHERE c1.ComponentType =? AND PCIeSlot=?`
		rows, err = DB.Query(Query, ComponentType, GraphicsCard)
		if err != nil {
			log.Print("Error Executing Query", err)
			// return data, err
		}
	}
	for rows.Next() {
		var data1 ComponentList1
		err := rows.Scan(&data1.ID, &data1.SpecName, &data1.SocketNumber, &data1.SupportedRam, &data1.GraphicsCard)
		if err != nil {
			log.Print("Error Executing Scanning the data", err)
		}
		data = append(data, data1)
	}
	if err = rows.Err(); err != nil {
		log.Print("Error in rows.Err()", err)
		return data, err
	}

	return data, nil
}
