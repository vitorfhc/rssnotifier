package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/vitorfhc/rssnotifier/pkg/db"
	"github.com/vitorfhc/rssnotifier/pkg/types"
)

type NotifierOpts struct {
	DiscordWebhookURL string
}

type Notifier struct {
	database *db.Database
	opts     *NotifierOpts
}

type NotifierOpt func(*NotifierOpts)

func New(database *db.Database, opts ...NotifierOpt) *Notifier {
	n := &Notifier{
		opts:     &NotifierOpts{},
		database: database,
	}

	for _, opt := range opts {
		opt(n.opts)
	}

	return n
}

func WithDiscordWebhookURL(url string) NotifierOpt {
	return func(opts *NotifierOpts) {
		opts.DiscordWebhookURL = url
	}
}

func (n *Notifier) Run() error {
	feeds := n.database.GetFeeds()
	parser := gofeed.NewParser()

	for _, feed := range feeds {
		parsedFeeds, err := parser.ParseURL(feed.Link)
		if err != nil {
			return err
		}

		feedLastUpdated := time.Now().AddDate(-10, 0, 0)
		if feed.LastUpdated != "" {
			timeLayout := time.RFC3339
			feedLastUpdated, err = time.Parse(timeLayout, feed.LastUpdated)
			if err != nil {
				return err
			}
		}

		newItems := []gofeed.Item{}
		for _, item := range parsedFeeds.Items {
			if item.UpdatedParsed != nil && item.UpdatedParsed.After(feedLastUpdated) {
				newItems = append(newItems, *item)
			}
		}

		feed.LastUpdated = time.Now().Format(time.RFC3339)
		n.database.UpdateFeed(feed)

		if len(newItems) == 0 {
			continue
		}

		if n.opts.DiscordWebhookURL != "" {
			if err := n.SendDiscordNotification(feed, newItems); err != nil {
				return err
			}
		}
	}

	return n.database.Save()
}

func (n *Notifier) SendDiscordNotification(feed types.Feed, items []gofeed.Item) error {
	content := "New items in " + feed.Name + ":\n\n"

	for _, item := range items {
		newContent := item.Title + "\n"
		newContent += item.Link + "\n"
		newContent += "\n"

		if len(content)+len(newContent) > 2000 {
			break
		}

		content += newContent
	}

	payload := map[string]string{
		"content": content,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(n.opts.DiscordWebhookURL, "application/json", bytes.NewReader(payloadBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
