package main

import (
	"fmt"
	"log"
	"net/http"

	"search"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		message := "Hi Gophers!"
		fmt.Fprintf(w, message)
	})

	http.HandleFunc("/gophers", func(w http.ResponseWriter, r *http.Request) {
		var gophers = search.KeywordSearch("What do Gophers eat?")
		fmt.Println("Gopher count", len(gophers))

		fmt.Fprintf(w, "%+v", gophers)
	})

	http.HandleFunc("/vector-gophers", func(w http.ResponseWriter, r *http.Request) {
		var gophers = search.VectorSearch("What do Gophers eat?")
		fmt.Println("Gopher count", len(gophers))

		fmt.Fprintf(w, "%+v", gophers)
	})

	http.HandleFunc("/embedding-vector-gophers", func(w http.ResponseWriter, r *http.Request) {
		var gophers = search.VectorSearch("What do Gophers eat?")
		fmt.Println("Gopher count", len(gophers))

		fmt.Fprintf(w, "%+v", gophers)
	})

	http.HandleFunc("/filtered-gophers", func(w http.ResponseWriter, r *http.Request) {
		var gophers = search.VectorSearchWithFilter("What do Gophers eat?")
		fmt.Println("Gopher count", len(gophers))

		fmt.Fprintf(w, "%+v", gophers)
	})

	http.HandleFunc("/generated-vector-gophers", func(w http.ResponseWriter, r *http.Request) {
		var gophers = search.VectorSearchWithGeneratedQueryVector("What do Gophers eat?")
		fmt.Println("Gopher count", len(gophers))

		fmt.Fprintf(w, "%+v", gophers)
	})

	http.HandleFunc("/hybrid-gophers", func(w http.ResponseWriter, r *http.Request) {
		var gophers = search.HybridSearchWithBoost("What do Gophers eat?")
		fmt.Println("Gopher count", len(gophers))

		fmt.Fprintf(w, "%+v", gophers)
	})

	http.HandleFunc("/rrf-gophers", func(w http.ResponseWriter, r *http.Request) {
		var gophers = search.HybridSearchWithRRF("What do Gophers eat?")
		fmt.Println("Gopher count", len(gophers))

		fmt.Fprintf(w, "%+v", gophers)
	})

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	const port = ":80"

	log.Println("starting server on port 80")
	http.ListenAndServe(port, nil)
}
