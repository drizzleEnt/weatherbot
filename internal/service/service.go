package service

import "context"

type WeatherService interface {
	GetWeather(ctx context.Context, city string, chatID int) error
}
