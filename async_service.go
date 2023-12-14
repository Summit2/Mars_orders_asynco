package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Payload struct {
	ID int `json:"id"`
	// Add other fields if needed
}

func main() {
	http.HandleFunc("/deliver/", handlePost)
	http.ListenAndServe(":8080", nil)
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Read the request body
		buf := new(bytes.Buffer)
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Process the request body
		var payload Payload
		err = json.Unmarshal(buf.Bytes(), &payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Access the ID from the payload
		id := payload.ID
		fmt.Println("Received POST request with ID:", id)

		// Trigger the handleDelete function asynchronously
		go handleDelete(w, r, id)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("POST request processed successfully"))
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func handleDelete(w http.ResponseWriter, r *http.Request, id int) {
	// Simulate waiting for 5 seconds
	time.Sleep(5 * time.Second)

	// Perform a DELETE request to another server with the ID
	deleteURL := fmt.Sprintf("http://127.0.0.1:8000/cargo/%d/", id)
	req, err := http.NewRequest(http.MethodDelete, deleteURL, nil)
	if err != nil {
		fmt.Println("Error creating DELETE request:", err)
		return
	}

	client := &http.Client{}
	deleteResp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error performing DELETE request:", err)
		return
	}

	defer deleteResp.Body.Close()

	fmt.Println("DELETE request to another server completed with status code:", deleteResp.Status)
}
