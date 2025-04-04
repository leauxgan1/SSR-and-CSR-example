package main

import (
	"fmt"
	"strconv"
	"log"
	"net/http"
	"os"
)

const (
	PORT = 8080
)	
var id int32 = 0

type TodoNote struct {
	ID int32 `json:"todoID"`
	Title string `json:"todoTitle"`
	Details string `json:"todoDetails"`
}

type RemovalRequest struct {
	ID string `json:"todoID"`
}

var store []TodoNote 

func main() {
	http.HandleFunc("/getTodo",handleGetTodo)
	http.HandleFunc("/submitTodo",handleSubmitTodo)
	http.HandleFunc("/removeTodo",handleRemoveTodo)

	store = append(store, 
		TodoNote {
			ID: generateID(),
			Title: "Do the dishes", 
			Details: "The utensils need extra attention",
		}, 
		TodoNote {
			ID: generateID(),
			Title: "Do the laundry", 
			Details: "Make sure to take any money out of pockets",
		},
	)

	addFileServer()
	log.Printf("Starting server on localhost:%d",PORT);
	if err := http.ListenAndServe(fmt.Sprintf(":%d",PORT),nil); err != nil {
		log.Fatalf("Failed to initialize server: %s",err.Error())
	}
}

func handleGetTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w,"Incorrect method, expecting GET for this endpoint", http.StatusMethodNotAllowed)
		return
	}
	sendTodoList(w)
}

func handleSubmitTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w,"Incorrect method, expecting GET for this endpoint", http.StatusMethodNotAllowed)
		return
	}
	submittedTitle := r.FormValue("titleTextArea")
	// Handle invalid Title
	if submittedTitle == "" {
		http.Error(w,"Invalid title for todo", http.StatusBadRequest)
		return
	}
	submittedDetails := r.FormValue("detailTextArea")
	// Handle invalid Details
	if submittedDetails == "" {
		http.Error(w,"Invalid details for todo", http.StatusBadRequest)
		return
	}
	log.Printf("title: %s",submittedTitle)
	log.Printf("details: %s",submittedDetails)
	newTodo := TodoNote {ID: generateID(),Title: submittedTitle, Details: submittedDetails}
	store = append(store, newTodo)
	sendTodoList(w)
}


// Expects an index of a specific element from store to remove,then removes it and sends back html of list
func handleRemoveTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w,"Incorrect method, expecting DELETE for this endpoint", http.StatusMethodNotAllowed)
		return
	}

	// log.Printf("Request body: %s",r.Body)
	// var data RemovalRequest
	// err := json.NewDecoder(r.Body).Decode(&data)
	// if err != nil {
	// 	http.Error(w,"Unable to process removal request, id is malformed",http.StatusBadRequest)
	// 	return
	// }
	queryParams := r.URL.Query()
	id := queryParams.Get("todoID")
	log.Printf("ID: %s",id)
	if id == "" {
		http.Error(w,"Unable to process removal request, id is malformed",http.StatusBadRequest)
		return
	}

	targetID, err := strconv.ParseInt(id,10,16)
	if err != nil {
		http.Error(w,"Unable to process removal index, request failed", http.StatusInternalServerError)
		return
	}	
	log.Printf("Received form value for id: %d",targetID)
	w.WriteHeader(http.StatusOK)
	for i := 0; i < len(store); i++ {
		if store[i].ID == int32(targetID) {
			store = removeAtIndex(store,i)
		}
	}
	log.Println(store)
	sendTodoList(w)
}

func sendTodoList(w http.ResponseWriter) {
	w.Header().Set("Content-Type","text/html")
	for i := 0; i < len(store); i++ {
		raw := `
		<li class='m-2 rounded-lg border-blue-400 border-solid border-2' >
					<div class='flex flex-row items-center justify-center'>
						<div class='flex-grow'>
							<h1 class='m-2 font-bold text-2xl'>%s</h1>
							<p class='text-md mx-4 my-2'>%s</p>
						</div>
						<button 
							class='text-6xl text-center w-14 h-full bg-blue-200 rounded-md border-solid border-2 border-blue-300 p-1 mx-1' 
							hx-delete='/removeTodo'
							hx-vals='{"todoID":"%d"}'
							hx-trigger='click'
							hx-target='#todolist'
							hx-swap='innerHTML'
						>
							X
						</button>
					</div>
				</li>
			`
		currTodo := store[i]
		formatted := fmt.Sprintf(raw,currTodo.Title,currTodo.Details,currTodo.ID)
		// log.Println(formatted)
		fmt.Fprint(w,formatted)//,currTodo.ID)
	}
}

func addFileServer() {
	directoryPath := "./client"

	_, err := os.Stat(directoryPath)
	if os.IsNotExist(err) {
			fmt.Printf("Directory '%s' not found.\n", directoryPath)
			return
	}

	// Create a file server handler to serve the directory's contents
	fileServer := http.FileServer(http.Dir(directoryPath))

	// Create a new HTTP server and handle requests
	http.Handle("/", fileServer)
}

func generateID() int32 {
	id += 1
	return id - 1
}

func removeAtIndex[T comparable](slice []T, index int) []T {
	return append(slice[:index], slice[index+1:]...)
}
