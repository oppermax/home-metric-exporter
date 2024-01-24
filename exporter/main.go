package main

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/oppermax/shelly-metric-exporter/pkg/deconz"
	"github.com/oppermax/shelly-metric-exporter/pkg/mappings"
	"github.com/oppermax/shelly-metric-exporter/pkg/shelly"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Spec struct {
	Host              string `envconfig:"HOST" required:"true"`
	Port              string `envconfig:"PORT" default:"7979"`
	DeconzHost        string `envconfig:"DECONZ_HOST" required:"true"`
	DeconzApiKey      string `envconfig:"DECONZ_API_KEY" required:"true"`
	ShellyUser        string `envconfig:"SHELLY_USER" required:"true"`
	ShellyPass        string `envconfig:"SHELLY_PASS" required:"true"`
	MappingsDirectory string `envconfig:"MAPPINGS_DIR" required:"false" default:"mappings/"`
}

func main() {
	var spec Spec
	err := envconfig.Process("home-metric-exporter", &spec)
	if err != nil {
		log.Panic().Err(err).Msg("could not process env")
	}

	log.Info().Msgf("starting home-metric-exporter on port %s", spec.Port)

	shellyMappings, err := mappings.GenerateMappings(mappings.PathConfig{ShellyMappingsFilepath: fmt.Sprintf("%s/shelly.json", spec.MappingsDirectory)})
	if err != nil {
		log.Error().Err(err).Msg("could not generate mappings")
	}

	deconzCollector := deconz.NewCollector(spec.DeconzHost, spec.DeconzApiKey)
	shellyPlugCollector := shelly.NewCollector(spec.ShellyUser, spec.ShellyPass, spec.Host, spec.Port, shellyMappings)

	humGauge, tempGauge := shelly.RegisterHTGauges()

	shellyWriter := shelly.MetricsWriter{Mappings: shellyMappings, HumGauge: humGauge, TempGauge: tempGauge}

	prometheus.MustRegister(deconzCollector, shellyPlugCollector, humGauge, tempGauge)

	http.HandleFunc("/shelly", shellyWriter.Write)
	http.HandleFunc("/shelly-plug-s/on", shelly.PlugSOn)
	http.HandleFunc("/shelly-plug-s/off", shelly.PlugSOff)
	http.Handle("/metrics", promhttp.Handler())

	err = http.ListenAndServe(fmt.Sprintf(":%s", spec.Port), nil)
	if err != nil {
		log.Panic().Err(err).Msg("could not listen and serve")
	}

	log.Info().Msg("collection successful")
}
