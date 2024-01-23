package shelly

import (
	"github.com/oppermax/shelly-metric-exporter/pkg/mappings"
	"github.com/prometheus/client_golang/prometheus"
)

type MetricsWriter struct {
	Mappings  []mappings.ShellyMapping
	HumGauge  *prometheus.GaugeVec
	TempGauge *prometheus.GaugeVec
}

type PlugSettingsResponse struct {
	Device Device `json:"device"`
}

type Device struct {
	Type     string `json:"type,omitempty"`
	Mac      string `json:"mac,omitempty"`
	Hostname string `json:"hostname"`
	Ip       string
	Room     string
}

type PlugStatusResponse struct {
	WifiSta         WifiSta      `json:"wifi_sta,omitempty"`
	Cloud           Cloud        `json:"cloud,omitempty"`
	Mqtt            Mqtt         `json:"mqtt,omitempty"`
	Time            string       `json:"time,omitempty"`
	Unixtime        int          `json:"unixtime,omitempty"`
	Serial          int          `json:"serial,omitempty"`
	HasUpdate       bool         `json:"has_update,omitempty"`
	Mac             string       `json:"mac,omitempty"`
	CfgChangedCnt   int          `json:"cfg_changed_cnt,omitempty"`
	ActionsStats    ActionsStats `json:"actions_stats,omitempty"`
	Relays          []Relays     `json:"relays,omitempty"`
	Meters          []Meters     `json:"meters,omitempty"`
	Temperature     float64      `json:"temperature,omitempty"`
	Overtemperature bool         `json:"overtemperature,omitempty"`
	Tmp             Tmp          `json:"tmp,omitempty"`
	Update          Update       `json:"update,omitempty"`
	RAMTotal        int          `json:"ram_total,omitempty"`
	RAMFree         int          `json:"ram_free,omitempty"`
	FsSize          int          `json:"fs_size,omitempty"`
	FsFree          int          `json:"fs_free,omitempty"`
	Uptime          int          `json:"uptime,omitempty"`
}

type WifiSta struct {
	Connected bool   `json:"connected,omitempty"`
	Ssid      string `json:"ssid,omitempty"`
	IP        string `json:"ip,omitempty"`
	Rssi      int    `json:"rssi,omitempty"`
}

type Cloud struct {
	Enabled   bool `json:"enabled,omitempty"`
	Connected bool `json:"connected,omitempty"`
}

type Mqtt struct {
	Connected bool `json:"connected,omitempty"`
}

type ActionsStats struct {
	Skipped int `json:"skipped,omitempty"`
}

type Relays struct {
	Ison           bool   `json:"ison,omitempty"`
	HasTimer       bool   `json:"has_timer,omitempty"`
	TimerStarted   int    `json:"timer_started,omitempty"`
	TimerDuration  int    `json:"timer_duration,omitempty"`
	TimerRemaining int    `json:"timer_remaining,omitempty"`
	Overpower      bool   `json:"overpower,omitempty"`
	Source         string `json:"source,omitempty"`
}

type Meters struct {
	Power     float64   `json:"power,omitempty"`
	Overpower float64   `json:"overpower,omitempty"`
	IsValid   bool      `json:"is_valid,omitempty"`
	Timestamp int       `json:"timestamp,omitempty"`
	Counters  []float64 `json:"counters,omitempty"`
	Total     int       `json:"total,omitempty"`
}

type Tmp struct {
	TC      float64 `json:"tC,omitempty"`
	TF      float64 `json:"tF,omitempty"`
	IsValid bool    `json:"is_valid,omitempty"`
}

type Update struct {
	Status      string `json:"status,omitempty"`
	HasUpdate   bool   `json:"has_update,omitempty"`
	NewVersion  string `json:"new_version,omitempty"`
	OldVersion  string `json:"old_version,omitempty"`
	BetaVersion string `json:"beta_version,omitempty"`
}
