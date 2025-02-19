# RSS Notifier

RSS Notifier is a command-line tool that polls RSS feeds and sends notifications when new items are found. Notifications can be sent via Discord webhook.

## Features
- Add RSS feeds to a local database
- Poll feeds for new items
- Send notifications to a Discord webhook

## Installation

```sh
# Install using Golang
go install github.com/vitorfhc/rssnotifier@latest
```

## Usage

### Add a New RSS Feed

```sh
rssnotifier add --url <RSS_FEED_URL> --name <FEED_NAME> --database <DATABASE_FILE>
```

- `--url`, `-u`: URL of the RSS feed (required)
- `--name`, `-n`: Name of the RSS feed (required)
- `--database`, `-d`: Path to the database file (default: `rssnotifier.json`)

### Poll RSS Feeds and Send Notifications

```sh
rssnotifier poll --discord-webhook <WEBHOOK_URL> --database <DATABASE_FILE>
```

- `--discord-webhook`, `-w`: Discord webhook URL for notifications (required)
- `--database`, `-d`: Path to the database file (default: `rssnotifier.json`)
