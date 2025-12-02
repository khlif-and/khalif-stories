package repository

import (
	"context"

	"github.com/meilisearch/meilisearch-go"

	"khalif-stories/internal/domain"

)

type SearchRepo struct {
	client *meilisearch.Client
}

func NewSearchRepository(client *meilisearch.Client) *SearchRepo {
	return &SearchRepo{client: client}
}

func (r *SearchRepo) SearchCategories(query string) ([]domain.Category, error) {
	searchRes, err := r.client.Index("categories").Search(query, &meilisearch.SearchRequest{
		Limit: 10,
	})
	if err != nil {
		return nil, err
	}

	var categories []domain.Category
	if err := searchRes.UnmarshalHits(&categories); err != nil {
		return nil, err
	}

	return categories, nil
}