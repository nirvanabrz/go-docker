package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"encoding/json"

	"github.com/gorilla/mux"
	"gopkg.in/natefinch/lumberjack.v2"
)

type AppResponse struct {
	HostName string `json:"HostName"`
	Version  string `json:"Version"`
	Language string `json:"Language"`
}

var appResponse AppResponse

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Requisição recebida com sucesso!")
	json.NewEncoder(w).Encode(appResponse)

}

func main() {
	language := "Golang"
	version := "V3"
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	appResponse = AppResponse{HostName: hostname, Version: version, Language: language}

	// Create Server and Route Handlers
	r := mux.NewRouter()

	r.HandleFunc("/", handler)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Configure Logging
	LOG_FILE_LOCATION := os.Getenv("LOG_FILE_LOCATION")
	if LOG_FILE_LOCATION != "" {
		log.SetOutput(&lumberjack.Logger{
			Filename:   LOG_FILE_LOCATION,
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28,   //days
			Compress:   true, // disabled by default
		})
	}

	// Start Server
	go func() {
		log.Println("Iniciando o servidor - Golang " + version + " em " + hostname)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful Shutdown
	waitForShutdown(srv)
}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Parando o serviço.")
	os.Exit(0)
}
