package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/moviedb/api/pkg/server"
)

type Movie struct {
	Title         map[string]string `json:",omitempty"`
	Url           string            `json:",omitempty"`
	ContentRating string            `json:",omitempty"`
	Type          string            `json:",omitempty"`
	Description   map[string]string `json:",omitempty"`
	Genre         string            `json:",omitempty"`
	Image         string            `json:",omitempty"`
	ReleaseDate   int64             `json:",omitempty"`
	Director      []string          `json:",omitempty"`
	Actors        []string          `json:",omitempty"`
	Trailer       []Trailer         `json:",omitempty"`
}

type Trailer struct {
	Name         map[string]string `json:",omitempty"`
	Description  map[string]string `json:",omitempty"`
	Url          string            `json:",omitempty"`
	ThumbnailUrl string            `json:",omitempty"`
}

func SearchHandler(srv server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var page int
		fmt.Sscan(r.URL.Query().Get("p"), &page)
		query := r.URL.Query().Get("q")
		country := r.URL.Query().Get("c")

		url := srv.Config.ApiUrl
		b, err := httpGet(fmt.Sprintf("%s/search/?query=%s&country=%s&page=%d", url, query, country, page))
		if err != nil {
			log.Printf("Failed to http get %s - %v\n", fmt.Sprintf("%s/search/?query=%s&country=%s&page=%d", url, query, country, page), err)
			srv.JSON(w, http.StatusInternalServerError, "failed to http get")
			return
		}

		var movieList []Movie
		if err := json.Unmarshal([]byte(b), &movieList); err != nil {
			log.Printf("Failed to Unmarshall %s", b)
			srv.JSON(w, http.StatusInternalServerError, "failed to unmarshall")
			return
		}

		srv.JSON(w, http.StatusOK, movieList)
	}
}

func httpGet(url string) ([]byte, error) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Content-Type": []string{"application/json"},
		"ApiKey":       []string{"178fb68d-3500-4b1d-96d7-6c0bf549b045"},
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error Url: %s - Status Code: %d", req.URL, resp.StatusCode)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed: %s", resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}
