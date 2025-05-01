package weather

import (
	"weatherbot/internal/clients"
	"weatherbot/internal/repository"
)

type repo struct {
	db clients.DBClient
}

var _ repository.WeatherRepository = (*repo)(nil)

func NewRepository(client clients.DBClient) *repo {
	return &repo{
		db: client,
	}
}

// GetGeodata implements repository.WeatherRepository.
func (r *repo) GetGeodata() {
	panic("unimplemented")
}
