package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/alitto/pond"
	"github.com/blevesearch/bleve/v2"
	sitemap "github.com/oxffaa/gopher-parse-sitemap"
	"log"
	"net/http"
	"os"
	"shopifyspider/pkg"
	"time"
)

func main() {

	go func() {
		// healthz
		http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	}()

	// init index
	dataPath := "data"
	index, err := bleve.Open(dataPath)
	if err == bleve.ErrorIndexPathDoesNotExist {

		appsMapping := bleve.NewDocumentMapping()
		appsMapping.AddFieldMappingsAt("lastUpdated", bleve.NewDateTimeFieldMapping())

		indexMapping := bleve.NewIndexMapping()
		indexMapping.AddDocumentMapping("apps", appsMapping)

		index, err = bleve.New(dataPath, indexMapping)
		if err != nil {
			log.Fatalf("Failed to create index: %v\n", err)
		}
	}
	defer index.Close()

	// index pool
	pool := pond.New(1, 100)
	defer pool.StopAndWait()

	sitemapStats := make(map[string]int)

	// parse sitemap
	if err := sitemap.ParseFromSite("https://apps.shopify.com/sitemap.xml", func(entry sitemap.Entry) error {
		urlType := pkg.GetUrlType(entry.GetLocation())

		if _, ok := sitemapStats[urlType]; !ok {
			sitemapStats[urlType] = 1
		} else {
			sitemapStats[urlType]++
		}

		switch urlType {
		case "apps":
			pool.Submit(func() {

				query := bleve.NewMatchQuery(entry.GetLocation())
				query.SetField("url")

				searchQuery := bleve.NewSearchRequest(query)
				searchQuery.Fields = []string{"name", "lastUpdated"}

				searchResult, err := index.Search(searchQuery)
				if err != nil {
					panic(err)
				}

				if searchResult.Total > 0 {
					searchedLastUpdated, _ := time.Parse(time.RFC3339, searchResult.Hits[0].Fields["lastUpdated"].(string))
					if searchedLastUpdated.Equal(*entry.GetLastModified()) {
						fmt.Printf("App: %v (skipped)\n", entry.GetLocation())
						return
					}
				}

				app, err := ParseApp(entry.GetLocation())
				if err != nil {
					panic(err)
				}

				app.LastUpdated = entry.GetLastModified()
				if err := index.Index(app.ClientId, app); err != nil {
					panic(err)
				}
				fmt.Printf("App: %v\n", entry.GetLocation())
				time.Sleep(3 * time.Second)
			})
		case "collections":
			pool.Submit(func() {
				//fmt.Printf("Collections: %v\n", urlType)
			})
		}
		return nil
	}); err != nil {
		log.Fatalf("Failed to parse sitemap: %v\n", err)
	}

	if err := index.Index("sitemapStats", sitemapStats); err != nil {
		log.Fatalf("Failed to index sitemap stats: %v\n", err)
	}
}

func getDocument(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func ParseApp(url string) (*pkg.App, error) {
	app := &pkg.App{
		ClientId:                 "",
		LastUpdated:              nil,
		Name:                     "",
		Url:                      url,
		IconUrl:                  "",
		PrimaryCategory:          "",
		Languages:                nil,
		Introduction:             "",
		Details:                  "",
		Features:                 "",
		DemoStoreUrl:             "",
		FeatureMediaUrl:          "",
		Screenshots:              nil,
		Integrations:             nil,
		SupportEmail:             "",
		SupportPhone:             "",
		PrivacyPolicyUrl:         "",
		DeveloperWebsiteUrl:      "",
		FaqUrl:                   "",
		ChangelogUrl:             "",
		SupportPortalUrl:         "",
		TutorialUrl:              "",
		AdditionalAppDocumentUrl: "",
		Pricing:                  nil,
		Subtitle:                 "",
	}

	doc, err := getDocument(url)
	if err != nil {
		return nil, err
	}

	app.ClientId = doc.Find("input[name='id']").AttrOr("value", "")
	if app.ClientId == "" {
		return nil, fmt.Errorf("failed to find client id")
	}

	return app, nil
}
