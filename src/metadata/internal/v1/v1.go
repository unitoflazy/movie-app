package v1

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"movie-app/gen"
	"movie-app/metadata/internal/v1/controller/metadata"
	grpchandler "movie-app/metadata/internal/v1/handler/grpc"
	"movie-app/metadata/internal/v1/handler/http"
	"movie-app/metadata/internal/v1/repository/memory"
	"net"
)

type GRPCApp struct {
	address string
	opts    []grpc.ServerOption
}

type GinApp struct {
	port string
}

func NewGRPCApp(address string, opts ...grpc.ServerOption) *GRPCApp {
	return &GRPCApp{
		address: address,
		opts:    opts,
	}
}

func (a *GRPCApp) Run() error {
	repo := memory.New()
	ctrl := metadata.New(repo)
	grpcHandler := grpchandler.NewGRPCHandler(ctrl)

	lis, err := net.Listen("tcp", a.address)
	if err != nil {
		return err
	}

	svr := grpc.NewServer(a.opts...)
	gen.RegisterMetadataServiceServer(svr, grpcHandler)

	return svr.Serve(lis)
}

func NewGinApp(port string) *GinApp {
	return &GinApp{
		port: port,
	}
}

func (a *GinApp) Run() error {
	repo := memory.New()
	ctrl := metadata.New(repo)
	ginHandler := http.NewGinHandler(ctrl)

	r := gin.Default()

	v1Router := r.Group("api/v1/metadata")
	v1Router.GET("/:id", ginHandler.GetMetadata)

	return r.Run(":" + a.port)
}
