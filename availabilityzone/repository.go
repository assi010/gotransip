package availabilityzone

import (
	"github.com/assi010/gotransip/v6/repository"
	"github.com/assi010/gotransip/v6/rest"
)

// Repository should be used to query the AvailabilityZones you want to order in
// when for example ordering a Vps
type Repository repository.RestRepository

// GetAll returns a list of AvailabilityZones
func (r *Repository) GetAll() ([]AvailabilityZone, error) {
	var response availabilityZonesResponse
	avRequest := rest.Request{Endpoint: "/availability-zones"}
	err := r.Client.Get(avRequest, &response)

	return response.AvailabilityZones, err
}
