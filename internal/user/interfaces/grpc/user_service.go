package grpc

import (
	"context"
	"github.com/Cospk/go-mall/internal/user/application/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/Cospk/go-mall/api/rpc/gen/go/user"
)

// UserServiceServer 用户服务gRPC实现
type UserServiceServer struct {
	pb.UnimplementedStreamGreeterServer
	userService *service.UserService
}

// NewUserServiceServer 创建用户服务gRPC服务器
func NewUserServiceServer(userService *service.UserService) *UserServiceServer {
	return &UserServiceServer{
		userService: userService,
	}
}

// GetUserInfo 获取用户信息
func (s *UserServiceServer) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
	// 调用应用服务
	user, err := s.userService.GetUserInfo(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "获取用户信息失败: %v", err)
	}

	if user == nil {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	// 转换响应
	return &pb.GetUserInfoResponse{
		UserId:    user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Phone:     user.Phone,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
