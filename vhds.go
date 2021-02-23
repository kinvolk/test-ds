package main

import (
	"errors"
	"log"

	//coreconfig "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	routeservice "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	//discoveryservice "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
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
	request, err := dvhs.Recv()
	if err != nil {
		v.logger.Printf("DeltaVirtualHosts:: failed to receive a request: %v", err)
		return err
	}
	v.logger.Printf("DeltaVirtualHosts:: request: %s", Dump(request))

	return errors.New("not implemented")
	/*
		response := &discoveryservice.DeltaDiscoveryResponse{}

		response.Nonce = request.ResponseNonce
		response.ControlPlane = &configcore.ControlPlane{
			Identifier: v.controlPlaneIdentifier,
		}
		dvhs.Send(response)
	*/
}
