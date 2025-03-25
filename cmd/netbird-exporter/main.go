package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gocloudio/netbird-exporter/internal/collector"
	"github.com/gocloudio/netbird-exporter/internal/netbird"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	listenAddress    = flag.String("web.listen-address", ":9101", "Address to listen on for web interface and telemetry")
	metricsPath      = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics")
	socketAddr       = flag.String("netbird.socket", "unix:///var/run/netbird.sock", "Address of NetBird daemon socket")
	ignoreFQDNDomain = flag.String("ignore-fqdn-domain", "", "Domain suffix to remove from FQDNs (e.g. '.example.com')")
)

func main() {
	flag.Parse()

	// Create NetBird client
	client := netbird.NewClient(*socketAddr)

	// Create and register collector
	collector := collector.NewNetBirdCollector(client, *ignoreFQDNDomain)
	prometheus.MustRegister(collector)

	// Setup HTTP server
	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>NetBird Exporter</title></head>
			<body>
			<h1>NetBird Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	// Start server in a goroutine
	go func() {
		log.Printf("Starting NetBird exporter on %s", *listenAddress)
		if err := http.ListenAndServe(*listenAddress, nil); err != nil {
			log.Fatalf("Error starting HTTP server: %s", err)
		}
	}()

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down NetBird exporter")
}
