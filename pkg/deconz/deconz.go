package deconz

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	"io"
	"math"
	"net/http"
)

const (
	allLightsUrl = "http://%s/api/%s/lights"
)

type LightState struct {
	Name  string
	State float64
	Level float64
}

type Collector struct {
	// Possible metric descriptions.
	LightStates *prometheus.Desc
	LightLevels *prometheus.Desc

	Host         string
	DeconzApiKey string
}

// Describe implements prometheus.Collector.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	// Gather metadata about each metric.
	ch <- c.LightStates
}

func NewCollector(host, apiKey string) *Collector {
	return &Collector{
		LightStates: prometheus.NewDesc(
			prometheus.BuildFQName("home_metric_exporter", "", "light_states"),
			"State of lights",
			[]string{"device_name", "device_type"},
			nil,
		),
		LightLevels: prometheus.NewDesc(
			prometheus.BuildFQName("home_metric_exporter", "", "light_levels"),
			"Level of lights",
			[]string{"device_name", "device_type"},
			nil,
		),
		Host:         host,
		DeconzApiKey: apiKey,
	}
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	log.Info().Msg("Starting to collecting metrics.")

	log.Info().Msg("Collecting metrics on light states.")

	states, err := c.getLightStates()
	if err != nil {
		// If an error occurs, send an invalid metric to notify
		// Prometheus of the problem.
		log.Error().Err(err).Msgf("cannot get light states")
		ch <- prometheus.NewInvalidMetric(prometheus.NewDesc("failed getting light states", "", []string{"device_name", "device_type"}, nil), err)

		return
	}
	for _, state := range states {
		ch <- prometheus.MustNewConstMetric(c.LightStates, prometheus.GaugeValue, state.State, state.Name, "light")
		ch <- prometheus.MustNewConstMetric(c.LightLevels, prometheus.GaugeValue, state.Level, state.Name, "light")
	}
}

func generateLights(data []byte) ([]Light, error) {
	response := map[string]Light{}
	err := json.Unmarshal(data, &response)
	if err != nil {
		log.Error().Err(err).Msg("could not unmarshal deconz /lights response")

		return nil, err
	}

	log.Info().Msgf("generated %d lights", len(response))

	out := []Light{}

	for _, val := range response {
		out = append(out, val)
	}

	return out, nil
}

func (c *Collector) getAllLights() ([]Light, error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf(allLightsUrl, c.Host, c.DeconzApiKey), nil)
	if err != nil {
		log.Error().Err(err).Msg("could not generate request")

		return nil, err
	}

	client := http.Client{}
	res, err := client.Do(request)
	log.Info().Msg("firing request GET /lights")
	if err != nil || res.StatusCode != http.StatusOK {
		//log.Error().Err(err).Msgf("request unsuccessful GET /lights %s", res.StatusCode)
		log.Error().Err(err).Msg("request unsuccessful GET /lights")

		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error().Err(err).Msg("could not read response body")

		return nil, err
	}

	return generateLights(body)
}

func (c *Collector) getLightStates() ([]LightState, error) {
	lights, err := c.getAllLights()
	if err != nil {
		log.Error().Err(err).Msg("could not get all lights")

		return nil, err
	}

	var out []LightState

	for _, light := range lights {
		out = append(out, LightState{
			Name:  light.Name,
			State: convertState(light.State.On),
			Level: convertLevel(light.State),
		})
	}

	log.Info().Msgf("got state of %d lights", len(out))

	return out, nil
}

func convertState(on bool) float64 {
	if on {
		return 1
	}

	return 0
}

func convertLevel(state State) float64 {
	if !state.On {
		return 0
	}

	return math.Round(float64((float64(state.Bri) / 254) * 100))
}
