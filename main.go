package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Simple programm to query the github API and search all repositories
//

const SEARCH_URL = "https://api.github.com/search/repositories"

type Response struct {
	TotalCount        int  `json:"total_count"`
	IncompleteResults bool `json:"incomplete_results"`
	Items             []*SearchResultItem
}

type SearchResultItem struct {
	Id    int
	Name  string
	Url   string
	Stars int `json:"stargazers_count"`
}

func RunQuery(searchTerm string) (*Response, error) {
	q := url.QueryEscape(searchTerm)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", SEARCH_URL, nil)
	searchParam := req.URL.Query()
	searchParam.Set("q", q)
	searchParam.Set("sort", "stars")
	searchParam.Set("order", "desc")

	req.URL.RawQuery = searchParam.Encode()

	fmt.Println(req.URL.Path)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to make request, got code: " + resp.Status)
	}

	var result Response

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Please provide a search term")
		return
	}
	searchTerm := strings.Join(os.Args[1:], " ")

	fmt.Printf("Searching for repositories: \n\n" + searchTerm)

	resp, err := RunQuery(searchTerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Search result: Found %d repositories\n", resp.TotalCount)
	for _, item := range resp.Items {
		fmt.Printf("%s \n[%d] %s  \n\n", item.Name, item.Stars, item.Url)
	}
}
