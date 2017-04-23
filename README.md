# tvmaze
--
    import "github.com/rickyninja/tvmaze"

Package tvmaze is an HTTP Client for the tvmaze API.

## Usage

#### type Candidate

```go
type Candidate struct {
	Score float64
	Show  Show
}
```

Candidate represents a search candidate.

#### type Client

```go
type Client struct {
	Debug     bool
	BaseURI   string
	Region    string
	Cache     *cache.Cache
	CacheFile string
	UseCache  bool
	*http.Client
	// UserAgent may be set to identify your application.
	UserAgent string
}
```

Client is a tvmaze client.

#### func  NewClient

```go
func NewClient(cachefile string) (*Client, error)
```
NewClient returns a ready to use Client.

#### func (*Client) FindShow

```go
func (c *Client) FindShow(showname string) (Show, error)
```
FindShow searches tvmaze for showname, and returns it as a Show if a match is
found.

#### func (*Client) GetEpisodes

```go
func (c *Client) GetEpisodes(showID int64) ([]Episode, error)
```
GetEpisodes queries tvmaze, and returns a list of Episodes.

#### func (*Client) GetShow

```go
func (c *Client) GetShow(show string) ([]Candidate, error)
```
GetShow queries tvmaze for show, and returns Candidates that may be a match.

#### func (*Client) Go

```go
func (c *Client) Go(uri *url.URL) ([]byte, error)
```
Go does an HTTP GET to tvmaze with the provided uri, and returns the response
body. It will cache response if UseCache is true.

#### func (*Client) WriteCache

```go
func (c *Client) WriteCache() error
```
WriteCache writes cache contents to disk.

#### type Country

```go
type Country struct {
	Name     string
	Code     string
	TimeZone string
}
```

Country represents the country the show aired in.

#### type Episode

```go
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
```

Episode represents a tv episode.

#### type External

```go
type External struct {
	TVRage  int64
	TheTVDB int64
}
```

External represents a 3rd party tv api.

#### type Image

```go
type Image struct {
	Medium   string
	Original string
}
```

Image represents an image.

#### type Link

```go
type Link struct {
	Href string
}
```

Link represents a uri link.

#### type Links

```go
type Links struct {
	Self            Link
	PreviousEpisode Link
}
```

Links represents Episode links.

#### type Network

```go
type Network struct {
	ID      int
	Name    string
	Country Country
}
```

Network represents the tv network airing the show.

#### type Rating

```go
type Rating struct {
	Average float64
}
```

Rating represents a tv show rating.

#### type Schedule

```go
type Schedule struct {
	Time string
	Days []string
}
```

Schedule represents a Show schedule.

#### type Show

```go
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
```

Show represents a tv show.
