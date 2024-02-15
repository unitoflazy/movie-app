package v1

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"movie-app/gen"
	"movie-app/rating/internal/v1/controller/rating"
	grpchandler "movie-app/rating/internal/v1/handler/grpc"
	"movie-app/rating/internal/v1/handler/http"
	"movie-app/rating/internal/v1/repository/memory"
	"net"
)

type GRPCApp struct {
	address string
	opts    []grpc.ServerOption
}

func NewGRPCApp(address string, opts ...grpc.ServerOption) *GRPCApp {
	return &GRPCApp{
		address: address,
		opts:    opts,
	}
}

func (g *GRPCApp) Run() error {
	repo := memory.New()
	ctrl := rating.New(repo)
	handler := grpchandler.NewGRPCHandler(ctrl)

	lis, err := net.Listen("tcp", g.address)
	if err != nil {
		return err
	}

	srv := grpc.NewServer(g.opts...)
	gen.RegisterRatingServiceServer(srv, handler)
	return srv.Serve(lis)
}

type GinApp struct {
	port string
}

func NewGinApp(port string) *GinApp {
	return &GinApp{
		port: port,
	}
}

func (g *GinApp) Run() error {
	repo := memory.New()
	ctrl := rating.New(repo)
	handler := http.NewGinHandler(ctrl)

	r := gin.Default()
	v1Router := r.Group("api/v1/rating")

	v1Router.GET("/:id", handler.GetAggregatedRating)
	v1Router.PUT("/:id", handler.PutRating)

	return r.Run(":" + g.port)
}
