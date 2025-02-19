package notifier

import (
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

// func (n *Notifier) Run() error {
// 	feeds := n.database.GetFeeds()
// 	now := time.Now()

// 	for _, feed := range feeds {
// 	}

// 	return nil
// }
