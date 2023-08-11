package main

import (
	"fmt"
	bleve "github.com/blevesearch/bleve/v2"
)

func main() {
	dataPath := "data.bleve"
	index, err := bleve.Open(dataPath)
	if err == bleve.ErrorIndexPathDoesNotExist {
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(dataPath, mapping)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	defer index.Close()

	text, _ := index.Stats().MarshalJSON()
	fmt.Println(string(text))

	data := struct {
		Name string
	}{
		Name: "demo",
	}

	// index some data
	index.Index("id", data)

	query := bleve.NewMatchQuery("text")
	search := bleve.NewSearchRequest(query)
	searchResults, err := index.Search(search)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(searchResults)
}
