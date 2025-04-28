package models

import (
	"database/sql"
	"log"
)

type Components struct {
	ID            int    `json:"id"`
	ComponentType string `json:"componentType"`
	ComponentName string `json:"componentName"`
	SpecName      string `json:"specName"`
	SocketNumber  string `json:"socketNumber"`
	Platform      string `json:"platform"`
}

// NEW Updated Code For Inserting Not Repeating page
func ComponentInsert(data Components) (Components, error) {
	// Query to check if the component already exists
	query := `SELECT id FROM component WHERE ComponentType = ? AND ComponentName = ?`
	var existingID int
	err := DB.QueryRow(query, data.ComponentType, data.ComponentName).Scan(&existingID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Component does not exist, insert it
			componentQuery := `INSERT INTO component (ComponentType, ComponentName) VALUES (?, ?)`
			res, insertErr := DB.Exec(componentQuery, data.ComponentType, data.ComponentName)
			if insertErr != nil {
				log.Print("Error executing component insertion query: ", insertErr)
				return data, insertErr
			}

			// Retrieve the last inserted ID
			lastID, idErr := res.LastInsertId()
			if idErr != nil {
				log.Print("Error getting last insert ID: ", idErr)
				return data, idErr
			}
			// Assign the last inserted ID to the Components struct
			data.ID = int(lastID)
		} else {
			// An unexpected error occurred
			log.Print("Error executing select query: ", err)
			return data, err
		}
	} else {
		log.Print("Already Exsists The Component Spec")
		data.ID = existingID
	}

	// Query to insert into the  table
	var existingComponentName string	
	var existingSpecName string
	comQuery := `SELECT p.ComponentName, c.specName FROM componentdata p JOIN componentspec c ON p.id = c.cid WHERE P.componentName=? AND c.specName =?`
	err1 := DB.QueryRow(comQuery, data.ComponentName, data.SpecName).Scan(&existingComponentName, &existingSpecName)
	log.Printf("data %s,another data %s", existingComponentName, existingSpecName)
	if err1 != nil {
		if err1 == sql.ErrNoRows {
			// Combination does not exist, proceed with insertion
			specQuery := `
				INSERT INTO  (cid, specName, socketNumber, Platform)
				VALUES (?, ?, ?, ?)`

			_, specErr := DB.Exec(specQuery, data.ID, data.SpecName, data.SocketNumber, data.Platform)
			if specErr != nil {
				log.Printf("Error executing  insertion query: %v", specErr)
				return data, specErr
			}

			log.Println("New specification successfully inserted.")
			return data, nil
		}

		// Unexpected error during query execution
		log.Printf("Error querying component-spec combination: %v", err)
		return data, err
	}
	// Combination already exists
	log.Println("This combination already exists. Try a different specification.")
	return data, nil
}
