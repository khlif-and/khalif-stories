package search

import (
	"github.com/meilisearch/meilisearch-go"

)

func NewMeiliClient(addr, apiKey string) *meilisearch.Client {
	return meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   addr,
		APIKey: apiKey,
	})
}