package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"code-search/auth-service/entity/config"
	pb "code-search/proto/auth"
	"google.golang.org/grpc"
)

var (
	serverPort int

	redisHost string
	redisPort int
	redisPass string
	redisDB   int

	mysqlHost     string
	mysqlPort     int
	mysqlUser     string
	mysqlPassword string
	mysqlDB       string
)

func init() {
	flag.IntVar(&serverPort, "port", 50051, "The server port")

	flag.StringVar(&redisHost, "redishost", "localhost", "Redis host")
	flag.IntVar(&redisPort, "redisport", 6379, "Redis port")
	flag.StringVar(&redisPass, "redispass", "", "Redis password")
	flag.IntVar(&redisDB, "redisdb", 0, "Redis db")

	flag.StringVar(&mysqlHost, "mysqlhost", "localhost", "MySQL host")
	flag.IntVar(&mysqlPort, "mysqlport", 3306, "MySQL port")
	flag.StringVar(&mysqlUser, "mysqluser", "root", "MySQL user")
	flag.StringVar(&mysqlPassword, "mysqlpassword", "", "MySQL password")
	flag.StringVar(&mysqlDB, "mysqldb", "code_search_user", "MySQL db")
}

func main() {
	flag.Parse()
	config.DefaultConfig.Services.Auth.Port = serverPort

	config.DefaultConfig.MySQL.Host = mysqlHost
	config.DefaultConfig.MySQL.Port = mysqlPort
	config.DefaultConfig.MySQL.User = mysqlUser
	config.DefaultConfig.MySQL.Password = mysqlPassword
	config.DefaultConfig.MySQL.DBName = mysqlDB

	config.DefaultConfig.Redis.Host = redisHost
	config.DefaultConfig.Redis.Port = redisPort
	config.DefaultConfig.Redis.Password = redisPass
	config.DefaultConfig.Redis.DB = redisDB

	// 创建gRPC服务器
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.DefaultConfig.Services.Auth.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, newAuthServiceServer())
	log.Printf("Auth service listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
