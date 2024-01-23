package shelly

import (
	"encoding/json"
	"fmt"
	"github.com/oppermax/shelly-metric-exporter/pkg/mappings"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

const (
	settingsUrl = "http://%s/settings/"
	statusUrl   = "http://%s/status"
)

type PlugStatus struct {
	Name   string
	Usages []float64
	Device Device
}

type Collector struct {
	// Possible metric descriptions.
	PowerLevels *prometheus.Desc

	Username string
	Password string
	Host     string
	Port     string
	mappings []mappings.ShellyMapping
}

// Describe implements prometheus.Collector.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	// Gather metadata about each metric.
	ch <- c.PowerLevels
}

func NewCollector(username, password, host, port string, mappings []mappings.ShellyMapping) *Collector {
	return &Collector{
		PowerLevels: prometheus.NewDesc(
			prometheus.BuildFQName("home_metric_exporter", "", "power_levels"),
			"Power consumption of Shelly plugs",
			[]string{"device_name", "room", "device_type"},
			nil,
		),
		Username: username,
		Password: password,
		Host:     host,
		Port:     port,
		mappings: mappings,
	}
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	log.Info().Msg("Starting to collecting Shelly plug metrics.")

	log.Info().Msg("Collecting metrics on light states.")

	statuses, err := c.GetPlugStatuses()
	if err != nil {
		// If an error occurs, send an invalid metric to notify
		// Prometheus of the problem.
		log.Error().Err(err).Msgf("cannot shelly plug statuses")
		ch <- prometheus.NewInvalidMetric(prometheus.NewDesc("failed getting plug statuses", "", []string{"device_name", "room", "device_type"}, nil), err)

		return
	}
	for _, status := range statuses {
		for _, usage := range status.Usages {
			ch <- prometheus.MustNewConstMetric(c.PowerLevels, prometheus.GaugeValue, usage, status.Name, status.Device.Room, "plug")
		}
	}

	log.Info().Msgf("Collected statuses of %d Shelly plugs", len(statuses))
}

func (c *Collector) GetPlugStatuses() ([]*PlugStatus, error) {
	out := []*PlugStatus{}

	for _, mapping := range c.mappings {
		if mapping.DeviceType == typePlug {
			settings, err := c.getPlugSettings(mapping.Ip)
			if err != nil {
				return nil, fmt.Errorf("could not get settings of shelly plug at %s: %w", mapping.Ip, err)
			}

			settings.Device.Room = mapping.Room

			status, err := c.getPlugStatus(settings)
			if err != nil {
				return nil, fmt.Errorf("could not get status of shelly plug at %s: %w", mapping.Ip, err)
			}

			out = append(out, status)
		}

	}

	return out, nil
}

func (c *Collector) getPlugSettings(deviceIp string) (*PlugSettingsResponse, error) {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(settingsUrl, deviceIp), nil)
	if err != nil {
		return nil, fmt.Errorf("could not get /settings of %s: %w", deviceIp, err)
	}

	req.SetBasicAuth(c.Username, c.Password)

	res, err := client.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request to get plug settings failed: %w", err)
	}

	settings := &PlugSettingsResponse{}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body: %w", err)
	}

	err = json.Unmarshal(data, settings)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal json response: %w", err)
	}

	settings.Device.Ip = deviceIp

	return settings, nil
}

func (c *Collector) getPlugStatus(settings *PlugSettingsResponse) (*PlugStatus, error) {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(statusUrl, settings.Device.Ip), nil)
	if err != nil {
		return nil, fmt.Errorf("could not get /status of %s: %w", settings.Device.Ip, err)
	}

	req.SetBasicAuth(c.Username, c.Password)

	res, err := client.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request to get plug settings failed: %w status: %d", err, res.StatusCode)
	}

	status := &PlugStatusResponse{}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body: %w", err)
	}

	err = json.Unmarshal(data, status)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal json response: %w", err)
	}

	usages := make([]float64, len(status.Meters))

	for i := 0; i < len(status.Meters); i++ {
		usages[i] = status.Meters[i].Power
	}

	out := &PlugStatus{
		Name:   mappings.GetPlugNameForHostname(settings.Device.Hostname, c.mappings),
		Usages: usages,
		Device: settings.Device,
	}

	return out, nil

}

func PlugSOn(w http.ResponseWriter, req *http.Request) {
	log.Info().Msg("plug turned on")
	fmt.Println(req)
}

func PlugSOff(w http.ResponseWriter, req *http.Request) {
	log.Info().Msg("plug turned off")
	fmt.Println(req)
}
