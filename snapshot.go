package main

import (
	"time"

	"github.com/golang/protobuf/ptypes"

	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
)

func GetSnapshot(clusterName string, localPort int) (cache.Snapshot, error) {
	version := "1"
	snapshot := cache.NewSnapshot(
		version,
		nil, // no endpoints, cluster has them
		[]types.Resource{
			getCluster(clusterName, localPort),
		},
		nil, // no routes
		nil, // no listeners
		nil, // no runtimes
		nil, // no secrets
	)
	if err := snapshot.Consistent(); err != nil {
		return cache.Snapshot{}, err
	}
	return snapshot, nil
}

func getCluster(clusterName string, localPort int) *cluster.Cluster {
	return &cluster.Cluster{
		Name:           clusterName,
		ConnectTimeout: ptypes.DurationProto(5 * time.Second),
		ClusterDiscoveryType: &cluster.Cluster_Type{
			Type: cluster.Cluster_LOGICAL_DNS,
		},
		LbPolicy:        cluster.Cluster_ROUND_ROBIN,
		LoadAssignment:  getEndpoint(clusterName, localPort),
		DnsLookupFamily: cluster.Cluster_V4_ONLY,
	}
}

func getEndpoint(clusterName string, localPort int) *endpoint.ClusterLoadAssignment {
	return &endpoint.ClusterLoadAssignment{
		ClusterName: clusterName,
		Endpoints: []*endpoint.LocalityLbEndpoints{{
			LbEndpoints: []*endpoint.LbEndpoint{{
				HostIdentifier: &endpoint.LbEndpoint_Endpoint{
					Endpoint: &endpoint.Endpoint{
						Address: &core.Address{
							Address: &core.Address_SocketAddress{
								SocketAddress: &core.SocketAddress{
									Protocol: core.SocketAddress_TCP,
									Address:  "127.0.0.1",
									PortSpecifier: &core.SocketAddress_PortValue{
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
