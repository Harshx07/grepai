package search

import (
	"context"

	"github.com/yoanbernabeu/grepai/embedder"
	"github.com/yoanbernabeu/grepai/store"
)

type Searcher struct {
	store    store.VectorStore
	embedder embedder.Embedder
}

func NewSearcher(st store.VectorStore, emb embedder.Embedder) *Searcher {
	return &Searcher{
		store:    st,
		embedder: emb,
	}
}

func (s *Searcher) Search(ctx context.Context, query string, limit int) ([]store.SearchResult, error) {
	// Embed the query
	queryVector, err := s.embedder.Embed(ctx, query)
	if err != nil {
		return nil, err
	}

	// Search the store
	return s.store.Search(ctx, queryVector, limit)
}
