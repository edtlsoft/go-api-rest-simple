package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type task struct {
	Id      int    `json:Id`
	Name    string `json:Name`
	Content string `json:Content`
}

type allTask []task

var tasks = allTask{
	{
		Id:      1,
		Name:    "Task One",
		Content: "Some content",
	},
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask task
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "Insert a valid task")
	}

	json.Unmarshal(reqBody, &newTask)

	newTask.Id = len(tasks) + 1

	tasks = append(tasks, newTask)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	taskId, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Fprintf(w, "The field id is required")
	} else {
		task, _, error := findById(taskId)

		if error {
			fmt.Fprintf(w, "The task with Id %v has not been found", taskId)
		} else {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
		}
	}

}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	taskId, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Fprintf(w, "Invalid Id")
	}

	_, taskIndice, error := findById(taskId)

	if error {
		fmt.Fprintf(w, "The task with Id %v has not been found", taskId)
	} else {
		tasks = append(tasks[:taskIndice], tasks[taskIndice+1:]...)
		fmt.Fprintf(w, "The task with Id %v has been remove successfuly", taskId)
	}
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	taskId, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Fprintf(w, "Invalid Id")
	}

	var updateTask task

	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "Please enter valid data")
	}

	json.Unmarshal(reqBody, &updateTask)

	_, taskIndice, error := findById(taskId)

	if error {
		fmt.Fprintf(w, "The task with Id %v has not been found", taskId)
	} else {
		tasks = append(tasks[:taskIndice], tasks[taskIndice+1:]...)
		updateTask.Id = taskId
		tasks = append(tasks, updateTask)
		fmt.Fprintf(w, "The task with Id %v has been updated successfuly", taskId)
	}
}

func findById(taskId int) (task, int, bool) {
	var task task
	var taskIndice int
	var error bool = true

	for i, t := range tasks {
		if t.Id == taskId {
			task = t
			taskIndice = i
			error = false
		}
	}

	return task, taskIndice, error
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my REST API 1.0")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", indexRoute)

	router.HandleFunc("/tasks", getTasks).Methods("GET")

	router.HandleFunc("/tasks", createTask).Methods("POST")

	router.HandleFunc("/tasks/{id}", getTask).Methods("GET")

	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")

	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")

	log.Fatal(http.ListenAndServe(":3000", router))
}
