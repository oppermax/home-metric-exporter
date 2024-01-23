package mappings

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Mappings struct {
	ShellyMappings []ShellyMapping `json:"shelly_mappings"`
}

type ShellyMapping struct {
	DeviceId   string `json:"device_id,omitempty"`
	DeviceType string `json:"device_type"`
	Room       string `json:"room,omitempty"`
	Ip         string `json:"ip,omitempty"`
	DeviceName string `json:"device_name,omitempty"`
}

type PathConfig struct {
	ShellyMappingsFilepath string
}

func GenerateMappings(paths PathConfig) ([]ShellyMapping, error) {
	file, err := os.Open(paths.ShellyMappingsFilepath)
	if err != nil {
		return nil, fmt.Errorf("could not open file %s: %w", paths.ShellyMappingsFilepath, err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("could not read data in file %s: %w", paths.ShellyMappingsFilepath, err)
	}

	mappings := &Mappings{}

	err = json.Unmarshal(data, mappings)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal data in file %s: %w", paths.ShellyMappingsFilepath, err)
	}

	return mappings.ShellyMappings, nil
}

func GetPlugNameForHostname(id string, mappings []ShellyMapping) string {
	for _, mapping := range mappings {
		if mapping.DeviceId == id {
			return mapping.DeviceName
		}
	}

	return ""
}
