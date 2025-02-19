package db

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/vitorfhc/rssnotifier/pkg/types"
)

type databaseData struct {
	Feeds []types.Feed `json:"feeds"`
}

type Database struct {
	filepath string
	data     databaseData
}

func NewFromJSON(filepath string) (*Database, error) {
	d := &Database{
		filepath: filepath,
		data:     databaseData{},
	}

	if err := d.load(); err != nil {
		return nil, err
	}

	return d, nil
}

func (d *Database) AddFeed(feed types.Feed) error {
	for _, f := range d.data.Feeds {
		if f.Link == feed.Link {
			return fmt.Errorf("feed with link %q already exists", feed.Link)
		}
	}

	d.data.Feeds = append(d.data.Feeds, feed)
	return nil
}

func (d *Database) UpdateFeed(feed types.Feed) error {
	for i, f := range d.data.Feeds {
		if f.Link == feed.Link {
			d.data.Feeds[i] = feed
			return nil
		}
	}

	return fmt.Errorf("feed with link %q does not exist", feed.Link)
}

func (d *Database) GetFeeds() []types.Feed {
	return d.data.Feeds
}

func (d *Database) Save() error {
	data, err := json.MarshalIndent(d.data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(d.filepath, data, 0644)
}

func (d *Database) load() error {
	dat, err := os.ReadFile(d.filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return json.Unmarshal(dat, &d.data)
}
