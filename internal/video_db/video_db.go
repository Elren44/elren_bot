package video_db

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Elren44/elren_bot/configs"
	"github.com/gocolly/colly"
)

type Movie struct {
	Result bool   `json:"result"`
	Data   []Data `json:"data"`
}

type Data struct {
	Id          uint   `json:"id"`
	Title       string `json:"ru_title"`
	KinopoiskID string `json:"kinopoisk_id"`
	ImdbID      string `json:"imdb_id"`
	Year        string `json:"year"`
	Link        string `json:"iframe_src"`
	IFrame      string `json:"iframe"`
	Poster      string
}

//Videocdn get video from db
func Videocdn(title string, cfg *configs.Config) (*Movie, error) {
	var movie Movie
	resp, err := http.Get("https://videocdn.tv/api/movies?api_token=" + cfg.VideocdnToken + "&query=" + title)
	if err != nil {
		return nil, fmt.Errorf("get response error: %w", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&movie)
	if err != nil {
		return nil, fmt.Errorf("failed to decode in struct: %w", err)
	}

	return &movie, nil
}

//GrabPoster find link on film poster
func GrabPoster(id string) (string, error) {
	var link string
	c := colly.NewCollector()

	c.OnHTML(".ipc-poster--baseAlt", func(h *colly.HTMLElement) {
		h.ForEach(".ipc-media > img", func(_ int, h *colly.HTMLElement) {
			link = h.Attr("src")
		})
	})

	err := c.Visit("https://www.imdb.com/title/" + id + "/")
	if err != nil {
		return link, err
	}
	return link, nil
}
