/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vitorfhc/rssnotifier/pkg/db"
	"github.com/vitorfhc/rssnotifier/pkg/notifier"
	"github.com/vitorfhc/rssnotifier/pkg/types"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rssnotifier",
	Short: "RSS Notifier polls RSS feeds and sends notifications if there are new items",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	dbpath := rootCmd.PersistentFlags().StringP("database", "d", "rssnotifier.json", "Path to the database file")

	addCmd := cobra.Command{
		Use:   "add",
		Short: "Add a new RSS feed",
		RunE: func(cmd *cobra.Command, args []string) error {
			if dbpath == nil {
				return fmt.Errorf("database path is required")
			}

			database, err := db.NewFromJSON(*dbpath)
			if err != nil {
				return fmt.Errorf("failed to read database: %w", err)
			}

			url, err := cmd.Flags().GetString("url")
			if err != nil {
				return fmt.Errorf("failed to get URL: %w", err)
			}

			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return fmt.Errorf("failed to get name: %w", err)
			}

			if url == "" {
				return fmt.Errorf("URL is required")
			}

			if name == "" {
				return fmt.Errorf("name is required")
			}

			database.AddFeed(types.Feed{
				Name: name,
				Link: url,
			})

			if err := database.Save(); err != nil {
				return fmt.Errorf("failed to save database: %w", err)
			}

			return nil
		},
	}

	addCmd.Flags().StringP("url", "u", "", "URL of the RSS feed")
	addCmd.MarkFlagRequired("url")

	addCmd.Flags().StringP("name", "n", "", "Name of the RSS feed")
	addCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(&addCmd)

	pollCmd := cobra.Command{
		Use:   "poll",
		Short: "Poll the RSS feeds",
		RunE: func(cmd *cobra.Command, args []string) error {
			if dbpath == nil {
				return fmt.Errorf("database path is required")
			}

			database, err := db.NewFromJSON(*dbpath)
			if err != nil {
				return fmt.Errorf("failed to read database: %w", err)
			}

			webhook, err := cmd.Flags().GetString("discord-webhook")
			if err != nil {
				return fmt.Errorf("failed to get Discord webhook: %w", err)
			}

			if webhook == "" {
				return fmt.Errorf("webhook for Discord is required")
			}

			notfr := notifier.New(database, notifier.WithDiscordWebhookURL(webhook))
			if err := notfr.Run(); err != nil {
				return fmt.Errorf("failed to run notifier: %w", err)
			}

			return nil
		},
	}

	pollCmd.Flags().StringP("discord-webhook", "w", "", "Discord webhook URL")
	pollCmd.MarkFlagRequired("discord-webhook")

	rootCmd.AddCommand(&pollCmd)
}
