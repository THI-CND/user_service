package grpcserver

import (
	"context"
	"github.com/BieggerM/userservice/pkg/service/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"

	"github.com/BieggerM/userservice/pkg/adapter/out/broker"
	"github.com/BieggerM/userservice/pkg/adapter/out/database"
	"github.com/BieggerM/userservice/pkg/adapter/out/logger"
	"github.com/BieggerM/userservice/pkg/models"
	"github.com/BieggerM/userservice/proto/user"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcServer interface {
	StartGRPCServer(MB broker.MessageBroker, DB database.Database, rlog logger.Logger, auth auth.AuthService)
}
type UserServiceServer struct {
	user.UnimplementedUserServiceServer
	DB   database.Database
	MB   broker.MessageBroker
	rlog logger.Logger
	auth auth.AuthService
}

func (s *UserServiceServer) StartGRPCServer(MB broker.MessageBroker, DB database.Database, rlog logger.Logger, auth auth.AuthService) {
	s.DB = DB
	s.MB = MB
	s.rlog = rlog
	s.auth = auth
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		logrus.Fatalf("Failed to listen: %v", err)
	}
	server := grpc.NewServer()
	user.RegisterUserServiceServer(server, s)
	reflection.Register(server)
	logrus.Infoln("GRPC Server started")
	if err := server.Serve(lis); err != nil {
		logrus.Fatalf("Failed to serve: %v", err)
	}
}

func (s *UserServiceServer) ListUsers(ctx context.Context, req *user.Empty) (*user.UserListResponse, error) {
	users := s.DB.ListUsers()
	var userList []*user.User
	for _, u := range users {
		userList = append(userList, &user.User{
			Username:  u.Username,
			Firstname: u.FirstName,
			Lastname:  u.LastName,
		})
	}
	return &user.UserListResponse{Users: userList}, nil
}

func (s *UserServiceServer) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.UserResponse, error) {
	u, err := s.DB.GetUser(req.Username)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	return &user.UserResponse{User: &user.User{
		Username:  u.Username,
		Firstname: u.FirstName,
		Lastname:  u.LastName,
	}}, nil
}

func (s *UserServiceServer) CreateUser(ctx context.Context, req *user.User) (*user.UserResponse, error) {
	newUser := models.User{
		Username:  req.Username,
		FirstName: req.Firstname,
		LastName:  req.Lastname,
	}
	if err := s.DB.SaveUser(newUser); err != nil {
		return nil, err
	}
	s.rlog.Info("User created", "username", newUser.Username)
	return &user.UserResponse{User: req}, nil
}

func (s *UserServiceServer) UpdateUser(ctx context.Context, req *user.User) (*user.UserResponse, error) {
	updatedUser := models.User{
		Username:  req.Username,
		FirstName: req.Firstname,
		LastName:  req.Lastname,
	}
	s.DB.UpdateUser(updatedUser)
	return &user.UserResponse{User: req}, nil
}

func (s *UserServiceServer) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (*user.DeleteUserResponse, error) {
	s.DB.DeleteUser(req.Username)
	return &user.DeleteUserResponse{Message: "user deleted"}, nil
}

func (s *UserServiceServer) Auth(ctx context.Context, req *user.AuthRequest) (*user.AuthResponse, error) {
	token := req.Token
	if token == "" {
		return nil, status.Errorf(codes.Unauthenticated, "Authorization token not provided")
	}

	valid, err := s.auth.ValidateJWT(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Failed to validate JWT: %v", err)
	}
	if !valid {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid JWT")
	}

	return &user.AuthResponse{Message: "valid JWT"}, nil
}
