package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
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

var store []TodoNote 

func main() {
	http.HandleFunc("/getTodo",handleGetTodo)
	http.HandleFunc("/submitTodo",handleSubmitTodo)
	http.HandleFunc("/removeTodo",handleRemoveTodo)

	store = append(store, 
		TodoNote {
			ID: generateID(),
			Title: "An example Todo", 
			Details: "Go read a book or something...",
		}, 
	)

	addFileServer()
	log.Printf("Starting server on localhost:%d",PORT);
	if err := http.ListenAndServe(fmt.Sprintf(":%d",PORT),nil); err != nil {
		log.Fatalf("Failed to initialize server: %s",err.Error())
	}
}

///
func handleGetTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w,"Incorrect method, expecting GET for this endpoint", http.StatusMethodNotAllowed)
		return
	}
	sendTodoList(w)
}

/// 
func handleSubmitTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w,fmt.Sprintf("Incorrect method, expecting POST for this endpoint, received %v",r.Method), http.StatusMethodNotAllowed)
		return
	}
	var data struct {
		TodoTitle   string `json:"todoTitle"`
		TodoDetails string `json:"todoDetails"`
	};
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w,"Invalid request body from client", http.StatusBadRequest)
		return
	}
	trimmedTitle := strings.Trim(data.TodoTitle,"\n\t")
	trimmedDetails := strings.Trim(data.TodoDetails,"\n\t")
	newTodo := TodoNote {ID: generateID(),Title: trimmedTitle, Details: trimmedDetails}
	store = append(store, newTodo)
	sendNewTodo(w,newTodo)
}

/// Expects an index of a specific element from store to remove,then removes it and sends back html of list
func handleRemoveTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w,"Incorrect method, expecting DELETE for this endpoint", http.StatusMethodNotAllowed)
		return
	}

	receivedID := r.Header.Get("ElementID")
	intID,err := strconv.ParseInt(receivedID,10,32)
	if err != nil {
		http.Error(w,"Unable to parse provided element id", http.StatusBadRequest)
		return
	}


	for i, elem := range store {
		if(elem.ID == int32(intID)) {
			store = slices.Delete(store,i,i+1)
		}
	}
	log.Printf("Removing at : %d",intID)
	log.Printf("Store: %v",store)
	w.Write([]byte("Delete Succeeded"))
}

///
func sendTodoList(w http.ResponseWriter) {
	w.Header().Set("Content-Type","application/json")
	encoded, err := json.Marshal(store)
	if err != nil {
		http.Error(w,"Error encoding the server store into json format", http.StatusInternalServerError)
		return
	}
	w.Write(encoded)	
}

///
func sendNewTodo(w http.ResponseWriter, newTodo TodoNote) {
	w.Header().Set("Content-Type","application/json")
	err := json.NewEncoder(w).Encode(newTodo)
	if err != nil {
		http.Error(w,"Error encoding new todo into json format", http.StatusInternalServerError)
		return	
	}
}

///
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

///
func generateID() int32 {
	id += 1
	return id - 1
}

