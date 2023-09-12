package main

import (
	"edSample/flow"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func endPoint(response http.ResponseWriter, request *http.Request) { //json decoding handler
	var source []byte
	json.NewDecoder(request.Body).Decode(&source) //there should be an error handler in case of the request arrives with more items than the expected
	result := flow.Process(source)
	if result == nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	response.Header().Set("Content=Type", "application/json")
	json.NewEncoder(response).Encode(result)
}

func main() { //Main

	http.HandleFunc("/keyData", endPoint)

	fmt.Println("RUNNING ON PORT 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
