package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Students struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Semester int     `json:"semester"`
	GPA      float64 `json:"gpa"`
}

var stds []Students

func getStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stds)
}

func getStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range stds {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Students{})
}

func delStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range stds {
		if item.ID == params["id"] {
			stds = append(stds[:index], stds[index+1:]...)
			break
		}
	}
}

func createStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var std Students
	_ = json.NewDecoder(r.Body).Decode(&std)
	std.ID = strconv.Itoa(rand.Intn(100000000))
	stds = append(stds, std)
	json.NewEncoder(w).Encode(std)
}

func updateStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range stds {
		if item.ID == params["id"] {
			stds = append(stds[:index], stds[index+1:]...)
			var std Students
			_ = json.NewDecoder(r.Body).Decode(&std)
			std.ID = params["id"]
			stds = append(stds, std)
			json.NewEncoder(w).Encode(std)
			return
		}
	}
}

func main() {
	r := mux.NewRouter()

	stds = append(stds, Students{ID: "1", Name: "Muhammad", Semester: 8, GPA: 3.53})
	stds = append(stds, Students{ID: "2", Name: "Soomro", Semester: 7, GPA: 3.20})

	r.HandleFunc("/students", getStudents).Methods("GET")
	r.HandleFunc("/students/{id}", getStudent).Methods("GET")
	r.HandleFunc("/students/{id}", delStudent).Methods("DELETE")
	r.HandleFunc("/students", createStudent).Methods("POST")
	r.HandleFunc("/students/{id}", updateStudent).Methods("PUT")

	fmt.Println("Server listining at port 8080")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
