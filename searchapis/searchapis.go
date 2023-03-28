package searchapis

import (
	"context"
	"encoding/json"
	"net/http"
    "errors"
    "strings"
)

const (
	PixabayapiURL = "https://pixabay.com/api/"
	TenorapiURL = "https://tenor.googleapis.com/v2/search"
)

var (
    ErrNoKey error = errors.New("No api key")
    ErrNoSearch error = errors.New("No search query")
)

type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// Config
	config *Config
}

type Config struct {
	// Base URL for API requests.
	BaseURL string
}

func NewPixabayClient(apikey string) (*Client, error) {
	if apikey == "" {
		return nil, ErrNoKey
	}

	return &Client{
		client: &http.Client{},
		config: &Config{
			BaseURL: PixabayapiURL + "?key=" + apikey,
		},
	}, nil
}

func NewTenorClient(apikey string) (*Client, error) {
	if apikey == "" {
		return nil, ErrNoKey
	}

	return &Client{
		client: &http.Client{},
		config: &Config{
			BaseURL: TenorapiURL + "?key=" + apikey,
		},
	}, nil
}

type PixabayHit struct {
    Collections     int             `json:"collections"`
    Comments        int             `json:"comments"`
    Downloads       int             `json:"downloads"`
	Id              int             `json:"id"`
    ImageHeight     int             `json:"imageHeight"`
    ImageSize       int             `json:"imageSize"`
    ImageWidth      int             `json:"imageWidth"`
    LargeImageURL   string          `json:"largeImageURL"`
    Likes           int             `json:"likes"`
    PageURL         string          `json:"pageURL"`
    PreviewHeight   int             `json:"previewHeight"`
    PreviewURL      string          `json:"previewURL"`
    PreviewWidth    int             `json:"previewWidth"`
    Tags            string          `json:"tags"`
    Type            string          `json:"type"`
    User            string          `json:"user"`
    UserImageURL    string          `json:"userImageURL"`
    User_Id         int             `json:"user_id"`
    Views           int             `json:"views"`
    WebformatHeight int             `json:"webformatHeight"`
    WebformatURL    string          `json:"webformatURL"`
    WebformatWidth  int             `json:"webformatWidth"`
}

type PixabayResult struct {
    Hits            []PixabayHit    `json:"hits"`  
    Total           int             `json:"total"`
    TotalHits       int             `json:"totalHits"`
}

func (c *Client) PixabayImageById(ctx context.Context, search string) (PixabayResult, error) {
	var (
        Err error
        empty PixabayResult
    )
	if !(len(search) > 0)  {
	    return empty, ErrNoSearch
	}
	search = strings.ReplaceAll(search, " ", "+")
	search = c.config.BaseURL + "&q=" + search 
	httpResp, err := http.Get(search)
	if err != nil {
	    return empty, err
	}
	defer httpResp.Body.Close()
	
	var result PixabayResult
	err = json.NewDecoder(httpResp.Body).Decode(&result); 
	if err != nil {
		return empty, err
	}
	
	return result, Err
}

type TenorResult struct {
	Results []struct {
		ID           string `json:"id"`
		Title        string `json:"title"`
		MediaFormats struct {
			Tinywebm struct {
				URL      string `json:"url"`
				Duration float64    `json:"duration"`
				Preview  string `json:"preview"`
				Dims     []int  `json:"dims"`
				Size     int    `json:"size"`
			} `json:"tinywebm"`
			Loopedmp4 struct {
				URL      string `json:"url"`
				Duration float64    `json:"duration"`
				Preview  string `json:"preview"`
				Dims     []int  `json:"dims"`
				Size     int    `json:"size"`
			} `json:"loopedmp4"`
			Mp4 struct {
				URL      string `json:"url"`
				Duration float64    `json:"duration"`
				Preview  string `json:"preview"`
				Dims     []int  `json:"dims"`
				Size     int    `json:"size"`
			} `json:"mp4"`
			Gif struct {
				URL      string `json:"url"`
				Duration float64    `json:"duration"`
				Preview  string `json:"preview"`
				Dims     []int  `json:"dims"`
				Size     int    `json:"size"`
			} `json:"gif"`
			Webm struct {
				URL      string `json:"url"`
				Duration float64   `json:"duration"`
				Preview  string `json:"preview"`
				Dims     []int  `json:"dims"`
				Size     int    `json:"size"`
			} `json:"webm"`
			Nanogif struct {
				URL      string `json:"url"`
				Duration float64 `json:"duration"`
				Preview  string `json:"preview"`
				Dims     []int  `json:"dims"`
				Size     int    `json:"size"`
			} `json:"nanogif"`
			Tinygifpreview struct {
				URL      string `json:"url"`
				Duration float64 `json:"duration"`
				Preview  string `json:"preview"`
				Dims     []int  `json:"dims"`
				Size     int    `json:"size"`
			} `json:"tinygifpreview"`
			Nanogifpreview struct {
				URL      string `json:"url"`
				Duration float64 `json:"duration"`
				Preview  string `json:"preview"`
				Dims     []int  `json:"dims"`
				Size     int    `json:"size"`
			} `json:"nanogifpreview"`
			Gifpreview struct {
				URL      string `json:"url"`
				Duration float64 `json:"duration"`
				Preview  string `json:"preview"`
				Dims     []int  `json:"dims"`
				Size     int    `json:"size"`
			} `json:"gifpreview"`
			Tinymp4 struct {
				URL      string  `json:"url"`
				Duration float64 `json:"duration"`
				Preview  string  `json:"preview"`
				Dims     []int   `json:"dims"`
				Size     int     `json:"size"`
			} `json:"tinymp4"`
			Nanomp4 struct {
				URL      string `json:"url"`
				Duration float64    `json:"duration"`
				Preview  string `json:"preview"`
				Dims     []int  `json:"dims"`
				Size     int    `json:"size"`
			} `json:"nanomp4"`
			Tinygif struct {
				URL      string `json:"url"`
				Duration float64    `json:"duration"`
				Preview  string `json:"preview"`
				Dims     []int  `json:"dims"`
				Size     int    `json:"size"`
			} `json:"tinygif"`
			Mediumgif struct {
				URL      string `json:"url"`
				Duration float64    `json:"duration"`
				Preview  string `json:"preview"`
				Dims     []int  `json:"dims"`
				Size     int    `json:"size"`
			} `json:"mediumgif"`
			Nanowebm struct {
				URL      string `json:"url"`
				Duration float64    `json:"duration"`
				Preview  string `json:"preview"`
				Dims     []int  `json:"dims"`
				Size     int    `json:"size"`
			} `json:"nanowebm"`
		} `json:"media_formats"`
		Created            float64  `json:"created"`
		ContentDescription string   `json:"content_description"`
		Itemurl            string   `json:"itemurl"`
		URL                string   `json:"url"`
		Tags               []string `json:"tags"`
		Flags              []any    `json:"flags"`
		Hasaudio           bool     `json:"hasaudio"`
	} `json:"results"`
	Next string `json:"next"`
}


func (c *Client) TenorGifById(ctx context.Context, search string) (TenorResult, error) {
	var (
        Err error
        empty TenorResult
    )
	if !(len(search) > 0)  {
	    return empty, ErrNoSearch
	}
	search = strings.ReplaceAll(search, " ", "+")
	search = c.config.BaseURL + "&q=" + search + "&limit=5"
	httpResp, err := http.Get(search)
	if err != nil {
	    return empty, err
	}
	defer httpResp.Body.Close()
	
	var result TenorResult
	err = json.NewDecoder(httpResp.Body).Decode(&result); 
	if err != nil {
		return empty, err
	}
	
	return result, Err
}