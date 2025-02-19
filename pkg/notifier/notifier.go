package notifier

import (
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/vitorfhc/rssnotifier/pkg/db"
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
			if item.UpdatedParsed.After(feedLastUpdated) {
				newItems = append(newItems, *item)
			}
		}

		feed.LastUpdated = time.Now().Format(time.RFC3339)
		n.database.UpdateFeed(feed)

		if len(newItems) == 0 {
			continue
		}

		// TODO: Send the new items to the Discord webhook
	}

	return n.database.Save()
}
