package giphy

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/oklahomer/go-sarah"
	"github.com/oklahomer/go-sarah/slack"
	"github.com/oklahomer/golack/webapi"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

var matchPattern = regexp.MustCompile(`^\.giphy\s*`)
var SlackCommand = sarah.NewCommandBuilder().
	Identifier("giphy").
	InputExample(`".giphy" shows trending gifs. ".giphy FOO" shows translated gif for FOO.`).
	MatchPattern(matchPattern).
	Func(SlackCommandFunc).
	MustBuild()

func SlackCommandFunc(ctx context.Context, input sarah.Input) (*sarah.CommandResponse, error) {
	text := sarah.StripMessage(matchPattern, input.Message())

	var attachments []*webapi.MessageAttachment
	var reqErr error
	if text == "" {
		attachments, reqErr = trend(ctx)
	} else {
		attachments, reqErr = translate(ctx, text)
	}

	if reqErr != nil {
		return nil, reqErr
	}

	if len(attachments) == 0 {
		return nil, errors.New("No trending gif found.")
	}

	return slack.NewPostMessageResponse(input, "", attachments), nil
}

func trend(ctx context.Context) ([]*webapi.MessageAttachment, error) {
	query := &url.Values{}
	query.Set("limit", strconv.Itoa(6))

	response := &trendingResponse{}
	err := request(ctx, "/v1/gifs/trending", query, response)
	if err != nil {
		return nil, err
	}

	attachments := []*webapi.MessageAttachment{}
	for _, gif := range response.Data {
		attachments = append(attachments, &webapi.MessageAttachment{
			Fallback:  gif.URL,
			Title:     "Trending gif",
			TitleLink: gif.URL,
			ImageURL:  gif.Images.FixedWidth.URL,
		})
	}

	return attachments, nil
}

func translate(ctx context.Context, text string) ([]*webapi.MessageAttachment, error) {
	query := &url.Values{}
	query.Set("s", text)

	response := &translateResponse{}
	err := request(ctx, "/v1/gifs/translate", query, response)
	if err != nil {
		return nil, err
	}

	return []*webapi.MessageAttachment{
		{
			Fallback:  response.Data.URL,
			Title:     "Translation gif",
			TitleLink: response.Data.URL,
			ImageURL:  response.Data.Images.FixedWidth.URL,
		},
	}, nil
}

func request(ctx context.Context, path string, query *url.Values, response interface{}) error {
	query.Set("api_key", "dc6zaTOxFJmzC")
	endpoint := &url.URL{
		Scheme:   "http",
		Host:     "api.giphy.com",
		Path:     path,
		RawQuery: query.Encode(),
	}

	reqCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := ctxhttp.Get(reqCtx, http.DefaultClient, endpoint.String())
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response status error. status: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, response); err != nil {
		return err
	}

	return nil
}

type translateResponse struct {
	Data *Gif `json:"data"`
}

type trendingResponse struct {
	Data []*Gif `json:"data"`
}

type Gif struct {
	Type        string `json:"type"`
	Id          string `json:"id"`
	URL         string `json:"url"`
	Tags        string `json:"tags"`
	BitlyGifURL string `json:"bitly_gif_url"`
	Images      struct {
		Original               ImageProps `json:"original"`
		OriginalStill          ImageProps `json:"original_still"`
		FixedHeight            ImageProps `json:"fixed_height"`
		FixedHeightStill       ImageProps `json:"fixed_height_still"`
		FixedHeightDownsampled ImageProps `json:"fixed_height_downsampled"`
		FixedWidth             ImageProps `json:"fixed_width"`
		FixedWidthStill        ImageProps `json:"fixed_width_still"`
		FixedWidthDownsampled  ImageProps `json:"fixed_width_downsampled"`
	} `json:"images"`
}

type ImageProps struct {
	URL    string `json:"url"`
	Width  string `json:"width"`
	Height string `json:"height"`
}
