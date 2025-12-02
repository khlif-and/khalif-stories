package main

import (
	"github.com/meilisearch/meilisearch-go"

	"khalif-stories/internal/config"
	"khalif-stories/pkg/search"

)

func NewMeiliClientFromConfig(cfg *config.Config) *meilisearch.Client {
	return search.NewMeiliClient(cfg.MeiliAddr, cfg.MeiliKey)
}