# NetBird Exporter

A Prometheus exporter for [NetBird](https://netbird.io/) peer metrics.

## Overview

NetBird Exporter communicates with the NetBird daemon via gRPC and exports various peer connection metrics in Prometheus format. It's designed to help monitor the health and performance of your NetBird network.

## Features

- Exposes NetBird peer connection status
- Tracks connection latency between peers
- Monitors data transfer (bytes sent/received)
- Reports connection types (P2P or relay)
- Supports custom FQDN domain suffix filtering

## Prerequisites

- NetBird daemon running with socket at `/var/run/netbird.sock`
- Linux OS (currently the only supported platform)

## Installation

### Using Go

```bash
go install github.com/gocloudio/netbird-exporter/cmd/netbird-exporter@latest
```

### Using Docker from GitHub Container Registry

```bash
docker pull ghcr.io/gocloudio/netbird-exporter:latest
docker run -p 9101:9101 -v /var/run/netbird.sock:/var/run/netbird.sock ghcr.io/gocloudio/netbird-exporter:latest
```

### From Source

```bash
git clone https://github.com/gocloudio/netbird-exporter.git
cd netbird-exporter
go build -o netbird-exporter ./cmd/netbird-exporter
./netbird-exporter
```

## Usage

```bash
./netbird-exporter [flags]
```

### Flags

- `--web.listen-address`: Address to listen on for web interface and telemetry (default: `:9101`)
- `--web.telemetry-path`: Path under which to expose metrics (default: `/metrics`)
- `--netbird.socket`: Address of NetBird daemon socket (default: `unix:///var/run/netbird.sock`)
- `--ignore-fqdn-domain`: Domain suffix to remove from FQDNs (e.g. `.example.com`)

## Available Metrics

| Metric Name | Description | Labels |
|-------------|-------------|--------|
| `netbird_connection` | Connection status (1=connected, 0=disconnected) | `local`, `remote`, `status`, `connection_type`, `local_ip`, `remote_ip` |
| `netbird_connection_latency` | Connection latency in milliseconds | `local`, `remote` |
| `netbird_connection_received_bytes` | Bytes received from peer | `local`, `remote` |
| `netbird_connection_sent_bytes` | Bytes sent to peer | `local`, `remote` |

### Label Details

- `local` and `remote`: FQDN of the local and remote peers (with domain suffix optionally removed)
- `status`: Status of the connection (e.g., "connected", "disconnected")
- `connection_type`: Either "p2p" for direct connections or "relay" for relayed connections
- `local_ip` and `remote_ip`: IP addresses of the local and remote peers (CIDR notation removed if present)

## Prometheus Configuration

Add the following to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'netbird'
    static_configs:
      - targets: ['localhost:9101']
```

## Grafana Dashboard

A sample Grafana dashboard is available in the `dashboards` directory.

## Docker Image Tags

The Docker images are available from GitHub Container Registry with the following tag patterns:

- `ghcr.io/gocloudio/netbird-exporter:latest` - Latest stable build
- `ghcr.io/gocloudio/netbird-exporter:v1.0.0` - Specific version
- `ghcr.io/gocloudio/netbird-exporter:1.0` - Major.Minor version

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 