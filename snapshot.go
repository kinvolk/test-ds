package main

import (
	"time"

	"github.com/golang/protobuf/ptypes"

	clusterconfig "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	coreconfig "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpointconfig "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	routeconfig "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	cache "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
)

func GetSnapshot(clusterName string, localPort int, routeConfigName, xdsClusterName string) (cache.Snapshot, error) {
	version := "1"
	snapshot := cache.NewSnapshot(
		version,
		nil, // no endpoints, cluster has them
		[]types.Resource{
			getCluster(clusterName, localPort),
		},
		[]types.Resource{
			getRoute(routeConfigName, xdsClusterName),
		},
		nil, // no listeners
		nil, // no runtimes
		nil, // no secrets
	)
	/*
	if err := snapshot.Consistent(); err != nil {
		return cache.Snapshot{}, err
	}
	*/
	return snapshot, nil
}

func getCluster(clusterName string, localPort int) *clusterconfig.Cluster {
	return &clusterconfig.Cluster{
		Name:           clusterName,
		ConnectTimeout: ptypes.DurationProto(5 * time.Second),
		ClusterDiscoveryType: &clusterconfig.Cluster_Type{
			Type: clusterconfig.Cluster_LOGICAL_DNS,
		},
		LbPolicy:        clusterconfig.Cluster_ROUND_ROBIN,
		LoadAssignment:  getEndpoint(clusterName, localPort),
		DnsLookupFamily: clusterconfig.Cluster_V4_ONLY,
	}
}

func getEndpoint(clusterName string, localPort int) *endpointconfig.ClusterLoadAssignment {
	return &endpointconfig.ClusterLoadAssignment{
		ClusterName: clusterName,
		Endpoints: []*endpointconfig.LocalityLbEndpoints{{
			LbEndpoints: []*endpointconfig.LbEndpoint{{
				HostIdentifier: &endpointconfig.LbEndpoint_Endpoint{
					Endpoint: &endpointconfig.Endpoint{
						Address: &coreconfig.Address{
							Address: &coreconfig.Address_SocketAddress{
								SocketAddress: &coreconfig.SocketAddress{
									Protocol: coreconfig.SocketAddress_TCP,
									Address:  "127.0.0.1",
									PortSpecifier: &coreconfig.SocketAddress_PortValue{
										PortValue: (uint32)(localPort),
									},
								},
							},
						},
					},
				},
			}},
		}},
	}
}

func getRoute(routeConfigName, xdsClusterName string) *routeconfig.RouteConfiguration {
	return &routeconfig.RouteConfiguration{
		Name: routeConfigName,
		Vhds: &routeconfig.Vhds{
			ConfigSource: &coreconfig.ConfigSource{
				ResourceApiVersion: coreconfig.ApiVersion_V3,
				ConfigSourceSpecifier: &coreconfig.ConfigSource_ApiConfigSource{
					ApiConfigSource: &coreconfig.ApiConfigSource{
						TransportApiVersion:       coreconfig.ApiVersion_V3,
						ApiType:                   coreconfig.ApiConfigSource_DELTA_GRPC,
						SetNodeOnFirstMessageOnly: true,
						GrpcServices: []*coreconfig.GrpcService{{
							TargetSpecifier: &coreconfig.GrpcService_EnvoyGrpc_{
								EnvoyGrpc: &coreconfig.GrpcService_EnvoyGrpc{
									ClusterName: xdsClusterName,
								},
							},
						}},
					},
				},
			},
		},
		/*
		VirtualHosts: []*routeconfig.VirtualHost{{
			Name:    "local_service",
			Domains: []string{"*"},
			Routes: []*routeconfig.Route{{
				Match: &routeconfig.RouteMatch{
					PathSpecifier: &routeconfig.RouteMatch_Prefix{
						Prefix: "/",
					},
				},
				Action: &routeconfig.Route_Route{
					Route: &routeconfig.RouteAction{
						ClusterSpecifier: &routeconfig.RouteAction_Cluster{
							Cluster: clusterName,
						},
						HostRewriteSpecifier: &routeconfig.RouteAction_HostRewriteLiteral{
							HostRewriteLiteral: UpstreamHost,
						},
					},
				},
			}},
		}},
		*/
	}
}
