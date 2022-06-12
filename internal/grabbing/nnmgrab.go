package grabbing

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gocolly/colly"
)

type SearchParameter struct {
	Type  string
	Value string
}

func NewSearchParameter(t string, v string) *SearchParameter {
	return &SearchParameter{
		Type:  t,
		Value: v,
	}
}

const (
	baseNnmURL     = "https://nnmclub.ro/forum/"
	proxy          = "HTTP://proxy-nossl.antizapret.prostovpn.org:29976"
	baseSearchLine = "https://nnmclub.ro/forum/tracker.php?"

	//search codes
	NewMovies         = "f=954&"
	Screener          = "f=217&"
	ForeignCartoons   = "f=1339&"
	ForeignCartoons3d = "f=1339&"
	RussianCartoons   = "f=1332&"
	AllCartoons       = "c=26&"
	ForeignMovie      = "f=227&"
	RussianMovie      = "f=221&"
	Game              = "c=17&"
	AllMovie          = "c=14&"
	Serials           = "c=27&"
	AllFiles          = ""
)

type Film struct {
	addInfo   string
	baseLink  string
	date      string
	magnet    string
	size      string
	tFileLink string
	Title     string
}

//SprintFilm return format string for telegram message with movie info
func SprintFilm(film Film) (string, error) {
	if film.Title == "" {
		return "", errors.New("not found")
	}
	return fmt.Sprintf("Название: %s\nРазмер: %s\nНа торренте с: %s\nДоп. инфо: %s\nМагнет ссылка: %s\nТоррент файл: %s\n",
		film.Title, film.size, film.date, film.addInfo, film.magnet, film.tFileLink), nil
}

//findDuplicate find duplicates in result search
func findDuplicate(arr []Film) []Film {
	var occurred = make(map[string]bool)
	var result []Film
	for e := range arr {
		if !occurred[arr[e].magnet] {
			occurred[arr[e].magnet] = true
			result = append(result, arr[e])
		}
	}
	return result

}

//printToStdout printing result search to stdout
func printToStdout(result []Film) {
	for _, el := range result {
		fmt.Println(el.Title)
		fmt.Println("\t--size\t:", el.size)
		fmt.Println("\t--date\t:", el.date)
		fmt.Println("\t--audio\t:", el.addInfo)
		fmt.Println("\t--magnet\t:", el.magnet)
		fmt.Println("\t--torrent-file\t:", el.tFileLink)
		fmt.Println("***********************")
	}
}

//Find func running search movie
func Find(searchString string, searchType string) ([]Film, error) {
	var ch = make(chan Film)
	films, err := nnmAPI(searchString, searchType, ch)
	if err != nil {
		return nil, err
	}
	result := findDuplicate(films)
	//fmt.Println(len(result))
	//printToStdout(result)
	return result, nil
}

//nnmAPI get links on movie from channel and grab each movie info
func nnmAPI(searchString string, searchType string, ch chan Film) ([]Film, error) {
	var filmMap []Film

	if err := searchForm(searchString, searchType, ch); err != nil {
		return filmMap, err
	}

	c := colly.NewCollector()
	//err := c.SetProxy(proxy)
	//if err != nil {
	//	return filmMap, err
	//}
	c.MaxDepth = 1

	err := c.Limit(&colly.LimitRule{
		DomainGlob:  "nnmclub.ro/*",
		Parallelism: 2,

		// Set a delay between requests to these domains
		Delay: 1 * time.Second,
		// Add a random delay
		RandomDelay: 1 * time.Second,
	})
	if err != nil {
		return filmMap, err
	}
	for chEl := range ch {

		c.OnHTML(".wrap", func(h *colly.HTMLElement) {
			var film Film
			film.date = chEl.date
			film.addInfo = chEl.addInfo

			h.ForEach(".btTbl", func(_ int, h *colly.HTMLElement) {
				h.ForEach("span[title]", func(_ int, h *colly.HTMLElement) {
					if strings.Contains(h.Attr("title"), "Размер") {
						film.size = h.Text
					}
				})
				h.ForEach("a", func(_ int, h *colly.HTMLElement) {
					if h.Attr("title") == "Примагнититься" {
						film.magnet = h.Attr("href")
					}
					if h.Attr("href")[:8] == "download" {
						film.tFileLink = baseNnmURL + h.Attr("href")
					}
				})
			})

			h.ForEach(".maintitle", func(_ int, h *colly.HTMLElement) {
				film.Title = h.Text
			})

			filmMap = append(filmMap, film)
		})
		err := c.Visit(chEl.baseLink)
		if err != nil {
			return filmMap, err
		}
	}
	return filmMap, nil
}

//searchForm post request on search form and send link on each result in channel
func searchForm(searchString string, searchType string, ch chan Film) error {
	if utf8.RuneCountInString(searchString) < 3 {
		return errors.New("too few symbols for search")
	}

	cFile := colly.NewCollector()
	//err := cFile.SetProxy(proxy)
	//if err != nil {
	//	return err
	//}
	err := cFile.Limit(&colly.LimitRule{
		Parallelism: 2,
		// Filter domains affected by this rule
		DomainGlob: "nnmclub.ro/*",
		// Set a delay between requests to these domains
		Delay: 1 * time.Second,
		// Add a random delay
		RandomDelay: 1 * time.Second,
	})
	if err != nil {
		return err
	}

	go func(channel chan Film) {
		var tempFilm Film

		cFile.OnHTML(".tablesorter", func(h *colly.HTMLElement) {
			h.ForEach("tbody", func(_ int, h *colly.HTMLElement) {
				h.ForEach("tr", func(_ int, h *colly.HTMLElement) {
					h.ForEach("td[title]", func(_ int, h *colly.HTMLElement) {
						if h.Attr("title") == "Добавлено" {
							tempFilm.date = h.Text[11:21]
						}
					})
					tempFilm.addInfo = h.ChildText(".opened")
					tempFilm.baseLink = baseNnmURL + h.ChildAttr(".topictitle", "href")

					channel <- tempFilm
				})
			})
		})

		cFile.OnError(func(_ *colly.Response, e error) {
			fmt.Println(e.Error())
		})

		searchString = "nm=" + searchString
		err := cFile.Visit(baseSearchLine + searchType + searchString)
		if err != nil {
			fmt.Println(err)
		}
		close(ch)
	}(ch)

	return nil
}
