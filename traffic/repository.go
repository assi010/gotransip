package traffic

import (
	"fmt"
	"github.com/assi010/gotransip/v6/repository"
	"github.com/assi010/gotransip/v6/rest"
)

// Repository allows you to get information about your usage in your traffic pool
// you can retrieve this information globally or per vps
type Repository repository.RestRepository

// GetTrafficPool returns all the traffic of your VPSes combined, overusage will also be billed based on this information
func (r *Repository) GetTrafficPool() (Information, error) {
	var response wrapper
	restRequest := rest.Request{Endpoint: "/traffic"}
	err := r.Client.Get(restRequest, &response)

	return response.TrafficInformation, err
}

// GetTrafficInformationForVps allows you to get specific traffic information for a given VPS
func (r *Repository) GetTrafficInformationForVps(vpsName string) (Information, error) {
	var response wrapper
	restRequest := rest.Request{Endpoint: fmt.Sprintf("/traffic/%s", vpsName)}
	err := r.Client.Get(restRequest, &response)

	return response.TrafficInformation, err
}
