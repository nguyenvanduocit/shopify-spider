package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/ThreeDotsLabs/watermill/message"
	sitemap "github.com/oxffaa/gopher-parse-sitemap"
	"net/http"
)

func getDocument(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
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

func GetAppClientID(url string) (string, error) {

	doc, err := getDocument(url)
	if err != nil {
		return "", err
	}

	clientId := doc.Find("input[name='id']").AttrOr("value", "")
	if clientId == "" {
		return "", fmt.Errorf("failed to find client id")
	}

	return clientId, nil
}

func WalkSitemap(publisher message.Publisher) {
	// parse sitemap
	if err := sitemap.ParseFromSite("https://apps.shopify.com/sitemap.xml", func(entry sitemap.Entry) error {
		urlType := GetUrlType(entry.GetLocation())

		switch urlType {
		case "apps":
			payload, _ := json.Marshal(SitemapEntry{
				Location:           entry.GetLocation(),
				ParsedLastModified: entry.GetLastModified(),
			})
			return publisher.Publish(EvtSiteMapAppFound, message.NewMessage(entry.GetLocation(), payload))
		case "collections":

		}
		return nil
	}); err != nil {
	}
}
