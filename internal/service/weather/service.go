package weather

import (
	"weatherbot/internal/repository"
	"weatherbot/internal/service"
)

type weatherService struct {
	r repository.WeatherRepository
}

var _ service.WeatherService = (*weatherService)(nil)

func NewService(r repository.WeatherRepository) *weatherService {
	return &weatherService{
		r: r,
	}
}

// GetWeather implements service.WeatherService.
func (w *weatherService) GetWeather() {
	panic("unimplemented")
}

func GetGeodata() {

}
