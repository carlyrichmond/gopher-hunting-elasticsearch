package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"search"
)

func main() {
	client, err := search.GetElasticsearchClient()
	if err != nil {
		log.Println(err)
		os.Exit(3)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		message := "Hi Gophers!"
		fmt.Fprintf(w, message)
	})

	http.HandleFunc("/gophers", func(w http.ResponseWriter, r *http.Request) {
		if gophers, err := search.KeywordSearch(client, "What do Gophers eat?"); err == nil {
			fmt.Println("Gopher count", len(gophers))
			json.NewEncoder(w).Encode(gophers)
		} else {
			log.Println(err)
			internalServerError(w)
		}
	})

	http.HandleFunc("/vector-gophers", func(w http.ResponseWriter, r *http.Request) {
		if gophers, err := search.VectorSearch(client, "What do Gophers eat?"); err == nil {
			fmt.Println("Gopher count", len(gophers))
			json.NewEncoder(w).Encode(gophers)
		} else {
			log.Println(err)
			internalServerError(w)
		}
	})

	http.HandleFunc("/embedding-vector-gophers", func(w http.ResponseWriter, r *http.Request) {
		if gophers, err := search.VectorSearch(client, "What do Gophers eat?"); err == nil {
			fmt.Println("Gopher count", len(gophers))
			json.NewEncoder(w).Encode(gophers)
		} else {
			log.Println(err)
			internalServerError(w)
		}
	})

	http.HandleFunc("/filtered-gophers", func(w http.ResponseWriter, r *http.Request) {
		if gophers, err := search.VectorSearchWithFilter(client, "What do Gophers eat?"); err == nil {
			fmt.Println("Gopher count", len(gophers))
			json.NewEncoder(w).Encode(gophers)
		} else {
			log.Println(err)
			internalServerError(w)
		}
	})

	http.HandleFunc("/generated-vector-gophers", func(w http.ResponseWriter, r *http.Request) {
		if gophers, err := search.VectorSearchWithGeneratedQueryVector(client, "What do Gophers eat?"); err == nil {
			fmt.Println("Gopher count", len(gophers))
			json.NewEncoder(w).Encode(gophers)
		} else {
			log.Println(err)
			internalServerError(w)
		}
	})

	http.HandleFunc("/hybrid-gophers", func(w http.ResponseWriter, r *http.Request) {
		if gophers, err := search.HybridSearchWithBoost(client, "What do Gophers eat?"); err == nil {
			fmt.Println("Gopher count", len(gophers))
			json.NewEncoder(w).Encode(gophers)
		} else {
			log.Println(err)
			internalServerError(w)
		}
	})

	http.HandleFunc("/rrf-gophers", func(w http.ResponseWriter, r *http.Request) {
		if gophers, err := search.HybridSearchWithRRF(client, "What do Gophers eat?"); err == nil {
			fmt.Println("Gopher count", len(gophers))
			json.NewEncoder(w).Encode(gophers)
		} else {
			log.Println(err)
			internalServerError(w)
		}
	})

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	const port = ":8080"

	log.Println("starting server on port 80")
	http.ListenAndServe(port, nil)
}

func internalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Internal server error"))
}
