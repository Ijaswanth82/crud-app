package main

import (
	"databaseconnection/controller"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", controller.ServeHome).Methods("GET")
	r.HandleFunc("/getallrecords", controller.GetAllRecords).Methods("GET")
	r.HandleFunc("/addrecord", controller.AddARecord).Methods("POST")               //validation string value instead of int
	r.HandleFunc("/deleterecord/{name}", controller.DeleteByName).Methods("DELETE") // No data available response
	r.HandleFunc("/getrecord/{name}", controller.GetByName).Methods("GET")
	r.HandleFunc("/updaterecord/{name}", controller.UpdateByName).Methods("PUT")
	log.Fatal(http.ListenAndServe("0.0.0.0:8888",r))
}

//data access layers
//connection pooling
