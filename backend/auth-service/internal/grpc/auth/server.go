package auth

import (
	"context"
	"errors"
	"google.golang.org/grpc/metadata"

	pb "github.com/Blxssy/social-media/auth-service/api/auth"
	"google.golang.org/grpc"
)

const (
	emptyValue = 0
)

type Auth interface {
	Register(ctx context.Context, username, email, password string) (string, string, error)
	Login(ctx context.Context, email, password string) (string, string, error)
	IsAdmin(ctx context.Context, userID int) (bool, error)
}

type ServerAPI struct {
	pb.UnsafeAuthServiceServer
	auth Auth
}

func Register(grpcServer *grpc.Server, auth Auth) {
	pb.RegisterAuthServiceServer(grpcServer, &ServerAPI{auth: auth})
}

func (s *ServerAPI) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := s.auth.Register(ctx, req.GetUsername(), req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}
	md := metadata.Pairs("authorization", "Bearer "+accessToken)
	metadata.NewOutgoingContext(context.Background(), md)
	return &pb.RegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *ServerAPI) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *ServerAPI) IsAdmin(ctx context.Context, req *pb.IsAdminRequest) (*pb.IsAdminResponse, error) {
	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, int(req.GetUserId()))
	if err != nil {
		return nil, err
	}

	return &pb.IsAdminResponse{IsAdmin: isAdmin}, nil
}

func validateRegister(req *pb.RegisterRequest) error {
	if req.GetEmail() == "" {
		return errors.New("missing email")
	}

	if req.GetPassword() == "" {
		return errors.New("missing password")
	}

	return nil
}

func validateLogin(req *pb.LoginRequest) error {
	if req.GetEmail() == "" {
		return errors.New("missing email")
	}

	if req.GetPassword() == "" {
		return errors.New("missing password")
	}

	return nil
}

func validateIsAdmin(req *pb.IsAdminRequest) error {
	if req.GetUserId() == emptyValue {
		return errors.New("missing auth id")
	}

	return nil
}
