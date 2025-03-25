package collector

import (
	"context"
	"testing"
	"time"

	"github.com/netbirdio/netbird/client/proto"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/protobuf/types/known/durationpb"
)

// MockNetBirdClient implements NetBirdClient interface for testing
type MockNetBirdClient struct {
	resp *proto.StatusResponse
	err  error
}

func (m *MockNetBirdClient) GetPeers(ctx context.Context) (*proto.StatusResponse, error) {
	return m.resp, m.err
}

func TestNetBirdCollector(t *testing.T) {
	// Create mock response
	mockResp := &proto.StatusResponse{
		Status:        "running",
		DaemonVersion: "0.0.1",
		FullStatus: &proto.FullStatus{
			LocalPeerState: &proto.LocalPeerState{
				IP:     "100.64.0.1",
				PubKey: "self-key",
				Fqdn:   "self.netbird.io",
			},
			Peers: []*proto.PeerState{
				{
					IP:         "100.64.0.1",
					PubKey:     "self-key",
					ConnStatus: "connected",
					Fqdn:       "self.netbird.io",
				},
				{
					IP:         "100.64.0.2",
					PubKey:     "peer-key",
					ConnStatus: "connected",
					Fqdn:       "peer.netbird.io",
					BytesRx:    3000,
					BytesTx:    2000,
					Latency:    durationpb.New(100 * time.Millisecond),
				},
			},
		},
	}

	// Create mock client
	mockClient := &MockNetBirdClient{
		resp: mockResp,
		err:  nil,
	}

	// Create collector with mock client
	collector := NewNetBirdCollector(mockClient, "")

	// Test Describe method
	ch := make(chan *prometheus.Desc, 4)
	collector.Describe(ch)

	if len(ch) != 4 {
		t.Errorf("Expected 4 metric descriptions, got %d", len(ch))
	}

	// Register the collector in a fresh registry
	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)
}
