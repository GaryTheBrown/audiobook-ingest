package apis

import (
	"audiobook-ingest/data"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

var asinValidatorRegex = regexp.MustCompile(`^[A-Z0-9]{10}$`)

type AudibleCatalogProduct struct {
	Asin string `json:"asin"`
}

type AudibleCatalogResponse struct {
	Products []AudibleCatalogProduct `json:"products"`
}

var audibleRegionTLD = map[string]string{
	"us": ".com", "ca": ".ca", "uk": ".co.uk", "au": ".com.au",
	"fr": ".fr", "de": ".de", "jp": ".co.jp", "it": ".it",
	"in": ".in", "es": ".es",
}

func AudibleSearch(titleKeyword, authorKeyword, region string) ([]data.MetadataInterface, error) {
	region = strings.ToLower(strings.TrimSpace(region))
	tld, ok := audibleRegionTLD[region]
	if !ok {
		tld = ".com"
	}

	titleClean := strings.TrimSpace(titleKeyword)
	authorClean := strings.TrimSpace(authorKeyword)

	var targetAsins []string

	queryObj := url.Values{}
	queryObj.Set("num_results", "10")
	queryObj.Set("products_sort_by", "Relevance")
	queryObj.Set("title", titleClean)
	if authorClean != "" {
		queryObj.Set("author", authorClean)
	}

	catalogURL := fmt.Sprintf("https://api.audible%s/1.0/catalog/products?%s", tld, queryObj.Encode())

	req, err := http.NewRequest("GET", catalogURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "AudiobookIngestStation/2.0.0")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("audible catalog search connection failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var catData AudibleCatalogResponse
		if err := json.NewDecoder(resp.Body).Decode(&catData); err == nil {
			for _, prod := range catData.Products {
				if prod.Asin != "" {
					targetAsins = append(targetAsins, prod.Asin)
				}
			}
		}
	}

	if len(targetAsins) == 0 {
		return []data.MetadataInterface{}, nil
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	finalMatches := make([]data.MetadataInterface, 0, len(targetAsins))

	for _, asin := range targetAsins {
		wg.Add(1)
		go func(targetAsin string) {
			defer wg.Done()
			match, err := FetchBookDetailsByID(targetAsin)
			if err != nil {
				return
			}
			mu.Lock()
			finalMatches = append(finalMatches, match)
			mu.Unlock()
		}(asin)
	}

	wg.Wait()
	return finalMatches, nil
}
