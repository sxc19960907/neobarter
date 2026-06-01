package es

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

const ItemIndex = "neobarter_items"

type Client struct {
	es *elasticsearch.Client
}

func NewClient(addresses []string) (*Client, error) {
	cfg := elasticsearch.Config{
		Addresses: addresses,
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create ES client: %w", err)
	}

	// 验证连接
	res, err := es.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ES: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("ES info error: %s", res.String())
	}

	log.Println("Elasticsearch connected")
	return &Client{es: es}, nil
}

// EnsureIndex 确保索引存在。优先用 IK 分词映射；若 IK 插件不可用则降级到内置分析器。
func (c *Client) EnsureIndex() error {
	res, err := c.es.Indices.Exists([]string{ItemIndex})
	if err != nil {
		return err
	}
	res.Body.Close()

	if res.StatusCode == 200 {
		return nil // 已存在
	}

	// 优先尝试 IK 中文分词映射
	if err := c.createIndex(ItemMapping); err != nil {
		log.Printf("IK 分词映射创建失败，降级到内置分析器: %v", err)
		// 降级到内置分析器映射
		if ferr := c.createIndex(ItemMappingFallback); ferr != nil {
			return fmt.Errorf("fallback index create failed: %w", ferr)
		}
		log.Printf("Created index %s (内置分析器降级模式)", ItemIndex)
		return nil
	}

	log.Printf("Created index: %s (IK 中文分词)", ItemIndex)
	return nil
}

func (c *Client) createIndex(mapping string) error {
	res, err := c.es.Indices.Create(
		ItemIndex,
		c.es.Indices.Create.WithBody(strings.NewReader(mapping)),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("%s", res.String())
	}
	return nil
}

// IndexItem 索引一个物品
func (c *Client) IndexItem(id int64, doc map[string]interface{}) error {
	body, _ := json.Marshal(doc)
	req := esapi.IndexRequest{
		Index:      ItemIndex,
		DocumentID: fmt.Sprintf("%d", id),
		Body:       strings.NewReader(string(body)),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), c.es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("index error: %s", res.String())
	}
	return nil
}

// DeleteItem 从索引中删除物品
func (c *Client) DeleteItem(id int64) error {
	req := esapi.DeleteRequest{
		Index:      ItemIndex,
		DocumentID: fmt.Sprintf("%d", id),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), c.es)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

// Search 搜索物品
func (c *Client) Search(query string) (*SearchResult, error) {
	res, err := c.es.Search(
		c.es.Search.WithIndex(ItemIndex),
		c.es.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.String())
	}

	var result SearchResult
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SearchResult ES 搜索结果
type SearchResult struct {
	Hits struct {
		Total struct {
			Value int64 `json:"value"`
		} `json:"total"`
		Hits []SearchHit `json:"hits"`
	} `json:"hits"`
	Suggest map[string][]SuggestEntry `json:"suggest"`
}

type SearchHit struct {
	ID        string                 `json:"_id"`
	Score     float64                `json:"_score"`
	Source    map[string]interface{} `json:"_source"`
	Highlight map[string][]string    `json:"highlight"`
}

type SuggestEntry struct {
	Text    string           `json:"text"`
	Options []SuggestOption  `json:"options"`
}

type SuggestOption struct {
	Text string `json:"text"`
}

func (c *Client) ES() *elasticsearch.Client {
	return c.es
}
