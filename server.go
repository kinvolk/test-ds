package main

import (
	"context"
	"log"
	"net"

	clusterservice "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discoveryservice "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	routeservice "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer         *grpc.Server
	vhdsServer         *VHDS
	controlPlaneServer server.Server
}

func NewServer(ctx context.Context, scache cache.SnapshotCache, logger *log.Logger, controlPlaneName string) Server {
	srv := Server{
		grpcServer:         grpc.NewServer(),
		vhdsServer:         NewVHDSServer(logger, controlPlaneName),
		controlPlaneServer: server.NewServer(ctx, scache, newLogCallbacks(logger)),
	}
	srv.registerGRPC()
	return srv
}

type logCallbacks struct {
	logger *log.Logger
}

var _ server.Callbacks = (*logCallbacks)(nil)

func (cb *logCallbacks) OnFetchRequest(ctx context.Context, request *discoveryservice.DiscoveryRequest) error {
	cb.logger.Printf("OnFetchRequest:: request: %s", cb.dump(request))
	return nil
}

func (cb *logCallbacks) OnFetchResponse(request *discoveryservice.DiscoveryRequest, response *discoveryservice.DiscoveryResponse) {
	cb.logger.Printf("OnFetchResponse:: request: %s, response: %s", cb.dump(request), cb.dump(response))
}

func (cb *logCallbacks) OnStreamOpen(ctx context.Context, streamID int64, typeURL string) error {
	cb.logger.Printf("OnStreamOpen:: stream ID: %d, type URL: %s", streamID, typeURL)
	return nil
}

func (cb *logCallbacks) OnStreamClosed(streamID int64) {
	cb.logger.Printf("OnStreamClosed:: stream ID: %d", streamID)
}

func (cb *logCallbacks) OnStreamRequest(streamID int64, request *discoveryservice.DiscoveryRequest) error {
	cb.logger.Printf("OnStreamRequest:: stream ID: %d, request: %s", streamID, cb.dump(request))
	return nil
}

func (cb *logCallbacks) OnStreamResponse(streamID int64, request *discoveryservice.DiscoveryRequest, response *discoveryservice.DiscoveryResponse) {
	cb.logger.Printf("OnStreamResponse:: stream ID: %d, request: %s, response: %s", streamID, cb.dump(request), cb.dump(response))
}

func (cb *logCallbacks) dump(thing interface{}) string {
	return Dump(thing)
}

func newLogCallbacks(logger *log.Logger) server.Callbacks {
	return &logCallbacks{
		logger: logger,
	}
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
	discoveryservice.RegisterAggregatedDiscoveryServiceServer(s.grpcServer, s.controlPlaneServer)
	clusterservice.RegisterClusterDiscoveryServiceServer(s.grpcServer, s.controlPlaneServer)
	routeservice.RegisterVirtualHostDiscoveryServiceServer(s.grpcServer, s.vhdsServer)
}
