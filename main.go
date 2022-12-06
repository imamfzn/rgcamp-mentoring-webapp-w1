package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
)

type GramediaBook struct {
	Name  string  `json:"name"`
	Image string  `json:"thumbnail"`
	Price float64 `json:"basePrice"`
}

type Book struct {
	Name  string  `json:"name"`
	Image string  `json:"image"`
	Price float64 `json:"price"`
}

type Data struct {
	Books []Book
}

var GramediaAPIBaseUrl = "https://www.gramedia.com/api/algolia/search/product?q=belajar"

func GetGramediaBooks() []Book {
	resp, _ := http.Get(GramediaAPIBaseUrl)

	var gbooks []GramediaBook
	err := json.NewDecoder(resp.Body).Decode(&gbooks)
	if err != nil {
		panic(err)
	}

	var books []Book
	for _, b := range gbooks {
		book := Book{b.Name, b.Image, b.Price}
		books = append(books, book)
	}

	return books
}

func main() {
	tmpl, err := template.ParseFiles("./books.html")
	if err != nil {
		panic(err)
	}

	books := GetGramediaBooks()

	http.HandleFunc("/api/books", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		err := json.NewEncoder(w).Encode(books)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Server error :("))
		}
	})

	http.HandleFunc("/page/books", func(w http.ResponseWriter, r *http.Request) {
		data := Data{Books: books}
		err = tmpl.Execute(w, data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Server error :("))
		}
	})

	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
