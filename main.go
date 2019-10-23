package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

type impeachedResponse struct {
	State bool
}

var impeached = false

func isImpeached(w http.ResponseWriter, req *http.Request) {
	s := impeachedResponse{
		State: impeached,
	}
	json.NewEncoder(w).Encode(s)
}

func fetchState() {
	// parameters
	sleep_duration := 10 * time.Second
	rss_nyt_url := "https://rss.nytimes.com/services/xml/rss/nyt/Politics.xml"
	// Some date that is definitely far in the past, so all rss feed
	// data is considered at the beginning.
	last_checked := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	p := gofeed.NewParser()
	for true {
		feed, _ := p.ParseURL(rss_nyt_url)
		// Dont consider items in the feed that were already checked.
		for _, item := range feed.Items {
			if !last_checked.Before(*item.PublishedParsed) {
				break
			}
			t := strings.ToLower(item.Title)
			if strings.Contains(t, "impeached") &&
				strings.Contains(t, "trump") &&
				!strings.Contains(t, "?") {
				fmt.Println(t)
				fmt.Println("Impeached!!!")
				impeached = true
				break
			}
			fmt.Printf("%s\n\n\n\n", t)
		}
		last_checked = *feed.Items[0].PublishedParsed
		fmt.Printf("It's %s and Trump is still not impeached...\n", time.Now().Truncate(time.Second))
		time.Sleep(sleep_duration)
		fmt.Println("Let's check again...")
	}
}

func main() {
	go fetchState()
	http.HandleFunc("/impeached", isImpeached)
	http.ListenAndServe(":8090", nil)
}
