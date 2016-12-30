package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	apiVersion = "v1"
	baseURL    = "https://api.unsplash.com"
)

// Client manages all searches.  All searches are performed from an instance of a client.
// It is the top level object used to perform a search and should be created through the
// createClient API.
type Client struct {
	Options *AuthOptions
	Client  *http.Client
}

// AuthOptions provides any keys and settings needed to use the API
type AuthOptions struct {
	ApplicationID string
}

// Search lets you query for photos
func (c *Client) Search(query string, page int, perPage int) (PhotoResponse, error) {
	const path = "/search/photos/"

	req, err := http.NewRequest("GET", fmt.Sprintf(baseURL+"%s", path), nil)
	req.Header.Set("Accept-Version", apiVersion)

	// Headers: `Authorization: Client-ID YOUR_APPLICATION_ID`
	// As query param: `https://api.unsplash.com/photos/?client_id=YOUR_APPLICATION_ID`
	req.Header.Set("Authorization", "Client-ID "+c.Options.ApplicationID)

	req.URL.Query().Set("query", query)
	req.URL.Query().Set("page", string(page))
	req.URL.Query().Set("per_page", string(perPage))

	resp, err := c.Client.Do(req)
	if err != nil {
		return PhotoResponse{}, err
	}

	defer resp.Body.Close()
	var r PhotoResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return PhotoResponse{}, err
	}

	return r, nil
}

// PhotoResponse is thing
type PhotoResponse struct {
	Total      int     `json:"total"`
	TotalPages int     `json:"total_pages"`
	Results    []Photo `json:"results"`
}

// Photo is thing
type Photo struct {
	ID     string   `json:"id"`
	Images ImageSet `json:"urls"`
	// "created_at": "2016-09-09T08:59:35-04:00",
	// "categories": [],
}

// ImageSet is thing
type ImageSet struct {
	Raw     string `json:"raw"`
	Full    string `json:"full"`
	Regular string `json:"regular"`
	Small   string `json:"small"`
	Thumb   string `json:"thumb"`
}

// New will create a new yelp search client.  All search operations should go through this API.
func New(options *AuthOptions, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		Options: options,
		Client:  httpClient,
	}
}
