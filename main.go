package main

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
	"encoding/json"
	"io"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("The environment PORT must be defined")
	}

	api := http.NewServeMux()

	api.HandleFunc("GET /health", health)

    api.HandleFunc("POST /qr-code", createQrCode)

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

type QrCodeRequest struct {
    Source string `json:"source"`
    Size *int `json:"size"`
}

/**
 * validates te request body
 * source is a URL from generates QR-QrCode
 * and size is QR-CODE size, by default is 256
 */
func validateQrCodeRequest(request QrCodeRequest) error {
    parsedUrl, err := url.Parse(request.source)
    if err != nil || (parsedUrl.Schema != "http" && != parsedUrl.Schema != "https") || parsedUrl.Hostname() == "" {
        return errors.New("Source must be a valid URL, with http or https")
    }
    if request.Size != nil {
        if *request.Size < 128 || *request.Size > 1024 {
            return errors.New("Size must be a integer between 128 and 1024 pixels")
        }
    }
    return nil
}

func createQrCode(w http.ResponseWriter, r *http.Request) {
    /** bytes limit for body - string source and int size **/
    const maxBodySizeKiB = 8 * 1024

    body, err := io.ReadAll(io.LimitReader(r.Body, maxBodySizeKiB + 1)
    statusCode := http.StatusUnprocessableEntity
    message := "Invalid body, send a valid string source and optional integer size"

    if err != nil {
        message = "Could not read request body"
        statusCode = http.StatusBadRequest
        ResponseApiJson(
            w,
            statusCode,
            map[string]string{"error": message},
            nil
        )
        return
    }

    if len(body) > maxBodySizeKiB {
        message = "The body must not exceed 8 KiB"
        statusCode = http.StatusRequestEntityTooLarge
        ResponseApiJson(
                w,
                statusCode,
                map[string]string{"error": message},
                nil
            )
        return
    }

    var request QrCodeRequest
    err = json.Unmarshal(body, &request)
    if err != nil {
        ResponseApiJson(
            w,
            statusCode,
            map[string]string{"error": message}
            nil,
        )
        return
    }

    err = validateQrCodeRequest(request)
    if err != nil {
        ResponseApiJson(
            w,
            statusCode,
            map[string]string{"error": err.Error()}
            nil,
        )
        return
    }

    size := 256
    if request.Size != nil {
        size = *request.Size
    }

    // TODO generate QR CODE
}

func health(w http.ResponseWriter, r *http.Request) {
    ResponseApiJson(
        w,
	    http.StatusOk,
	    map[string]string{
	        "status": "ok",
	        "message": "Success"
	    },
        nil,
    )
}

func ResponseApiJson(
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