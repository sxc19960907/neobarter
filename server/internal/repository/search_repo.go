package repository

import (
	"encoding/json"
	"fmt"

	"github.com/neobarter/server/internal/pkg/es"
)

type SearchRepository struct {
	esClient *es.Client
}

func NewSearchRepository(esClient *es.Client) *SearchRepository {
	return &SearchRepository{esClient: esClient}
}

type SearchQuery struct {
	Keyword    string
	CategoryID int
	Condition  string
	MinValue   float64
	MaxValue   float64
	Location   string
	SortBy     string // _score / created_at / view_count / estimated_value
	Page       int
	PageSize   int
}

type SearchResultItem struct {
	ID               int64             `json:"id"`
	Title            string            `json:"title"`
	Description      string            `json:"description"`
	CategoryID       int               `json:"category_id"`
	CategoryName     string            `json:"category_name"`
	EstimatedValue   float64           `json:"estimated_value"`
	Condition        string            `json:"condition"`
	Images           []string          `json:"images"`
	Location         string            `json:"location"`
	ViewCount        int               `json:"view_count"`
	UserNickname     string            `json:"user_nickname"`
	CreatedAt        string            `json:"created_at"`
	Highlight        map[string][]string `json:"highlight,omitempty"`
}

type SearchResponse struct {
	Items []SearchResultItem `json:"items"`
	Total int64              `json:"total"`
}

func (r *SearchRepository) Search(q SearchQuery) (*SearchResponse, error) {
	query := r.buildQuery(q)
	queryJSON, _ := json.Marshal(query)

	result, err := r.esClient.Search(string(queryJSON))
	if err != nil {
		return nil, err
	}

	response := &SearchResponse{
		Total: result.Hits.Total.Value,
		Items: make([]SearchResultItem, 0, len(result.Hits.Hits)),
	}

	for _, hit := range result.Hits.Hits {
		item := SearchResultItem{
			Highlight: hit.Highlight,
		}
		sourceBytes, _ := json.Marshal(hit.Source)
		json.Unmarshal(sourceBytes, &item)
		response.Items = append(response.Items, item)
	}

	return response, nil
}

func (r *SearchRepository) Suggest(prefix string) ([]string, error) {
	query := map[string]interface{}{
		"suggest": map[string]interface{}{
			"title_suggest": map[string]interface{}{
				"prefix": prefix,
				"completion": map[string]interface{}{
					"field":           "title.suggest",
					"size":            5,
					"skip_duplicates": true,
				},
			},
		},
	}

	queryJSON, _ := json.Marshal(query)
	result, err := r.esClient.Search(string(queryJSON))
	if err != nil {
		return nil, err
	}

	suggestions := make([]string, 0)
	if entries, ok := result.Suggest["title_suggest"]; ok && len(entries) > 0 {
		for _, opt := range entries[0].Options {
			suggestions = append(suggestions, opt.Text)
		}
	}
	return suggestions, nil
}

func (r *SearchRepository) buildQuery(q SearchQuery) map[string]interface{} {
	must := make([]interface{}, 0)
	filter := make([]interface{}, 0)

	// 全文搜索
	if q.Keyword != "" {
		must = append(must, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  q.Keyword,
				"fields": []string{"title^3", "description", "want_items"},
				"type":   "best_fields",
			},
		})
	}

	// 只搜索上架物品
	filter = append(filter, map[string]interface{}{
		"term": map[string]interface{}{"status": "active"},
	})

	// 分类筛选
	if q.CategoryID > 0 {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{"category_id": q.CategoryID},
		})
	}

	// 成色筛选
	if q.Condition != "" {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{"condition": q.Condition},
		})
	}

	// 地区筛选
	if q.Location != "" {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{"location": q.Location},
		})
	}

	// 估值范围
	if q.MinValue > 0 || q.MaxValue > 0 {
		rangeQ := map[string]interface{}{}
		if q.MinValue > 0 {
			rangeQ["gte"] = q.MinValue
		}
		if q.MaxValue > 0 {
			rangeQ["lte"] = q.MaxValue
		}
		filter = append(filter, map[string]interface{}{
			"range": map[string]interface{}{"estimated_value": rangeQ},
		})
	}

	// 构建 bool query
	boolQuery := map[string]interface{}{}
	if len(must) > 0 {
		boolQuery["must"] = must
	} else {
		boolQuery["must"] = []interface{}{map[string]interface{}{"match_all": map[string]interface{}{}}}
	}
	if len(filter) > 0 {
		boolQuery["filter"] = filter
	}

	// 排序
	var sort []interface{}
	switch q.SortBy {
	case "created_at":
		sort = []interface{}{map[string]interface{}{"created_at": "desc"}}
	case "view_count":
		sort = []interface{}{map[string]interface{}{"view_count": "desc"}}
	case "estimated_value":
		sort = []interface{}{map[string]interface{}{"estimated_value": "desc"}}
	default:
		if q.Keyword != "" {
			sort = []interface{}{map[string]interface{}{"_score": "desc"}}
		} else {
			sort = []interface{}{map[string]interface{}{"created_at": "desc"}}
		}
	}

	// 分页
	from := (q.Page - 1) * q.PageSize
	if from < 0 {
		from = 0
	}

	body := map[string]interface{}{
		"query": map[string]interface{}{"bool": boolQuery},
		"sort":  sort,
		"from":  from,
		"size":  q.PageSize,
		"highlight": map[string]interface{}{
			"fields": map[string]interface{}{
				"title":       map[string]interface{}{},
				"description": map[string]interface{}{"fragment_size": 150},
			},
			"pre_tags":  []string{"<em>"},
			"post_tags": []string{"</em>"},
		},
	}

	fmt.Printf("[ES Query] %+v\n", body)
	return body
}
