package urlgrpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	service "github.com/VadimShara/url_shortening_service/api/gen/go"
)

type serverApi struct {
	service.UnimplementedUrlServer
	url Url
}

type Url interface {
	SaveUrl(
		ctx context.Context,
		url string,
	) (alias string, err error)
	RedirectUrl(
		ctx context.Context,
		alias string,
	) (url string, err error)
}

func Register(gRPCServer *grpc.Server, url Url) {
	service.RegisterUrlServer(gRPCServer, &serverApi{url: url})
}

func (s *serverApi) SaveUrl(
	ctx context.Context,
	in *service.SaveUrlRequest,
) (*service.SaveUrlResponse, error) {
	if in.Url == "" {
		return nil, status.Error(codes.InvalidArgument, "url is required")
	}

	alias, err := s.url.SaveUrl(ctx, in.Url)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &service.SaveUrlResponse{Alias: alias}, nil
}

func (s *serverApi) RedirectUrl(
	ctx context.Context,
	in *service.RedirectUrlRequest,
) (*service.RedirectUrlResponse, error) {
	if in.Alias == "" {
		return nil, status.Error(codes.InvalidArgument, "alias is required")
	}

	url, err := s.url.RedirectUrl(ctx, in.Alias)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &service.RedirectUrlResponse{Url: url}, nil
}
