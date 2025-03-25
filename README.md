# NetBird Exporter

A Prometheus exporter for [NetBird](https://netbird.io/) peer metrics.

## Overview

NetBird Exporter connects to the NetBird daemon via a Unix socket and exports various peer connection metrics in Prometheus format. It's designed to help monitor the health and performance of your NetBird network.

## Features

- Exposes NetBird peer connection status
- Tracks connection latency between peers
- Monitors data transfer (bytes sent/received)
- Reports connection types (P2P or relay)
- Provides peer identification via public keys

## Prerequisites

- NetBird daemon running with socket at `/var/run/netbird.sock`
- Linux OS (currently the only supported platform)

## Installation

### Using Go

```bash
go install github.com/gocloudio/netbird-exporter/cmd/netbird-exporter@latest
```

### Using Docker

```bash
docker run -p 9101:9101 -v /var/run/netbird.sock:/var/run/netbird.sock gocloudio/netbird-exporter
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
- `--netbird.socket`: Path to NetBird Unix socket (default: `/var/run/netbird.sock`)

## Available Metrics

| Metric Name | Description | Labels |
|-------------|-------------|--------|
| `netbird_connection` | Connection status (1=connected, 0=disconnected) | `local`, `remote`, `status`, `name`, `connection_type`, `local_public_key`, `remote_public_key` |
| `netbird_connection_latency` | Connection latency in milliseconds | `local`, `remote` |
| `netbird_connection_received_bytes` | Bytes received from peer | `local`, `remote` |
| `netbird_connection_sent_bytes` | Bytes sent to peer | `local`, `remote` |

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

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 