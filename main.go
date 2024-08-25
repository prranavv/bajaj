package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"unicode"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {
	srv := http.Server{
		Addr:    ":8000",
		Handler: routes(),
	}
	fmt.Println("Server Running on 8000")
	err := srv.ListenAndServe()
	log.Println(err)
}

type Response struct {
	OperationCode int `json:"operation_code"`
}

type PostResponse struct {
	Status           bool     `json:"status"`
	UserID           string   `json:"user_id"`
	CollegeEmailId   string   `json:"email"`
	RollNumber       string   `json:"roll_number"`
	Numbers          []string `json:"numbers"`
	Aplhabets        []string `json:"alphabets"`
	HighestAlphabets string   `json:"highest_lowercase_alphabet"`
}

type PostRequest struct {
	Data []string `json:"data"`
}

func routes() *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	mux.Get("/bfhl", func(w http.ResponseWriter, r *http.Request) {
		response := Response{
			OperationCode: 1,
		}

		// Convert the struct to JSON
		jsonData, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set Content-Type as application/json
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Write the JSON response
		w.Write(jsonData)
	})

	mux.Post("/bfhl", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var reqData PostRequest
		// Unmarshal the JSON into the struct
		err = json.Unmarshal(body, &reqData)
		if err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}
		var Numbers []string
		var alphabets []string
		var loweralpha []string
		for i := 0; i < len(reqData.Data); i++ {
			if isAlphabet(reqData.Data[i]) {
				alphabets = append(alphabets, reqData.Data[i])
				if isAllLowercase(reqData.Data[i]) {
					loweralpha = append(loweralpha, reqData.Data[i])
				}
			} else if isInteger(reqData.Data[i]) {
				Numbers = append(Numbers, reqData.Data[i])
			}
		}
		post := PostResponse{
			Status:           true,
			UserID:           "pranav_vk_24062003",
			CollegeEmailId:   "pranav.veerendra2021@vitstudent.ac.in",
			RollNumber:       "21BCE3927",
			Numbers:          Numbers,
			Aplhabets:        alphabets,
			HighestAlphabets: loweralpha[len(loweralpha)-1],
		}

		jsonData, err := json.Marshal(post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set Content-Type as application/json
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Write the JSON response
		w.Write(jsonData)
	})
	return mux
}

func isAlphabet(s string) bool {
	for _, char := range s {
		if !unicode.IsLetter(char) {
			return false
		}
	}
	return true
}

func isInteger(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func isAllLowercase(s string) bool {
	for _, char := range s {
		if unicode.IsLetter(char) && !unicode.IsLower(char) {
			return false
		}
	}
	return true
}
