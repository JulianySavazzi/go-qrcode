package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"time"
	"encoding/json"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("The environment PORT must be defined")
	}

	api := http.NewServeMux()

	api.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		responseApiJson(
		    w,
        	http.StatusOk,
        	map[string]string{
        	    "status": "ok",
        	    "message": "Success"
        	},
            nil,
        )
	})
    api.HandleFunc("POST /qr-code", func(w http.ResponseWriter, r *http.Request) {

    })

	mux := http.NewServeMux()
	mux.Handle("/v1/", http.StripPrefix("/v1", api))

	server := &http.Server{
		Addr:              net.JoinHostPort("", port),
		Handler:           mux,
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	log.Printf("The sever HTTP starts on port %s", port)
	log.Fatal(server.ListenAndServe())
}

func responseApiJson(
    writer http.ResponseWriter,
    statusCode int,
    data any,
    additionalHeaders map[string]string,
) {
    /** serialize data to json **/
    responseBody, err := json.Marshal(data)
    if err != nil {
        statusCode = http.StatusInternalServerError
        /** catch the json Marshal error **/
        responseBody = []byte(`{error: "Response error"}`)
    }

    /** set additional headers **/
    for name, value := range additionalHeaders {
        writer.Header().Set(name, value)
    }

    /** added response headers **/
    writer.Header().Set("Content-Type", "application/json")
    writer.WriteHeader(statusCode)

    /** use _ for ignore bytes written returned by writer.Write, and writerErr for error handler **/
    _, writerErr := writer.Write(responseBody)
    if writerErr != nil {
        log.Printf("Error to send response: %v", writerErr)
    }
}