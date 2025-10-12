package subnet

import (
	"context"

	"github.com/IvanKondrashkov/go-shortener/internal/config"
	"github.com/IvanKondrashkov/go-shortener/internal/utils/admin"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// TrustedSubnet проверяет, принадлежит ли IP адрес указанной подсети для gRPC
func TrustedSubnet(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.PermissionDenied, "Missing metadata")
	}

	subnetHeaders := md.Get(subnetHeader)
	if len(subnetHeaders) == 0 {
		return nil, status.Errorf(codes.PermissionDenied, "Missing X-Real-IP header")
	}

	ipStr := subnetHeaders[0]
	isTrusted, err := admin.IsIPInSubnet(ipStr, config.TrustedSubnet)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Subnet check error")
	}

	if !isTrusted {
		return nil, status.Errorf(codes.PermissionDenied, "Access denied")
	}
	return handler(ctx, req)
}
