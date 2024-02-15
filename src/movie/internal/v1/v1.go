package v1

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"movie-app/gen"
	"movie-app/movie/internal/v1/controller"
	metagrpc "movie-app/movie/internal/v1/gateway/metadata/grpc"
	metahttp "movie-app/movie/internal/v1/gateway/metadata/http"
	ratinggrpc "movie-app/movie/internal/v1/gateway/rating/grpc"
	ratinghttp "movie-app/movie/internal/v1/gateway/rating/http"
	grpchandler "movie-app/movie/internal/v1/handler/grpc"
	"movie-app/movie/internal/v1/handler/http"
	"movie-app/pkg/discovery"
	"net"
)

type gatewayOptionMap map[string][]grpc.DialOption

type GRPCApp struct {
	address  string
	registry *discovery.Registry
	options  gatewayOptionMap
}

func NewGRPCApp(address string, registry *discovery.Registry, options gatewayOptionMap) *GRPCApp {
	return &GRPCApp{
		address:  address,
		registry: registry,
		options:  options,
	}
}

func (g *GRPCApp) Run() error {
	metadataGW := metagrpc.New(g.registry, g.options["metadata"]...)
	ratingGW := ratinggrpc.New(g.registry, g.options["rating"]...)
	ctrl := controller.New(ratingGW, metadataGW)

	handler := grpchandler.NewGRPCHandler(ctrl)
	lis, err := net.Listen("tcp", g.address)
	if err != nil {
		return err
	}
	srv := grpc.NewServer()
	gen.RegisterMovieServiceServer(srv, handler)
	return srv.Serve(lis)
}

type GinApp struct {
	port     string
	registry *discovery.Registry
}

func NewGinApp(port string, registry *discovery.Registry) *GinApp {
	return &GinApp{
		port:     port,
		registry: registry,
	}
}

func (g *GinApp) Run() error {
	metadataGW := metahttp.New(g.registry)
	ratingGW := ratinghttp.New(g.registry)
	ctrl := controller.New(ratingGW, metadataGW)
	handler := http.NewGinHandler(ctrl)

	r := gin.Default()
	v1Router := r.Group("api/v1/api")
	v1Router.GET("/:id", handler.GetMovieDetails)

	return r.Run(":" + g.port)
}
