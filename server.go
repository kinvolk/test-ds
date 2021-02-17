package main

import (
	"context"
	"net"

	"google.golang.org/grpc"

	clusterservice "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discoverygrpc "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
)

type Server struct {
	grpcServer         *grpc.Server
	controlPlaneServer server.Server
}

func NewServer(ctx context.Context, scache cache.SnapshotCache) Server {
	srv := Server{
		grpcServer:         grpc.NewServer(),
		controlPlaneServer: server.NewServer(ctx, scache, nil),
	}
	srv.registerGRPC()
	return srv
}

func (s *Server) Run(ctx context.Context, port int) error {
	addr := net.TCPAddr{
		Port: port,
	}
	lcfg := net.ListenConfig{}
	listener, err := lcfg.Listen(ctx, "tcp", addr.String())
	if err != nil {
		return err
	}
	return s.grpcServer.Serve(listener)
}

func (s *Server) registerGRPC() {
	discoverygrpc.RegisterAggregatedDiscoveryServiceServer(s.grpcServer, s.controlPlaneServer)
	clusterservice.RegisterClusterDiscoveryServiceServer(s.grpcServer, s.controlPlaneServer)
}
