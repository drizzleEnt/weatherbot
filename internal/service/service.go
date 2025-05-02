package service

import (
	"context"
	"weatherbot/internal/domain"
)

type WeatherService interface {
	GetWeather(ctx context.Context, userInfoIncome domain.UserInfo) error
}
