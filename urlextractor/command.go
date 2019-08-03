package urlextractor

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/mvdan/xurls"
	"github.com/oklahomer/go-sarah/v2"
	"github.com/oklahomer/go-sarah/v2/slack"
	"github.com/oklahomer/golack/webapi"
	"golang.org/x/net/context/ctxhttp"
	"net/http"
	"strings"
	"time"
)

func init() {
	props := sarah.NewCommandPropsBuilder().
		BotType(slack.SLACK).
		Identifier("urlextractor").
		Instruction(`"This is my page http://example.com/foo" to get the information about http://example.com/foo`).
		MatchPattern(xurls.Strict).
		Func(slackCommandFunc).
		MustBuild()
	sarah.RegisterCommandProps(props)
}

type Document struct {
	URL         string
	Title       string
	Description string
	ImageURL    string
}

// CommandFunc provides the core function of url extractor
func slackCommandFunc(ctx context.Context, input sarah.Input) (*sarah.CommandResponse, error) {
	urls := xurls.Strict.FindAllString(input.Message(), -1)

	var docs []*Document
	for _, url := range urls {
		doc, err := extractContent(ctx, url)
		if err != nil {
			break
		}

		docs = append(docs, doc)
	}

	var attachments []*webapi.MessageAttachment
	for _, doc := range docs {
		attachments = append(attachments, &webapi.MessageAttachment{
			Fallback:   fmt.Sprintf("title: %s. description: %s.", doc.Title, doc.Description),
			Title:      doc.Title,
			TitleLink:  doc.URL,
			ImageURL:   doc.ImageURL,
			ThumbURL:   doc.ImageURL,
			AuthorLink: doc.URL,
			Text:       doc.Description,
		})
	}

	if len(attachments) == 0 {
		return nil, fmt.Errorf("error on fetching URL content(s): %s", strings.Join(urls, ","))
	}

	return slack.NewResponse(input, "", slack.RespWithAttachments(attachments))
}

func extractContent(ctx context.Context, url string) (*Document, error) {
	reqCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := ctxhttp.Get(reqCtx, http.DefaultClient, url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status error. status: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	title, _ := doc.Find(`meta[property="og:title"]`).First().Attr("content")
	if title == "" {
		title = doc.Find("title").First().Text()
	}

	description, _ := doc.Find(`meta[property="og:description"]`).First().Attr("content")
	if description == "" {
		description = doc.Find("description").First().Text()
	}

	image, _ := doc.Find(`meta[property="og:image"]`).First().Attr("content")

	return &Document{
		URL:         url,
		Title:       title,
		Description: description,
		ImageURL:    image,
	}, nil
}
