package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"

	"github.com/whyeasy/sonarcloud-exporter/lib/client"
)

//Collector struct for holding Prometheus Desc and Exporter Client
type Collector struct {
	up     *prometheus.Desc
	client *client.ExporterClient

	projectInfo *prometheus.Desc

	linesOfCode     *prometheus.Desc
	codeCoverage    *prometheus.Desc
	vulnerabilities *prometheus.Desc
	bugs            *prometheus.Desc
	codeSmells      *prometheus.Desc
}

//New creates a new Collecotor with Prometheus descriptors
func New(c *client.ExporterClient) *Collector {
	log.Info("Creating collector")
	return &Collector{
		up:     prometheus.NewDesc("sonarcloud_up", "Whether Sonarcloud scrape was successfull", nil, nil),
		client: c,

		projectInfo: prometheus.NewDesc("sonarcloud_project_info", "General information about projects", []string{"project_name", "project_qualifier", "project_key", "project_organization"}, nil),

		linesOfCode:     prometheus.NewDesc("sonarcloud_lines_of_code", "Lines of code within a project in SonarCloud", []string{"project_key"}, nil),
		codeCoverage:    prometheus.NewDesc("sonarcloud_code_coverage", "Code coverage within a project in SonarCloud", []string{"project_key"}, nil),
		vulnerabilities: prometheus.NewDesc("sonarcloud_vulnerabilities", "Amount of vulnerabilities within a project in SonarCloud", []string{"project_key"}, nil),
		bugs:            prometheus.NewDesc("sonarcloud_bugs", "Amount of bugs within a project in SonarCloud", []string{"project_key"}, nil),
		codeSmells:      prometheus.NewDesc("sonarcloud_code_smells", "Amount of code smells within a project in SonarCloud", []string{"project_key"}, nil),
	}
}

//Describe the metrics that are collected
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up

	ch <- c.projectInfo

	ch <- c.linesOfCode
	ch <- c.codeCoverage
	ch <- c.bugs
	ch <- c.vulnerabilities
	ch <- c.codeSmells
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	log.Info("Running scrape")

	if stats, err := c.client.GetStats(); err != nil {
		log.Error(err)
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)
	} else {
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 1)

		collectProjectInfo(c, ch, stats)

		collectMeasurements(c, ch, stats)

		log.Info("Scrape Complete")
	}
}

func collectProjectInfo(c *Collector, ch chan<- prometheus.Metric, stats *client.Stats) {
	for _, project := range *stats.Projects {
		ch <- prometheus.MustNewConstMetric(c.projectInfo, prometheus.GaugeValue, 1, project.Name, project.Qualifier, project.Key, project.Organization)
	}
}

func collectMeasurements(c *Collector, ch chan<- prometheus.Metric, stats *client.Stats) {
	for _, measurement := range *stats.Measurements {
		value, err := strconv.ParseFloat(measurement.Value, 64)
		if err != nil {
			log.Error(err)
		}
		switch {
		case measurement.Metric == "ncloc":
			ch <- prometheus.MustNewConstMetric(c.linesOfCode, prometheus.GaugeValue, value, measurement.Key)
		case measurement.Metric == "coverage":
			ch <- prometheus.MustNewConstMetric(c.codeCoverage, prometheus.GaugeValue, value, measurement.Key)
		case measurement.Metric == "vulnerabilities":
			ch <- prometheus.MustNewConstMetric(c.vulnerabilities, prometheus.GaugeValue, value, measurement.Key)
		case measurement.Metric == "bugs":
			ch <- prometheus.MustNewConstMetric(c.bugs, prometheus.GaugeValue, value, measurement.Key)
		case measurement.Metric == "violations":
			ch <- prometheus.MustNewConstMetric(c.codeSmells, prometheus.GaugeValue, value, measurement.Key)
		}
	}
}
