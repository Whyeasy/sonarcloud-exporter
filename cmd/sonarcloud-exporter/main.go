package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"github.com/whyeasy/sonarcloud-exporter/internal"
	"github.com/whyeasy/sonarcloud-exporter/lib/client"
	"github.com/whyeasy/sonarcloud-exporter/lib/collector"
)

var config internal.Config

func init() {
	flag.StringVar(&config.Token, "scToken", os.Getenv("SC_TOKEN"), "Token to access SonarCloud API")
	flag.StringVar(&config.ListenAddress, "listenAddress", os.Getenv("LISTEN_ADDRESS"), "Port address of exporter to run on")
	flag.StringVar(&config.ListenPath, "listenPath", os.Getenv("LISTEN_PATH"), "Path where metrics will be exposed")
	flag.StringVar(&config.Organization, "organization", os.Getenv("SC_ORGANIZATION"), "Organization to query within SonarCloud")
}

func main() {
	if err := parseConfig(); err != nil {
		log.Error(err)
		flag.Usage()
		os.Exit(2)
	}

	log.Info("Starting SonarCloud Exporter")

	client := client.New(config)
	coll := collector.New(client)

	prometheus.MustRegister(coll)

	log.Info("Start serving metrics")

	http.Handle(config.ListenPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`<html>
			<head><title>SonarCloud Exporter</title></head>
			<body>
			<h1>SonarCloud Exporter</h1>
			<p><a href="` + config.ListenPath + `">Metrics</a></p>
			</body>
			</html>`))
		if err != nil {
			log.Error(err)
		}
	})
	log.Fatal(http.ListenAndServe(":"+config.ListenAddress, nil))
}

func parseConfig() error {
	flag.Parse()
	required := []string{"scToken"}
	var err error
	flag.VisitAll(func(f *flag.Flag) {
		for _, r := range required {
			if r == f.Name && (f.Value.String() == "" || f.Value.String() == "0") {
				err = fmt.Errorf("%v is empty", f.Usage)
			}
		}
		if f.Name == "listenAddress" && (f.Value.String() == "" || f.Value.String() == "0") {
			err = f.Value.Set("8080")
			if err != nil {
				log.Error(err)
			}
		}
		if f.Name == "listenPath" && (f.Value.String() == "" || f.Value.String() == "0") {
			err = f.Value.Set("/metrics")
			if err != nil {
				log.Error(err)
			}
		}

	})
	return err
}
