package main

import (
	pb "code-search/proto/search"
	"code-search/search-service/repo/searchdb"
	"context"
	"log"
)

func newSearchService() pb.SearchServiceServer {
	return &server{
		esClient: searchdb.New(),
	}
}

type server struct {
	pb.UnimplementedSearchServiceServer
	esClient searchdb.ESClient
}

func (s *server) Search(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	// 构建Elasticsearch查询
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  req.Query,
				"fields": []string{"content", "repo_name", "language"},
			},
		},
		"highlight": map[string]interface{}{
			"fields": map[string]interface{}{
				"content": map[string]interface{}{
					"number_of_fragments": 3,
					"fragment_size":       150,
				},
			},
		},
	}
	result, err := s.esClient.Search(ctx, query)
	if err != nil {
		return nil, err
	}
	log.Printf("Elasticsearch result: %v\n", result)
	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	results := make([]*pb.SearchResult, 0, len(hits))
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		highlight := hit.(map[string]interface{})["highlight"].(map[string]interface{})

		var highlightText string
		if highlights, ok := highlight["content"].([]interface{}); ok && len(highlights) > 0 {
			highlightText = highlights[0].(string)
		}
		result := &pb.SearchResult{
			License:   source["lic"].(string),
			Stars:     int32(source["stars"].(float64)),
			Highlight: highlightText,
			RepoName:  source["file_name"].(string),
			RepoPath:  source["file_name"].(string),
			Language:  source["lang"].(string),
		}
		results = append(results, result)
	}
	return &pb.SearchResponse{
		Results: results,
	}, nil
}
