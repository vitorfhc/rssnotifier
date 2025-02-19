package types

type Feed struct {
	Name          string `json:"name"`
	URL           string `json:"url"`
	LastPublished string `json:"last_published"`
}
