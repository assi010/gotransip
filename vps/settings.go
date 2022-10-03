package vps

import (
	"fmt"

	"github.com/assi010/gotransip/v6/repository"
	"github.com/assi010/gotransip/v6/rest"
)

// SettingRepository allows you to get and modify settings a VPS. These include
// the `blockVpsMailPorts` and `tcpMonitoringAvailable` settings
type SettingRepository repository.RestRepository

const (
	// SettingBlockVPSMailPorts describes if the mail ports (25,465,465) are blocked for a VPS
	SettingBlockVPSMailPorts = "blockVpsMailPorts"
	// SettingTCPMonitoringAvailable is true when the TCP monitoring feature is enabled
	SettingTCPMonitoringAvailable = "tcpMonitoringAvailable"
)

// These constants define which data types can be returned by the api
const (
	SettingDataTypeString  = "string"
	SettingDataTypeBoolean = "boolean"
)

// Setting is a struct that describes a vps setting
type Setting struct {
	Name     string       `json:"name"`
	DataType string       `json:"dataType"`
	ReadOnly bool         `json:"readOnly"`
	Value    SettingValue `json:"value"`
}

// SettingValue contains the value of a setting. Only one the fields will contain the
// value. Which field that is can be determined by checking the DataType of a Setting.
type SettingValue struct {
	ValueBoolean bool   `json:"valueBoolean"`
	ValueString  string `json:"valueString"`
}

// settingsWrapper struct contains Settings in it,
// this is solely used for unmarshalling
type settingsWrapper struct {
	Settings []Setting `json:"settings"`
}

// settingWrapper struct contains Settings in it,
// this is solely used for unmarshalling
type settingWrapper struct {
	Setting Setting `json:"setting"`
}

// GetAll returns all the Settings for a vps
func (r *SettingRepository) GetAll(vpsName string) ([]Setting, error) {
	var response settingsWrapper
	restRequest := rest.Request{Endpoint: fmt.Sprintf("/vps/%s/settings", vpsName)}
	err := r.Client.Get(restRequest, &response)

	return response.Settings, err
}

// GetByName returns a setting by name
func (r *SettingRepository) GetByName(vpsName string, settingName string) (Setting, error) {
	var response settingWrapper
	restRequest := rest.Request{Endpoint: fmt.Sprintf("/vps/%s/settings/%s", vpsName, settingName)}
	err := r.Client.Get(restRequest, &response)
	return response.Setting, err
}

// Update updates a setting for a vps
func (r *SettingRepository) Update(vpsName string, setting Setting) error {
	requestBody := settingWrapper{setting}
	restRequest := rest.Request{Endpoint: fmt.Sprintf("/vps/%s/settings/%s", vpsName, setting.Name), Body: requestBody}
	return r.Client.Put(restRequest)
}
