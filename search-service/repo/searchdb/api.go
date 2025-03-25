package searchdb

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"code-search/search-service/entity/config"
	"github.com/elastic/go-elasticsearch/v8"
)

func New() ESClient {
	cfg := config.DefaultConfig
	// 创建Elasticsearch客户端
	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: cfg.Elasticsearch.Addresses,
	})
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	return &impl{
		esClient: esClient,
	}
}

type ESClient interface {
	Search(ctx context.Context, query map[string]interface{}) (map[string]interface{}, error)
}

type impl struct {
	esClient *elasticsearch.Client
}

func (i *impl) Search(ctx context.Context, query map[string]interface{}) (map[string]interface{}, error) {
	res, err := i.esClient.Search(
		i.esClient.Search.WithContext(ctx),
		i.esClient.Search.WithIndex("code"),
		i.esClient.Search.WithBody(strings.NewReader(mustEncodeJSON(query))),
	)
	if err != nil {
		log.Printf("Error getting response: %s", err)
		return nil, err
	}
	defer res.Body.Close()
	// 解析响应
	var result map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Printf("Error parsing the response body: %s", err)
		return nil, err
	}
	return result, err
}

func mustEncodeJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(data)
}
