package service

import (
	"github.com/neobarter/server/internal/repository"
)

type SearchService struct {
	searchRepo *repository.SearchRepository
}

func NewSearchService(searchRepo *repository.SearchRepository) *SearchService {
	return &SearchService{searchRepo: searchRepo}
}

func (s *SearchService) Search(q repository.SearchQuery) (*repository.SearchResponse, error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 || q.PageSize > 50 {
		q.PageSize = 20
	}
	return s.searchRepo.Search(q)
}

func (s *SearchService) Suggest(prefix string) ([]string, error) {
	if prefix == "" {
		return []string{}, nil
	}
	return s.searchRepo.Suggest(prefix)
}
