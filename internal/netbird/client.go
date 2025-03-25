package netbird

import (
	"context"
	"fmt"
	"time"

	"github.com/netbirdio/netbird/client/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const (
	defaultSocketAddr = "unix:///var/run/netbird.sock"
	defaultTimeout    = 5 * time.Second
)

// Client represents a NetBird client
type Client struct {
	socketAddr string
}

// NewClient creates a new NetBird client
func NewClient(socketAddr string) *Client {
	if socketAddr == "" {
		socketAddr = defaultSocketAddr
	}

	return &Client{
		socketAddr: socketAddr,
	}
}

// GetPeers returns the list of peers from the NetBird daemon
func (c *Client) GetPeers(ctx context.Context) (*proto.StatusResponse, error) {
	conn, err := dialClientGRPCServer(ctx, c.socketAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to daemon: %w", err)
	}
	defer conn.Close()

	resp, err := proto.NewDaemonServiceClient(conn).Status(ctx, &proto.StatusRequest{GetFullPeerStatus: true})
	if err != nil {
		return nil, fmt.Errorf("status request failed: %v", status.Convert(err).Message())
	}

	return resp, nil
}

// dialClientGRPCServer connects to the daemon gRPC server
func dialClientGRPCServer(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	return grpc.DialContext(
		timeoutCtx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
}
