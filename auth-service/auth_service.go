package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"

	"code-search/auth-service/repo/userdb"
	"code-search/auth-service/repo/usertoken"
	pb "code-search/proto/auth"
)

func newAuthServiceServer() pb.AuthServiceServer {
	return &server{
		userDBClient:    userdb.NewClient(),
		userTokenClient: usertoken.NewClient(),
	}
}

type server struct {
	pb.UnimplementedAuthServiceServer
	userDBClient    userdb.Client
	userTokenClient usertoken.Client
}

// Login 登陆
func (s *server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if err := s.userDBClient.CheckUserPassword(ctx, req.GetUsername(), req.GetPassword()); err != nil {
		log.Printf("check user password error: %v", err)
		return nil, err
	}
	token := generateToken()
	if err := s.userTokenClient.SetUserToken(ctx, token, req.GetUsername()); err != nil {
		log.Printf("set user token error: %v", err)
		return nil, err
	}
	return &pb.LoginResponse{
		Success: true,
		Token:   token,
		Message: "Login successful",
	}, nil
}

func (s *server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if err := s.userDBClient.CreateUser(ctx, req.GetUsername(), req.GetEmail(), req.GetPassword()); err != nil {
		log.Printf("create user error: %v", err)
		return nil, err
	}
	return &pb.RegisterResponse{
		Success: true,
		Message: "Registration successful",
	}, nil
}

func (s *server) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	username, err := s.userTokenClient.ValidateToken(ctx, req.GetToken())
	if err != nil {
		return &pb.ValidateTokenResponse{
			Valid:    false,
			Username: "",
			Message:  "validate error",
		}, err
	}
	return &pb.ValidateTokenResponse{
		Valid:    true,
		Username: username,
		Message:  "validate pass",
	}, nil
}

func generateToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
