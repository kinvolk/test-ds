package main

import (
	"io"
	"log"
	"strconv"

	routeservice "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	discoveryservice "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
)

type VHDS struct {
	logger                 *log.Logger
	controlPlaneIdentifier string
}

var (
	_ routeservice.VirtualHostDiscoveryServiceServer = (*VHDS)(nil)
)

func NewVHDSServer(logger *log.Logger, controlPlaneIdentifier string) *VHDS {
	return &VHDS{
		logger:                 logger,
		controlPlaneIdentifier: controlPlaneIdentifier,
	}
}

func (v *VHDS) DeltaVirtualHosts(dvhs routeservice.VirtualHostDiscoveryService_DeltaVirtualHostsServer) error {
	streamNonce := (uint64)(0)
	for {
		request, err := dvhs.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			v.logger.Printf("DeltaVirtualHosts:: failed to receive a request: %v", err)
			return err
		}
		v.logger.Printf("DeltaVirtualHosts:: request: %s", Dump(request))
		response := &discoveryservice.DeltaDiscoveryResponse{}

		response.TypeUrl = request.TypeUrl
		streamNonce++
		response.Nonce = strconv.FormatUint(streamNonce, 10)
		if err := dvhs.Send(response); err != nil {
			return err
		}
	}
}
