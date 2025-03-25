package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "code-search/proto/search"
	"code-search/search-service/entity/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	serverPort             int
	elasticsearchAddresses string
	elasticsearchPort      int
)

func init() {
	flag.IntVar(&serverPort, "port", 50052, "The server port")
	flag.StringVar(&elasticsearchAddresses, "esaddr", "127.0.0.1", "The elasticsearch addresses")
	flag.IntVar(&elasticsearchPort, "esport", 9200, "The elasticsearch port")
}

func main() {
	flag.Parse()
	config.DefaultConfig.Elasticsearch.Addresses = []string{fmt.Sprintf("http://%s:%d", elasticsearchAddresses, elasticsearchPort)}
	config.DefaultConfig.Services.Search.Port = serverPort
	// 创建gRPC服务器
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.DefaultConfig.Services.Search.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	log.Printf("Search service listening at %v", lis.Addr())
	reflection.Register(s)
	pb.RegisterSearchServiceServer(s, newSearchService())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
