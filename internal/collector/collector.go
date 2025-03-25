package collector

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/netbirdio/netbird/client/proto"
	"github.com/prometheus/client_golang/prometheus"
)

// NetBirdClient defines the interface for the NetBird client
type NetBirdClient interface {
	GetPeers(ctx context.Context) (*proto.StatusResponse, error)
}

// NetBirdCollector implements the prometheus.Collector interface.
type NetBirdCollector struct {
	client           NetBirdClient
	ignoreFQDNDomain string

	connectionDesc          *prometheus.Desc
	connectionLatencyDesc   *prometheus.Desc
	connectionRecvBytesDesc *prometheus.Desc
	connectionSentBytesDesc *prometheus.Desc
}

// NewNetBirdCollector creates a new NetBirdCollector
func NewNetBirdCollector(client NetBirdClient, ignoreFQDNDomain string) *NetBirdCollector {
	return &NetBirdCollector{
		client:           client,
		ignoreFQDNDomain: ignoreFQDNDomain,
		connectionDesc: prometheus.NewDesc(
			"netbird_connection",
			"Status of NetBird connection",
			[]string{"local", "remote", "status", "connection_type", "local_ip", "remote_ip"},
			nil,
		),
		connectionLatencyDesc: prometheus.NewDesc(
			"netbird_connection_latency",
			"Latency of NetBird connection in milliseconds",
			[]string{"local", "remote"},
			nil,
		),
		connectionRecvBytesDesc: prometheus.NewDesc(
			"netbird_connection_received_bytes",
			"Bytes received through NetBird connection",
			[]string{"local", "remote"},
			nil,
		),
		connectionSentBytesDesc: prometheus.NewDesc(
			"netbird_connection_sent_bytes",
			"Bytes sent through NetBird connection",
			[]string{"local", "remote"},
			nil,
		),
	}
}

// Describe implements the prometheus.Collector interface.
func (c *NetBirdCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.connectionDesc
	ch <- c.connectionLatencyDesc
	ch <- c.connectionRecvBytesDesc
	ch <- c.connectionSentBytesDesc
}

// Collect implements the prometheus.Collector interface.
// This is called on each Prometheus scrape, creating fresh metrics
// with the latest values from the NetBird daemon.
func (c *NetBirdCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	status, err := c.client.GetPeers(ctx)
	if err != nil {
		log.Printf("Error getting peers: %v", err)
		return
	}

	if status.FullStatus == nil {
		log.Printf("No full status available")
		return
	}

	// Get local peer info
	local := status.FullStatus.LocalPeerState
	if local == nil {
		log.Printf("No local peer state available")
		return
	}

	// Process peer connections
	for _, peer := range status.FullStatus.Peers {
		// Skip if it's the local peer
		if peer.PubKey == local.PubKey {
			continue
		}

		// Process FQDNs to remove domain suffix if needed
		localFQDN := ProcessFQDN(local.Fqdn, c.ignoreFQDNDomain)
		remoteFQDN := ProcessFQDN(peer.Fqdn, c.ignoreFQDNDomain)

		// Clean IP addresses by removing CIDR notation if present
		localIP := CleanIPAddress(local.IP)
		remoteIP := CleanIPAddress(peer.IP)

		// Determine connection status (1 = connected, 0 = disconnected)
		connectionStatus := 0.0
		if strings.ToLower(peer.ConnStatus) == "connected" {
			connectionStatus = 1.0
		}

		// Determine connection type (p2p or relay)
		connType := "p2p"
		if peer.Relayed {
			connType = "relay"
		}

		ch <- prometheus.MustNewConstMetric(
			c.connectionDesc,
			prometheus.GaugeValue,
			connectionStatus,
			localFQDN, remoteFQDN, peer.ConnStatus, connType, localIP, remoteIP,
		)

		// Only report metrics for connected peers
		if strings.ToLower(peer.ConnStatus) == "connected" {
			var latencyMs int64
			if peer.Latency != nil {
				latencyMs = peer.Latency.AsDuration().Milliseconds()
			}

			ch <- prometheus.MustNewConstMetric(
				c.connectionLatencyDesc,
				prometheus.GaugeValue,
				float64(latencyMs),
				localFQDN, remoteFQDN,
			)

			ch <- prometheus.MustNewConstMetric(
				c.connectionRecvBytesDesc,
				prometheus.CounterValue,
				float64(peer.BytesRx),
				localFQDN, remoteFQDN,
			)

			ch <- prometheus.MustNewConstMetric(
				c.connectionSentBytesDesc,
				prometheus.CounterValue,
				float64(peer.BytesTx),
				localFQDN, remoteFQDN,
			)
		}
	}
}
