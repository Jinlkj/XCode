package main

import (
	pb "code-search/proto/search"
	"code-search/search-service/repo/searchdb"
	"context"
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
	//query := map[string]interface{}{
	//	"query": map[string]interface{}{
	//		"multi_match": map[string]interface{}{
	//			"query":  req.Query,
	//			"fields": []string{"content", "repo_name", "language"},
	//		},
	//	},
	//	"highlight": map[string]interface{}{
	//		"fields": map[string]interface{}{
	//			"content": map[string]interface{}{
	//				"number_of_fragments": 3,
	//				"fragment_size":       150,
	//			},
	//		},
	//	},
	//}
	//result, err := s.esClient.Search(ctx, query)
	//if err != nil {
	//	return nil, err
	//}
	//log.Printf("Elasticsearch result: %v\n", result)
	//hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	//results := make([]*pb.SearchResult, 0, len(hits))
	//for _, hit := range hits {
	//	source := hit.(map[string]interface{})["_source"].(map[string]interface{})
	//	highlight := hit.(map[string]interface{})["highlight"].(map[string]interface{})
	//
	//	var highlightText string
	//	if highlights, ok := highlight["content"].([]interface{}); ok && len(highlights) > 0 {
	//		highlightText = highlights[0].(string)
	//	}
	//	result := &pb.SearchResult{
	//		License:   source["lic"].(string),
	//		Stars:     int32(source["stars"].(float64)),
	//		Highlight: highlightText,
	//		RepoName:  source["file_name"].(string),
	//		RepoPath:  source["file_name"].(string),
	//		Language:  source["lang"].(string),
	//	}
	//	results = append(results, result)
	//}
	fakeResult := []*pb.SearchResult{
		{
			License:   "GPL-3.0 license",
			Stars:     655,
			Highlight: "Sorting algorithms visualized using the Blender Python API.",
			RepoName:  "ForeignGods/Sorting-Algorithms-Blender",
			RepoPath:  "https://github.com/ForeignGods/Sorting-Algorithms-Blender",
			Language:  "Python",
		},
		{
			License:   "MIT",
			Stars:     24,
			Highlight: "hello world",
			RepoName:  "B3ns44d/Python_Sorting_Algorithms",
			RepoPath:  "https://github.com/B3ns44d/Python_Sorting_Algorithms",
			Language:  "Python",
		},
		{
			License:   "MIT",
			Stars:     151,
			Highlight: "Basic sorting algorithms in python. Simplicity.",
			RepoName:  "ztgu/sorting_algorithms_py",
			RepoPath:  "test",
			Language:  "Python",
		},
		{
			License:   "MIT",
			Stars:     239,
			Highlight: "Sorting algorithm visualisation with Cairo\n",
			RepoName:  "cortesi/sortvis",
			RepoPath:  "test",
			Language:  "Python",
		},
		{
			License:   "MIT",
			Stars:     65,
			Highlight: "Python sorting algorithm visualizer.\n",
			RepoName:  "techwithtim/Sorting-Algorithm-Visualizer",
			RepoPath:  "techwithtim/Sorting-Algorithm-Visualizer",
			Language:  "Python",
		},
		{
			License:   "MIT",
			Stars:     596,
			Highlight: "realtime multiple people tracking (centerNet based person detector + deep sort algorithm with pytorch)",
			RepoName:  "kimyoon-young/centerNet-deep-sort",
			RepoPath:  "test",
			Language:  "Python",
		},
	}
	return &pb.SearchResponse{
		Results: fakeResult,
	}, nil
}
