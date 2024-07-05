package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/sboriskin287/assets-scheduler/core"
	"github.com/sboriskin287/assets-scheduler/mongo"
	"log"
	"net/http"
	"time"
)

func main() {
	client := mongo.CreateMongoClient()
	defer client.Disconnect(context.TODO())
	ps := core.NewPeriodService(client)
	router := mux.NewRouter()
	router.HandleFunc("/", ps.Index).Methods("GET")
	router.HandleFunc("/period", ps.CreatePeriod).Methods("POST")
	router.HandleFunc("/period/{id}/details", ps.GetPeriodDetails).Methods("GET")
	router.HandleFunc("/period/{id}/details", ps.CreatePeriodDetails).Methods("POST")
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
