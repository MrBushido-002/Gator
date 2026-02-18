package main

import(
	"context"
	"net/http"
	"io"
	"encoding/xml"
	"html"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	var data []byte
	var myFeed RSSFeed
	
	// creates a request object that is not sent to the server yet
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")

	// uses to request object to get and store a response from the server
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	// closes the http request once the function returns to prevent memory leaks
	defer res.Body.Close()

	//reads all of the body data from the http response
	data, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	//converts the raw data to fit the RSSFeed struct
	err = xml.Unmarshal(data, &myFeed)
	if err != nil {
		return nil, err
	}

	//loops through the channel struct and cleans up the title and description of unwanted characters with the html.UnescapeString
	myFeed.Channel.Title = html.UnescapeString(myFeed.Channel.Title)
	myFeed.Channel.Description = html.UnescapeString(myFeed.Channel.Description)
	
	for i, item := range myFeed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		myFeed.Channel.Item[i] = item
	}

	return &myFeed, nil
}