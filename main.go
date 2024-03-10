package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var apiURL string = "https://api.openai.com/v1/images/generations"

var apiKey string = "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"

type Request struct {
	Prompt         string `json:"prompt"`
	N              int    `json:"n"`
	Size           string `json:"size"`
	ResponseFormat string `json:"response_format"`
}

const (
	low    string = "256x256"
	medium string = "512x512"
	high   string = "1024x1024"
)

type Response struct {
	Created int `json:"created"`
	Data    []struct {
		Url string `json:"url"`
	} `json:"data"`
}

func generateImage(w http.ResponseWriter, req *http.Request) {

	text := req.URL.Query().Get("text")

	body := Request{Prompt: text, N: 1, ResponseFormat: "url", Size: low}

	bodyBytes, _ := json.Marshal(body)

	r, _ := http.NewRequest("POST", apiURL, bytes.NewReader(bodyBytes))
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", "Bearer "+apiKey)
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return
	}

	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	output := Response{}

	err = json.Unmarshal(responseBody, &output)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output.Data[0].Url)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/generate", generateImage).Methods("GET")

	port := 8080
	fmt.Printf("Server running on :%d...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port),
		handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTION"}),
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		)(router))
}
