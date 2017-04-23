// Package tvmaze is an HTTP Client for the tvmaze API.
package tvmaze

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	cache "github.com/robfig/go-cache"
)

// Client is a tvmaze client.
type Client struct {
	Debug     bool
	BaseURI   string
	Region    string
	Cache     *cache.Cache
	CacheFile string
	UseCache  bool
	*http.Client
}

// NewClient returns a ready to use Client.
func NewClient(cachefile string) (*Client, error) {
	c := cache.New(time.Minute*60*24*7, time.Minute*60)
	if _, err := os.Stat(cachefile); err == nil {
		err := c.LoadFile(cachefile)
		if err != nil {
			return nil, err
		}
	}

	timeout := time.Duration(180 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}

	return &Client{
		Cache:     c,
		CacheFile: cachefile,
		BaseURI:   "http://api.tvmaze.com",
		Client:    client,
	}, nil
}

// WriteCache writes cache contents to disk.
func (c *Client) WriteCache() error {
	err := c.Cache.SaveFile(c.CacheFile)
	if err != nil {
		return err
	}
	return nil
}

// FindShow searches tvmaze for showname, and returns it as a Show if a match is found.
func (c *Client) FindShow(showname string) (Show, error) {
	candidates, err := c.GetShow(showname)
	if err != nil {
		return Show{}, err
	}
	// Lost Girl is listed as CA country in tvrage, and so my configured region of US
	// will not match.  Lost Girl has only aired in CA though, so you can get a string
	// equality match on the 2nd retry (the old behavior).
	for retry := 0; retry <= 1; retry++ {
		for _, cand := range candidates {
			show := cand.Show
			if c.Region != "" && retry < 1 {
				if c.Region == show.Network.Country.Code && strings.HasPrefix(show.Name, showname) {
					return show, nil
				}
			} else {
				if strings.ToLower(show.Name) == strings.ToLower(showname) {
					return show, nil
				} else if retry > 0 && strings.HasPrefix(show.Name, showname) {
					return show, nil
				}
			}
		}
	}

	return Show{}, errors.New("Failed to match show in tvmaze!")
}

// GetShow queries tvmaze for show, and returns Candidates that may be a match.
func (c *Client) GetShow(show string) ([]Candidate, error) {
	route := "/search/shows"
	uri, err := url.Parse(c.BaseURI + route)
	if err != nil {
		log.Fatal(err)
	}

	query := url.Values{}
	query.Add("q", show)
	uri.RawQuery = query.Encode()

	var candidates []Candidate
	jsondata, err := c.Go(uri)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsondata, &candidates)
	if err != nil {
		return nil, err
	}

	return candidates, nil
}

// GetEpisodes queries tvmaze, and returns a list of Episodes.
func (c *Client) GetEpisodes(showID int64) ([]Episode, error) {
	route := fmt.Sprintf("/shows/%d/episodes", showID)
	uri, err := url.Parse(c.BaseURI + route)
	if err != nil {
		return nil, err
	}

	jsondata, err := c.Go(uri)
	if err != nil {
		return nil, err
	}

	var episodes []Episode
	err = json.Unmarshal(jsondata, &episodes)
	if err != nil {
		return nil, err
	}

	return episodes, nil
}

// Go does an HTTP GET to tvmaze with the provided uri, and returns the response body.
// It will cache response if UseCache is true.
func (c *Client) Go(uri *url.URL) ([]byte, error) {
	data, found := c.Cache.Get(uri.String())

	if !found || !c.UseCache {
		if c.Debug {
			log.Print("cache miss: " + uri.String() + "\n")
		}
		request, err := http.NewRequest("GET", uri.String(), nil)
		if err != nil {
			return nil, err
		}

		resp, err := c.Do(request)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != 200 {
			return nil, errors.New(fmt.Sprintf("Request failed: %s", http.StatusText(resp.StatusCode)))
		}

		data, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		c.Cache.Set(uri.String(), data, 0)
	} else {
		if c.Debug {
			fmt.Printf("cache hit: %s\n", uri.String())
		}
	}

	return data.([]byte), nil
}

// Episode represents a tv episode.
type Episode struct {
	ID       int64
	URL      string
	Name     string
	Season   int
	Number   int
	AirDate  string
	AirTime  string
	AirStamp string
	Runtime  int
	//Image
	Summary string
	Links   Links `json:"_links"`
}

// Candidate represents a search candidate.
type Candidate struct {
	Score float64
	Show  Show
}

// Show represents a tv show.
type Show struct {
	ID        int64
	URL       string
	Name      string
	Type      string
	Language  string
	Genres    []string
	Status    string
	Runtime   int
	Premiered string
	Schedule  Schedule
	Rating    Rating
	Weight    int
	Network   Network
	//WebChannel
	Externals External
	Image     Image
	Summary   string
	Updated   int64
	Links     Links
}

// Links represents Episode links.
type Links struct {
	Self            Link
	PreviousEpisode Link
}

// Link represents a uri link.
type Link struct {
	Href string
}

// Image represents an image.
type Image struct {
	Medium   string
	Original string
}

// External represents a 3rd party tv api.
type External struct {
	TVRage  int64
	TheTVDB int64
}

// Schedule represents a Show schedule.
type Schedule struct {
	Time string
	Days []string
}

// Rating represents a tv show rating.
type Rating struct {
	Average float64
}

// Network represents the tv network airing the show.
type Network struct {
	ID      int
	Name    string
	Country Country
}

// Country represents the country the show aired in.
type Country struct {
	Name     string
	Code     string
	TimeZone string
}
