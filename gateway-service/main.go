package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"code-search/gateway-service/config"
	pb "code-search/proto/auth"
	searchpb "code-search/proto/search"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SearchRequest struct {
	Query string `json:"query"`
}

var (
	gatewayPort = 8000
	authPort    = 50051
	searchPort  = 50052
)

func init() {
	flag.IntVar(&gatewayPort, "gatewayport", gatewayPort, "gateway port")
	flag.IntVar(&authPort, "authport", authPort, "auth service port")
	flag.IntVar(&searchPort, "searchport", searchPort, "search service port")
}

func main() {
	flag.Parse()

	config.DefaultConfig.Services.Gateway.Port = gatewayPort
	config.DefaultConfig.Services.Auth.Port = authPort
	config.DefaultConfig.Services.Search.Port = searchPort

	cfg := config.DefaultConfig

	// 连接认证服务
	authConn, err := grpc.NewClient(
		fmt.Sprintf("localhost:%d", cfg.Services.Auth.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to auth service: %v", err)
	}
	defer authConn.Close()

	// 连接搜索服务
	searchConn, err := grpc.NewClient(
		fmt.Sprintf("localhost:%d", cfg.Services.Search.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to search service: %v", err)
	}
	defer searchConn.Close()

	authClient := pb.NewAuthServiceClient(authConn)
	searchClient := searchpb.NewSearchServiceClient(searchConn)

	r := gin.Default()
	r.Use(CORSMiddleware())

	// 登录接口
	r.POST("/login", func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
			return
		}

		resp, err := authClient.Login(context.Background(), &pb.LoginRequest{
			Username: req.Username,
			Password: req.Password,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Login failed"})
			return
		}

		if !resp.Success {
			c.JSON(http.StatusUnauthorized, gin.H{"message": resp.Message})
			return
		}

		c.JSON(http.StatusOK, gin.H{"usertoken": resp.Token})
	})

	// 注册接口
	r.POST("/register", func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
			return
		}

		resp, err := authClient.Register(context.Background(), &pb.RegisterRequest{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Registration failed"})
			return
		}

		if !resp.Success {
			c.JSON(http.StatusBadRequest, gin.H{"message": resp.Message})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
	})

	// 搜索接口
	r.POST("/search", func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			log.Printf("failed to get usertoken: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"message": "No usertoken provided"})
			return
		}

		// 验证token
		token = strings.TrimPrefix(token, "Bearer ")
		validateResp, err := authClient.ValidateToken(context.Background(), &pb.ValidateTokenRequest{
			Token: token,
		})
		if err != nil {
			log.Printf("failed to validate token: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Token validation failed"})
			return
		}

		if !validateResp.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid usertoken"})
			return
		}

		var req SearchRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("failed to bind json: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
			return
		}
		resp, err := searchClient.Search(context.Background(), &searchpb.SearchRequest{
			Query: req.Query,
		})
		if err != nil {
			fmt.Printf("failed to search: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Search failed"})
			return
		}
		c.JSON(http.StatusOK, resp.Results)
	})

	log.Printf("Gateway service listening at :%d", cfg.Services.Gateway.Port)
	if err = r.Run(fmt.Sprintf(":%d", cfg.Services.Gateway.Port)); err != nil {
		log.Fatal(err)
	}
}

// CORSMiddleware 处理跨域请求
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization, Accept, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
