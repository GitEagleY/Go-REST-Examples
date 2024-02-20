package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Item is an abstract struct with fields ID, Name, and Data.
type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Data string `json:"data"`
}

func main() {

	//MARSHAL UNMARSHAL
	// Create an instance of the Item struct
	item := Item{
		ID:   1,
		Name: "Example Item",
		Data: "Initial data",
	}

	// Convert struct to JSON
	jsonData, err := json.Marshal(item)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	// Print the JSON string
	fmt.Println("JSON String:")
	fmt.Println(string(jsonData))

	// Parse JSON back into a struct
	var newItem Item
	err = json.Unmarshal(jsonData, &newItem)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	// Modify struct data
	newItem.Data = "Modified data"

	// Convert modified struct back to JSON
	modifiedJSON, err := json.Marshal(newItem)
	if err != nil {
		log.Fatalf("Error marshaling modified JSON: %v", err)
	}

	// Print the modified JSON string
	fmt.Println("\nModified JSON String:")
	fmt.Println(string(modifiedJSON))

	//ENCODE DECODE

	// Encode struct to JSON and write to a file
	file, err := os.Create("output.json")
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(item)
	if err != nil {
		log.Fatalf("Error encoding JSON: %v", err)
	}
	fmt.Println("Struct encoded and written to 'output.json'")

	// Decode JSON from the file back into a struct
	var newItem Item
	file, err = os.Open("output.json")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&newItem)
	if err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}

	// Print the decoded struct
	fmt.Println("\nDecoded Struct:")
	fmt.Printf("ID: %d\nName: %s\nData: %s\n", newItem.ID, newItem.Name, newItem.Data)
}
