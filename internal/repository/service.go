package repository

import "context"

type WeatherRepository interface {
	GetGeodata()
	GetUserCity(ctx context.Context, chatID int) (string, error)
}