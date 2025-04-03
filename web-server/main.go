package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Response struct {
	Message string `json:"message"`
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf("%s %s - %v\n", r.Method, r.URL.Path, duration)
	})
}

func GetHello(w http.ResponseWriter, req *http.Request) {
	msg := Response{Message: "hello"}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

func GetStatus(w http.ResponseWriter, req *http.Request) {
	response := map[string]string{
		"status": "ok",
		"uptime": fmt.Sprintf("%v", time.Since(startTime)),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

var startTime time.Time

const port = 8090

func main() {
	startTime = time.Now()

	mux := http.NewServeMux()

	logRouter := loggingMiddleware(mux)

	mux.HandleFunc("/hello", GetHello)
	mux.HandleFunc("/status", GetStatus)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: logRouter,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		fmt.Printf("Server starting on http://localhost:%d\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	<-stop
	fmt.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown with error: %v", err)
	}

	fmt.Println("Successfully shut down.")
}
