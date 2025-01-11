package grpcserver

import (
	"context"
	"net"

	"github.com/BieggerM/userservice/pkg/adapter/out/broker"
	"github.com/BieggerM/userservice/pkg/adapter/out/database"
	"github.com/BieggerM/userservice/pkg/models"
	"github.com/BieggerM/userservice/proto/user"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type UserServiceServer struct {
	user.UnimplementedUserServiceServer
	DB database.Postgres
	MB broker.RabbitMQ
}

func (s *UserServiceServer) StartGRPCServer(MB broker.RabbitMQ, DB database.Postgres) {
	s.DB = DB
	s.MB = MB
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
	u := s.DB.GetUser(req.Username)
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
