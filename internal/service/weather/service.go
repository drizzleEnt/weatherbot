package weather

import (
	"context"
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
func (w *weatherService) GetWeather(ctx context.Context, city string, chatID int) error {
	if city == "" {
		mainCity, err := w.r.GetUserCity(ctx, chatID)
		if err != nil {
			return err
		}
		city = mainCity
	}

	w.r.GetGeodata()

	panic("unimplemented")
}

func GetGeodata() {

}
