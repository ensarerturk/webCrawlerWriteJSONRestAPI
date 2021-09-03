package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

type Fact struct {
	ID          string    `json:"id"`
	Description string `json:"description"`
}

var Articles []Fact

func homePage(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w,"Welcome the Home Page")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	json.NewEncoder(w).Encode(Articles)
}

func handleRequest() {

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/",homePage)
	myRouter.HandleFunc("/all",returnAllArticles)

	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func allFuc() {

	collector := colly.NewCollector(
		colly.AllowedDomains("imdb.com", "www.imdb.com"),
	)
	collector.OnHTML(".lister-item-header a ", func(element *colly.HTMLElement) {
		factId:= element.Attr("href")

		factDesc := element.Text

		fact := Fact{
			ID:          factId,
			Description: factDesc,
		}
		Articles = append(Articles, fact)
	})

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL.String())
	})

	collector.Visit("https://www.imdb.com/search/title/?groups=top_100&sort=user_rating,desc")

	writeJSON(Articles)
}

func writeJSON(data []Fact) {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("Unable to create json file")
		return
	}

	_ = ioutil.WriteFile("index.json", file, 0644)
}

func main() {

	allFuc()
	handleRequest()
}