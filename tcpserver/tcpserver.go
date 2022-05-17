package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/jonasngs/go_entry_task/grpc"
	"github.com/jonasngs/go_entry_task/tcpserver/config"
	"github.com/jonasngs/go_entry_task/tcpserver/manager"
	"google.golang.org/grpc"
)

type tcpserver struct {
	pb.UnimplementedUserManagerServer
}

var user_manager *manager.Manager

func (s *tcpserver) LoadProfile(ctx context.Context, req *pb.ProfileRequest) (*pb.ProfileResponse, error) {
	session_token := req.GetSessionToken()
	user, err := user_manager.Session.VerifySession(session_token)
	if err != nil {
		return nil, err
	}
	return &pb.ProfileResponse{
		User: user,
	}, nil
}

func (s *tcpserver) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	nickname := req.GetNickname()
	profile_picture := req.GetProfilePicture() // Type image
	session_token := req.GetSessionToken()
	user, err := user_manager.Session.VerifySession(session_token)
	if err != nil {
		return nil, err
	}

	err = user_manager.User.UpdateProfilePicture(user.GetUsername(), session_token, profile_picture.GetData(), profile_picture.GetFileExtension())
	if err != nil {
		return &pb.UpdateResponse{
		}, err
	}

	err = user_manager.User.UpdateNickname(user.GetUsername(), session_token, nickname)
	if err != nil {
		return &pb.UpdateResponse{
		}, err
	}
	return &pb.UpdateResponse{
	}, nil
}

func (s *tcpserver) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	username := req.GetUsername()
	password := req.GetPassword()
	err := user_manager.Authentication.AuthenticateLogin(username, password)
	if err != nil {
		return &pb.LoginResponse{
			SessionToken: "",
		}, err
	}
	new_token := user_manager.Session.GenerateSessionToken()
	fmt.Println(new_token)
	user_manager.Session.CreateSession(username, new_token)
	return &pb.LoginResponse{
		SessionToken: new_token,
	}, nil
}

func main() {
	config.OpenConfigFile()

	//lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 4040))
	lis, err := net.Listen("tcp", config.GetTCPServer())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	user_manager = manager.InitializeManager()

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterUserManagerServer(grpcServer, &tcpserver{})
	grpcServer.Serve(lis)

}
