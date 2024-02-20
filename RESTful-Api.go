package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Task represents a task entity.
type Task struct {
	TaskContent string `json:"task_content"`
	NumInList   int    `json:"num_in_list"`
}

var db = make(map[int]Task)

func main() {
	// Set up the server routes
	http.HandleFunc("/tasks/display", displayAllTasks)
	http.HandleFunc("/tasks/addById/", addByID)
	http.HandleFunc("/tasks/update/", updateTask)
	http.HandleFunc("/tasks/delete/", deleteTask)
	http.HandleFunc("/tasks/getTaskById/", getByID)

	// Start the server
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

// GET
func displayAllTasks(w http.ResponseWriter, r *http.Request) {

	tempTasksDB := make([]Task, 0, len(db))

	switch r.Method {
	case http.MethodGet: //if http method that used to acces /tasks/display route GET

		for _, task := range db {
			tempTasksDB = append(tempTasksDB, task) //adding to temp task db tasks from actual db
		}
		respondWithJSON(w, tempTasksDB) //sending json
	default: //if http not GET
		respondWithError(w, "Method not allowed") //if err sending err json
	}
}

// POST
func addByID(w http.ResponseWriter, r *http.Request) {
	var taskToAdd Task // Temporary variable to hold the task data.

	switch r.Method {
	case http.MethodPost: // Handle only HTTP POST requests.
		id := getTaskNumFromRequest(r) // Extract the task ID from the request.

		// Decode the JSON payload from the request body.
		err := json.NewDecoder(r.Body).Decode(&taskToAdd)
		if err != nil {
			respondWithError(w, "Invalid request payload")
			return
		}

		// No error while decoding, so add the task to the database.
		db[id] = taskToAdd

		//responding with succsess json
		respondWithJSON(w, taskToAdd)
	default:
		// Respond with an error for methods other than POST.
		respondWithError(w, "Method not allowed")
	}
}

// PUT
func updateTask(w http.ResponseWriter, r *http.Request) {
	switch r.Method { //checking for http method
	case http.MethodPut: //if PUT
		taskID := getTaskNumFromRequest(r) //getting id of task
		task, exists := db[taskID]         //taking task from db and checking for existance
		if exists {                        //checking if task exists in db
			var updatedTask Task                                                 //temp task for holding updated task data
			if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil { //trying to decode json from body
				respondWithError(w, "Invalid request payload") //return eror json if there an err
				return
			}

			//sucsess
			task.TaskContent = updatedTask.TaskContent //updating TaskContent from just received updatedTask
			db[taskID] = task                          //updating db[taskID] by new task with updated data
			//responding with succsess json
			respondWithJSON(w, task)
		} else {
			respondWithError(w, "Task not found") //if task not found respoond with error json
		}
	default:
		respondWithError(w, "Method not allowed") //if inapropriate method respoond with error json
	}
}

// DELETE
func deleteTask(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		taskID := getTaskNumFromRequest(r)
		_, exists := db[taskID] //checking if task existing
		if exists {             //if true
			delete(db, taskID)                                               //deleting from db
			respondWithJSON(w, map[string]string{"message": "Task deleted"}) //respond with succsess message
		} else {
			respondWithError(w, "Task not found") //respond error json if there no such task
		}
	default:
		respondWithError(w, "Method not allowed") //respond with error if inaproproate method
	}
}

// GET
func getByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method { //checking for method
	case http.MethodGet: //if GET
		taskID := getTaskNumFromRequest(r) //getting id of task from request
		task, exists := db[taskID]         //taking task by id and checking for existance
		if exists {                        //if exists return it
			respondWithJSON(w, task)
		} else { //if not retunr not found
			respondWithError(w, "Task not found")
		}
	default: //return err if not appropriate method
		respondWithError(w, "Method not allowed")
	}
}

// ///
// ///
// ///
// //////HELP FUNCS
// ///
// ///
// ///
// respondWithJSON is a utility function that sends a JSON response with the specified payload.
func respondWithJSON(w http.ResponseWriter, payload interface{}) {

	//SETTING

	w.Header().Set("Content-Type", "application/json") //setting header settings
	response, err := json.Marshal(payload)             //marshalling payload to json responce

	//ERROR CHECKING
	if err != nil { //if err return err
		log.Fatalf("Error marshaling JSON: %v", err)
		respondWithError(w, "Internal Server Error")
		return
	}

	//SENDING
	w.Write(response) //write json respoce to writer
}

// respondWithError is a utility function that sends a JSON response with an error message.
func respondWithError(w http.ResponseWriter, message string) {

	//SETTING
	w.Header().Set("Content-Type", "application/json") //setting header settings

	response, err := json.Marshal(map[string]string{"error": message}) //marshalling error message to json responce

	//ERROR CHECKING
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//writing error
	http.Error(w, string(response), http.StatusMethodNotAllowed)
}

// getTaskNumFromRequest is a utility function to extract the task number from the URL.
func getTaskNumFromRequest(r *http.Request) int {
	// Extract the task number from the URL
	parts := strings.Split(r.URL.Path, "/") //splitting to parts
	if len(parts) < 4 {                     //checking url for validness
		log.Printf("Invalid URL format: %s", r.URL.Path)
		return 0
	}

	id := parts[len(parts)-1]      //getting last part
	idInt, err := strconv.Atoi(id) //converting to int
	if err != nil {                //error checking
		log.Printf("Error converting task number to integer: %v", err)
		return 0
	}
	return idInt //returning id
}
