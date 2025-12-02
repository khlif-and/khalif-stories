package main

import (
	"github.com/meilisearch/meilisearch-go"

	"khalif-stories/internal/config"
	"khalif-stories/pkg/search"

)

// Wrapper function agar Wire hanya melihat function ini, bukan package aslinya
func NewMeiliClientFromConfig(cfg *config.Config) *meilisearch.Client {
	return search.NewMeiliClient(cfg.MeiliAddr, cfg.MeiliKey)
}