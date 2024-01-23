package shelly

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
	"strings"
)

const (
	name        = "id"
	temperature = "temp"
	humidity    = "hum"
	typeHT      = "shellyht"
	typePlug    = "shellyplug"
)

type HtReading struct {
	Name        string
	Room        string
	Temperature float64
	Humidity    float64
	Type        string
}

func RegisterHTGauges() (*prometheus.GaugeVec, *prometheus.GaugeVec) {
	var tempGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   "home_metric_exporter",
			Name:        "temperature",
			ConstLabels: nil,
		},
		[]string{"device_name", "room", "device_type"},
	)

	var humGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   "home_metric_exporter",
			Name:        "humidity",
			ConstLabels: nil,
		},
		[]string{"device_name", "room", "device_type"},
	)

	return tempGauge, humGauge
}

func (s *MetricsWriter) Write(w http.ResponseWriter, req *http.Request) {
	reading, err := GenerateReading(req)
	if err != nil {
		log.Panic().Err(err).Msg("could not generate reading")
	}

	room := "default"

	for _, mapping := range s.Mappings {
		if mapping.DeviceId == reading.Name {
			room = mapping.Room
		}
	}

	log.Info().Fields("Temperature").Fields("Humidity").Fields("Name").Msg("got reading from shelly ht")

	s.TempGauge.WithLabelValues(reading.Name, room, reading.Type).Set(reading.Temperature)
	s.HumGauge.WithLabelValues(reading.Name, room, reading.Type).Set(reading.Humidity)
}

func GenerateReading(r *http.Request) (*HtReading, error) {
	name := r.URL.Query().Get(name)
	hum, err := strconv.ParseFloat(r.URL.Query().Get(humidity), 64)
	if err != nil {
		log.Error().Err(err).Msgf("could not cast Humidity into int. device: %s", name)
		return nil, err
	}
	temp, err := strconv.ParseFloat(r.URL.Query().Get(temperature), 64)
	if err != nil {
		log.Error().Err(err).Msgf("could not cast Humidity into int. device: %s", name)
		return nil, err
	}
	return &HtReading{
		Name:        name,
		Temperature: temp,
		Humidity:    hum,
		Type:        determineDeviceType(name),
	}, nil
}

func determineDeviceType(name string) string {
	if strings.HasPrefix(name, typeHT) {
		return typeHT
	} else if strings.HasPrefix(name, typePlug) {
		return typePlug
	} else {
		return name
	}
}
